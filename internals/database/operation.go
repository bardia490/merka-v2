package database

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var activeMessage *tview.TextView
var OperationNameChan = make(chan string)
var OperationErrorChan = make(chan error)

// TOOD: use the following struct instead of using two channels for the Name and Err
// Single struct channel prevents race conditions and message desyncs
type OperationStatus struct {
	Name string
	Err  error
}

// NOTE: this function needs to be called in a seperate goroutine
// to handle updates at the very begining before calling App.app.Run()
func UpdateStatusMessage() {
	for {
		var operationName string
		var err error
		operationName = <-OperationNameChan
		err = <-OperationErrorChan
		App.app.QueueUpdateDraw(func() {
			if err == nil {
				// Success State
				activeMessage.Box.
					SetBorderColor(colSuccess).
					SetTitle(" [ #a6e3a1]Success[#cdd6f4] ] ").
					SetTitleColor(colSuccess)

				textMsg := fmt.Sprintf(
					"\n\n[#a6e3a1]✔ Operation successful[-]\n\n[#cdd6f4]Name: [#89dceb]%s[-]\n\n[#6c7086]Press 'q', Esc, or Enter to go back[-]",
					tview.Escape(operationName),
				)
				activeMessage.SetText(textMsg)
			} else {
				// Failure / Error State
				activeMessage.Box.
					SetBorderColor(colError).
					SetTitle(" [ #f38ba8]Error[#cdd6f4] ] ").
					SetTitleColor(colError)

				textMsg := fmt.Sprintf(
					"\n[#f38ba8]✖ Operation failed[-]\n\n[#cdd6f4]Name: [#89dceb]%s[-]\n[#f38ba8]Error: %s[-]\n\n[#6c7086]Press 'q', Esc, or Enter to go back[-]",
					tview.Escape(operationName),
					tview.Escape(err.Error()),
				)
				activeMessage.SetText(textMsg)
			}
		})
	}
}

func addOperationStatusPage(pageName string) {
	activeMessage = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	activeMessage.SetBackgroundColor(colCardBg)
	activeMessage.Box.SetBorder(true)

	// Navigation: Return to "Start" page on 'q', Esc, or Enter
	activeMessage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyCtrlC {
			App.pages.SwitchToPage("Start")
			return nil
		}
		return event
	})

	// Center the card nicely in the middle of the screen
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(nil, 0, 1, false).
				AddItem(activeMessage, 54, 0, true). // Card width
				AddItem(nil, 0, 1, false),
			13, 0, true, // Card height
		).
		AddItem(nil, 0, 1, false)

	layout.SetBackgroundColor(colBg)

	// Register page
	App.pages.AddPage(pageName, layout, true, false)
}
