package database

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"merka-v2/pkg/questions"
	"merka-v2/pkg/utility"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var db *DB

func init() {
	db = newDB()
}

// this function will initialize the [db] global variable when the library is imported
func newDB() *DB {
	return &DB{
		Works:           make(map[string]Work),
		Codes:           make(map[string]Price),
		OtherMaterials:  make(map[string]Price),
		PriceTags:       make(map[string]float64),
		Time:            make(map[string]Price),
		AdditionalCosts: make(map[string]Price),
	}
}

// this function should be the first thing called when wanting to work with the libray
// the [path] is the relative path to the configs\works.json file
func LoadDataBase(path string) error {
	if !utility.DoesFileExist(path) {
		path = "configs/work_template.json"
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, db)
	if err != nil {
		// check for a syntax Error
		if syntaxErr, ok := errors.AsType[*jsontext.SyntacticError](err); ok {
			line, col := calculateLineAndColumnJsonError(contents, syntaxErr.ByteOffset)
			return fmt.Errorf("there was a problem whith parsing the json file at line: %d, column: %d.\nthe ERROR: %s", line, col, err.Error())
		}
		// check for symantic Error
		if semErr, ok := errors.AsType[*json.SemanticError](err); ok {
			// SemanticError in v2 includes JSONPointer path information!
			if semErr.ByteOffset > 0 {
				line, col := calculateLineAndColumnJsonError(contents, semErr.ByteOffset)
				return fmt.Errorf("JSON type mismatch at line %d, column %d (at path %q): %w",
					line, col, semErr.JSONPointer, semErr)
			}
			return fmt.Errorf("JSON type mismatch at path %q: %w", semErr.JSONPointer, semErr)
		}
		return err
	}
	//runChecks(checkMonjogs | checkTimes | checkMaterials | checkTime | checkPrices | checkDefaultPrice)
	return nil
}

// prints all the work names + monjog + material + time costs
func PrintFullDataBase() {
	printWorks(db.Works, true)
}

// this should be used when we want to print a work based on its name and we also want
// to check if the work exists beforehand (returns error if the work doesn't exist)
func PrintWork(workname string) error {
	// first make sure these things exist:
	// if they dont exist error messages for each will be printed onto the screen
	// and the function terminates
	if errors := runChecks(checkTime | checkDefaultPrice); len(errors) != 0 {
		fmt.Printf("the function cannot run due to the following errors:")
		for error_number, err := range errors {
			fmt.Printf("%d: %s", error_number+1, err.Error())
			return fmt.Errorf("function terminated due to ERRORS!")
		}
	}

	// if the work existed inside the database
	if work, ok := db.Works[workname]; ok {
		fmt.Println("work:", workname)
		work.print(false)
		return nil
	}
	return fmt.Errorf("could not find work: %s in database", workname)
}

// prints all the works in database using the Print method for each work
// if verbose is set to true then the whole information about each monjog
// (including count, price and total price) is displayed
func printWorks(works map[string]Work, verbose bool) {
	if len(works) == 0 {
		fmt.Println("there was no data here!")
		return
	}
	// make sure Time price and Default Price have been set, if not then terminate early
	if errors := runChecks(checkTime | checkDefaultPrice); len(errors) != 0 {
		fmt.Printf("the function cannot run due to the following errors:")
		for error_number, err := range errors {
			fmt.Printf("%d: %s", error_number+1, err.Error())
			return
		}
	}

	for key, work := range works {
		fmt.Printf("work: %s\n", key)
		work.print(verbose)
		fmt.Println("<><><><><><><><><><><><><><><><><><><><><>")
	}
}

// if a work exists inside the database we can use this function to print its
// monjogs, materials and time (if unsure about the existence of the work,
// should use DataBase.PrintWork function)
// if verbose is set to true then the whole information about each monjog
// (including count, price and total price) is displayed
func (work *Work) print(verbose bool) {
	total_price, monjog_price, material_price, time_price := 0.0, 0.0, 0.0, 0.0
	var err error
	// print the complete monjogs prices and counts
	if monjog_price, err = work.calcMonjogPrices(verbose); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the monjogs, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += monjog_price
	// print the materials
	if material_price, err = work.calcMaterialPrices(verbose); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the materials, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += material_price
	// print the time
	if time_price, err = work.calcTimePrice(verbose); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the time, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += time_price
}

