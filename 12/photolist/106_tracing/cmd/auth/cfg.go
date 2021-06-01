package main

import (
	"strings"

	"github.com/spf13/viper"
)

func Read(appName string, defaults map[string]interface{}, cfg interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(appName)
	v.AddConfigPath("/etc/")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if cfg != nil {
		err := v.Unmarshal(cfg)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}
