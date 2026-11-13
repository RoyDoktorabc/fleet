package ui

import (
	"testing"

	"github.com/brizzai/fleet/internal/session"
)

// An errored row reads inactive like a suspended one: dim, not underlined.
func TestErrorRowTitleIsDim(t *testing.T) {
	got := TitleStyleForStatus(session.StatusError)
	want := TitleStyleForStatus(session.StatusSuspended)
	if got.GetForeground() != want.GetForeground() || got.GetUnderline() {
		t.Fatalf("error title fg=%v underline=%v, want suspended fg=%v and no underline",
			got.GetForeground(), got.GetUnderline(), want.GetForeground())
	}
}
