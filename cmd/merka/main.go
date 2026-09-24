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
	"log/slog"
	database "merka-v2/internals/database"
	"merka-v2/pkg/questions"
)

func main() {
	path := "configs/example_works.json"
	err := database.LoadDataBase(path) // db is a pointer to the database struct
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//database.PrintFullDataBase()
	_questions := []string{"bardia", "nika"}
	index := questions.GetAnswer("something", _questions[:])
	slog.Error("you have chosen:", _questions[index], "")
	//if err := database.PrintWork("goshvare_mosalasi2"); err != nil {
	//	fmt.Println("couldn't print work")
	//}
	//if err = database.SaveDataBaseToFile(); err != nil {
	//	fmt.Println(err.Error())
	//}
}
