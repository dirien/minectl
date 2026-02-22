package ui

import (
	"errors"
	"testing"

	"github.com/dirien/minectl/internal/logging"
)

func newTestUI(t *testing.T, headless bool) *UI {
	t.Helper()
	l, err := logging.NewLogging("info", "console", headless)
	if err != nil {
		t.Fatalf("failed to create logging: %v", err)
	}
	return NewUI(headless, l)
}

func TestNewSpinner(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		headless bool
	}{
		{"headless spinner", "creating server...", true},
		{"interactive spinner", "creating server...", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			s := NewSpinner(tt.message, u)

			if s.Message != tt.message {
				t.Errorf("got Message %q, want %q", s.Message, tt.message)
			}
			if s.headless != tt.headless {
				t.Errorf("got headless=%v, want %v", s.headless, tt.headless)
			}
			if s.started {
				t.Error("expected started to be false for new spinner")
			}
		})
	}
}

func TestSpinnerHeadlessStartStop(t *testing.T) {
	tests := []struct {
		name         string
		stopErr      error
		wantStarted  bool
		wantHeadless bool
	}{
		{"stop without error", nil, false, true},
		{"stop with error", errors.New("server failed"), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, true)
			s := NewSpinner("creating server...", u)
			s.FinalMessage = "server created"
			s.ErrorMessage = "server failed"

			s.Start()
			if !s.headless {
				t.Error("expected headless mode")
			}
			// In headless mode, started remains false (no goroutine launched)
			if s.started != tt.wantStarted {
				t.Errorf("got started=%v, want %v", s.started, tt.wantStarted)
			}
			s.Stop(tt.stopErr)
		})
	}
}

func TestSpinnerStopWithoutStart(t *testing.T) {
	tests := []struct {
		name     string
		headless bool
	}{
		{"headless mode", true},
		{"interactive mode", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			s := NewSpinner("test", u)
			s.FinalMessage = "done"

			// Stop without Start should not panic
			s.Stop(nil)

			if s.started {
				t.Error("expected started to be false")
			}
		})
	}
}
