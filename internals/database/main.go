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
