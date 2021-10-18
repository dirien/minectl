package common

import (
	"github.com/fatih/color"
)

const InstanceTag = "minectl"

const NameRegex = "^[a-z-0-9]+$"

func Green(value string) string {
	return color.GreenString(value)
}
