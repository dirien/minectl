package manifest

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/minectl/pkg/common"

	"github.com/minectl/pkg/model"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

type MinecraftServerManifest struct {
	MinecraftServer *model.MinecraftResource
}

const (
	MinecraftProxy  = "MinecraftProxy"
	MinecraftServer = "MinecraftServer"
)

//go:embed server.json
var server string

//go:embed proxy.json
var proxy string

func validate(manifest []byte) error {
	var schemaLoader gojsonschema.JSONLoader
	if strings.Contains(string(manifest), MinecraftProxy) {
		schemaLoader = gojsonschema.NewStringLoader(proxy)
	} else if strings.Contains(string(manifest), MinecraftServer) {
		schemaLoader = gojsonschema.NewStringLoader(server)
	}
	yaml, err := yaml.YAMLToJSON(manifest)
	if err != nil {
		return err
	}
	documentLoader := gojsonschema.NewStringLoader(string(yaml))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
			return errors.New("validation error")
		}
	}
	return nil
}

func checkNamePattern(serverName string) error {
	match, _ := regexp.MatchString(common.NameRegex, serverName)
	if !match {
		return errors.New("the name of your Minecraft server must consist of lower case alphanumeric characters or '-'")
	}
	return nil
}

func NewMinecraftResource(manifestPath string) (*model.MinecraftResource, error) {
	var server model.MinecraftResource
	manifestFile, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	err = validate(manifestFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(manifestFile, &server)
	if err != nil {
		return nil, err
	}
	err = checkNamePattern(server.GetName())
	if err != nil {
		return nil, err
	}
	return &server, nil
}
