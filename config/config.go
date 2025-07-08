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

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Skipping .env loading, assuming Render-provided env vars")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGODB_URI is not set!")
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("❌ SECRET_KEY is not set!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // fallback for local
	}

	return Config{
		MongoURI:   mongoURI,
		Port:       port,
		SECRET_KEY: secret,
	}
}
