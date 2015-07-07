package elasticsearch

import (
	"encoding/json"
)

type Aggregation interface {
	json.Marshaler
	Name() string
}

type DateHistogramAggregation struct {
	Field    string   `json:"field"`
	Interval Duration `json:"interval"`
}

func (d *DateHistogramAggregation) Name() string {
	return "date_histogram"
}
