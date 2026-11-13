package session

import (
	"testing"
	"time"

	"github.com/brizzai/fleet/internal/debuglog"
)

// The agent's last frame (an idle ❯ box) stays on screen above the shell prompt
// after it exits, so the pane scraper alone reads a dead session as finished.
// With the shell in front, the row must read as error whatever the hook last said.
func TestAgentExitedToShellIsError(t *testing.T) {
	debuglog.Init()
	lastFrame := "Done!\n\n❯ \n\nroy@mac proj % "
	for _, tc := range []struct {
		name string
		old  Status
		hook string
	}{
		{"no hook, was finished", StatusFinished, ""},
		{"no hook, was idle", StatusIdle, ""},
		{"stale running hook", StatusRunning, "running"},
		{"stale finished hook", StatusFinished, "finished"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &Session{ID: "exited", Status: tc.old, paneCapturer: &mockPane{content: lastFrame, shellFg: true}}
			if tc.hook != "" {
				s.UpdateHookStatus(&HookStatus{Status: tc.hook, UpdatedAt: time.Now().Add(-time.Hour)}, false)
			}
			s.UpdateStatus()
			if got := s.GetStatus(); got != StatusError {
				t.Fatalf("status = %s, want error", got)
			}
		})
	}
}

// A launch puts the shell in front for a moment before the agent starts; that
// must not read as an exit.
func TestStartingWithShellInFrontStaysStarting(t *testing.T) {
	debuglog.Init()
	s := &Session{ID: "launching", Status: StatusStarting, paneCapturer: &mockPane{content: "roy@mac proj % ", shellFg: true}}
	s.UpdateStatus()
	if got := s.GetStatus(); got == StatusError {
		t.Fatalf("status = error during launch")
	}
}
