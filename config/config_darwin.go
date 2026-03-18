package config

import "github.com/spf13/viper"

func InitPlatform() {
	viper.SetDefault("runtime", "podman")
}
