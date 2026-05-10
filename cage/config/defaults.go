package config

import "github.com/spf13/viper"

var (
	AppName = "cage"
)

func initPlatformShared() {
	viper.SetDefault("secrets.backend", "keyring")
}
