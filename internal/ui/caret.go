package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// paintCaret draws the focus caret on column x of a rendered line. Past the end
// of the line it pads with spaces so the caret still lands on the right cell.
func paintCaret(line string, x int) string {
	cur := FocusCaret()
	lw := lipgloss.Width(line)
	if x < 0 || x >= lw {
		pad := x - lw
		if pad < 0 {
			pad = 0
		}
		return line + strings.Repeat(" ", pad) + cur.Render(" ")
	}
	// The grapheme under the caret may be double-width (CJK/emoji/wide box
	// drawing). A 1-cell cut of a wide glyph returns "" (it can't fit), which
	// would blank the caret and leave the glyph in the right-hand remainder,
	// garbling the row. Detect the wide case and slice by the full width.
	cellW := 1
	if ansi.Strip(ansi.Cut(line, x, x+1)) == "" &&
		ansi.StringWidth(ansi.Cut(line, x, x+2)) == 2 {
		cellW = 2
	}
	cell := ansi.Strip(ansi.Cut(line, x, x+cellW))
	if cell == "" {
		cell = " "
	}
	return ansi.Truncate(line, x, "") + cur.Render(cell) + ansi.TruncateLeft(line, x+cellW, "")
}
