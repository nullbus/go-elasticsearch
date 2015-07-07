package elasticsearch

import (
	"encoding/json"
	"fmt"
)

type Aggregation struct {
	Name           string
	Aggregator     Aggregator
	SubAggregation *Aggregation
}

func (a *Aggregation) MarshalJSON() ([]byte, error) {
	if nil == a.Aggregator {
		return nil, fmt.Errorf("empty aggregator")
	}

	doc := map[string]interface{}{
		a.Name: map[string]interface{}{
			a.Aggregator.Name(): a.Aggregator.Aggregate(),
		},
	}

	if nil != a.SubAggregation {
		if a == a.SubAggregation {
			return nil, fmt.Errorf("cannot aggregate recursive struct")
		}

		doc["aggs"] = a.SubAggregation
	}

	return json.Marshal(doc)
}

type Aggregator interface {
	Aggregate() json.RawMessage
	Name() string
}

func aggregateSelf(agg Aggregator) json.RawMessage {
	jsonDoc, jsonErr := json.Marshal(agg)
	if jsonErr != nil {
		return json.RawMessage("{}")
	}

	return json.RawMessage(jsonDoc)
}

type DateHistogramAggregator struct {
	Field    string   `json:"field"`
	Interval Duration `json:"interval"`
}

func (d *DateHistogramAggregator) Name() string {
	return "date_histogram"
}

func (d *DateHistogramAggregator) Aggregate() json.RawMessage {
	return aggregateSelf(d)
}

type SumAggregator struct {
	Field  string `json:"field,omitempty"`
	Script string `json:"script,omitempty"`
}

func (s *SumAggregator) Name() string {
	return "sum"
}

func (s *SumAggregator) Aggregate() json.RawMessage {
	return aggregateSelf(s)
}
