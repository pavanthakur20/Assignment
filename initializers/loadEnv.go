package initializers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file, using system environment variables")
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
