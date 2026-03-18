package config

import (
	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

func Init(cfgFile string) {
	InitPlatform()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cage" (without extension).
		viper.AddConfigPath(xdg.ConfigHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cage")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SafeWriteConfig()
}
