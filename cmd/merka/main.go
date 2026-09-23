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
	err := database.LoadDataBase(path) // db is a pointer to the database struct
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//database.PrintFullDataBase()
	database.GetAnswer("something", []string{"bardia", "nika"})
	//if err := database.PrintWork("goshvare_mosalasi2"); err != nil {
	//	fmt.Println("couldn't print work")
	//}
	//if err = database.SaveDataBaseToFile(); err != nil {
	//	fmt.Println(err.Error())
	//}
}
