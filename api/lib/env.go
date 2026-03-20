package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type OAuthProvider struct {
	Name         string   `mapstructure:"name"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	Scopes       []string `mapstructure:"scopes"`
}

type Env struct {
	ServerPort     string `mapstructure:"SERVER_PORT"`
	LogLevel       string `mapstructure:"SERVER_LOG_LEVEL"`
	Environment    string `mapstructure:"ENV"`
	DBUsername     string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASS"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBName         string `mapstructure:"DB_NAME"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	AuthHeader     string `mapstructure:"AUTH_HEADER"`
	OAuthProviders []OAuthProvider
}

func (e Env) GetOAuthProvider(name string) (OAuthProvider, bool) {
	for _, p := range e.OAuthProviders {
		if p.Name == name {
			return p, true
		}
	}
	return OAuthProvider{}, false
}

func NewEnv() Env {
	godotenv.Load()

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

	env.Environment = firstNonEmpty(
		getConfigValue("server.environment", "", ""),
		os.Getenv("ENV"),
		os.Getenv("NODE_ENV"),
		"development",
	)
	env.LogLevel = firstNonEmpty(
		getConfigValue("server.log.level", "", ""),
		os.Getenv("SERVER_LOG_LEVEL"),
		os.Getenv("LOG_LEVEL"),
		"info",
	)
	env.ServerPort = firstNonEmpty(
		getConfigValue("server.port", "", ""),
		os.Getenv("SERVER_PORT"),
		"8080",
	)
	env.DBHost = firstNonEmpty(
		getConfigValue("database.host", "", ""),
		os.Getenv("DB_HOST"),
		"localhost",
	)
	env.DBPort = firstNonEmpty(
		getConfigValue("database.port", "", ""),
		os.Getenv("DB_PORT"),
		"5432",
	)
	env.DBName = firstNonEmpty(
		getConfigValue("database.database", "", ""),
		os.Getenv("DB_NAME"),
	)
	env.DBUsername = firstNonEmpty(
		getConfigValue("database.username", "", ""),
		os.Getenv("DB_USER"),
	)
	env.DBPassword = firstNonEmpty(
		getConfigValue("database.password", "", ""),
		os.Getenv("DB_PASS"),
	)
	env.JWTSecret = firstNonEmpty(
		getConfigValue("server.jwt_secret", "", ""),
		os.Getenv("JWT_SECRET"),
	)
	env.AuthHeader = firstNonEmpty(
		getConfigValue("server.auth_header", "", ""),
		os.Getenv("AUTH_HEADER"),
	)

	var providers []OAuthProvider
	if err := viper.UnmarshalKey("sso.providers", &providers); err == nil {
		env.OAuthProviders = providers
	}

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
