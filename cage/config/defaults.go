package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	AppName = "cage"
)

func initPlatformShared() {
	viper.SetDefault("secrets.backend", "keyring")
}

func CageDirName() string {
	return fmt.Sprintf(".%v", AppName)
}
