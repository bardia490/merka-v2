package database

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func addTodoPage(page_name string) {
	// Create a styled TextView for the central "card" content
	todoWidget := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n\n[#cba6f7]🚧 TODO PAGE 🚧[-]\n\n[#6c7086]Feature under construction[-]\n\n[#89dceb]Press 'q' or Esc to go back[-]")

	todoWidget.SetBackgroundColor(colCardBg)
	todoWidget.Box.SetBorder(true).
		SetBorderColor(colAccent).
		SetTitle(" [ [#89dceb]Status[#cdd6f4] ] ").
		SetTitleColor(colText)

	// Handle key events to let the user exit back to the start page
	todoWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			App.pages.SwitchToPage("Start")
			return nil
		}
		return event
	})

	ok_button := tview.NewButton("Ok")
	// Style the button colors to match your Catppuccin theme
	ok_button.SetBackgroundColor(colCardBg)          // Normal background
	ok_button.SetBackgroundColorActivated(colAccent) // Background when focused/hovered
	ok_button.SetLabelColor(colText)                 // Normal text color
	ok_button.SetLabelColorActivated(colBg)          // Text color when focused (dark text on accent background looks super clean)

	ok_button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			App.pages.SwitchToPage("Start")
			return nil
		}
		return event
	})

	ok_button.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {

		if action == tview.MouseLeftClick {
			App.pages.SwitchToPage("Start")
			return tview.MouseConsumed, nil
		}
		return action, event
	})

	// Center the card nicely in the middle of the screen using nested Flex layouts
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Vertical spacer top
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(nil, 0, 1, false).        // Horizontal spacer left
				AddItem(todoWidget, 50, 0, true). // Card width: 50 columns
				AddItem(nil, 0, 1, false),        // Horizontal spacer right
			11, 0, true, // Card height: 11 rows
		).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(nil, 0, 1, false).        // Horizontal spacer left
				AddItem(ok_button, 10, 0, false). // Card width: 10 columns
				AddItem(nil, 0, 1, false),        // Horizontal spacer right
			1, 0, true, // button height: 2 rows
		).
		AddItem(nil, 0, 1, false) // Vertical spacer bottom

	// Optional: Ensure the background of the outer container matches your theme
	layout.SetBackgroundColor(colBg)

	// Add the page to your application's pages router
	App.pages.AddPage(page_name, layout, true, false)
}
