package template

import (
	"bytes"
	"embed"
	_ "embed"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/minectl/pkg/model"
)

type Template struct {
	Template *template.Template
	Values   *templateValues
}

type templateValues struct {
	*model.MinecraftResource
	Mount      string
	Properties []string
}

type TemplateName string

const (
	TemplateBash               TemplateName = "bash"
	TemplateCloudConfig        TemplateName = "cloud-config"
	TemplateJavaBinary         TemplateName = "java-binary"
	TemplateBedrockBinary      TemplateName = "bedrock-binary"
	TemplateSpigotBukkitBinary TemplateName = "spigotbukkit-binary"
	TemplateFabricBinary       TemplateName = "fabric-binary"
	TemplateForgeBinary        TemplateName = "forge-binary"
	TemplatePaperMCBinary      TemplateName = "papermc-binary"
	TemplateProxyCloudConfig   TemplateName = "proxy-cloud-config"
	TemplateProxyBash          TemplateName = "proxy-bash"
	TemplateBungeeCordBinary   TemplateName = "bungeecord-binary"
	TemplateWaterfallBinary    TemplateName = "waterfall-binary"
	TemplateNukkitBinary       TemplateName = "nukkit-binary"
)

func GetUpdateTemplate() *Template {
	bash := template.Must(template.New("base").Funcs(sprig.TxtFuncMap()).ParseFS(templateBash, "templates/bash/*"))
	return &Template{
		Template: bash,
		Values:   &templateValues{},
	}
}

func (t *Template) DoUpdate(model *model.MinecraftResource, name TemplateName) (string, error) {
	return t.GetTemplate(model, "", name)
}

func (t *Template) GetTemplate(model *model.MinecraftResource, mount string, name TemplateName) (string, error) {
	var buff bytes.Buffer

	t.Values.MinecraftResource = model
	t.Values.Mount = mount
	t.Values.Properties = strings.Split(model.GetProperties(), "\n")

	err := t.Template.ExecuteTemplate(&buff, string(name), t.Values)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func GetTemplateCloudConfigName(isProxy bool) TemplateName {
	if isProxy {
		return TemplateProxyCloudConfig
	}
	return TemplateCloudConfig
}

func GetTemplateBashName(isProxy bool) TemplateName {
	if isProxy {
		return TemplateProxyBash
	}
	return TemplateBash
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
