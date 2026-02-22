package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dirien/minectl/internal/logging"
	"go.uber.org/zap"
)

var (
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

// UI provides a unified interface for terminal output that handles both
// interactive and headless modes. In interactive mode, it uses Charm components
// for polished output. In headless mode, it falls back to structured logging.
type UI struct {
	headless bool
	logging  *logging.MinectlLogging
}

// NewUI creates a new UI instance.
func NewUI(headless bool, log *logging.MinectlLogging) *UI {
	return &UI{
		headless: headless,
		logging:  log,
	}
}

// IsHeadless returns true if the UI is running in headless mode (CI/CD).
func (u *UI) IsHeadless() bool {
	return u.headless
}

// Logging returns the underlying logging instance.
func (u *UI) Logging() *logging.MinectlLogging {
	return u.logging
}

// Info prints an informational message with a dim/neutral style.
func (u *UI) Info(msg string) {
	if u.headless {
		zap.S().Infow(msg)
		return
	}
	fmt.Println(infoStyle.Render("ℹ " + msg))
}

// Success prints a success message with a green checkmark prefix.
func (u *UI) Success(msg string) {
	if u.headless {
		zap.S().Infow(msg)
		return
	}
	fmt.Println(successStyle.Render("✓ " + msg))
}

// ErrorMsg prints an error message with a red cross prefix.
func (u *UI) ErrorMsg(err error) {
	if u.headless {
		zap.S().Errorw(err.Error())
		return
	}
	fmt.Println(errorStyle.Render("✗ " + err.Error()))
}

// Warn prints a warning message with a yellow warning prefix.
func (u *UI) Warn(msg string) {
	if u.headless {
		zap.S().Warnw(msg)
		return
	}
	fmt.Println(warnStyle.Render("⚠ " + msg))
}
