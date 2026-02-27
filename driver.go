package search

import . "github.com/bamgoo/base"

type (
	Capabilities struct {
		SyncIndex bool
		Clear     bool
		Upsert    bool
		Delete    bool
		Search    bool
		Count     bool
		Suggest   bool

		Sort      bool
		Facets    bool
		Highlight bool

		FilterOps []string
	}

	Driver interface {
		Connect(*Instance) (Connection, error)
	}

	Connection interface {
		Open() error
		Close() error

		Capabilities() Capabilities
		SyncIndex(name string, index Index) error
		Clear(index string) error
		Upsert(index string, rows []Map) error
		Delete(index string, ids []string) error
		Search(index string, query Query) (Result, error)
		Count(index string, query Query) (int64, error)
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
		Prefix    bool
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
