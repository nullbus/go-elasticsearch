package elasticsearch

import (
	"encoding/json"
	"fmt"
)

type FilteredQuery struct {
	Filter Filterer
	Query  Query
}

func (f *FilteredQuery) Name() string {
	return "filtered"
}

func (f *FilteredQuery) MarshalJSON() ([]byte, error) {
	if f.Filter == nil {
		return nil, fmt.Errorf("empty filter")
	}

	doc := map[string]interface{}{
		"filter": map[string]interface{}{
			f.Filter.Name(): f.Filter.Filter(),
		},
	}

	if f.Query != nil {
		doc["query"] = f.Query
	}

	return json.Marshal(doc)
}

type Filterer interface {
	Filter() json.RawMessage
	Name() string
}

type AndFilter []Filterer

func (a AndFilter) Name() string {
	return "and"
}

func (a AndFilter) Filter() json.RawMessage {
	filters := make([]*json.RawMessage, 0, len(a))
	for _, f := range a {
		jsonStr, jsonErr := json.Marshal(map[string]Filterer{
			f.Name(): f,
		})
		if jsonErr != nil {
			return json.RawMessage("[]")
		}

		filters = append(filters, (*json.RawMessage)(&jsonStr))
	}

	jsonStr, _ := json.Marshal(filters)
	return json.RawMessage(jsonStr)
}

type OrFilter struct {
	// uses same marshaling logit of AndFilter
	AndFilter
}

func (o *OrFilter) Name() string {
	return "or"
}

type RangeFilter struct {
	Field string
	From  interface{}
	To    interface{}
}

func (r *RangeFilter) Name() string {
	return "range"
}

func (r *RangeFilter) Filter() json.RawMessage {
	doc := map[string]interface{}{}
	if r.From != nil {
		doc["from"] = r.From
	}
	if r.To != nil {
		doc["to"] = r.To
	}

	jsonStr, jsonErr := json.Marshal(doc)
	if jsonErr != nil {
		return json.RawMessage("{}")
	}

	return json.RawMessage(jsonStr)
}
