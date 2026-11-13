package database

import "github.com/rivo/tview"

// used to specify the return value of the widget (I may change this to use interfaces or something else)
type Result struct {
	ValString  string
	ValInt     int
	ValFloat   float64
	IsValueSet bool
}

type Application struct {
	app     *tview.Application // the actual application responsible for running the program
	pages   *tview.Pages       // this field will be set at the start of the program to know what pages should the program draw
	Results Result             // sets the final result value of the widget
}

var App Application

func SwitchToPage(page_name string) {
	App.pages.SwitchToPage(page_name)
}

func Run() error {
	return App.app.SetRoot(App.pages, true).EnableMouse(true).Run()
}
