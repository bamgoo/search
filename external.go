package search

import . "github.com/infrago/base"

func Clear(index string) error {
	return module.Clear(index)
}

func GetCapabilities(index string) Capabilities {
	return module.Capabilities(index)
}

func ListCapabilities() map[string]Capabilities {
	return module.ListCapabilities()
}

func Upsert(index string, rows ...Map) error {
	return module.Upsert(index, rows...)
}

func Delete(index string, ids []string) error {
	return module.Delete(index, ids)
}

func Search(index, keyword string, args ...Any) (Result, error) {
	return module.Search(index, keyword, args...)
}

func Count(index, keyword string, args ...Any) (int64, error) {
	return module.Count(index, keyword, args...)
}

func Signature(index, keyword string, args ...Any) string {
	return QuerySignature(index, BuildQuery(keyword, args...))
}
