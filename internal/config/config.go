package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	App      AppConfig      `mapstructure:"app"`
	AI       AIConfig       `mapstructure:"ai"`
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

type AIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
	Timeout int    `mapstructure:"timeout"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("$HOME/.config/rules_engine_api")
	viper.AddConfigPath("/etc/rules_engine_api")

	viper.SetEnvPrefix("") // optional
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("database.port", 5432)

	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.name", "DATABASE_NAME")
	viper.BindEnv("app.host", "APP_HOST")
	viper.BindEnv("app.port", "APP_PORT")
	viper.BindEnv("ai.api_key", "AI_API_KEY")
	viper.BindEnv("ai.base_url", "AI_BASE_URL")
	viper.BindEnv("ai.model", "AI_MODEL")
	viper.BindEnv("ai.timeout", "AI_TIMEOUT")

	// AI (LLM) defaults
	viper.SetDefault("ai.base_url", "")
	viper.SetDefault("ai.model", "")
	viper.SetDefault("ai.timeout", 30)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("WARNING: config file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return &cfg
}
