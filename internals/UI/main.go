package UI

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Catppuccin / Modern Terminal Palette
var (
	colBg      = tcell.NewHexColor(0x1e1e2e) // Dark slate background
	colCardBg  = tcell.NewHexColor(0x181825) // Inner card background
	colAccent  = tcell.NewHexColor(0xcba6f7) // Mauve / Purple
	colCyan    = tcell.NewHexColor(0x89dceb) // Cyan label
	colText    = tcell.NewHexColor(0xcdd6f4) // Primary text
	colSuccess = tcell.NewHexColor(0xa6e3a1) // Emerald Green
	colError   = tcell.NewHexColor(0xf38ba8) // Coral Red
	colSubtext = tcell.NewHexColor(0x6c7086) // Muted gray
)

// used to specify the return value of the widget (I may change this to use interfaces or something else)
type Result struct {
	valString  string
	valInt     int
	valFloat   float64
	isValueSet bool
}

type Application struct {
	app    *tview.Application // the actual application responsible for running the program
	pages  *tview.Pages       // this field will be set at the start of the program to know what pages should the program draw
	result Result             // sets the final result value of the widget
}

var App Application

// this should be called Immediatly
func Init() {
	App.app = tview.NewApplication()
	App.pages = tview.NewPages()
	// TODO: rename the following names
	addStartingPage("Start")
	addPromptingPage("Select")
	addTablePage("Table")
	App.pages.SwitchToPage("Table")
}

func SwitchToPage(page_name string) {
	App.pages.SwitchToPage(page_name)
}

func Run() error {
	return App.app.SetRoot(App.pages, true).EnableMouse(true).Run()
}
