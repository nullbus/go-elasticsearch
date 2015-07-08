package es

import (
	"encoding/json"
)

type Sort struct {
	Field        string
	AscOrder     bool
	Mode         string
	Unit         string
	DistanceType string
}

func (s Sort) MarshalJSON() ([]byte, error) {
	order := "desc"
	if s.AscOrder {
		order = "asc"
	}

	var value struct {
		Order        string `json:"order,omitempty"`
		Mode         string `json:"mode,omitempty"`
		Unit         string `json:"unit,omitempty"`
		DistanceType string `json:"distance_type,omitempty"`
	}

	value.Order = order
	value.Mode = s.Mode
	value.DistanceType = s.DistanceType
	value.Unit = s.Unit

	return json.Marshal(map[string]interface{}{
		s.Field: value,
	})
}
