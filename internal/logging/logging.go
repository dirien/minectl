package logging

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MinectlLogging struct {
	headless bool
}

func NewLogging(verbose, logEncoding string, headless bool) (*MinectlLogging, error) {
	var level zapcore.Level
	err := level.Set(verbose)
	if err != nil {
		return nil, err
	}
	cfg := zap.Config{
		Encoding: logEncoding,
		Level:    zap.NewAtomicLevelAt(level),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	if verbose != "" {
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
	}

	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	return &MinectlLogging{
		headless: headless,
	}, nil
}

func (l *MinectlLogging) Error(msg error) {
	if l.headless {
		zap.S().Error(msg)
	} else {
		fmt.Println(msg)
	}
}

func (l *MinectlLogging) RawMessage(msg string) {
	if l.headless {
		zap.S().Infow(strings.ReplaceAll(msg, "\n", ""))
	} else {
		fmt.Println(msg)
	}
}

func (l *MinectlLogging) PrintMixedGreen(format, value string) {
	if l.headless {
		zap.S().Infow(strings.ReplaceAll(fmt.Sprintf(format, value), "\n", ""))
	} else {
		greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		fmt.Printf(format+"\n", greenStyle.Render(value))
	}
}

// IsHeadless returns true if the logging instance is running in headless mode.
func (l *MinectlLogging) IsHeadless() bool {
	return l.headless
}
