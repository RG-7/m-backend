package routes

import (
	"github.com/RG-7/m-backend/controllers"
	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", controllers.Register).Methods("POST")
	router.HandleFunc("/auth/login", controllers.Login).Methods("POST")
	router.HandleFunc("/auth/validate", controllers.ValidateToken).Methods("GET")
	router.HandleFunc("/auth/{id}", controllers.DeleteUser).Methods("GET")
}
