package config

import "os"

type Config struct {
	APIURL string
	Port   string
}

func NewConfig() *Config {
	return &Config{
		APIURL: getEnv("API_URL", "https://api.chargeblast.io/bin/"),
		Port:   getEnv("PORT", "9000"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
