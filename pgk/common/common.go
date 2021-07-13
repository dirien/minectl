package common

import (
	"fmt"
	"github.com/fatih/color"
)

const InstanceTag = "minectl"

func PrintMixedGreen(format string, value string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf(format, green(value))
}

func Green(value string) string {
	return color.GreenString(value)
}
