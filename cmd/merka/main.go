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
	path := "configs/works.json"
	err := database.LoadDataBase(path)
	if err != nil {
		fmt.Println(err.Error())
		return
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
