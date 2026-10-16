package database

import (
	"encoding/json/v2"
	"fmt"
)

type PriceKind int // so I can diffrentiate between an actual price value and an PriceTag
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

// the items each work should be able to hold
// time is float64 because its the minutes it takes to finish the job
type Work struct {
	Monjogs   map[string]float64 `json:"monjogs"`
	Materials map[string]float64 `json:"materials"`
	Time      float64            `json:"time"`
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
	Works             map[string]Work    `json:"works"`
	Time              map[string]Price   `json:"time"`
	OtherMaterials    map[string]Price   `json:"other_materials"`
	Codes             map[string]Price   `json:"codes"`
	AdditionalCosts   map[string]Price   `json:"additional_costs"`
	PriceTags         map[string]float64 `json:"price tags"`
	DefaultMultiplier float64            `json:"default multiplier"`
}

// NOTE: tags work as follows:
// tag_nam1: float64 (tags cannot reference each other)

// NOTE: DefaultMultiplier is only an float64 because its only one value and
// probably doesn't need a Tag value

// contains information on how to format and write the information about
// each work
type workWriterSettings struct {
	// tells the write or calc functions that use
	// these settings to give full information about
	// every step of the way
	verbose bool
	// writes the data in a non-verbose table format
	// perfect for table UI design.
	// NOTE: in order to use this option the [verbose] flag
	// must be set to false
	tableFormat bool
	// usefull when there is a lot
	//of work to be printed or written to file
	writeOneDiscountValue bool
	// if the writeOneDiscountValue is set the
	// following value is used for the multiplier
	defaultMultiplierValue float64
}
