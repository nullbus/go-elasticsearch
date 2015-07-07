package elasticsearch

import (
	"encoding/json"
)

type Aggregation map[string]Aggregator

func (a Aggregation) MarshalJSON() ([]byte, error) {
	doc := map[string]map[string]interface{}{}
	for name, agg := range a {
		doc[name] = map[string]interface{}{
			agg.Name(): agg.Aggregate(),
		}

		if pAgg, ok := agg.(ParentAggregator); ok && pAgg.ChildAggregation() != nil {
			doc[name]["aggs"] = pAgg.ChildAggregation()
		}
	}

	return json.Marshal(doc)
}

type Aggregator interface {
	Name() string
	Aggregate() *json.RawMessage
}

type ParentAggregator interface {
	ChildAggregation() Aggregation
}

func aggregateSelf(agg Aggregator) *json.RawMessage {
	jsonDoc, jsonErr := json.Marshal(agg)
	if jsonErr != nil {
		return nil
	}

	return (*json.RawMessage)(&jsonDoc)
}

type DateHistogramAggregator struct {
	Field          string      `json:"field"`
	Interval       Duration    `json:"interval"`
	SubAggregation Aggregation `json:"-"`
}

func (d *DateHistogramAggregator) Name() string {
	return "date_histogram"
}

func (d *DateHistogramAggregator) Aggregate() *json.RawMessage {
	return aggregateSelf(d)
}

func (d *DateHistogramAggregator) ChildAggregation() Aggregation {
	return d.SubAggregation
}

type SumAggregator struct {
	Field  string `json:"field,omitempty"`
	Script string `json:"script,omitempty"`
}

func (s *SumAggregator) Name() string {
	return "sum"
}

func (s *SumAggregator) Aggregate() *json.RawMessage {
	return aggregateSelf(s)
}

type SingleJSONMap struct {
	Key   string
	Value string
}

func (m *SingleJSONMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		m.Key: m.Value,
	})
}

type TermsAggregator struct {
	Field          string         `json:"field"`
	MinDocCount    int            `json:"min_doc_count,omitempty"`
	Size           int            `json:"size,omitempty"`
	Order          *SingleJSONMap `json:"order,omitempty"`
	Include        string         `json:"include,omitempty"`
	Exclude        string         `json:"exclude,omitempty"`
	SubAggregation Aggregation    `json:"-"`
}

func (t *TermsAggregator) Name() string {
	return "terms"
}

func (t *TermsAggregator) Aggregate() *json.RawMessage {
	return aggregateSelf(t)
}

func (t *TermsAggregator) ChildAggregation() Aggregation {
	return t.SubAggregation
}
