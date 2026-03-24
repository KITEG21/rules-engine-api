package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	App      AppConfig      `mapstructure:"app"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type AppConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("$HOME/.config/rules_engine_api")
	viper.AddConfigPath("/etc/rules_engine_api")

	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("database.port", 5432)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("WARNING: config file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return &cfg
}
