package ui

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
)

// Spinner provides an animated spinner for long-running operations.
// In headless mode, it falls back to simple log messages.
type Spinner struct {
	Message      string
	FinalMessage string
	ErrorMessage string
	headless     bool
	program      *tea.Program
	mu           sync.Mutex
	started      bool
	done         chan struct{}
}

// spinnerModel is the Bubble Tea model for the spinner.
type spinnerModel struct {
	spinner  spinner.Model
	message  string
	quitting bool
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case quitMsg:
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// quitMsg is sent to stop the spinner.
type quitMsg struct{}

// NewSpinner creates a new animated spinner.
func NewSpinner(message string, u *UI) *Spinner {
	return &Spinner{
		Message:  message,
		headless: u.IsHeadless(),
		done:     make(chan struct{}),
	}
}

// Start begins the spinner animation.
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.headless {
		zap.S().Infow(s.Message)
		return
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	model := spinnerModel{
		spinner: sp,
		message: s.Message,
	}

	s.program = tea.NewProgram(model, tea.WithOutput(os.Stderr))
	s.started = true

	go func() {
		if _, err := s.program.Run(); err != nil {
			zap.S().Debugw("spinner program error", "error", err)
		}
		close(s.done)
	}()
}

// Stop stops the spinner and displays the appropriate final message.
func (s *Spinner) Stop(err error) {
	s.mu.Lock()
	var message string
	if err != nil {
		message = s.ErrorMessage
	} else {
		message = s.FinalMessage
	}

	if s.headless {
		s.mu.Unlock()
		if message != "" {
			if err != nil {
				zap.S().Errorw(message)
			} else {
				zap.S().Infow(message)
			}
		}
		return
	}

	prog := s.program
	started := s.started
	s.mu.Unlock()

	// Send quit and wait outside of the lock to avoid deadlock
	if started && prog != nil {
		prog.Send(quitMsg{})
		<-s.done
	}

	if message != "" {
		if err != nil {
			fmt.Fprintln(os.Stderr, errorStyle.Render("✗ "+message))
		} else {
			fmt.Fprintln(os.Stderr, successStyle.Render("✓ "+message))
		}
	}
}
