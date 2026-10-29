// Package ui provides terminal UI components for minectl using the Charm ecosystem.
package ui

import (
	"errors"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

// FormTheme returns a consistent theme for Huh forms.
func FormTheme() huh.Theme {
	return huh.ThemeFunc(formStyles)
}

// formStyles builds the form styles on top of the base theme. The overrides
// use fixed colors, so they look the same on light and dark terminals.
func formStyles(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)

	t.Focused.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Bold(true)
	t.Focused.Description = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))
	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229"))
	t.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	return t
}

// RunForm runs a Huh form with consistent styling.
// In headless mode, it will use accessible mode for better CI/CD compatibility.
func RunForm(form *huh.Form, headless bool) error {
	if headless {
		// Accessible mode prints styled prompts straight to its writer.
		// lipgloss.Writer downsamples colors to what stdout supports, so
		// piped CI output stays free of escape codes.
		form = form.WithAccessible(true).WithOutput(lipgloss.Writer)
	}
	return form.WithTheme(FormTheme()).Run()
}

// Confirm displays a yes/no confirmation prompt.
// Returns true if the user confirmed, false otherwise.
func Confirm(title string) (bool, error) {
	confirmed := true

	confirm := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Inline(true).
		Value(&confirmed)

	err := confirm.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, nil
		}
		return false, err
	}
	return confirmed, nil
}
