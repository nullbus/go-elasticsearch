package elasticsearch

import (
	"encoding/json"
	"fmt"
)

type QueryHead struct {
	Query Query
}

func (qh *QueryHead) MarshalJSON() ([]byte, error) {
	if nil == qh.Query {
		return nil, fmt.Errorf("empty query")
	}

	return json.Marshal(map[string]interface{}{
		qh.Query.Name(): qh.Query,
	})
}

type Query interface {
	json.Marshaler
	Name() string
}

type TextValue string

func (q *TextValue) MarshalJSON() ([]byte, error) {
	return []byte(*q), nil
}

const (
	MATCH_OPERATER_OR = iota // default
	MATCH_OPERATER_AND
)

const (
	MATCH_TYPE_BOOL = iota // default
	MATCH_TYPE_PHRASE
	MATCH_TYPE_PHRASE_PREFIX
)

type MatchAllQuery struct{}

func (q *MatchAllQuery) Name() string {
	return "match_all"
}

func (q *MatchAllQuery) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

type MatchQuery struct {
	Field    string
	Query    string
	Type     int
	Operator int
}

func (q *MatchQuery) Name() string {
	return "match"
}

func (q *MatchQuery) MarshalJSON() ([]byte, error) {
	if q.Field == "" {
		return nil, fmt.Errorf("empty field")
	}

	format := struct {
		Query    string `json:"query"`
		Type     string `json:"type,omitempty"`
		Operator string `json:"operator"`
	}{
		Query: q.Query,
	}

	if q.Operator == MATCH_OPERATER_OR {
		format.Operator = "or"
	} else {
		format.Operator = "and"
	}

	if q.Type == MATCH_TYPE_PHRASE {
		format.Type = "phrase"
	} else if q.Type == MATCH_TYPE_PHRASE_PREFIX {
		format.Type = "phrase_prefix"
	}

	return json.Marshal(&format)
}
