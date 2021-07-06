package template

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed civo.sh.tmpl
var bash string

//go:embed cloud-config.yaml.tmpl
var cloudConfig string

type Template struct {
	Template *template.Template
	Values   templateValues
}

type templateValues struct {
	Properties string
	Edition    string
	Mount      string
}

type Templater interface {
	GetTemplate() string
}

func (t Template) GetTemplate() (string, error) {
	var buff bytes.Buffer
	err := t.Template.Execute(&buff, t.Values)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func NewTemplateCivo(properties, edition string) (*Template, error) {
	civo, err := template.New("civo").Parse(bash)
	if err != nil {
		return nil, err
	}
	return &Template{
		Template: civo,
		Values: templateValues{
			Properties: strings.Replace(properties, "\t", "", -1),
			Edition:    edition,
		},
	}, nil
}

func NewTemplateCloudConfig(properties, edition, mount string) (*Template, error) {
	cloudInit, err := template.New("cloud-init").Parse(cloudConfig)
	if err != nil {
		return nil, err
	}
	return &Template{
		Template: cloudInit,
		Values: templateValues{
			Properties: strings.Replace(properties, "\t", "      ", -1),
			Edition:    edition,
			Mount:      mount,
		},
	}, nil
}
