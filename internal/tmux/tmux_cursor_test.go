package tmux

import "testing"

func TestSplitCaptureCursor(t *testing.T) {
	cases := []struct {
		name        string
		out         string
		wantContent string
		wantCursor  PaneCursor
		wantOK      bool
	}{
		{"visible cursor", "line one\n❯ \n2 1 1\n", "line one\n❯ \n", PaneCursor{X: 2, Y: 1, Visible: true}, true},
		{"hidden cursor", "a\n0 0 0\n", "a\n", PaneCursor{X: 0, Y: 0, Visible: false}, true},
		{"blank pane keeps its rows", "\n\n\n5 2 1\n", "\n\n\n", PaneCursor{X: 5, Y: 2, Visible: true}, true},
		{"empty capture", "0 0 1\n", "", PaneCursor{X: 0, Y: 0, Visible: true}, true},
		{"cursor line missing", "just text\n", "", PaneCursor{}, false},
		{"garbled cursor line", "text\nx y z\n", "", PaneCursor{}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			content, cur, ok := splitCaptureCursor(c.out)
			if ok != c.wantOK || content != c.wantContent || cur != c.wantCursor {
				t.Fatalf("got (%q, %+v, %v), want (%q, %+v, %v)", content, cur, ok, c.wantContent, c.wantCursor, c.wantOK)
			}
		})
	}
}
