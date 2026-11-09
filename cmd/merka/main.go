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
	"merka-v2/internals/UI"
	database "merka-v2/internals/database"
	"os"
)

func main() {
	path := "configs/works.json"
	err := database.LoadDataBase(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := UI.Run(); err != nil {
		fmt.Printf("app couldn't run because of: %s", err.Error())
		os.Exit(1)
	}
	database.PrintWorkTable()
	//database.PrintFullDataBase()
	//database.WriteWorksToFile()
	//if err := database.PrintWork("goshvare_mosalasi2"); err != nil {
	//	fmt.Println("couldn't print work")
	//}
	//if err = database.SaveDataBaseToFile(); err != nil {
	//	fmt.Println(err.Error())
	//}
}
