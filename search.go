package elasticsearch

import (
	"io"
	"net/url"
	"strings"
)

const (
	SEARCH_TYPE_QUERY_THEN_FETCH SearchType = iota // default
	SEARCH_TYPE_COUNT
	SEARCH_TYPE_DFS_QUERY_FETCH
	SEARCH_TYPE_SCAN
	SEARCH_TYPE_SCROLL
)

type SearchType int

func (t SearchType) String() string {
	switch t {
	case SEARCH_TYPE_COUNT:
		return "count"
	case SEARCH_TYPE_DFS_QUERY_FETCH:
		return "dfs_query_then_fetch"
	case SEARCH_TYPE_SCAN:
		return "scan"
	}

	// include SEARCH_TYPE_QUERY_THEN_FETCH
	return "query_then_fetch"
}

type Search interface {
	API
	Type() SearchType
	Data() interface{}
}

type QueryData struct {
	Query       Query       `json:"query,omitempty"`
	Aggregation Aggregation `json:"aggs,omitempty"`
	Size        *int        `json:"size,omitempty"`
}

type DefaultSearch struct {
	Indices   []string
	Types     []string
	QueryData QueryData
}

func (s *DefaultSearch) Type() SearchType {
	return SEARCH_TYPE_QUERY_THEN_FETCH
}

func (s *DefaultSearch) AddIndex(name string) {
	s.Indices = append(s.Indices, name)
}

func (s *DefaultSearch) AddType(name string) {
	s.Types = append(s.Types, name)
}

func (s *DefaultSearch) SetSize(size int) {
	s.QueryData.Size = &size
}

func (s *DefaultSearch) Path() (path string) {
	if len(s.Indices) != 0 {
		path += "/" + strings.Join(s.Indices, ",")
	}

	if len(s.Types) != 0 {
		path += "/" + strings.Join(s.Types, ",")
	}

	path += "/_search"
	return path
}

func (s *DefaultSearch) Query() url.Values {
	return url.Values{}
}

func (s *DefaultSearch) Data() interface{} {
	return &s.QueryData
}

type CountSearch DefaultSearch

func (s *CountSearch) Type() SearchType {
	return SEARCH_TYPE_COUNT
}

func (s *CountSearch) Data() interface{} {
	return (*DefaultSearch)(s).Data()
}

type DFSSearch DefaultSearch

func (s *DFSSearch) Type() SearchType {
	return SEARCH_TYPE_DFS_QUERY_FETCH
}

func (s *DFSSearch) Data() interface{} {
	return (*DefaultSearch)(s).Data()
}

type ScanSearch struct {
	DefaultSearch
	ScrollTime string
}

func (s *ScanSearch) Type() SearchType {
	return SEARCH_TYPE_SCAN
}

func (s *ScanSearch) Query() url.Values {
	ret := s.DefaultSearch.Query()
	ret["scroll"] = []string{s.ScrollTime}
	return ret
}

type ScrollSearch struct {
	ScanSearch
	ScrollID string
}

func (s *ScrollSearch) Type() SearchType {
	return SEARCH_TYPE_SCROLL
}

func (s *ScrollSearch) Path() string {
	return "/_search/scroll"
}

func (s *ScrollSearch) Query() url.Values {
	return url.Values{
		"scroll":    []string{s.ScrollTime},
		"scroll_id": []string{s.ScrollID},
	}
}

func (s *ScrollSearch) Data() io.Reader {
	return nil
}
