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
	"merka-v2/internals/database"
	"os"
)

func main() {
	go database.UpdateTableContents()
	go database.UpdatePromptOptions()
	go database.UpdateStatusMessage()
	if err := database.Run(); err != nil {
		fmt.Printf("app couldn't run because of: %s", err.Error())
		os.Exit(1)
	}
	database.PrintResults() // for debugging
	//database.PrintWorkTable()
	//database.PrintFullDataBase()
	//database.WriteWorksToFile()
	//if err := database.PrintWork("goshvare_mosalasi2"); err != nil {
	//	fmt.Println("couldn't print work")
	//}
	//if err = database.SaveDataBaseToFile(); err != nil {
	//	fmt.Println(err.Error())
	//}
}