// calculates the total price for the monjogs in [work]
// if [verbose] is set to true, then the price for each monjog will
// be displayed during calculation + the total price at the end
// it returns the final price for monjogs or an error if any of the prices could not be
// resolved (see [resolve] function for more details on resolving prices)
func (work *Work) calcMonjogPrices(verbose bool) (total_price float64, resolve_error error) {
	total_price = 0
	for m_code, m_count := range work.Monjogs {
		p := db.Codes[m_code]
		price, err := p.resolve()
		if err != nil {
			fmt.Println("there was a problem with calculating the price of the monjogs")
			return 0, err
		}
		total_price += price * float64(m_count)
		if verbose {
			fmt.Printf("monjog code: %s| monjog count: %f| monjog price: %.2f| total price: %.2f\n", m_code, m_count, price, float64(m_count)*price)
		}
	}
	if verbose {
		fmt.Printf("the total price for monjogs is: %.2f\n", total_price)
	}
	return total_price, nil
}

// calculates the total price for the materials in [work]
// if [verbose] is set to true, then the price for each material will
// be displayed during calculation + the total price at the end
// it returns the final price for materials or an error if any of the prices could not be
// resolved (see [resolve] function for more details on resolving prices)
func (work *Work) calcMaterialPrices(verbose bool) (total_price float64, resolve_error error) {
	total_price = 0
	//Materials map[string]Price
	if len(work.Materials) != 0 {
		fmt.Println("materials:")
		for m_name, m_count := range work.Materials {
			p := db.OtherMaterials[m_name]
			price, err := p.resolve()
			if err != nil {
				fmt.Printf("there was a problem with calculating the price of the material: %s\n", m_name)
				return 0, err
			}
			if verbose {
				fmt.Printf("material name: %s| material price: %f| material price: %f.2| total price: %.2f\n", m_name, m_count, price, float64(m_count)*price)
			}
			total_price += price * float64(m_count)
		}
		if verbose {
			fmt.Printf("the total price for monjogs is: %.2f\n", total_price)
		}
	}
	return total_price, nil
}

// calculates the total price for the time in [work]
// if [verbose] is set to true, then more information about the time of work will
// be displayed during calculation + the total price at the end
// it returns the final price for materials or an error if the price could not be
// resolved (see [resolve] function for more details on resolving prices)
func (work *Work) calcTimePrice(verbose bool) (total_price float64, resolve_error error) {
	// work.Time: float64
	time := work.Time

	//Time map[string]Price
	p := db.Time["price"]
	price, err := p.resolve()
	if err != nil {
		fmt.Println("there was a problem with calculating the price of the time")
		return 0, err
	}
	total_price = time * price
	if verbose {
		fmt.Printf("total time for work: %.2f| time price: %.2f| total time price: %.2f\n", time, price, time*price)
	}
	return total_price, nil
}

// just in case the apps file is modified during runtime, this can be used to reload everything
func Reload() error {
	path := "configs/works.json"
	if !utility.DoesFileExist(path) {
		path = "configs/work_template.json"
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(contents), db)
	if err != nil {
		return err
	}
	return nil
}

// every price must be converted into an float64
// this method will find the correct tag value (if the type is an price tag)
// or the value (if the type is an value from the begining), it will also convert
// every -1 into the default price value (for legacy reasons)
// if none of the above matches it will return an error
func (price *Price) resolve() (float64, error) {
	if price.Kind == PriceValue {
		result := price.Value
		if result == -1.0 { // for legacy default value
			result = getDefaultPrice() // this never fails
		}
		return result, nil
	}
	// priceTag
	var ok bool
	result := 0.0
	result, ok = db.PriceTags[price.Tag]
	if !ok { // if the Tag wasn't set return an error
		return 0, fmt.Errorf("could not resolve the price: %v. the Tag: %s could not be found, please add it", *price, price.Tag)
	}

	return result, nil
}

// this function assumes a code with the name [default] always exists and its non negative
func getDefaultPrice() float64 {
	defaultPrice := db.Codes["default"]
	val, _ := defaultPrice.resolve()
	return val
}

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

func SaveDataBaseToFile() error {
	path := "configs/works.json"
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := json.MarshalWrite(file, *db, jsontext.WithIndent("  ")); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}

