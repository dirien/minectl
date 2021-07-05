package manifest

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"sigs.k8s.io/yaml"
	"strings"
)

// Spec
type Spec struct {
	Server    Server    `yaml:"server"`
	Minecraft Minecraft `yaml:"minecraft"`
}

// Server
type Server struct {
	Size       string `yaml:"size"`
	VolumeSize int    `yaml:"volumeSize"`
	Ssh        string `yaml:"ssh"`
	Cloud      string `yaml:"cloud"`
	Region     string `yaml:"region"`
}

// Minecraft
type Minecraft struct {
	Java       Java   `yaml:"java"`
	Properties string `yaml:"properties"`
}

// Java
type Java struct {
	Xmx string `yaml:"xmx"`
	Xms string `yaml:"xms"`
}

// Metadata
type Metadata struct {
	Name string `yaml:"name"`
}

// MinecraftServer
type MinecraftServer struct {
	Spec       Spec     `yaml:"spec"`
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
}

type MinecraftServerManifester interface {
	UpdateManifest(filename string)
	SetID(id string)
	GetID() string
	GetName() string
	GetCloud() string
	GetSSH() string
	GetRegion() string
	GetSize() string
}

func (m *MinecraftServer) GetName() string {
	return m.Metadata.Name
}

func (m *MinecraftServer) GetCloud() string {
	return m.Spec.Server.Cloud
}

func (m *MinecraftServer) GetSSH() string {
	return m.Spec.Server.Ssh
}

func (m *MinecraftServer) GetRegion() string {
	return m.Spec.Server.Region
}

func (m *MinecraftServer) GetSize() string {
	return m.Spec.Server.Size
}

func (m *MinecraftServer) GetProperties() string {
	text := strings.Replace(m.Spec.Minecraft.Properties, "\n", "\n\t", -1)
	return text
}

func (m *MinecraftServer) GetVolumeSize() int {
	return m.Spec.Server.VolumeSize
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

func NewMinecraftServer(manifestPath string) (*MinecraftServer, error) {

	var server MinecraftServer
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
	return &server, nil
}
