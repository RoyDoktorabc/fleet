package tmux

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// A multi-line cmd+v in split view or the drawer used to land only its first
// line: the control client takes one command per line, so the newline ended
// send-keys. Paste must deliver every line.
func TestPasteDeliversEveryLine(t *testing.T) {
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed")
	}
	// Not t.TempDir(): on macOS its path pushes the socket past the 104-byte limit.
	sock, err := os.MkdirTemp("", "ftp")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(sock) })
	t.Setenv("TMUX_TMPDIR", sock)
	t.Setenv("TMUX", "")
	os.Unsetenv("TMUX")

	s := NewSessionWithPrefix("fleettest_", "paste", t.TempDir())
	if out, err := exec.Command("tmux", "new-session", "-d", "-s", s.Name, "cat").CombinedOutput(); err != nil {
		t.Fatalf("new-session: %v: %s", err, out)
	}
	defer exec.Command("tmux", "kill-server").Run()

	if err := s.Paste("first line\nsecond line\nthird line\n"); err != nil {
		t.Fatalf("Paste: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	var pane string
	for time.Now().Before(deadline) {
		out, _ := exec.Command("tmux", "capture-pane", "-p", "-t", s.Name).Output()
		pane = string(out)
		if strings.Contains(pane, "third line") {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	for _, want := range []string{"first line", "second line", "third line"} {
		if !strings.Contains(pane, want) {
			t.Errorf("pane missing %q after paste; got:\n%s", want, pane)
		}
	}
}
