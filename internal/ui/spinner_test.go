package ui

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/dirien/minectl/internal/logging"
)

// testSpinnerMessage is the spinner message used across spinner tests.
const testSpinnerMessage = "creating server..."

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
		{"headless spinner", testSpinnerMessage, true},
		{"interactive spinner", testSpinnerMessage, false},
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
			s := NewSpinner(testSpinnerMessage, u)
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

func TestSpinnerModelUpdate(t *testing.T) {
	tests := []struct {
		name         string
		msg          tea.Msg
		wantQuitting bool
		wantQuit     bool
	}{
		{"ctrl+c quits", tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}, true, true},
		{"quit message quits", quitMsg{}, true, true},
		{"other key is ignored", tea.KeyPressMsg{Code: 'q', Text: "q"}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := spinnerModel{spinner: spinner.New(), message: testSpinnerMessage}

			updated, cmd := m.Update(tt.msg)

			got, ok := updated.(spinnerModel)
			if !ok {
				t.Fatalf("Update returned %T, want spinnerModel", updated)
			}
			if got.quitting != tt.wantQuitting {
				t.Errorf("got quitting=%v, want %v", got.quitting, tt.wantQuitting)
			}
			assertQuitCmd(t, cmd, tt.wantQuit)
		})
	}
}

// assertQuitCmd checks that cmd produces tea.QuitMsg when wantQuit is true,
// and that cmd is nil otherwise.
func assertQuitCmd(t *testing.T, cmd tea.Cmd, wantQuit bool) {
	t.Helper()
	if !wantQuit {
		if cmd != nil {
			t.Error("expected no command")
		}
		return
	}
	if cmd == nil {
		t.Fatal("expected a quit command, got nil")
	}
	if _, isQuit := cmd().(tea.QuitMsg); !isQuit {
		t.Error("expected command to produce tea.QuitMsg")
	}
}

func TestSpinnerModelView(t *testing.T) {
	tests := []struct {
		name     string
		quitting bool
		want     string
	}{
		{"running shows message", false, testSpinnerMessage},
		{"quitting clears output", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := spinnerModel{spinner: spinner.New(), message: testSpinnerMessage, quitting: tt.quitting}

			content := m.View().Content
			if tt.want == "" && content != "" {
				t.Errorf("expected empty view, got %q", content)
			}
			if !strings.Contains(content, tt.want) {
				t.Errorf("view %q does not contain %q", content, tt.want)
			}
		})
	}
}
