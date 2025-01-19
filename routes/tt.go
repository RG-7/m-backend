package routes

import (
	"github.com/RG-7/m-backend/controllers"
	"github.com/gorilla/mux"
)

func TTRoutes(router *mux.Router) {
	router.HandleFunc("/tt/linrange", controllers.CreateTimetableEntry).Methods("POST")
	router.HandleFunc("/tt/linrangedel", controllers.DeleteTimetableEntry).Methods("POST")
	router.HandleFunc("/tt/c", controllers.CreateCommonTimetableEntry).Methods("POST")
	router.HandleFunc("/tt/cdel", controllers.DeleteCommonTimetableEntry).Methods("POST")
	// router.HandleFunc("/tt/faculty/{facultyCode}/{date}", controllers.GetTimetableByFaculty).Methods("GET")
	// router.HandleFunc("/tt/subgroup/{subgroup}/{date}", controllers.GetTimetableBySubgroup).Methods("GET")
	// router.HandleFunc("/tt/room/{room}/{date}", controllers.GetTimetableByRoom).Methods("GET")
}
