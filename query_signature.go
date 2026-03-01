package search

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	. "github.com/infrago/base"
)

func QuerySignature(index string, q Query) string {
	parts := make([]string, 0, 12)
	parts = append(parts, "index="+strings.TrimSpace(index))
	parts = append(parts, "keyword="+strings.TrimSpace(q.Keyword))
	parts = append(parts, fmt.Sprintf("prefix=%t", q.Prefix))
	parts = append(parts, "filters="+filterSignature(q.Filters))
	parts = append(parts, "sorts="+sortSignature(q.Sorts))
	parts = append(parts, "fields="+strings.Join(q.Fields, ","))
	parts = append(parts, "facets="+strings.Join(q.Facets, ","))
	parts = append(parts, "highlight="+strings.Join(q.Highlight, ","))
	parts = append(parts, fmt.Sprintf("offset=%d", q.Offset))
	parts = append(parts, fmt.Sprintf("limit=%d", q.Limit))
	parts = append(parts, "raw="+stableAnySignature(q.Raw))
	parts = append(parts, "setting="+stableAnySignature(q.Setting))
	return strings.Join(parts, "|")
}

func filterSignature(filters []Filter) string {
	if len(filters) == 0 {
		return ""
	}
	parts := make([]string, 0, len(filters))
	for _, f := range filters {
		parts = append(parts, strings.Join([]string{
			strings.TrimSpace(f.Field),
			strings.TrimSpace(f.Op),
			stableAnySignature(f.Value),
			stableAnySignature(f.Values),
			stableAnySignature(f.Min),
			stableAnySignature(f.Max),
		}, ":"))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func sortSignature(in []Sort) string {
	parts := make([]string, 0, len(in))
	for _, s := range in {
		dir := "asc"
		if s.Desc {
			dir = "desc"
		}
		parts = append(parts, strings.TrimSpace(s.Field)+":"+dir)
	}
	return strings.Join(parts, ",")
}

func stableAnySignature(v Any) string {
	switch vv := v.(type) {
	case nil:
		return "null"
	case string:
		return "s:" + vv
	case bool:
		if vv {
			return "b:1"
		}
		return "b:0"
	case int:
		return fmt.Sprintf("i:%d", vv)
	case int64:
		return fmt.Sprintf("i64:%d", vv)
	case float64:
		return fmt.Sprintf("f:%g", vv)
	case time.Time:
		return "t:" + vv.UTC().Format(time.RFC3339Nano)
	case []Any:
		parts := make([]string, 0, len(vv))
		for _, one := range vv {
			parts = append(parts, stableAnySignature(one))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case []string:
		parts := make([]string, 0, len(vv))
		for _, one := range vv {
			parts = append(parts, "s:"+one)
		}
		return "[" + strings.Join(parts, ",") + "]"
	case Map:
		keys := make([]string, 0, len(vv))
		for k := range vv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+":"+stableAnySignature(vv[k]))
		}
		return "{" + strings.Join(parts, ",") + "}"
	}

	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return "null"
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		parts := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			parts = append(parts, stableAnySignature(rv.Index(i).Interface()))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case reflect.Map:
		keys := rv.MapKeys()
		keyStrings := make([]string, 0, len(keys))
		keyMap := make(map[string]reflect.Value, len(keys))
		for _, k := range keys {
			ks := fmt.Sprintf("%v", k.Interface())
			keyStrings = append(keyStrings, ks)
			keyMap[ks] = k
		}
		sort.Strings(keyStrings)
		parts := make([]string, 0, len(keyStrings))
		for _, ks := range keyStrings {
			val := rv.MapIndex(keyMap[ks]).Interface()
			parts = append(parts, ks+":"+stableAnySignature(val))
		}
		return "{" + strings.Join(parts, ",") + "}"
	default:
		return fmt.Sprintf("%T:%v", v, v)
	}
}
