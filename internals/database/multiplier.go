package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// this probably doesn't even need to be global since the multiplier doesn't get updated unlike the other widgets
var activeMultiplierPrompt *tview.Flex

// addMultiplier adds a page that displays an interactive input modal on App.pages.
// When valid input is entered or confirmed, onDone is called with the calculated multiplier.
// If defaultMult is 0 or less, it will fall back to 1.0
// next_action is a function that will be called when the multiplier value is submited
// and can be used to switch to a new page or do some work with the value. if the
// argument is not needed user can just pass in nil.
// the multiplier value can be acceced in the other function using the App.Results
// global variable. this argument is usefull for chaining operations if needed.
// Example use cases: printing the work table after getting the multiplier value
// or writing the work prices to file after geting the multiplier value.
func addMultiplierPage(pageName string, next_action func()) {
	// Text description
	infoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	infoView.SetBackgroundColor(colCardBg)

	defaultMult := 1.2
	defaultDisplay := fmt.Sprintf("%.2f", defaultMult)
	if defaultMult <= 0 {
		defaultDisplay = "1.00"
	}

	infoText := fmt.Sprintf(
		"[#cba6f7]Enter Price Multiplier[-]\n\n"+
			"[#cdd6f4]• [#a6e3a1]> 1.0[-] (e.g. [#89dceb]1.25[-] = +25%%)\n"+
			"• [#f38ba8]< 1.0[-] (e.g. [#89dceb]0.75[-] = -25%%)\n"+
			"• [#fab387]Enter[-] on empty for default ([#89dceb]%s[-])\n"+
			"• [#6c7086]0 or 1 keeps price unchanged.[-]",
		defaultDisplay,
	)
	infoView.SetText(infoText)

	// Status / Error Message
	statusView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusView.SetBackgroundColor(colCardBg)

	// Input Field
	inputField := tview.NewInputField().
		SetLabel("Multiplier: ").
		SetFieldWidth(12)

	inputField.SetFieldBackgroundColor(colBg).
		SetFieldTextColor(colText).
		SetLabelColor(colCyan)

	// Submit handler
	submit := func() {
		raw := strings.TrimSpace(inputField.GetText())
		var mult float64

		if len(raw) == 0 {
			mult = defaultMult
		} else {
			parsed, err := strconv.ParseFloat(raw, 64)
			if err != nil || parsed < 0 {
				statusView.SetText("[#f38ba8]⚠ Invalid number format[-]")
				return
			}
			mult = parsed
		}

		if mult <= 0 {
			mult = 1.0
		}

		App.Results.IsValueSet = true
		App.Results.ValFloat = mult
		inputField.SetText("")
		if next_action != nil {
			next_action()
		}
	}

	cancel := func() {
		App.pages.SwitchToPage("Start")
	}

	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			submit()
		case tcell.KeyEscape:
			cancel()
		}
	})

	// Buttons
	btnOk := tview.NewButton("Confirm").SetSelectedFunc(submit)
	btnCancel := tview.NewButton("Cancel").SetSelectedFunc(cancel)

	btnOk.SetBackgroundColor(colAccent)
	btnOk.SetLabelColor(colBg)
	btnCancel.SetBackgroundColor(colCardBg)
	btnOk.SetLabelColor(colText)

	btnRow := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(nil, 0, 1, false).
		AddItem(btnOk, 11, 0, false).
		AddItem(nil, 2, 0, false).
		AddItem(btnCancel, 11, 0, false).
		AddItem(nil, 0, 1, false)

	// Modal Card
	// Total height = 5 (info) + 1 (spacing) + 1 (input) + 1 (status) + 1 (buttons) + 2 (borders) = ~11 rows
	card := tview.NewFlex().
		SetDirection(tview.FlexColumnCSS).
		AddItem(infoView, 5, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(inputField, 1, 0, true).
		AddItem(statusView, 1, 0, false).
		AddItem(btnRow, 1, 0, false)

	card.SetBorder(true).
		SetTitle(" Multiplier ").
		SetTitleColor(colAccent).
		SetBorderColor(colAccent).
		SetBackgroundColor(colCardBg)

	activeMultiplierPrompt = tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumnCSS).
				AddItem(nil, 0, 1, false).
				AddItem(card, 16, 0, true).
				AddItem(nil, 0, 1, false),
			33, 0, true,
		).
		AddItem(nil, 0, 1, false)
	App.pages.AddPage(pageName, activeMultiplierPrompt, true, false)
}
