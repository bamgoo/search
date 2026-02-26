package search

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/bamgoo/base"
)

func BuildQuery(keyword string, args ...Any) Query {
	q := Query{Keyword: strings.TrimSpace(keyword), Offset: 0, Limit: 20, Raw: Map{}, Setting: Map{}}
	for _, arg := range args {
		switch v := arg.(type) {
		case Query:
			q = mergeQuery(q, v)
		case *Query:
			if v != nil {
				q = mergeQuery(q, *v)
			}
		case Map:
			q = mergeQueryMap(q, v)
		}
	}
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	return q
}

func mergeQuery(dst Query, src Query) Query {
	if strings.TrimSpace(src.Keyword) != "" {
		dst.Keyword = strings.TrimSpace(src.Keyword)
	}
	if src.Offset >= 0 {
		dst.Offset = src.Offset
	}
	if src.Limit > 0 {
		dst.Limit = src.Limit
	}
	if len(src.Filters) > 0 {
		dst.Filters = append(dst.Filters, src.Filters...)
	}
	if len(src.Sorts) > 0 {
		dst.Sorts = append([]Sort{}, src.Sorts...)
	}
	if len(src.Fields) > 0 {
		dst.Fields = append([]string{}, src.Fields...)
	}
	if len(src.Facets) > 0 {
		dst.Facets = append([]string{}, src.Facets...)
	}
	if len(src.Highlight) > 0 {
		dst.Highlight = append([]string{}, src.Highlight...)
	}
	if src.Raw != nil {
		dst.Raw = mergeMaps(dst.Raw, src.Raw)
	}
	if src.Setting != nil {
		dst.Setting = mergeMaps(dst.Setting, src.Setting)
	}
	return dst
}

func mergeQueryMap(dst Query, cfg Map) Query {
	if cfg == nil {
		return dst
	}
	if v, ok := pickString(cfg, "$keyword", "$q", "keyword", "q"); ok && strings.TrimSpace(v) != "" {
		dst.Keyword = strings.TrimSpace(v)
	}
	if v, ok := toInt(pickValue(cfg, "$offset", "offset")); ok {
		dst.Offset = v
	}
	if v, ok := toInt(pickValue(cfg, "$limit", "limit")); ok {
		dst.Limit = v
	}
	if v, ok := pickValueOK(cfg, "$fields", "$select", "fields", "select"); ok {
		dst.Fields = toStrings(v)
	}
	if v, ok := pickValueOK(cfg, "$facets", "facets"); ok {
		dst.Facets = toStrings(v)
	}
	if v, ok := pickValueOK(cfg, "$highlight", "highlight"); ok {
		dst.Highlight = toStrings(v)
	}
	if v, ok := pickMap(cfg, "$setting", "setting"); ok {
		dst.Setting = mergeMaps(dst.Setting, v)
	}
	if v, ok := pickMap(cfg, "$raw", "raw"); ok {
		dst.Raw = mergeMaps(dst.Raw, v)
	}

	if vv, ok := pickValueOK(cfg, "$sort", "sort", "sorts", "$sorts"); ok {
		dst.Sorts = parseSorts(vv)
	}

	if vv, ok := pickValueOK(cfg, "$filters", "$filter", "filters", "filter"); ok {
		dst.Filters = append(dst.Filters, parseFilters(vv)...)
	}

	// data-style top-level filters:
	// Map{"category":"tech", "$sort": Map{"score": DESC}}
	dst.Filters = append(dst.Filters, parseTopLevelFilters(cfg)...)

	return dst
}

func parseSorts(v Any) []Sort {
	out := make([]Sort, 0)
	switch vv := v.(type) {
	case Map:
		for field, dirVal := range vv {
			if strings.TrimSpace(field) == "" {
				continue
			}
			desc, ok := parseSortDirection(dirVal)
			if !ok {
				desc = false
			}
			out = append(out, Sort{Field: field, Desc: desc})
		}
	case string:
		for _, one := range toStrings(vv) {
			one = strings.TrimSpace(one)
			if one == "" {
				continue
			}
			desc := false
			if strings.HasPrefix(one, "-") {
				desc = true
				one = strings.TrimPrefix(one, "-")
			}
			out = append(out, Sort{Field: one, Desc: desc})
		}
	case []string:
		for _, one := range vv {
			out = append(out, parseSorts(one)...)
		}
	case []Any:
		for _, one := range vv {
			out = append(out, parseSorts(one)...)
		}
	case []Map:
		for _, one := range vv {
			for field, dirVal := range one {
				if strings.TrimSpace(field) == "" {
					continue
				}
				desc, ok := parseSortDirection(dirVal)
				if !ok {
					desc = false
				}
				out = append(out, Sort{Field: field, Desc: desc})
			}
		}
	}
	return out
}

