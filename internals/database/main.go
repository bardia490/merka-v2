package database

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"merka-v2/pkg/utility"
	"os"
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
			return fmt.Errorf("there was a problem whith parsing the json file at line: %d, column: %d. the ERROR: %s", line, col, err.Error())
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

// prints all the works names + monjog + material + time costs
func PrintFullDataBase() {
	printWorks(db.Works, true)
}

// this should be used when we want to print a work based on its name and we also want
// to check if the work exists beforehand (returns error if the work doesn't exist)
func PrintWork(workname string) error {
	// if the work existed inside the database
	if work, ok := db.Works[workname]; ok {
		fmt.Println("work:", workname)
		work.print(false)
		return nil
	}
	return fmt.Errorf("could not find work: %s in database", workname)
}

// prints all the works in database using the Print method for each work
func printWorks(works map[string]Work, verbose bool) {
	if len(works) == 0 {
		fmt.Println("there was no data here!")
		return
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
// if verbose is set to true then the whole information about each monjog (including count, price and total price) is displayed
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
			fmt.Printf("monjog code: %s| monjog count: %d| monjog price: %.2f| total price: %.2f\n", m_code, m_count, price, float64(m_count)*price)
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
				fmt.Printf("material name: %s| material price: %d| material price: %f.2| total price: %.2f\n", m_name, m_count, price, float64(m_count)*price)
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
	t := work.Time
	time, err := t.resolve()
	if err != nil {
		fmt.Println("there was a problem with calculating the amount of the time")
		return 0, err
	}

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
// every -1 into the default price value
// if none of the above matches it will return an error
func (price *Price) resolve() (float64, error) {
	switch price.Kind {
	case PriceValue:
		result := price.Value
		if result == -1.0 { // for legacy default value
			result = getDefaultPrice() // this never fails
		}
		return result, nil
	case PriceTag:
		result := 0.0
		if val, ok := db.PriceTags[price.Tag]; ok {
			result = val
		} else { // if the tag didn't exist or wasn't set, return an error
			return 0, fmt.Errorf("could not resolve the price: %v. the Tag: %s could not be found, please add it", *price, price.Tag)
		}
		return result, nil
	}
	return 0, fmt.Errorf("could not resolve the price: %v. please contact your boyfriend immediatley", *price)
}

// this function assumes a code with the name [default] always exists and its non negative
// NOTE: Should probably add a checker function to make sure this is always the case
func getDefaultPrice() float64 {
	defaultPrice := db.Codes["default"]
	val, _ := defaultPrice.resolve()
	return val
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
