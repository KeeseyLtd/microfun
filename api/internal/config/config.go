package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env  string `mapstructure:"APP_ENV"`
	Port int    `mapstructure:"APP_PORT"`
}

func LoadConfig(config string) *Config {
	viper.SetConfigFile(config)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("could not load config", err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("could not unmarshal config", err)
	}

	return &cfg
}
