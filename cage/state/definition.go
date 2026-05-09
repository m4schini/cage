package state

import (
	"cage/agent/secrets"
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

type CageConfig struct {
	Name     string   `yaml:"name"`
	Packages []string `yaml:"packages"`
	Agent    struct {
		AuthTokenSecretLabel string `yaml:"authTokenSecretLabel"`
		BaseURL              string `yaml:"baseURL"`
		Model                string `yaml:"model"`
		HaikuModel           string `yaml:"haikuModel"`
	} `yaml:"agent"`
}

func (c CageConfig) ClaudeCodeEnv() ([]string, error) {
	var env []string
	if c.Agent.AuthTokenSecretLabel != "" {
		s, err := secrets.Retrieve(c.Agent.AuthTokenSecretLabel)
		if err != nil {
			return nil, err
		}
		env = append(env, fmt.Sprintf("ANTHROPIC_AUTH_TOKEN=%v", s))
	}
	if c.Agent.BaseURL != "" {
		env = append(env, fmt.Sprintf("ANTHROPIC_BASE_URL=%v", c.Agent.BaseURL))
	}
	if c.Agent.BaseURL != "" {
		env = append(env, fmt.Sprintf("ANTHROPIC_MODEL=%v", c.Agent.Model))
	}
	if c.Agent.BaseURL != "" {
		env = append(env, fmt.Sprintf("ANTHROPIC_DEFAULT_HAIKU_MODEL=%v", c.Agent.HaikuModel))
	}
	env = append(env, "DISABLE_TELEMETRY=1")

	return env, nil
}

func Write(def CageDefinition, writer io.Writer) error {
	return yaml.NewEncoder(writer).Encode(&def)
}

func Read(reader io.Reader) (CageDefinition, error) {
	var def CageDefinition
	err := yaml.NewDecoder(reader).Decode(&def)
	return def, err
}
