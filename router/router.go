package router

import (
	"Task2/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/ems-api/v1").Subrouter()

	subRouter.HandleFunc("/employees", controllers.GetAllEmployees).Methods("GET")
	subRouter.HandleFunc("/employees/search", controllers.SearchEmployees).Methods("GET")
	subRouter.HandleFunc("/employees/{id}", controllers.GetEmployee).Methods("GET")
	subRouter.HandleFunc("/employees", controllers.CreateEmployee).Methods("POST")
	subRouter.HandleFunc("/employees/{id}", controllers.UpdateEmployee).Methods("PUT")
	subRouter.HandleFunc("/employees/{id}", controllers.DeleteEmployee).Methods("DELETE")

	return subRouter
}
