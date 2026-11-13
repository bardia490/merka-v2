package database

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NOTE: make sure to call addTablePage initialization of the program
// so the following variable can be set; otherwise it will cause a runtime panic
var activeTable *tview.Table
var TableRowsChan = make(chan [][]string)

// TODO: add comments
func addTablePage(page_name string) {
	activeTable = tview.NewTable().
		SetSelectable(true, false) // Select whole rows, not individual cells

	// for now were just choosing the column 2 (price with discount)
	activeTable.SetSelectedFunc(func(row, column int) {
		selectedCell := activeTable.GetCell(row, 2) // Column 1 is Price or whatever you need
		if selectedCell != nil {
			App.Results.ValString = selectedCell.Text
			App.Results.IsValueSet = true
		}
	})

	// 4. Handle quit (q or esc/ctrl+c)
	activeTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlC {
			App.pages.SwitchToPage("Start")
			return nil
		}
		return event
	})

	// Focus the table and set it as root
	activeTable.Select(1, 0).SetFixed(1, 0) // Start cursor on first data row
	activeTable.Box.SetBorder(true)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumnCSS).
				AddItem(nil, 0, 1, false).
				AddItem(activeTable, 10, 0, true).
				AddItem(nil, 0, 1, false),
			70, 0, true,
		).
		AddItem(nil, 0, 1, false)
	App.pages.AddPage(page_name, layout, true, false)
}

// this function is used to update the rows of the work table of neccasery
// the function is designed to run on a seperate goroutine before App.Run is called.
// IMPORTANT NOTE: the data must come from the TableRowsChan channel
// if no data is recieved it will result in a hanging channel and sleep goroutine
func UpdateTableContents() {
	for {
		var rows [][]string
		rows = <-TableRowsChan
		App.app.QueueUpdateDraw(func() {
			if activeTable == nil {
				return
			}
			// 1. Set Header Row (Row 0)
			headers := []string{
				"Work",
				"Price",
				"Price with multiplier",
			}
			for colIndex, header := range headers {
				cell := tview.NewTableCell(header).
					SetTextColor(tcell.ColorYellow).
					SetAlign(tview.AlignLeft).
					SetSelectable(false) // Headers shouldn't be selectable/focusable
				activeTable.SetCell(0, colIndex, cell)
			}

			// 2. Populate Data Rows (Starting at row 1)
			for rowIndex, rowData := range rows {
				for colIndex, val := range rowData {
					cell := tview.NewTableCell(val).
						SetTextColor(tcell.ColorWhite).
						SetAlign(tview.AlignLeft)

					// tview table takes 0-indexed rows where row 0 is our header,
					// so data starts at rowIndex + 1
					activeTable.SetCell(rowIndex+1, colIndex, cell)
				}
			}
		})
	}
}
