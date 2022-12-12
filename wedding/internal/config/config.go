package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env string `mapstructure:"APP_ENV"`

	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBDatabase string `mapstructure:"DB_DATABASE"`
	DBHostname string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBSSL      string `mapstructure:"DB_SSLMODE"`
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
