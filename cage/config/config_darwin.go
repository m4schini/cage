package config

import "github.com/spf13/viper"

func InitPlatform() {
	initPlatformShared()
	viper.SetDefault("runtime", "colima")
}
