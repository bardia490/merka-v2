package database

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"merka-v2/pkg/utility"
	"os"
)

// DataBase
type DB struct {
	data map[string]any
}

// the [path] is the relative path to the configs\works.json file
func Create(path string) (DB, error) {
	var db DB
	db.data = make(map[string]any, 10)

	if !utility.DoesFileExist(path) {
		path = "configs/work_template.json"
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return DB{}, err
	}
	err = json.Unmarshal([]byte(contents), &db.data)
	if err != nil {
		return DB{}, err
	}
	return db, nil
}

func (db DB) PrintFullDataBase() {
	printJson(db.data)
}

func printJson(data map[string]any) {
	if len(data) == 0 {
		fmt.Println("there was no data here!")
	}
	for key, val := range data {
		switch vv := val.(type) {
		case string:
			fmt.Println(key, "is string", vv)
		case float64:
			fmt.Println(key, "is float64", vv)
		case []any:
			fmt.Println(key, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case map[string]any: // for objects (maps in go)
			fmt.Println(key, "is an map with contents:")
			printJson(vv)
		default:
			fmt.Println(key, "is of a type I don't know how to handle")
		}
	}
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
	err = json.Unmarshal([]byte(contents), &db.data)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) SaveDataBaseToFile() error {
	path := "configs/works.json"
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := json.MarshalWrite(file, db.data, jsontext.WithIndent("  ")); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}

func (db DB) PrintWork(workname string) error {
	// if the work existed inside the database
	if _, ok := db.data[workname]; !ok {
		return fmt.Errorf("could not find work: %s in database", workname)
	}
	work_contents, work_assertion_ok := db.data[workname].(map[string]any)
	if !work_assertion_ok { // if work wasn't map[string]any for some wierd reason
		return fmt.Errorf("PLEASE CALL YOU'RE BOYFRIEND IMMIEDIETLY, work %s wasn't an map[string]any, it was: %v", workname, db.data[workname])
	}
	// check to see if monjogs are in work
	if monjogs, ok := work_contents["monjogs"]; ok {
		monjogs_contents, monjogs_assertion_ok := monjogs.(map[string]any)
		if !monjogs_assertion_ok { // if monjogs wasn't map[string]any for some wierd reason
			return fmt.Errorf("PLEASE CALL YOU'RE BOYFRIEND IMMIEDIETLY, monjogs wasn't an map[string]any in work: %s, it was: %v", workname, monjogs)
		}
		db.printMonjogs(monjogs_contents)
	}
	// check to see if materials are in work
	if materials, ok := work_contents["materials"]; ok {
		materials_contents, materials_assertion_ok := materials.(map[string]any)
		if !materials_assertion_ok { // if materials wasn't map[string]any for some wierd reason
			return fmt.Errorf("PLEASE CALL YOU'RE BOYFRIEND IMMIEDIETLY, materials wasn't an map[string]any in work: %s, it was: %v", workname, materials)
		}
		db.printMaterials(materials_contents)
	}
	// check to see if time is in work
	if times, ok := work_contents["times"]; ok {
		times_content, times_assertion_ok := times.(float64)
		if !times_assertion_ok { // if times wasn't map[string]any for some wierd reason
			return fmt.Errorf("PLEASE CALL YOU'RE BOYFRIEND IMMIEDIETLY, times wasn't an float64 in work: %s, it was: %v", workname, times)
		}
		db.printTime(times_content)
	}
	return nil
}

func (db *DB) printMonjogs(monjogs map[string]any) error {
	codes := db.data["codes"].(map[string]any) // NOTE: should not use any here, it should be float64|string, change it later

	for monjog_name, count_interface := range monjogs {
		var price float64

		// how many monjogs of this type do we have
		count := count_interface.(float64)

		// find the price for this monjog
		if price_interface, ok := codes[monjog_name]; ok {
			switch val := price_interface.(type) {
			case float64: // NOTE: for backwards compatibility it should check for -1 in case of default
				price = val
			case string:
				fmt.Println("not implemented yet, will use 1 for now")
				price = 1
			default:
				return fmt.Errorf("the price for monjog: %s was not correctly written in the database, please correct the price for this monjog in the \"codes\" part", monjog_name)
			}
		} else {
			return fmt.Errorf("the price for monjog: %s was not in the database, please add the price for this monjog in the \"codes\" part", monjog_name)
		}
		fmt.Printf("monjog code: %s\n count: %d\n price: %f, total price: %f", monjog_name, int(count), price, count*price)
	}
	return nil
}
func (db *DB) printMaterials(materials map[string]any) {
}
func (db *DB) printTime(time float64) {
}
