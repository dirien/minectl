package manifest

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/minectl/pkg/model"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

type MinecraftServerManifest struct {
	MinecraftServer *model.MinecraftServer
}

//go:embed schema.json
var schema string

func validate(manifest []byte) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	yaml, err := yaml.YAMLToJSON(manifest)
	if err != nil {
		log.Fatal(err)
	}
	documentLoader := gojsonschema.NewStringLoader(string(yaml))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatal(err)
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

func NewMinecraftServer(manifestPath string) (*MinecraftServerManifest, error) {
	var server MinecraftServerManifest
	manifestFile, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	err = validate(manifestFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(manifestFile, &server.MinecraftServer)
	if err != nil {
		return nil, err
	}
	return &server, nil
}
