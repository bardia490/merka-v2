package database

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"merka-v2/pkg/utility"
	"os"
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

// the items each work should be able to hold
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
	Works            map[string]Work  `json:"works"`
	Time             map[string]Price `json:"time"`
	OtherMaterials   map[string]Price `json:"other_materials"`
	Codes            map[string]Price `json:"codes"`
	Additional_costs map[string]Price `json:"additional_costs"`
	PriceTags        map[string]Price `json:"price tags"`
}

// the [path] is the relative path to the configs\works.json file
func Create(path string) (*DB, error) {
	var db DB

	if !utility.DoesFileExist(path) {
		path = "configs/work_template.json"
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(contents), &db)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// prints all the works names + monjog + material + time costs
func (db *DB) PrintFullDataBase() {
	printWorks(db.Works)
}

// this should be used when we want to print a work based on its name and we also want
// to check if the work exists beforehand
func (db *DB) PrintWork(workname string) error {
	// if the work existed inside the database
	if work, ok := db.Works[workname]; ok {
		fmt.Println("work:", workname)
		work.Print()
		return nil
	}
	return fmt.Errorf("could not find work: %s in database", workname)
}
func printWorks(works map[string]Work) {
	if len(works) == 0 {
		fmt.Println("there was no data here!")
		return
	}
	for key, work := range works {
		fmt.Printf("work: %s\n", key)
		work.Print()
	}
}

// if a work exists inside the database we can use this function to print its
// monjogs, materials and time (if unsure about the existence of the work, should use DataBase.PrintWork function)
func (work *Work) Print() {
	work.printMonjogs()
	work.printMaterials()
	work.printTime()
}

func (work *Work) printMonjogs() error {
	// NOTE: this needs to change cause currently I
	// don't access to the price tags from the database
	//Monjogs  map[string]Price
	if len(work.Monjogs) != 0 {
		fmt.Println("monjogs:")
		for m_name, m_count := range work.Monjogs {
			fmt.Printf("monjog code: %s, monjog count: %f\n", m_name, m_count)
		}
	} else {
		fmt.Println("there were no monjogs in this work!!, please make sure this was intentional")
	}
	return nil
}

func (work *Work) printMaterials() {
	// NOTE: this needs to change cause currently I
	// don't access to the price tags from the database
	//Materials map[string]float64
	if len(work.Materials) != 0 {
		fmt.Println("materials:")
		for m_name, m_price := range work.Materials {
			fmt.Printf("material name: %s, material price: %d\n", m_name, m_price)
		}
	}
}
func (work *Work) printTime() {
	//Time map[string]Price
	fmt.Printf("total time for work: %.2f\n", work.Time)
}

func (db *DB) Reload() error {
	path := "configs/works.json"
	if !utility.DoesFileExist(path) {
		path = "configs/work_template.json"
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(contents), db)
	if err != nil {
		return err
	}
	return nil
}

// every price must be converted into an float64
// this method will find the correct tag value (if the type is an price tag)
// or the value (if the type is an value from the begining), it will also convert
// every -1 into the default price value
// if none of the above matches it will return an error
func (price *Price) resolve() (float64, error) {
	switch price.Kind {
	case PriceValue:
		return price.Value, nil // NOTE: Add -1 check for default value
	case PriceTag:
		return 100, nil // NOTE: FIX this
	}
	return 0, fmt.Errorf("could not resolve the price: %v. please contact your boyfriend immediatley", *price)
}

func (db *DB) SaveDataBaseToFile() error {
	path := "configs/works.json"
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := json.MarshalWrite(file, *db, jsontext.WithIndent("  ")); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}
