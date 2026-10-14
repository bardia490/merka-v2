package database

// this file basically contains all the user getter and setter functions

import (
	"bufio"
	"fmt"
	"merka-v2/pkg/questions"
	"os"
	"strconv"
	"strings"
)

// prompts the user for the price. the question is used to convey to the user in what
// context is the price going to be used for (material_price, monjog_price or ...).
// it returns the price either as an price tag or an float64
// if the price is entered as a price tag, the function will make sure the price tag
// exists and prompts the user for the tag if it does not using the setPriceTagFromUser function.
func getPriceFromUser(question string) Price {
	var price Price
	// it just has to have an length bigger than 0
	condition := func(s string) bool { return len(s) > 0 }
	// using PromptString because we expect the Price to be either float64 (PriceValue) or string (PriceTag)
	p := questions.PromptString(question, "price can be either an value (such as 14562.32) or a Price Tag (such as tag1) it just has to non zero", condition)

	// check to see if its a float64 (PriceValue) or not
	if val, err := strconv.ParseFloat(p, 64); err != nil {
		price.Kind = PriceValue
		price.Value = val
	} else { // if it wasn't float64 we will assume it was an price tag
		price.Kind = PriceTag
		price.Tag = p
		// check to make sure the tag exists
		if _, ok := db.PriceTags[price.Tag]; !ok {
			setPriceTagFromUser(price.Tag)
		}
	}
	return price
}

// this is just an helper function for the getPriceFromUser and
// it sets the [tag_name] by prompting the user for a (new) value for [tag_name].
// if the [tag_name] already exists in db.PriceTags then it will replace it.
func setPriceTagFromUser(tag_name string) {
	var question strings.Builder
	question.WriteString("Please Enter The price for the PriceTag: ")
	question.WriteString(tag_name)
	condition := func(val float64) bool { return val > 0 }
	db.PriceTags[tag_name] = questions.PromptFloat(question.String(), "the price has to be greater than zero", condition)
}

// it prompts the user for monjog name and count values
// if the monjog name is not part of the Codes (which holds the monjogs prices)
// prompts the user for the prices which can be either price tag or price value
// if the price tag is not defined the user will be prompted for that as well
// at the end it modifies the monjogs map inplace
func getMonjogsFromUser(monjogs map[string]float64) {
	var question strings.Builder
	// this one is always true since if the length of the monjog_name is zero it means we are done
	monjog_name_condition := func(name string) bool { return true }
	monjog_count_condition := func(count float64) bool { return count > 0 }
	for {
		monjog_name := questions.PromptString("please Enter monjog name:", "the name can be anything you want but it needs to have at least one character long, when you were done you can just press Enter", monjog_name_condition)
		// we are done
		if len(monjog_name) == 0 {
			break
		}
		question.WriteString("Please Enter ")
		question.WriteString(monjog_name)
		question.WriteString(" count")
		monjog_count := questions.PromptFloat(question.String(), "the count must be a positive integer", monjog_count_condition)
		monjogs[monjog_name] = monjog_count
		question.Reset() // need to reset it for the next monjog name (code)
	}
	// check all monjogs to see which ones don't have an price or price tag
	for monjog_code := range monjogs {
		if _, is_in_codes := db.Codes[monjog_code]; !is_in_codes {
			question.WriteString("Please Enter the price for ")
			question.WriteString(monjog_code)
			db.Codes[monjog_code] = getPriceFromUser(question.String())
		}
	}
}

// it prompts the user for material name and count values
// if the material name is not part of the Codes (which holds the materials prices)
// prompts the user for the prices which can be either price tag or price value
// if the price tag is not defined the user will be prompted for that as well
// at the end it modifies the materials map inplace
func getMaterialsFromUser(materials map[string]float64) {
	var question strings.Builder
	// this one is always true since if the length of the material_name is zero it means we are done
	material_name_condition := func(name string) bool { return true }
	material_count_condition := func(count float64) bool { return count > 0 }
	for {
		material_name := questions.PromptString("please Enter material name:", "the name can be anything you want but it needs to have at least one character long, when you were done you can just press Enter", material_name_condition)
		// we are done
		if len(material_name) == 0 {
			break
		}
		question.WriteString("Please Enter ")
		question.WriteString(material_name)
		question.WriteString(" count")
		material_count := questions.PromptFloat(question.String(), "the count must be a positive integer", material_count_condition)
		materials[material_name] = material_count
		question.Reset() // need to reset it for the next material name (code)
	}
	// check all monjogs to see which ones don't have an price or price tag
	for material_name := range materials {
		if _, is_in_names := db.Codes[material_name]; !is_in_names {
			question.WriteString("Please Enter the price for ")
			question.WriteString(material_name)
			db.Codes[material_name] = getPriceFromUser(question.String())
		}
	}
}

// TODO: currently the getMaterialsFromUser and getMonjogsFromUser are exactly the same.
// if i don't change my mind about their implementation in the future I should probably
// use one function for both.
func AddWork() {
	var work Work
	name_condition := func(name string) bool { return len(name) > 0 } // it askes for a name to be used in work
	work_name := questions.PromptString("please Enter work name:", "the name can be anything you want but it needs to have at least one character long", name_condition)
	// get the monjogs for work (monjogs: map[string]float64)
	monjogs := make(map[string]float64)
	getMonjogsFromUser(monjogs)
	work.Monjogs = monjogs
	// get the materials for work (materials: map[string]int)
	materials := make(map[string]float64)
	getMaterialsFromUser(materials)
	work.Materials = materials
	// get the time for work
	time := questions.PromptTime()
	work.Time = time
	db.Works[work_name] = work
}

// prompts the user for the multiplier.
// this function basically prompts the user for a float64 value
// but the multiplier can also be the default value set in the works.json file
// in the form of "default multiplier": float64, in case the user wants
// the default value they must press Enter (empty string)
func getMultiplierValueFromUser() float64 {
	var answer float64
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Q. please enter the value for the multiplier")
		fmt.Println("NOTE: to increase the price enter a number bigger than zero in the format 1.xx (e.g 1.25 for 25% price increase)")
		fmt.Println("NOTE: to decrease the price enter a number smaller than zero in the format 0.xx (e.g 0.75 for 25% price discount)")
		fmt.Println("NOTE: for the default value just press Enter")
		fmt.Println("NOTE: pressing 0 or 1 will just use the current price")
		fmt.Printf("> ")
		scanner.Scan()
		output := strings.TrimSpace(scanner.Text())
		if len(output) == 0 {
			answer = db.DefaultMultiplier // if the default multiplier is not set it will return zero
			break
		}
		answer, err = strconv.ParseFloat(output, 64)
		if err == nil {
			break
		}
		fmt.Println("ERROR: Please enter a correct multiplier value")
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("this is probably not important but please report this: %s\n", err.Error())
	}
	return answer
}
