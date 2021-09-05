package utils

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func HumanDate(time time.Time) string {
	if time.IsZero() {
		return ""
	}

	return time.UTC().Format("02 Jan 2006 at 15:04")
}

func GetEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Couldn't load .env file.")
	}

	return os.Getenv(key)
}
