package database

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
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
	// resolve the multiplier value
	var multiplier float64
	multiplier = getMultiplierValueFromUser()
	settings := workWriterSettings{
		verbose:                true,
		writeOneDiscountValue:  true,
		defaultMultiplierValue: multiplier,
	}
	printWorks(db.Works, settings)
}

// this should be used when we want to print a work based on its name and we also want
// to check if the work exists beforehand (returns error if the work doesn't exist)
// [settings] defines additional information about how the work should be printed
// for for information about [workWriterSettings] check the definitions.go file
func PrintWork(workname string, settings workWriterSettings) error {
	// first make sure the following things exist
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
		var output strings.Builder
		output.Grow(512)

		// "work:", workname
		output.WriteString("work: ")
		output.WriteString(workname)
		output.WriteByte('\n')
		work.writeToBuffer(&output, settings)
		// Wrap the target file directly in tabwriter
		// minwidth=0, tabwidth=0, padding=4, padchar=' ', flags=0
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		tw.Write([]byte(output.String()))

		tw.Flush()
		return nil
	}
	return fmt.Errorf("could not find work: %s in database", workname)
}

// prints all the works in database using the Print method for each work
// if verbose is set to true then the whole information about each monjog
// (including count, price and total price) is displayed
// [settings] defines additional information about how the work should be printed
// for for information about [workWriterSettings] check the definitions.go file
// NOTE: for this function the [settings.verbose] flag should most probably be set to true
func printWorks(works map[string]Work, settings workWriterSettings) {
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

	var output strings.Builder
	output.Grow(1024)
	for workname, work := range works {
		// "work:", workname
		output.WriteString("work: ")
		output.WriteString(workname)
		output.WriteByte('\n')
		work.writeToBuffer(&output, settings)
	}
	// Wrap the target file directly in tabwriter
	// minwidth=0, tabwidth=0, padding=4, padchar=' ', flags=0
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	tw.Write([]byte(output.String()))

	tw.Flush()
}

// if a work exists inside the database we can use this function to print its
// monjogs, materials and time (if unsure about the existence of the work,
// should use DataBase.PrintWork function)
// if verbose is set to true then the whole information about each monjog
// (including count, price and total price) is displayed
// NOTE: this function literally prints the contents of the work
// to the screen using the [settings] provided by the [workWriterSettings] struct
// if this is not desired behavior, should probably check the writeToBuffer method
// for works
func (work *Work) print(output *strings.Builder, settings workWriterSettings) {
	work.writeToBuffer(output, settings)
	// Wrap the target file directly in tabwriter
	// minwidth=0, tabwidth=0, padding=4, padchar=' ', flags=0
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	tw.Write([]byte(output.String()))

	tw.Flush()
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

// if a work exists inside the database we can use this function to print its
// monjogs, materials and time (if unsure about the existence of the work,
// should use DataBase.PrintWork function which takes a work_name parameter)
// if verbose is set to true then the whole information about each monjog
// (including count, price and total price) is displayed
// NOTE: this function doesn't actually print anything, it recieves an output parameter
// which collects the impotant information inside an strings.Builder struct and the user
// can flush the information out later.
func (work *Work) writeToBuffer(output *strings.Builder, settings workWriterSettings) {
	total_price, monjog_price, material_price, time_price := 0.0, 0.0, 0.0, 0.0
	var err error
	// print the complete monjogs prices and counts
	if monjog_price, err = work.calcMonjogPrices(output, settings); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the monjogs, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += monjog_price
	if settings.verbose {
		output.WriteString("<><><><><><><><><><><><><><><><><><><><><>\n")
	}
	// print the materials
	if material_price, err = work.calcMaterialPrices(output, settings); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the materials, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += material_price
	if settings.verbose {
		output.WriteString("<><><><><><><><><><><><><><><><><><><><><>\n")
	}
	// print the time
	if time_price, err = work.calcTimePrice(output, settings); err != nil {
		fmt.Printf("ERROR: there was a problem can calculating the price of the time, please correct the mistake and try again later. the error: %s.\n Aborting Operation\n", err.Error())
		return
	}
	total_price += time_price
	if settings.verbose {
		output.WriteString("<><><><><><><><><><><><><><><><><><><><><>\n")
	}
	var numBuf [32]byte
	if settings.verbose {
		output.WriteString("the total price for work is: ")
	}
	output.Write(strconv.AppendFloat(numBuf[:0], total_price, 'f', 2, 64))
	if !settings.verbose {
		output.WriteByte('\t')
	}

	// resolve the multiplier value
	var multiplier float64
	// check to see if the multiplier value has been set
	if settings.writeOneDiscountValue {
		multiplier = settings.defaultMultiplierValue
	} else {
		multiplier = getMultiplierValueFromUser()
	}
	total_price *= multiplier
	// the price of work after discount
	if settings.verbose {
		output.WriteString("the total price for work after discount is: ")
	}
	output.Write(strconv.AppendFloat(numBuf[:0], total_price, 'f', 2, 64))
	output.WriteByte('\n')
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
func (work *Work) calcMonjogPrices(output *strings.Builder, settings workWriterSettings) (total_price float64, resolve_error error) {
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

			output.WriteString("\tmonjog count: ")
			// %f equivalent: -1 precision keeps default floating representation, or specify precision
			output.Write(strconv.AppendFloat(numBuf[:0], m_count, 'f', -1, 64))

			output.WriteString("\tmonjog price: ")
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
		output.WriteByte('\n')
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
func (work *Work) calcMaterialPrices(output *strings.Builder, settings workWriterSettings) (total_price float64, resolve_error error) {
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

				output.WriteString("material name: ")
				output.WriteString(m_name)

				output.WriteString("\tmaterial count: ")
				output.Write(strconv.AppendFloat(numBuf[:0], m_count, 'f', -1, 64))

				output.WriteString("\tmaterial price: ")
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
			output.WriteByte('\n')
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
func (work *Work) calcTimePrice(output *strings.Builder, settings workWriterSettings) (total_price float64, resolve_error error) {
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

		output.WriteString("\ttime price: ")
		output.Write(strconv.AppendFloat(numBuf[:0], price, 'f', 2, 64))

		output.WriteString("\ttotal time price: ")
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

	// set one defaultMultiplierValue value for all the works
	multiplier := getMultiplierValueFromUser()
	// initialize the settings
	settings := workWriterSettings{
		verbose:                false, // don't write full information about each work
		writeOneDiscountValue:  true,  // use only one discount value for all the works
		defaultMultiplierValue: multiplier,
	}
	// the header
	output.WriteString("work\t\tPrice\tprice with multiplier\n")
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
