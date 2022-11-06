package config

import "github.com/spf13/viper"

func InitViper() error {
	viper.AddConfigPath("internal/config")
	viper.SetConfigName("configs")
	return viper.ReadInConfig()
}