// NOTE: the following functions with the postfix v2 (*_v2) are going to
// replace the functions with similar names that don't have the _v2 postfix
// for now though they will be only used inside the WriteWorksToFile method
// until I'm Done experimenting with them

// TODO: replace the legacy functions with the new _v2 functions

// if a work exists inside the database we can use this function to print its
// monjogs, materials and time (if unsure about the existence of the work,
// should use DataBase.PrintWork function which takes a work_name parameter)
// if verbose is set to true then the whole information about each monjog
// (including count, price and total price) is displayed
// NOTE: this function doesn't actually print anything, it recieves an output parameter
// which collects the impotant information inside an strings.Builder struct and the user
// can flush the information out later.
// TODO: add a configuration struct to this function so costumisation and adding
// new settings could happen easier in the future
func (work *Work) writeToBuffer(output *strings.Builder, settings workWriterSettings) {
	total_price, monjog_price, material_price, time_price := 0.0, 0.0, 0.0, 0.0
	var err error
	// print the complete monjogs prices and counts
	if monjog_price, err = work.calcMonjogPrices_v2(output, settings); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the monjogs, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += monjog_price
	// print the materials
	if material_price, err = work.calcMaterialPrices_v2(output, settings); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the materials, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += material_price
	// print the time
	if time_price, err = work.calcTimePrice_v2(output, settings); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the time, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += time_price
	var numBuf [32]byte
	if settings.verbose {
		output.WriteString("the total price for work is: ")
	}
	output.Write(strconv.AppendFloat(numBuf[:0], total_price, 'f', 2, 64))
	output.WriteByte('\n')
	// TODO: add the discount mechanism after this
	// TODO: add a getDiscount function because the discount can be either:
	// 1) default value (probably 1.3)
	// 2) float value
	// 3) Tag value
}

// calculates the total price for the monjogs in [work]
// if [verbose] is set to true, then the price for each monjog will
// be displayed during calculation + the total price at the end
// it returns the final price for monjogs or an error if any of the prices could not be
// resolved (see [resolve] function for more details on resolving prices)
// NOTE: this function doesn't actually print anything, it recieves an output parameter
// which collects the impotant information inside an strings.Builder struct and the user
// can flush the information out later.
// NOTE: the ouput parameter will only be filled with information if the verbose flag
// is set to true, otherwise the final price will be returned so the user
// can decide what to do with it
func (work *Work) calcMonjogPrices_v2(output *strings.Builder, settings workWriterSettings) (total_price float64, resolve_error error) {
	total_price = 0
	for m_code, m_count := range work.Monjogs {
		p := db.Codes[m_code]
		price, err := p.resolve()
		if err != nil {
			fmt.Println("there was a problem with calculating the price of the monjogs")
			return 0, err
		}
		total_price += price * float64(m_count)
		if settings.verbose {
			// Reusable stack buffer for strconv formatting
			var numBuf [32]byte

			output.WriteString("monjog code: ")
			output.WriteString(m_code)

			output.WriteString("| monjog count: ")
			// %f equivalent: -1 precision keeps default floating representation, or specify precision
			output.Write(strconv.AppendFloat(numBuf[:0], m_count, 'f', -1, 64))

			output.WriteString("| monjog price: ")
			// %.2f equivalent: 2 decimal places
			output.Write(strconv.AppendFloat(numBuf[:0], price, 'f', 2, 64))

			output.WriteByte('\n')
		}
	}
	if settings.verbose {
		// Reusable stack buffer for strconv formatting
		var numBuf [32]byte
		output.WriteString("the total price for monjogs is: ")
		output.Write(strconv.AppendFloat(numBuf[:0], total_price, 'f', 2, 64))
	}
	return total_price, nil
}

