package globalconfig

import (
	"github.com/joho/godotenv"
	"os"
	"log"
)

// use godot package to load/read the .env file and
func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}