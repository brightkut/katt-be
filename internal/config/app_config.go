package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(env string) {

	// Load only not run on aws lambda
	if env == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}
}
