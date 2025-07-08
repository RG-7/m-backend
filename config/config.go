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

	mongoURI := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET_KEY")

	// Fail fast if critical env vars are missing
	if mongoURI == "" {
		log.Fatal("❌ MONGODB_URI is not set!")
	}
	if secret == "" {
		log.Fatal("❌ SECRET_KEY is not set!")
	}

	return Config{
		MongoURI:   mongoURI,
		Port:       port,
		SECRET_KEY: secret,
	}
}
