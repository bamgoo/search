package search

import . "github.com/bamgoo/base"

type (
	Driver interface {
		Connect(*Instance) (Connection, error)
	}

	Connection interface {
		Open() error
		Close() error

		CreateIndex(name string, index Index) error
		DropIndex(name string) error
		Upsert(index string, docs []Document) error
		Delete(index string, ids []string) error
		Search(index string, query Query) (Result, error)
		Count(index string, query Query) (int64, error)
		Suggest(index string, text string, limit int) ([]string, error)
	}

	Index struct {
		Name        string
		Desc        string
		Primary     string
		Attributes  Vars
		StrictWrite bool
		StrictRead  bool
		Fields      Map
		Language    string
		Analyzer    string
		Setting     Map
	}

	Indexes map[string]Index

	Document struct {
		ID      string `json:"id"`
		Payload Map    `json:"payload"`
	}

	Filter struct {
		Field  string
		Op     string
		Value  Any
		Values []Any
		Min    Any
		Max    Any
	}

	Sort struct {
		Field string
		Desc  bool
	}

	Facet struct {
		Field string
		Value string
		Count int64
	}

	Hit struct {
		ID        string  `json:"id"`
		Score     float64 `json:"score"`
		Payload   Map     `json:"payload"`
		Highlight Map     `json:"highlight,omitempty"`
	}

	Query struct {
		Keyword   string
		Filters   []Filter
		Sorts     []Sort
		Offset    int
		Limit     int
		Fields    []string
		Facets    []string
		Highlight []string
		Raw       Map
		Setting   Map
	}

	Result struct {
		Total  int64              `json:"total"`
		Took   int64              `json:"took"`
		Hits   []Hit              `json:"hits"`
		Facets map[string][]Facet `json:"facets,omitempty"`
		Raw    Any                `json:"raw,omitempty"`
	}
)
