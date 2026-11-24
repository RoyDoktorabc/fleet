package ui

import (
	"strings"
	"testing"

	"github.com/brizzai/fleet/internal/session"
	"github.com/brizzai/fleet/internal/tmux"
	"github.com/charmbracelet/x/ansi"
)

func TestPaintCaret(t *testing.T) {
	cases := []struct {
		name string
		line string
		x    int
		want string // the line's text once styling is stripped
	}{
		{"on a letter", "hello", 1, "hello"},
		{"past the end pads to the column", "ab", 5, "ab    "},
		{"right after the text", "❯ ", 2, "❯  "},
		{"on a wide glyph", "a世b", 1, "a世b"},
		{"on a styled line", "\x1b[31mred\x1b[0m text", 4, "red text"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := paintCaret(c.line, c.x)
			if plain := ansi.Strip(got); plain != c.want {
				t.Fatalf("text changed: got %q, want %q", plain, c.want)
			}
			if !strings.Contains(got, "\x1b[") {
				t.Fatalf("no styling applied: %q", got)
			}
		})
	}
}

func TestPaintCaretMarksTheRightCell(t *testing.T) {
	got := paintCaret("abc", 1)
	caret := FocusCaret().Render("b")
	if !strings.Contains(got, caret) {
		t.Fatalf("expected %q to hold the caret cell %q", got, caret)
	}
	if !strings.HasPrefix(ansi.Strip(strings.SplitN(got, caret, 2)[0]), "a") {
		t.Fatalf("caret is not after the first cell: %q", got)
	}
}

func TestRenderPreviewDrawsCaretOnCursorRow(t *testing.T) {
	s := session.NewSession("s", "/tmp/caret-test")
	content := "one\n❯ \nthree\n"
	caret := FocusCaret().Render(" ")

	withCaret := RenderPreview(s, content, nil, 30, 5, true, tmux.PaneCursor{X: 2, Y: 1, Visible: true})
	if !strings.Contains(withCaret, caret) {
		t.Fatalf("caret missing from the focused preview: %q", withCaret)
	}
	rows := strings.Split(withCaret, "\n")
	var caretRows []int
	for i, r := range rows {
		if strings.Contains(r, caret) {
			caretRows = append(caretRows, i)
		}
	}
	if len(caretRows) != 1 || !strings.Contains(ansi.Strip(rows[caretRows[0]]), "❯") {
		t.Fatalf("caret should sit on the ❯ row only, found rows %v in %q", caretRows, withCaret)
	}

	hidden := RenderPreview(s, content, nil, 30, 5, true, tmux.PaneCursor{X: 2, Y: 1, Visible: false})
	if strings.Contains(hidden, caret) {
		t.Fatalf("a hidden tmux cursor must not draw a caret: %q", hidden)
	}
}

func TestPreviewCaretGate(t *testing.T) {
	h, s := scrollTestHome(t)
	h.previewCursor[s.ID] = tmux.PaneCursor{X: 2, Y: 1, Visible: true}

	if got := h.previewCaret(s); got.Visible {
		t.Fatal("no caret while the split view is not focused")
	}
	h.focusMode = true
	if got := h.previewCaret(s); !got.Visible {
		t.Fatal("caret expected while focused")
	}
	h.previewScroll.offset = 3
	if got := h.previewCaret(s); got.Visible {
		t.Fatal("no caret while scrolled back")
	}
	h.previewScroll.offset = 0
	h.previewSel.dragging = true
	if got := h.previewCaret(s); got.Visible {
		t.Fatal("no caret while selecting")
	}
}

func TestCursorMoveRepaintsEvenWhenTextIsUnchanged(t *testing.T) {
	h, s := scrollTestHome(t)
	h.previewCache[s.ID] = "same"
	h.previewCursor[s.ID] = tmux.PaneCursor{X: 1, Y: 0, Visible: true}
	if !h.isIdlePreviewPoll(previewMsg{sessionID: s.ID, content: "same", cursor: tmux.PaneCursor{X: 1, Y: 0, Visible: true}}) {
		t.Fatal("identical text and cursor is an idle poll")
	}
	if h.isIdlePreviewPoll(previewMsg{sessionID: s.ID, content: "same", cursor: tmux.PaneCursor{X: 0, Y: 0, Visible: true}}) {
		t.Fatal("a moved cursor must repaint")
	}
}
