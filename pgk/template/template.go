package template

import (
	"bytes"
	_ "embed"
	"github.com/minectl/pgk/model"
	"strings"
	"text/template"
)

//go:embed civo.sh.tmpl
var bash string

//go:embed cloud-config.yaml.tmpl
var cloudConfig string

type Template struct {
	Template *template.Template
	Values   *templateValues
}

type templateValues struct {
	*model.MinecraftServer
	Mount      string
	Properties []string
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

func NewTemplateCivo(model *model.MinecraftServer) (*Template, error) {
	civo, err := template.New("civo").Parse(bash)
	if err != nil {
		return nil, err
	}
	return &Template{
		Template: civo,
		Values: &templateValues{
			MinecraftServer: model,
		},
	}, nil
}

func NewTemplateCloudConfig(model *model.MinecraftServer, mount string) (*Template, error) {
	cloudInit, err := template.New("cloud-init").Parse(cloudConfig)
	if err != nil {
		return nil, err
	}
	return &Template{
		Template: cloudInit,
		Values: &templateValues{
			MinecraftServer: model,
			Properties:      strings.Split(model.GetProperties(), "\n"),
			Mount:           mount,
		},
	}, nil
}