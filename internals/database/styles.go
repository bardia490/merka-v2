package database

import "github.com/gdamore/tcell/v2"

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
