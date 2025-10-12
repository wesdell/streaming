package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(variableName string) string {
	env := godotenv.Load(".env")
	if env != nil {
		log.Println("Warning: Unable to find .env file!")
	}
	return os.Getenv(variableName)
}
