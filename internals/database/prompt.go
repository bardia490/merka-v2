package database

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: at this point the prompting functions below are only used for
// getting user input based on some suggestions so it might be better to
// rename them to PromptOption or something

// TODO: this (probably) isn't needed anymore
// variable for setting the possible answers (options) for prompts
var validInputs []string

// A channel for recieving new version of the prompt options
var PromptOptionsChan = make(chan []string)

var suggestionsPane *tview.List
var input *tview.InputField

// updates what the user sees on the suggestions (right) pane
func updateSuggestions(currentText string) {
	trimmed := strings.TrimSpace(currentText)
	suggestionsPane.Clear()

	count := 0
	for _, lang := range validInputs {
		if trimmed == "" || strings.HasPrefix(strings.ToLower(lang), strings.ToLower(trimmed)) {
			selectedLang := lang
			// AddItem sets the text and a callback fired on Click or Enter
			suggestionsPane.AddItem("  • "+selectedLang, "", 0, func() {
				input.SetText(selectedLang) // Fill input with chosen option
				App.app.SetFocus(input)     // Return focus to input
			})
			count++
		}
	}

	if count == 0 {
		suggestionsPane.AddItem("  [#f38ba8]No matching Items[-]", "", 0, nil)
	}
}

// this function is used to update the suggestions of the prompting functions table of necessary
// the function is designed to run on a seperate goroutine before App.Run is called.
// IMPORTANT NOTE: the data must come from the UpdatePromptOptions channel
// if no data is recieved it will result in a hanging channel and a sleep goroutine
func UpdatePromptOptions() {
	for {
		validInputs = <-PromptOptionsChan
		App.app.QueueUpdateDraw(func() {
			// in the begining the input field is empty
			// so we don't suggest anything
			updateSuggestions("")
		})
	}
}

// TODO: add comments
func addPromptingPage(page_name string) {
	// Feedback / Status Label
	statusView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusView.SetBackgroundColor(colCardBg)
	statusView.SetText("[#6c7086]Type to search. Must match a valid input.[-]")

	// Input Field
	input = tview.NewInputField().
		SetLabel("  Language: ").
		SetFieldWidth(32).
		SetPlaceholder("").
		SetPlaceholderTextColor(tcell.Color121) //colSubtext

	// Styling the input field
	input.SetLabelColor(colCyan)
	input.SetFieldTextColor(colText)
	input.SetFieldBackgroundColor(tcell.NewHexColor(0x313244))
	input.SetBackgroundColor(colCardBg)

	mainStyle := tcell.StyleDefault.Foreground(colText).Background(colCardBg)
	selectedStyle := tcell.StyleDefault.Foreground(colAccent).Background(tcell.NewHexColor(0x313244))

	input.SetAutocompleteStyles(colCardBg, mainStyle, selectedStyle)

	suggestionsPane = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedTextColor(colAccent).
		SetSelectedBackgroundColor(tcell.NewHexColor(0x313244)).
		SetMainTextColor(colText)

	suggestionsPane.SetBorder(true).
		SetTitle(" Available suggestions (Click / Enter to choose) ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(colAccent).
		SetBorderColor(colAccent).
		SetBackgroundColor(colCardBg)

	suggestionsPane.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseScrollUp:
			cur := suggestionsPane.GetCurrentItem()
			if cur > 0 {
				suggestionsPane.SetCurrentItem(cur - 1)
			}
			return action, nil // Event consumed
		case tview.MouseScrollDown:
			cur := suggestionsPane.GetCurrentItem()
			if cur < suggestionsPane.GetItemCount()-1 {
				suggestionsPane.SetCurrentItem(cur + 1)
			}
			return action, nil // Event consumed
		}
		return action, event
	})
	suggestionsPane.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			App.app.SetFocus(input)
			return nil
		}
		if event.Key() == tcell.KeyUp && suggestionsPane.GetCurrentItem() == 0 {
			App.app.SetFocus(input)
			return nil
		}
		return event
	})
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
		SetTitle(" 🚀 SELECT AN OPTION ").
		SetTitleColor(colAccent).
		SetBorderColor(colAccent).
		SetBackgroundColor(colCardBg).
		SetBorderPadding(1, 1, 2, 2)

	//  Real-Time Validator (Fires on every keystroke)
	input.SetChangedFunc(func(text string) {
		updateSuggestions(text)
		if strings.TrimSpace(text) == "" {
			card.SetBorderColor(colAccent)
			statusView.SetText("[#6c7086]Type to search. Must match a valid Option.[-]")
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

	// Submit Validation (Fires on Enter)
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
			App.pages.SwitchToPage("TODO")
		}
	})
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyTab {
			if suggestionsPane.GetItemCount() > 0 {
				App.app.SetFocus(suggestionsPane)
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
		AddItem(suggestionsPane, 0, 1, false) // blank

	layout.SetBackgroundColor(colBg)
	// capture func(event *tcell.EventKey) *tcell.EventKey

	//func (p *Pages) AddPage(name string, item Primitive, resize, visible bool) *Pages
	App.pages.AddPage(page_name, layout, true, false)
}
