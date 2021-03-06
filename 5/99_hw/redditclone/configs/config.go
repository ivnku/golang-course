package configs

import "github.com/spf13/viper"

type Config struct {
	Token string `mapstructure:"TOKEN"`
}

var Conf Config

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	Conf = config

	return
}
