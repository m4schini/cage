package config

func InitPlatform() {
	viper.SetDefault("runtime", "docker")
}