func parseFilters(v Any) []Filter {
	out := make([]Filter, 0)
	switch vv := v.(type) {
	case Map:
		for key, val := range vv {
			out = append(out, parseFieldFilters(key, val)...)
		}
	case []Map:
		for _, one := range vv {
			for field, val := range one {
				out = append(out, parseFieldFilters(field, val)...)
			}
		}
	case []Any:
		for _, one := range vv {
			out = append(out, parseFilters(one)...)
		}
	}
	return out
}

func parseFieldFilters(field string, val Any) []Filter {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil
	}

	if mv, ok := val.(Map); ok {
		// support data-style operators:
		// {"score":{"$gt":100,"$lt":500}} or {"category":{"$eq":"123"}}
		out := make([]Filter, 0)
		handled := false
		for opKey, opVal := range mv {
			op := normalizeFilterOp(opKey)
			switch op {
			case "eq", "ne", "gt", "gte", "lt", "lte":
				out = append(out, Filter{Field: field, Op: op, Value: opVal})
				handled = true
			case "in", "nin":
				out = append(out, Filter{Field: field, Op: op, Values: toAnys(opVal)})
				handled = true
			case "range":
				if rv, ok := opVal.(Map); ok {
					out = append(out, Filter{Field: field, Op: "range", Min: rv["min"], Max: rv["max"]})
				} else {
					out = append(out, Filter{Field: field, Op: "range", Min: mv["min"], Max: mv["max"]})
				}
				handled = true
			case "op", "value", "values", "min", "max":
				// handled below
			}
		}
		if handled {
			return out
		}
		// map without supported operators is treated as eq value
		return []Filter{{Field: field, Op: "eq", Value: val}}
	}

	return []Filter{{Field: field, Op: "eq", Value: val}}
}

func parseSortDirection(v Any) (bool, bool) {
	switch vv := v.(type) {
	case int:
		if vv == DESC {
			return true, true
		}
		if vv == ASC {
			return false, true
		}
		return vv < 0, true
	case int64:
		if vv == int64(DESC) {
			return true, true
		}
		if vv == int64(ASC) {
			return false, true
		}
		return vv < 0, true
	case float64:
		if vv == float64(DESC) {
			return true, true
		}
		if vv == float64(ASC) {
			return false, true
		}
		return vv < 0, true
	case bool:
		return vv, true
	case string:
		s := strings.ToLower(strings.TrimSpace(vv))
		switch s {
		case "desc", "descending", "-1":
			return true, true
		case "asc", "ascending", "1":
			return false, true
		}
	}
	return false, false
}

func normalizeFilterOp(op string) string {
	s := strings.ToLower(strings.TrimSpace(op))
	s = strings.TrimPrefix(s, "$")
	switch s {
	case "=":
		return "eq"
	case "!=":
		return "ne"
	case ">":
		return "gt"
	case ">=":
		return "gte"
	case "<":
		return "lt"
	case "<=":
		return "lte"
	case "not_in":
		return "nin"
	default:
		return s
	}
}

func toAnys(v Any) []Any {
	switch vv := v.(type) {
	case []Any:
		return vv
	case []string:
		out := make([]Any, 0, len(vv))
		for _, one := range vv {
			out = append(out, one)
		}
		return out
	case []int:
		out := make([]Any, 0, len(vv))
		for _, one := range vv {
			out = append(out, one)
		}
		return out
	case []Map:
		out := make([]Any, 0, len(vv))
		for _, one := range vv {
			out = append(out, one)
		}
		return out
	default:
		return []Any{v}
	}
}

