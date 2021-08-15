package common

import (
	"github.com/fatih/color"
)

const InstanceTag = "minectl"

func Green(value string) string {
	return color.GreenString(value)
}
