/*
 *                                       /$$
 *                                      | $$
 *     /$$$$$$/$$$$   /$$$$$$   /$$$$$$ | $$   /$$  /$$$$$$
 *    | $$_  $$_  $$ /$$__  $$ /$$__  $$| $$  /$$/ |____  $$
 *    | $$ \ $$ \ $$| $$$$$$$$| $$  \__/| $$$$$$/   /$$$$$$$
 *    | $$ | $$ | $$| $$_____/| $$      | $$_  $$  /$$__  $$
 *    | $$ | $$ | $$|  $$$$$$$| $$      | $$ \  $$|  $$$$$$$
 *    |__/ |__/ |__/ \_______/|__/      |__/  \__/ \_______/
 *
 *
 *
 */

package main

import (
	"fmt"
	database "merka-v2/internals/database"
)

func main() {
	path := "configs/example_works.json"
	db, err := database.Create(path) // db is a pointer to the database struct
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db.PrintFullDataBase()
	if err := db.PrintWork("goshvare_mosalasi2"); err != nil {
		fmt.Println("couldn't print work")
	}
	if err = db.SaveDataBaseToFile(); err != nil {
		fmt.Println(err.Error())
	}
}
