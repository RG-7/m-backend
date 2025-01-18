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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!!!")
	}

	return Config{
		MongoURI:   os.Getenv("MONGODB_URI"),
		Port:       os.Getenv("PORT"),
		SECRET_KEY: os.Getenv("SECRET_KEY"),
	}
}
