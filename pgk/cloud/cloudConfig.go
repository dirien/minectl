package cloud

import (
	_ "embed"
	"strings"
)

//go:embed cloud-config.yaml
var CloudConfig string

func ReplaceServerProperties(source, value string) string {
	value = strings.Replace(value, "\t", "      ", -1)
	return strings.Replace(source, "<properties>", value, -1)
}
