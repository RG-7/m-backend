package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	Port       string
	SECRET_KEY string
}

// load Config
func LoadConfig() Config {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Skipping .env loading, using Render-provided env vars")
	}

	return Config{
		MongoURI:   os.Getenv("MONGODB_URI"),
		Port:       os.Getenv("PORT"),
		SECRET_KEY: os.Getenv("SECRET_KEY"),
	}
}
