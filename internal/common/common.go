package common

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const InstanceTag = "minectl"

const NameRegex = "^[a-z-0-9]+$"

func Green(value string) string {
	return color.GreenString(value)
}

func CreateServerNameWithTags(instanceName, label string) (id string) {
	return fmt.Sprintf("%s|%s", instanceName, label)
}

func ExtractFieldsFromServername(id string) (label string, err error) {
	fields := strings.Split(id, "|")
	if len(fields) == 3 {
		label = strings.Join([]string{fields[1], fields[2]}, ",")
	} else {
		err = fmt.Errorf("could not get fields from custom ID: fields: %v", fields)
		return "", err
	}
	return label, nil
}
