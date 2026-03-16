package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Env struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	LogLevel    string `mapstructure:"SERVER_LOG_LEVEL"`
	Environment string `mapstructure:"ENV"`
	DBUsername  string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASS"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBName      string `mapstructure:"DB_NAME"`
}

func NewEnv() Env {
	env := Env{}

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		wd, _ := os.Getwd()
		possiblePaths := []string{
			filepath.Join(wd, "config", "config.yaml"),
			filepath.Join(wd, "..", "config", "config.yaml"),
			filepath.Join(wd, "..", "..", "config", "config.yaml"),
		}
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			viper.SetConfigType("yaml")
			viper.SetConfigFile(configPath)
			viper.ReadInConfig()
		}
	}

	err := viper.Unmarshal(&env)
	if err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	env.Environment = getEnvOrDefault("NODE_ENV", getEnvOrDefault("ENV",
		viper.GetString("ENV")))
	env.LogLevel = getEnvOrDefault("LOG_LEVEL",
		getConfigValue("server.log.level", viper.GetString("LOG_LEVEL"), ""))
	env.ServerPort = getEnvOrDefault("SERVER_PORT", env.ServerPort)
	env.DBUsername = getEnvOrDefault("DB_USER", env.DBUsername)
	env.DBPassword = getEnvOrDefault("DB_PASS", env.DBPassword)
	env.DBHost = getEnvOrDefault("DB_HOST", env.DBHost)
	env.DBPort = getEnvOrDefault("DB_PORT", env.DBPort)
	env.DBName = getEnvOrDefault("DB_NAME", env.DBName)

	return env
}

func getConfigValue(key, envValue, defaultValue string) string {
	if envValue != "" {
		return envValue
	}
	if viper.IsSet(key) {
		val := viper.Get(key)
		if val != nil {
			switch v := val.(type) {
			case string:
				return v
			case int, int32, int64:
				return fmt.Sprintf("%d", v)
			case float64:
				return fmt.Sprintf("%.0f", v)
			default:
				return viper.GetString(key)
			}
		}
	}
	return defaultValue
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
