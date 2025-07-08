package main

import (
	"log"
	"net/http"

	"github.com/RG-7/m-backend/config"
	db "github.com/RG-7/m-backend/database"
	"github.com/RG-7/m-backend/helpers"
	"github.com/RG-7/m-backend/routes"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("🚀 Starting Backend Server...")
	// load config
	cfg := config.LoadConfig()
	helpers.SetSecret(cfg.SECRET_KEY)
	helpers.InitConfig()

	// connect to database
	db.ConnectDB(cfg.MongoURI)

	// initialize router
	router := mux.NewRouter()

	// register routes
	routes.RegisterRoutes(router)

	// get port
	port := cfg.Port
	if port == "" {
		port = "10000"
	}

	// Start server
	log.Println("✅ Server running on port:", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("❌ Server startup failed:", err)
	}
}
