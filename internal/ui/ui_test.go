package ui

import (
	"errors"
	"testing"
)

func TestUI(t *testing.T) {
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

			if got := u.IsHeadless(); got != tt.headless {
				t.Errorf("IsHeadless() = %v, want %v", got, tt.headless)
			}
			if u.Logging() == nil {
				t.Error("Logging() returned nil")
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name     string
		headless bool
	}{
		{"headless info", true},
		{"interactive info", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			// Should not panic
			u.Info("test message")
		})
	}
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		name     string
		headless bool
	}{
		{"headless success", true},
		{"interactive success", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			// Should not panic
			u.Success("server created")
		})
	}
}

func TestErrorMsg(t *testing.T) {
	tests := []struct {
		name     string
		headless bool
	}{
		{"headless error", true},
		{"interactive error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			// Should not panic
			u.ErrorMsg(errors.New("something went wrong"))
		})
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		name     string
		headless bool
	}{
		{"headless warn", true},
		{"interactive warn", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUI(t, tt.headless)
			// Should not panic
			u.Warn("this is a warning")
		})
	}
}
