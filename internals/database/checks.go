package database

import "fmt"

type setting uint8

const (
	checkMonjogs      setting = 1 << iota // checks if all works have non-empty monjogs field
	checkMaterials                        // checks if all works have non-empty materials field
	checkTime                             // checks to see if the time price has been set
	checkTimes                            // checks to see if all the works have a dedicated work time (minute) count
	checkPriceTags                        // makes sure everywork that has a price-tag, its tag can be resolved
	checkDefaultPrice                     // checks to see if a default price for price has been set
)

// functions for manipulating the settings
func (s setting) Has(flag setting) bool { return s&flag != 0 }
func (s *setting) Set(flag setting)     { *s |= flag }
func (s *setting) Clear(flag setting)   { *s &^= flag }
func (s *setting) Toggle(flag setting)  { *s ^= flag }

// for checking if the time:price field has been set and time_price <= 0
func runTimeCheck() {
	p := db.Time["price"]
	price, err := p.resolve()
	if err != nil {
		fmt.Printf("there is a problem with time:price => %s\n", err.Error())
	} else if price <= 0 {
		fmt.Printf("the price for time was 0 or negative: %f", price)
	}
}

// checks to see if all the works have a dedicated work time (minute) count
func (work *Work) runTimeChecks(name string) {
	time_price := work.Time
	if price, err := time_price.resolve(); err != nil {
		fmt.Printf("in work: %s, time could not be set correctly with the following error: %s\n", name, err.Error())
	} else if price <= 0 {
		fmt.Printf("in work: %s, the price for time was either 0 or negative: %f\n", name, price)
	}
}

// checks if all works have non-empty monjogs field
func (work *Work) runMonjogsCheck(name string) {
	if len(work.Monjogs) == 0 {
		fmt.Printf("there are no monjogs in the work: %s", name)
	}
}

// checks if all works have non-empty materials field
func (work *Work) runMaterialsCheck(name string) {
	if len(work.Materials) == 0 {
		fmt.Printf("there are no materials in the work: %s", name)
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

func runChecks(s setting) {
}
