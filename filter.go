package es

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

type FilterAggregaion struct {
	Filter         Filterer
	SubAggregation Aggregation
}

func (f *FilterAggregaion) Name() string {
	return "filter"
}

func (f *FilterAggregaion) Aggregate() *json.RawMessage {
	if f.Filter == nil {
		emptyJSON := []byte("{}")
		return (*json.RawMessage)(&emptyJSON)
	}

	jsonStr, _ := json.Marshal(map[string]interface{}{
		f.Filter.Name(): f.Filter.Filter(),
	})

	return (*json.RawMessage)(&jsonStr)
}

func (f *FilterAggregaion) ChildAggregation() Aggregation {
	return f.SubAggregation
}

type Filterer interface {
	Filter() *json.RawMessage
	Name() string
}

type AndFilter []Filterer

func (a AndFilter) Name() string {
	return "and"
}

func (a AndFilter) Filter() *json.RawMessage {
	filters := make([]map[string]*json.RawMessage, 0, len(a))
	for _, f := range a {
		filters = append(filters, map[string]*json.RawMessage{
			f.Name(): f.Filter(),
		})
	}

	jsonStr, _ := json.Marshal(filters)
	return (*json.RawMessage)(&jsonStr)
}

// uses same marshaling logit of AndFilter
type OrFilter AndFilter

func (o OrFilter) Name() string {
	return "or"
}

func (o OrFilter) Filter() *json.RawMessage {
	return AndFilter(o).Filter()
}

type RangeFilter struct {
	Field string
	From  interface{}
	To    interface{}
}

func (r *RangeFilter) Name() string {
	return "range"
}

func (r *RangeFilter) Filter() *json.RawMessage {
	doc := map[string]map[string]interface{}{
		r.Field: map[string]interface{}{},
	}
	if r.From != nil {
		doc[r.Field]["from"] = r.From
	}
	if r.To != nil {
		doc[r.Field]["to"] = r.To
	}

	jsonStr, jsonErr := json.Marshal(doc)
	if jsonErr != nil {
		return nil
	}

	return (*json.RawMessage)(&jsonStr)
}
