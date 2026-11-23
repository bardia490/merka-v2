package database

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: add comments
func addStartingPage(page_name string) {
	// 1. Header banner
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[#cba6f7::b]Navigation Dashboard[-:-:-]\n[#6c7086]Click an item or use Arrow keys + Enter[-]")
	header.SetBackgroundColor(colCardBg)

	// 2. Interactive List (Mouse Clickable & Keyboard Accessible)
	menuList := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetMainTextColor(colText).
		SetSelectedTextColor(colAccent).
		SetSelectedBackgroundColor(colSelected)

	menuList.SetBackgroundColor(colCardBg)

	// Target pages to switch to
	pagesList := []struct {
		ID    string
		Label string
	}{
		{ID: "MultiplierPrompt", Label: "Show work prices"},
		{ID: "B", Label: "Page B - Overview & Stats"},
		{ID: "C", Label: "Page C - Settings & Config"},
		{ID: "D", Label: "Page D - Logs & Activity"},
		{ID: "Exit", Label: "Exit the App"},
	}

	for _, p := range pagesList {
		targetID := p.ID
		menuList.AddItem(fmt.Sprintf("  ❯ %s", p.Label), "", 0, func() {
			if targetID == "Exit" {
				App.app.Stop() // TODO: add a clean up function for saving and stuff if needed
			}
			App.pages.SwitchToPage(targetID)
		})
	}

	// Global key overrides inside the card
	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			App.app.Stop()
			return nil
		}
		// Quick navigation with right arrow (e.g. to default prompt)
		if event.Key() == tcell.KeyRight {
			App.pages.SwitchToPage("MultiplierPrompt")
			return nil
		}
		return event
	})
	// add our own mouse capture
	menuList.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseScrollUp:
			cur := menuList.GetCurrentItem()
			if cur > 0 {
				menuList.SetCurrentItem(cur - 1)
			}
			return tview.MouseConsumed, nil

		case tview.MouseScrollDown:
			cur := menuList.GetCurrentItem()
			if cur < menuList.GetItemCount()-1 {
				menuList.SetCurrentItem(cur + 1)
			}
			return tview.MouseConsumed, nil
		}
		return action, event
	})
	// 3. Inner Card containing header + interactive items
	cardContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 4, 1, false).
		AddItem(menuList, len(pagesList)+1, 1, true)
	cardContent.SetBackgroundColor(colCardBg)

	cardContent.SetBorder(true).
		SetBorderColor(colSubtext).
		SetTitle("[#89dceb::b] Home [-:-:-]").
		SetTitleAlign(tview.AlignCenter)

	// 4. Centered Layout (Rows + Columns Flex)
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Top spacer
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).         // Left spacer
			AddItem(cardContent, 46, 1, true). // Centered Card (width 46)
			AddItem(nil, 0, 1, false),         // Right spacer
			len(pagesList)+7, 1, true, // Height dynamically scaled to items
		).
		AddItem(nil, 0, 1, false) // Bottom spacer

	layout.SetBackgroundColor(colBg)

	App.pages.AddPage(page_name, layout, true, true)
	App.pages.SwitchToPage(page_name)
}
