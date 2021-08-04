package template

import (
	"bytes"
	"embed"
	_ "embed"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
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

func GetUpdateTemplate() *Template {
	bash := template.Must(template.ParseFS(templateBash, "templates/bash/*"))
	return &Template{
		Template: bash,
		Values:   &templateValues{},
	}
}

func (t *Template) DoUpdate(model *model.MinecraftServer, name TemplateName) (string, error) {
	return t.GetTemplate(model, "", name)
}

func (t *Template) GetTemplate(model *model.MinecraftServer, mount string, name TemplateName) (string, error) {
	var buff bytes.Buffer

	t.Values.MinecraftServer = model
	t.Values.Mount = mount
	t.Values.Properties = strings.Split(model.GetProperties(), "\n")

	err := t.Template.ExecuteTemplate(&buff, string(name), t.Values)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

//go:embed templates
var templateBash embed.FS

func NewTemplateBash() (*Template, error) {
	bash := template.Must(template.New("base").Funcs(sprig.TxtFuncMap()).ParseFS(templateBash, "templates/bash/*"))
	return &Template{
		Template: bash,
		Values:   &templateValues{},
	}, nil
}

//go:embed templates
var templateCloudConfig embed.FS

func NewTemplateCloudConfig() (*Template, error) {
	cloudInit := template.Must(template.New("base").Funcs(sprig.TxtFuncMap()).ParseFS(templateCloudConfig, "templates/cloud-init/*"))
	return &Template{
		Template: cloudInit,
		Values:   &templateValues{},
	}, nil
}
