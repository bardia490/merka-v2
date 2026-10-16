package UI

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TableMain renders columns and rows using tview
func TableMain(headers []string, rows [][]string) {
	app := tview.NewApplication()

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false) // Select whole rows, not individual cells

	// 1. Set Header Row (Row 0)
	for colIndex, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false) // Headers shouldn't be selectable/focusable
		table.SetCell(0, colIndex, cell)
	}

	// 2. Populate Data Rows (Starting at row 1)
	for rowIndex, rowData := range rows {
		for colIndex, val := range rowData {
			cell := tview.NewTableCell(val).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft)

			// tview table takes 0-indexed rows where row 0 is our header,
			// so data starts at rowIndex + 1
			table.SetCell(rowIndex+1, colIndex, cell)
		}
	}

	// 3. Handle selection/Enter key (equivalent to your "Let's go to..." action)
	table.SetSelectedFunc(func(row, column int) {
		selectedCell := table.GetCell(row, 1) // Column 1 is Price or whatever you need
		if selectedCell != nil {
			// You can trigger actions or print info here.
			// tview doesn't have tea.Printf directly, but we can write to app or exit with info.
			app.Stop()
			fmt.Printf("Let's go to %s!\n", selectedCell.Text)
		}
	})

	// 4. Handle quit (q or esc/ctrl+c)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyCtrlC {
			app.Stop()
			os.Exit(0)
		}
		return event
	})

	// Focus the table and set it as root
	table.Select(1, 0).SetFixed(1, 0) // Start cursor on first data row

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
