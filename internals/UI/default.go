package UI

import "github.com/rivo/tview"

// TODO: add comments
func addStartingPage(page_name string) {
	main_widget := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n\n[#cba6f7]This is the default Page[-]\n\n[#6c7086]Press Right Arrow to go to selection[-]")
	main_widget.SetBackgroundColor(colCardBg)
	//main_widget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	//	return event
	//})
	layout := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(nil, 0, 1, false). // blank
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumnCSS). // set the main axis to columns
			AddItem(nil, 0, 1, false).
			AddItem(main_widget, 7, 1, true).
			AddItem(nil, 0, 1, false),
			0, 1, true,
		).
		AddItem(nil, 0, 1, false) // blank
	layout.SetBackgroundColor(colBg)
	//func (p *Pages) AddPage(name string, item Primitive, resize, visible bool) *Pages
	App.pages.AddPage(page_name, layout, true, true)
}
