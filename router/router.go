package router

import (
	"Task2/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Router() http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/ems-api/v1").Subrouter()

	subRouter.HandleFunc("/employees", controllers.GetAllEmployees).Methods("GET")
	subRouter.HandleFunc("/employees/search", controllers.SearchEmployees).Methods("GET")
	subRouter.HandleFunc("/employees/{id}", controllers.GetEmployee).Methods("GET")
	subRouter.HandleFunc("/employees", controllers.CreateEmployee).Methods("POST")
	subRouter.HandleFunc("/employees/{id}", controllers.UpdateEmployee).Methods("PUT")
	subRouter.HandleFunc("/employees/{id}", controllers.DeleteEmployee).Methods("DELETE")

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://shreyasawankar.github.io"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})

	// Wrap your subRouter with the CORS middleware
	handler := c.Handler(router)

	return handler
}
