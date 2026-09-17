package main

import (
	"fmt"
	database "merka-v2/internals/database"
)

func main() {
	path := "configs/works.json"
	db, err := database.Create(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db.PrintFullDataBase()
	if err = db.SaveDataBaseToFile(); err != nil {
		fmt.Println(err.Error())
	}
}
