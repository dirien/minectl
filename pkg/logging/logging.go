package logging

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

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

	if len(verbose) > 0 {
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
		zap.S().Infow(strings.Replace(msg, "\n", "", -1))
	} else {
		fmt.Println(msg)
	}
}

func (l *MinectlLogging) PrintMixedGreen(format string, value string) {
	if l.headless {
		zap.S().Infow(strings.Replace(fmt.Sprintf(format, value), "\n", "", -1))
	} else {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(fmt.Sprintf(format, green(value)))
	}
}
