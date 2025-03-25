package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var Env string
var DatabasePath string
var ApiUrl string
var Debug bool

func LoadConfig() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // Default is "dev"
	}

	// If not exist this code, godotenv tries to find from go project root dir
	envPath := filepath.Join("config", ".env."+env)

	err := godotenv.Load(envPath)
	if err != nil {
		log.Println("Warning: No .env file found for environment:", env)
	}

	Env = os.Getenv("ENV")
	DatabasePath = os.Getenv("DATABASE_PATH")
	ApiUrl = os.Getenv("API_URL")
	Debug = os.Getenv("DEBUG") == "true"
}