// calculates the total price for the materials in [work]
// if [verbose] is set to true, then the price for each material will
// be displayed during calculation + the total price at the end
// it returns the final price for materials or an error if any of the prices could not be
// resolved (see [resolve] function for more details on resolving prices)
// NOTE: this function doesn't actually print anything, it recieves an output parameter
// which collects the impotant information inside an strings.Builder struct and the user
// can flush the information out later.
// NOTE: the ouput parameter will only be filled with information if the verbose flag
// is set to true, otherwise the final price will be returned so the user
// can decide what to do with it
func (work *Work) calcMaterialPrices_v2(output *strings.Builder, settings workWriterSettings) (total_price float64, resolve_error error) {
	total_price = 0
	//Materials map[string]Price
	if len(work.Materials) != 0 {
		for m_name, m_count := range work.Materials {
			p := db.OtherMaterials[m_name]
			price, err := p.resolve()
			if err != nil {
				fmt.Printf("there was a problem with calculating the price of the material: %s\n", m_name)
				return 0, err
			}
			if settings.verbose {
				// Reusable stack buffer for strconv formatting
				var numBuf [32]byte

				output.WriteString("monjog code: ")
				output.WriteString(m_name)

				output.WriteString("| material count: ")
				output.Write(strconv.AppendFloat(numBuf[:0], m_count, 'f', -1, 64))

				output.WriteString("| material price: ")
				output.Write(strconv.AppendFloat(numBuf[:0], price, 'f', 2, 64))

				output.WriteByte('\n')
			}
			total_price += price * float64(m_count)
		}
		if settings.verbose {
			// Reusable stack buffer for strconv formatting
			var numBuf [32]byte
			output.WriteString("the total price for materials is: ")
			output.Write(strconv.AppendFloat(numBuf[:0], total_price, 'f', 2, 64))
		}
	}
	return total_price, nil
}

// calculates the total price for the time in [work]
// if [verbose] is set to true, then more information about the time of work will
// be displayed during calculation + the total price at the end
// it returns the final price for materials or an error if the price could not be
// resolved (see [resolve] function for more details on resolving prices)
// NOTE: this function doesn't actually print anything, it recieves an output parameter
// which collects the impotant information inside an strings.Builder struct and the user
// can flush the information out later.
// NOTE: the ouput parameter will only be filled with information if the verbose flag
// is set to true, otherwise the final price will be returned so the user
// can decide what to do with it
func (work *Work) calcTimePrice_v2(output *strings.Builder, settings workWriterSettings) (total_price float64, resolve_error error) {
	// work.Time: float64
	time := work.Time

	//Time map[string]Price
	p := db.Time["price"]
	price, err := p.resolve()
	if err != nil {
		fmt.Println("there was a problem with calculating the price of the time")
		return 0, err
	}
	total_price = time * price
	if settings.verbose {
		// Reusable stack buffer for strconv formatting
		var numBuf [32]byte

		output.WriteString("total time for work: ")
		output.Write(strconv.AppendFloat(numBuf[:0], time, 'f', 2, 64))

		output.WriteString("| time price: ")
		output.Write(strconv.AppendFloat(numBuf[:0], price, 'f', 2, 64))

		output.WriteString("| total time price: ")
		output.Write(strconv.AppendFloat(numBuf[:0], total_price, 'f', 2, 64))

		output.WriteByte('\n')
	}
	return total_price, nil
}

func WriteWorksToFile() {
	if len(db.Works) == 0 {
		fmt.Println("there was no data here!")
		return
	}
	// make sure Time price and Default Price have been set, if not then terminate early
	if errors := runChecks(checkTime | checkDefaultPrice); len(errors) != 0 {
		fmt.Printf("the function cannot run due to the following errors:")
		for error_number, err := range errors {
			fmt.Printf("%d: %s", error_number+1, err.Error())
			return
		}
	}

	// try to allocate enough buffer (the size may change in the future)
	var output strings.Builder
	output.Grow(1024)
	// initialize the settings
	settings := workWriterSettings{
		verbose:               false, // don't write full information about each work
		writeOneDiscountValue: true,  // use only one discount value for all the works
	}
	// the header
	output.WriteString("work\t\tPrice\n")
	for key, work := range db.Works {
		output.WriteString(key)
		output.WriteString("\t\t")
		work.writeToBuffer(&output, settings) // don't write every detail only the final value
	}
	file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Could not create file, the output due to:\n%s\nexiting the function for now", err.Error())
		return
	}
	defer file.Close()
	// Wrap the target file directly in tabwriter
	// minwidth=0, tabwidth=0, padding=4, padchar=' ', flags=0
	tw := tabwriter.NewWriter(file, 0, 0, 4, ' ', 0)
	tw.Write([]byte(output.String()))

	tw.Flush()
}
