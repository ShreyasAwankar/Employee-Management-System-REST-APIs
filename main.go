package main

import (
	"Task2/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Employee Management System APIs")
	r := router.Router()
	fmt.Println("Listening to port 4000...")
	log.Fatal(http.ListenAndServe(":4000", r))
}
