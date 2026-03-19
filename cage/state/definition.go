package state

import (
	"cage/cage/config"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

var DefinitionFileName = fmt.Sprintf("%v.yaml", config.AppName)

type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type CageDefinition struct {
	Image    string
	Shell    string   `yaml:"shell"`
	Packages []string `yaml:"packages"`
	Env      []EnvVar `yaml:"env"`
}

func Write(def CageDefinition, writer io.Writer) error {
	return yaml.NewEncoder(writer).Encode(&def)
}

func Read(reader io.Reader) (CageDefinition, error) {
	var def CageDefinition
	err := yaml.NewDecoder(reader).Decode(&def)
	return def, err
}
