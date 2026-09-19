package database

import (
	"encoding/json/v2"
	"fmt"
)

type PriceKind int // to I can diffrentiate between an actual price value and an PriceTag
const (
	PriceValue PriceKind = iota
	PriceTag
)

// price can either be an float64 or a "price tag" which is an string
type Price struct {
	Kind  PriceKind
	Value float64 // for PriceKind == PriceValue
	Tag   string  // for PriceKind == PriceTag
}

// the items each work should be able to hold (NOTE: change Time from Price to something else)
type Work struct {
	Monjogs   map[string]int `json:"monjogs"`
	Materials map[string]int `json:"materials"`
	Time      Price          `json:"times"`
}

// Implement the Unmarshal interface so i can read Price values easier
func (price *Price) UnmarshalJSON(b []byte) error {
	var num float64
	if err := json.Unmarshal(b, &num); err == nil {
		price.Kind = PriceValue
		price.Value = num
		price.Tag = ""
		return nil
	}
	var tag string
	if err := json.Unmarshal(b, &tag); err == nil {
		price.Kind = PriceTag
		price.Value = 0
		price.Tag = tag
		return nil
	}
	return fmt.Errorf("price must be a number or a string tag, got: %s", string(b))
}

// Implement the Unmarshal interface so i can write Price values easier
func (p Price) MarshalJSON() ([]byte, error) {
	switch p.Kind {
	case PriceValue:
		return json.Marshal(p.Value)
	case PriceTag:
		return json.Marshal(p.Tag)
	default:
		return nil, fmt.Errorf("unknown price kind: %v", p.Kind)
	}
}

// DataBase
type DB struct {
	Works           map[string]Work    `json:"works"`
	Time            map[string]Price   `json:"time"`
	OtherMaterials  map[string]Price   `json:"other_materials"`
	Codes           map[string]Price   `json:"codes"`
	AdditionalCosts map[string]Price   `json:"additional_costs"`
	PriceTags       map[string]float64 `json:"price tags"`
}
