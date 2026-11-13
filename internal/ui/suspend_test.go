package ui

import (
	"testing"
	"time"

	"github.com/brizzai/fleet/internal/config"
	"github.com/brizzai/fleet/internal/perfwatch"
	"github.com/brizzai/fleet/internal/session"
)

// suspendIdleThreshold encodes the aggressiveness-mode table. Lock it down so a
// tweak to one mode can't silently change another.
func TestSuspendIdleThreshold(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		level   int
		wantAct bool
		wantMin time.Duration
	}{
		{"off never acts (even critical)", config.SuspendOff, perfwatch.PressureCritical, false, 0},
		{"light ignores normal", config.SuspendLight, perfwatch.PressureNormal, false, 0},
		{"light ignores warning", config.SuspendLight, perfwatch.PressureWarning, false, 0},
		{"light acts on critical", config.SuspendLight, perfwatch.PressureCritical, true, 24 * time.Hour},
		{"balanced ignores normal", config.SuspendBalanced, perfwatch.PressureNormal, false, 0},
		{"balanced acts at warning", config.SuspendBalanced, perfwatch.PressureWarning, true, 4 * time.Hour},
		{"balanced acts at critical", config.SuspendBalanced, perfwatch.PressureCritical, true, 4 * time.Hour},
		{"aggressive ignores normal", config.SuspendAggressive, perfwatch.PressureNormal, false, 0},
		{"aggressive acts at warning", config.SuspendAggressive, perfwatch.PressureWarning, true, time.Hour},
		{"aggressive acts at critical", config.SuspendAggressive, perfwatch.PressureCritical, true, time.Hour},
		{"unknown pressure: light stays off", config.SuspendLight, perfwatch.PressureUnknown, false, 0},
		{"unknown pressure: aggressive stays off", config.SuspendAggressive, perfwatch.PressureUnknown, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotAct := suspendIdleThreshold(tt.mode, tt.level)
			if gotAct != tt.wantAct {
				t.Fatalf("act = %v, want %v", gotAct, tt.wantAct)
			}
			if tt.wantAct && gotMin != tt.wantMin {
				t.Fatalf("minIdle = %v, want %v", gotMin, tt.wantMin)
			}
		})
	}
}

// sweepDue: the inactivity limit suspends any at-rest row with no pressure; the
// pressure threshold keeps its old idle-only scope; 0 turns the limit off.
func TestSweepDue(t *testing.T) {
	const day = 24 * time.Hour
	tests := []struct {
		name          string
		st            session.Status
		idle          time.Duration
		pressureMin   time.Duration
		pressureAct   bool
		inactiveAfter time.Duration
		want          bool
	}{
		{"idle past limit, no pressure", session.StatusIdle, 49 * time.Hour, 0, false, 2 * day, true},
		{"finished past limit, no pressure", session.StatusFinished, 49 * time.Hour, 0, false, 2 * day, true},
		{"error past limit, no pressure", session.StatusError, 49 * time.Hour, 0, false, 2 * day, true},
		{"error under limit, no pressure", session.StatusError, 47 * time.Hour, 0, false, 2 * day, false},
		{"limit off, no pressure", session.StatusIdle, 100 * day, 0, false, 0, false},
		{"pressure suspends idle sooner", session.StatusIdle, 2 * time.Hour, time.Hour, true, 2 * day, true},
		{"pressure leaves finished alone", session.StatusFinished, 2 * time.Hour, time.Hour, true, 2 * day, false},
		{"pressure leaves error alone", session.StatusError, 2 * time.Hour, time.Hour, true, 0, false},
		{"pressure, idle too fresh", session.StatusIdle, 30 * time.Minute, time.Hour, true, 2 * day, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sweepDue(tt.st, tt.idle, tt.pressureMin, tt.pressureAct, tt.inactiveAfter); got != tt.want {
				t.Fatalf("sweepDue = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtRest(t *testing.T) {
	for st, want := range map[session.Status]bool{
		session.StatusIdle: true, session.StatusFinished: true, session.StatusError: true,
		session.StatusRunning: false, session.StatusWaiting: false, session.StatusStarting: false,
		session.StatusSuspended: false,
	} {
		if got := atRest(st); got != want {
			t.Errorf("atRest(%s) = %v, want %v", st, got, want)
		}
	}
}
