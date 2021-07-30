package template

import (
	"bytes"
	"embed"
	_ "embed"
	"strings"
	"text/template"

	"github.com/minectl/pgk/model"
)

type Template struct {
	Template *template.Template
	Values   *templateValues
}

type templateValues struct {
	*model.MinecraftServer
	Mount      string
	Properties []string
}

type TemplateName string

const (
	TemplateBash               TemplateName = "bash"
	TemplateCloudConfig        TemplateName = "cloud-config"
	TemplateJavaBinary         TemplateName = "java-binary"
	TemplateBedrockBinary      TemplateName = "bedrock-binary"
	TemplatesSigotbukkitBinary TemplateName = "spigotbukkit-binary"
	TemplatesFabricBinary      TemplateName = "fabric-binary"
	TemplatesForgeBinary       TemplateName = "forge-binary"
	TemplatesPaperMCBinary     TemplateName = "papermc-binary"
)

func (t Template) GetUpdateTemplate(model *model.MinecraftServer, name TemplateName) (string, error) {
	return "", nil
}

func (t Template) GetTemplate(model *model.MinecraftServer, name TemplateName) (string, error) {
	var buff bytes.Buffer

	t.Values.MinecraftServer = model
	t.Values.Properties = strings.Split(model.GetProperties(), "\n")

	err := t.Template.ExecuteTemplate(&buff, string(name), t.Values)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

//go:embed templates
var templateBash embed.FS

func NewTemplateBash(mount string) (*Template, error) {
	bash := template.Must(template.ParseFS(templateBash, "templates/bash/*"))
	return &Template{
		Template: bash,
		Values: &templateValues{
			Mount: mount,
		},
	}, nil
}

//go:embed templates
var templateCloudConfig embed.FS

func NewTemplateCloudConfig(mount string) (*Template, error) {
	cloudInit := template.Must(template.ParseFS(templateCloudConfig, "templates/cloud-init/*"))
	return &Template{
		Template: cloudInit,
		Values: &templateValues{
			Mount: mount,
		},
	}, nil
}
