package database

import "fmt"

type setting uint8

const (
	checkMonjogs      setting = 1 << iota // checks if all works have non-empty monjogs field
	checkMaterials                        // checks if all works have non-empty materials field
	checkTime                             // checks to see if the time price has been set
	checkTimes                            // checks to see if all the works have a dedicated work time (minute) count
	checkPrices                           // makes sure every field that can be a price, can be coereced into a price, either be parsing to a number or resolving the price tag
	checkDefaultPrice                     // checks to see if a default price for price has been set
)

// functions for manipulating the settings
func (s setting) Has(flag setting) bool { return s&flag != 0 }
func (s *setting) Set(flag setting)     { *s |= flag }
func (s *setting) Clear(flag setting)   { *s &^= flag }
func (s *setting) Toggle(flag setting)  { *s ^= flag }

// for checking if the time:price field has been set and time_price >= 0
func runTimePriceCheck() {
	p := db.Time["price"]
	price, err := p.resolve()
	if err != nil {
		fmt.Printf("there is a problem with time:price => %s\n", err.Error())
	} else if price <= 0 {
		fmt.Printf("the price for time was 0 or negative: %f\n", price)
	}
}

// checks to see if the work has a dedicated work time (minute) count
func (work *Work) runTimePriceChecks(name string) {
	time_price := work.Time
	if price, err := time_price.resolve(); err != nil {
		fmt.Printf("in work: %s, time could not be set correctly with the following error: %s\n", name, err.Error())
	} else if price <= 0 {
		fmt.Printf("in work: %s, the price for time was either 0 or negative: %f\n", name, price)
	}
}

// checks if the work has a non-empty monjogs field
func (work *Work) runMonjogsCheck(name string) {
	if len(work.Monjogs) == 0 {
		fmt.Printf("there are no monjogs in the work: %s\n", name)
	}
}

// checks if the work has a non-empty materials field
func (work *Work) runMaterialsCheck(name string) {
	if len(work.Materials) == 0 {
		fmt.Printf("there are no materials in the work: %s\n", name)
	}
}

// checks to see if a default price for price has been set
func runCheckDefaultPrice() {
	_, code_ok := db.Codes["default"]
	_, tag_ok := db.PriceTags["default"]
	if !code_ok && !tag_ok {
		fmt.Println("could not find the default value in the database")
	}
}

func runChecks(s setting) { // NOTE: add seperators
	f_checkMonjogs := s.Has(checkMonjogs)
	f_checkMaterials := s.Has(checkMaterials)
	f_checkTimes := s.Has(checkTimes)
	f_checkPrices := s.Has(checkPrices)
	f_checkTime := s.Has(checkTime)
	f_checkDefaultPrice := s.Has(checkDefaultPrice)

	if f_checkDefaultPrice {
		runCheckDefaultPrice()
	}
	if f_checkTime {
		runTimePriceCheck()
	}

	for work_name, work := range db.Works {
		if f_checkMonjogs {
			work.runMonjogsCheck(work_name)
		}
		if f_checkMaterials {
			work.runMaterialsCheck(work_name)
		}
		if f_checkTimes {
			work.runTimePriceChecks(work_name)
		}
	}
	if f_checkPrices {
		for code_name, price := range db.Codes {
			if _, err := price.resolve(); err != nil {
				fmt.Printf("could not determine the price for Code: %s. ERROR: %s\n", code_name, err.Error())
			}
		}
		for material_name, price := range db.OtherMaterials {
			if _, err := price.resolve(); err != nil {
				fmt.Printf("could not determine the price for Material: %s. ERROR: %s\n", material_name, err.Error())
			}
		}
		additional_costs_price := db.AdditionalCosts["price"]
		if _, err := additional_costs_price.resolve(); err != nil {
			fmt.Printf("could not determine the price for additional_costs:price. ERROR: %s\n", err.Error())
		}
	}
}

func calculateLineAndColumnJsonError(data []byte, offset int64) (line, col int) {
	line, col = 1, 1
	for i := range offset {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}