func toInt(v Any) (int, bool) {
	switch vv := v.(type) {
	case int:
		return vv, true
	case int64:
		return int(vv), true
	case float64:
		return int(vv), true
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(vv))
		return n, err == nil
	}
	return 0, false
}

func pickString(m Map, keys ...string) (string, bool) {
	for _, key := range keys {
		if v, ok := m[key].(string); ok {
			return v, true
		}
	}
	return "", false
}

func pickMap(m Map, keys ...string) (Map, bool) {
	for _, key := range keys {
		if v, ok := m[key].(Map); ok {
			return v, true
		}
	}
	return nil, false
}

func pickValue(m Map, keys ...string) Any {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return nil
}

func pickValueOK(m Map, keys ...string) (Any, bool) {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			return v, true
		}
	}
	return nil, false
}

func parseTopLevelFilters(cfg Map) []Filter {
	if cfg == nil {
		return nil
	}
	reserved := map[string]struct{}{
		"keyword": {}, "q": {}, "$keyword": {}, "$q": {},
		"offset": {}, "$offset": {}, "limit": {}, "$limit": {},
		"fields": {}, "$fields": {}, "select": {}, "$select": {},
		"facets": {}, "$facets": {},
		"highlight": {}, "$highlight": {},
		"setting": {}, "$setting": {},
		"raw": {}, "$raw": {},
		"sort": {}, "$sort": {}, "sorts": {}, "$sorts": {},
		"filters": {}, "$filters": {}, "filter": {}, "$filter": {},
	}

	out := make([]Filter, 0)
	for key, val := range cfg {
		if _, ok := reserved[key]; ok {
			continue
		}
		if strings.HasPrefix(key, "$") {
			continue
		}
		out = append(out, parseFieldFilters(key, val)...)
	}
	return out
}

func FilterMatch(filter Filter, payload Map) bool {
	if payload == nil {
		return false
	}
	val, ok := payload[filter.Field]
	if !ok {
		return false
	}
	op := strings.ToLower(strings.TrimSpace(filter.Op))
	if op == "" {
		op = "eq"
	}
	switch op {
	case "eq", "=":
		return compareEqual(val, filter.Value)
	case "ne", "!=":
		return !compareEqual(val, filter.Value)
	case "in":
		for _, one := range filter.Values {
			if compareEqual(val, one) {
				return true
			}
		}
		return false
	case "nin", "not_in":
		for _, one := range filter.Values {
			if compareEqual(val, one) {
				return false
			}
		}
		return true
	case "gt", ">":
		return compareNumber(val, filter.Value, ">")
	case "gte", ">=":
		return compareNumber(val, filter.Value, ">=")
	case "lt", "<":
		return compareNumber(val, filter.Value, "<")
	case "lte", "<=":
		return compareNumber(val, filter.Value, "<=")
	case "range":
		return compareRange(val, filter.Min, filter.Max)
	default:
		return compareEqual(val, filter.Value)
	}
}

func compareEqual(a, b Any) bool {
	if fa, oka := toFloat(a); oka {
		if fb, okb := toFloat(b); okb {
			return fa == fb
		}
	}
	return strings.EqualFold(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

func compareNumber(a, b Any, op string) bool {
	fa, oka := toFloat(a)
	fb, okb := toFloat(b)
	if !oka || !okb {
		return false
	}
	switch op {
	case ">":
		return fa > fb
	case ">=":
		return fa >= fb
	case "<":
		return fa < fb
	case "<=":
		return fa <= fb
	default:
		return false
	}
}

func compareRange(v, min, max Any) bool {
	fv, ok := toFloat(v)
	if !ok {
		return false
	}
	if min != nil {
		if fmin, ok := toFloat(min); ok {
			if fv < fmin {
				return false
			}
		}
	}
	if max != nil {
		if fmax, ok := toFloat(max); ok {
			if fv > fmax {
				return false
			}
		}
	}
	return true
}

func toFloat(v Any) (float64, bool) {
	switch vv := v.(type) {
	case int:
		return float64(vv), true
	case int64:
		return float64(vv), true
	case float64:
		return vv, true
	case float32:
		return float64(vv), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(vv), 64)
		return f, err == nil
	default:
		return 0, false
	}
}
