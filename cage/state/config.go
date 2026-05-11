package state

import (
	"cage/agent/secrets"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type CageConfig struct {
	Name     string          `yaml:"name"`
	Packages []string        `yaml:"packages"`
	Agent    CageAgentConfig `yaml:"agent"`
	Env      []EnvVar        `yaml:"env"`
}

type CageAgentConfig struct {
	AuthTokenSecretLabel string `yaml:"authTokenSecretLabel"`
	BaseURL              string `yaml:"baseURL"`
	Model                string `yaml:"model"`
	HaikuModel           string `yaml:"haikuModel"`
}

func (c CageConfig) ImageName() string {
	return "localhost/cage:" + strings.ToLower(c.Name)
}

func (c CageConfig) ClaudeCodeEnv(redactSecret ...bool) (env []string, err error) {
	if c.Agent.AuthTokenSecretLabel != "" {
		var s = fmt.Sprintf("<%v=%v>", viper.GetString("secrets.backend"), c.Agent.AuthTokenSecretLabel)
		if len(redactSecret) == 0 || !redactSecret[0] {
			s, err = secrets.Retrieve(c.Agent.AuthTokenSecretLabel)
			if err != nil {
				return nil, err
			}
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

	for _, envVar := range c.Env {
		env = append(env, fmt.Sprintf("%v=%v", envVar.Key, envVar.Value))
	}

	return env, nil
}
