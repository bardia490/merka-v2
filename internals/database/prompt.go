package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var validInputs []string

// used for setting valid inputs and suggestions for the promptOption function
func SetValidInputs(inputs []string) {
	validInputs = inputs
}

// TODO: add comments
func addPromptingPage(page_name string) {
	// --- 1. Feedback / Status Label ---
	statusView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusView.SetBackgroundColor(colCardBg)
	statusView.SetText("[#6c7086]Type to search. Must match a valid language.[-]")

	// --- 2. Input Field ---
	input := tview.NewInputField().
		SetLabel("  Language: ").
		SetFieldWidth(32).
		SetPlaceholder("").                     //"e.g. Go, Zig, Rust..."
		SetPlaceholderTextColor(tcell.Color121) //colSubtext

	// Styling the input field
	input.SetLabelColor(colCyan)
	input.SetFieldTextColor(colText)
	input.SetFieldBackgroundColor(tcell.NewHexColor(0x313244))
	input.SetBackgroundColor(colCardBg)

	// Correct Autocomplete Styles signature:
	// func (i *InputField) SetAutocompleteStyles(background tcell.Color, main, selected tcell.Style)
	mainStyle := tcell.StyleDefault.Foreground(colText).Background(colCardBg)
	selectedStyle := tcell.StyleDefault.Foreground(colAccent).Background(tcell.NewHexColor(0x313244))

	input.SetAutocompleteStyles(colCardBg, mainStyle, selectedStyle)

	right_pane := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(colAccent).
		SetSelectedBackgroundColor(tcell.NewHexColor(0x313244)).
		SetMainTextColor(colText)

	right_pane.SetBorder(true).
		SetTitle(" Available suggestions (Click / Enter to choose) ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(colAccent).
		SetBorderColor(colAccent).
		SetBackgroundColor(colCardBg)

	right_pane.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseScrollUp:
			cur := right_pane.GetCurrentItem()
			if cur > 0 {
				right_pane.SetCurrentItem(cur - 1)
			}
			return action, nil // Event consumed
		case tview.MouseScrollDown:
			cur := right_pane.GetCurrentItem()
			if cur < right_pane.GetItemCount()-1 {
				right_pane.SetCurrentItem(cur + 1)
			}
			return action, nil // Event consumed
		}
		return action, event
	})
	right_pane.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			App.app.SetFocus(input)
			return nil
		}
		if event.Key() == tcell.KeyUp && right_pane.GetCurrentItem() == 0 {
			App.app.SetFocus(input)
			return nil
		}
		return event
	})
	updateSuggestions := func(currentText string) {
		trimmed := strings.TrimSpace(currentText)
		right_pane.Clear()

		count := 0
		for _, lang := range validInputs {
			if trimmed == "" || strings.HasPrefix(strings.ToLower(lang), strings.ToLower(trimmed)) {
				selectedLang := lang
				// AddItem sets the text and a callback fired on Click or Enter
				right_pane.AddItem("  • "+selectedLang, "", 0, func() {
					input.SetText(selectedLang) // Fill input with chosen language
					App.app.SetFocus(input)     // Return focus to input
				})
				count++
			}
		}

		if count == 0 {
			right_pane.AddItem("  [#f38ba8]No matching languages[-]", "", 0, nil)
		}
	}
	updateSuggestions("")
	// Autocomplete suggestion generator
	input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		currentText = strings.TrimSpace(currentText)
		if currentText == "" {
			return nil
		}
		for _, lang := range validInputs {
			if strings.HasPrefix(strings.ToLower(lang), strings.ToLower(currentText)) {
				entries = append(entries, lang)
			}
		}
		return entries
	})

	// Outer card container to handle the dynamic border color
	card := tview.NewFlex().
		SetDirection(tview.FlexRow)

	card.SetBorder(true).
		SetTitle(" 🚀 SELECT COMPILER / RUNTIME ").
		SetTitleColor(colAccent).
		SetBorderColor(colAccent).
		SetBackgroundColor(colCardBg).
		SetBorderPadding(1, 1, 2, 2)

	// --- 3. Real-Time Validator (Fires on every keystroke) ---
	input.SetChangedFunc(func(text string) {
		updateSuggestions(text)
		if strings.TrimSpace(text) == "" {
			card.SetBorderColor(colAccent)
			statusView.SetText("[#6c7086]Type to search. Must match a valid language.[-]")
			return
		}

		isValid, msg := currentStringValidator(text)
		if isValid {
			card.SetBorderColor(colSuccess)
			statusView.SetText(fmt.Sprintf("[%s]%s[-]", colSuccess.TrueColor().String(), msg))
		} else {
			card.SetBorderColor(colError)
			statusView.SetText(fmt.Sprintf("[%s]%s[-]", colError.TrueColor().String(), msg))
		}
	})

	// --- 4. Submit Validation (Fires on Enter) ---
	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			text := input.GetText()
			isValid, msg := currentStringValidator(text)

			if !isValid {
				card.SetBorderColor(colError)
				statusView.SetText(fmt.Sprintf("[%s]Cannot submit! %s[-]", colError.TrueColor().String(), msg))
				return
			}

			App.Results.IsValueSet = true
			App.Results.ValString = strings.TrimSpace(text)
			App.app.Stop()

		case tcell.KeyEscape:
			App.app.Stop()
		}
	})
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyTab {
			if right_pane.GetItemCount() > 0 {
				App.app.SetFocus(right_pane)
				return nil
			}
		}
		return event
	})

	// Keyboard hint banner
	hints := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	hints.SetBackgroundColor(colCardBg)
	hints.SetText("[#6c7086][Tab/↓] Autocomplete   [Enter] Submit   [Esc] Quit[-]")

	// Assemble card elements
	card.
		AddItem(input, 1, 0, true).
		AddItem(tview.NewBox().SetBackgroundColor(colCardBg), 1, 0, false).
		AddItem(statusView, 1, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(colCardBg), 1, 0, false).
		AddItem(hints, 1, 0, false)

	// Center the card in the terminal using Flex columns & rows
	layout := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(nil, 0, 1, false). // blank
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumnCSS). // set the main axis to columns
			AddItem(nil, 0, 1, false).
			AddItem(card, 9, 0, true).
			AddItem(nil, 0, 1, false),
			0, 1, true,
		).
		AddItem(right_pane, 0, 1, false) // blank

	layout.SetBackgroundColor(colBg)
	// capture func(event *tcell.EventKey) *tcell.EventKey

	//func (p *Pages) AddPage(name string, item Primitive, resize, visible bool) *Pages
	App.pages.AddPage(page_name, layout, true, false)
}

var activeMultiplierPrompt *tview.Flex

// addMultiplier adds a page that displays an interactive input modal on App.pages.
// When valid input is entered or confirmed, onDone is called with the calculated multiplier.
// If defaultMult is 0 or less, it will fall back to 1.0.
func addMultiplierPage() {
	const pageName = "MultiplierPrompt"

	// 1. Text description
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

	// 2. Status / Error Message
	statusView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusView.SetBackgroundColor(colCardBg)

	// 3. Input Field
	inputField := tview.NewInputField().
		SetLabel("Multiplier: ").
		SetFieldWidth(12)

	inputField.SetFieldBackgroundColor(colBg).
		SetFieldTextColor(colText).
		SetLabelColor(colCyan)

	// 4. Submit handler
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
		PrintWorkTable(mult) // TODO: this name will change in the future
		App.pages.SwitchToPage("Table")
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

	// 5. Buttons
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

	// 6. Modal Card
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
	// Thread-safe update to make sure it immediately renders
	App.pages.AddPage("MultiplierPrompt", activeMultiplierPrompt, true, false)
}
