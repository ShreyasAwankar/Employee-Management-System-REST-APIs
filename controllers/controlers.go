package controllers

import (
	"Task2/models"
	"Task2/validations"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var logger = CreateLogger()
var mu sync.Mutex

// Controllers
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {

	// Reading data from DB (ems.csv)
	employees, err := ReadCSV("get")

	if err != nil {
		http.Error(w, "Error rading data", http.StatusInternalServerError)
		logger.Error("\nInternal Server Error : occured while reading the data from DB during controllers.GetAllEmployees function call")
		return
	}

	// Serialize the results to JSON and send the response
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(employees)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)

	empId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid Employee ID", http.StatusBadRequest)
		logger.Error("\nInvalid id was employeeId was provided on GetEmployee function call")
		return
	}

	employees, err := ReadCSV("get")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Error("\nCould not read data from DB during controllers.GetEmployee function call")
		return
	}

	empFound := false

	var employee models.Employee

	// Finding employee from employees slice created with ReadCSV() to update it...
	for _, emp := range employees {
		if empId == emp.ID {
			empFound = true
			employee = emp
			break
		}
	}

	if !empFound {
		http.Error(w, "Employee not found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(employee)
}

// Generating employeeId for creating new employee record. Generated employeeId will be maximum existing employeeId in DB + 1
var EmpId = GenerateEmployeeId()

func CreateEmployee(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	var employee models.Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	employee.ID = EmpId

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logger.Error("\nInvalid JSON input for type Employee during controllers.CreateEmployee function call")
		w.Header().Set("Content-Type", "application/json")
		return
	}

	err1 := validations.V.Struct(employee)

	if err1 != nil {
		http.Error(w, "Invalid input for employee deatails", http.StatusUnprocessableEntity)
		logger.Errorf("\nInvalid employee data input : occured while validating employee fields during controllers.CreateEmployee function call")
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Writting data to DB ("ems.csv")
	emp, err := WriteToCSV(employee)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Errorf("\nCould not write data to DB during controllers.CreateEmployee function call")
		return
	}

	// Setting headers and staus codes
	newEmployeeURL := fmt.Sprintf("http://localhost:5000/employees/%d", EmpId)
	w.Header().Set("Location", newEmployeeURL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)

	empId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid Employee ID", http.StatusBadRequest)
		logger.Errorf("\nInvalid id was employeeId was provided during controllers.UpdateEmployee function call")
		return
	}

	employees, err := ReadCSV("update")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Errorf("\nInternal Server Error : occured while reading the data from DB during controllers.UpdateEmployee function call")
		return
	}

	var employeeToUpdate models.Employee
	err = json.NewDecoder(r.Body).Decode(&employeeToUpdate)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logger.Errorf("\nInvalid JSON input for type Employee during controllers.UpdateEmployee function call")
		return
	}

	err1 := validations.V.Struct(employeeToUpdate)

	if err1 != nil {
		http.Error(w, "Invalid input for employee deatails", http.StatusUnprocessableEntity)
		logger.Errorf("\nInvalid employee data input : occured while validating employee fields during controllers.UpdateEmployee function call")
		return
	}

	// Following function is just to beautify the response body...
	func() {
		employeeToUpdate.FirstName = strings.TrimSpace(employeeToUpdate.FirstName)
		employeeToUpdate.FirstName = strings.ToUpper(employeeToUpdate.FirstName[0:1]) + strings.ToLower(employeeToUpdate.FirstName[1:])

		employeeToUpdate.LastName = strings.TrimSpace(employeeToUpdate.LastName)
		employeeToUpdate.LastName = strings.ToUpper(employeeToUpdate.LastName[0:1]) + strings.ToLower(employeeToUpdate.LastName[1:])
	}()

	empFound := false

	// Creating a new slice of type models.Employee so as to write updated records to DB (ems.csv)
	empSliceToWrite := make([]models.Employee, len(employees))

	// Finding employee from employees slice created with ReadCSV() to update it...
	for i, emp := range employees {
		if empId == emp.ID {
			employeeToUpdate.ID = empId
			empFound = true
			copy(empSliceToWrite[:i], employees[:i])
			empSliceToWrite[i] = employeeToUpdate
			copy(empSliceToWrite[i+1:], employees[i+1:])
			break
		}
	}

	if !empFound {
		http.Error(w, "Employee not found!", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Truncating DB file to avoid duplicate entries
	_, err = os.Create("ems.csv")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Errorf("\nError occured while to truncating DB (ems.csv) during UpdateEmployee operation")
		return
	}

	// Writing each rmployee in []employees to DB (ems.csv)
	for _, emp := range empSliceToWrite {
		_, err := WriteToCSV(emp)

		if err != nil {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			logger.Errorf("\nCould not write data to DB during controllers.UpdateEmployee function call")
			return
		}
	}

	// Setting headers and staus codes over suucessful operation.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employeeToUpdate)
}

func SearchEmployees(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	query := r.URL.Query()

	// Get the search criteria from query parameters
	firstName := query.Get("firstName")
	lastName := query.Get("lastName")
	email := query.Get("email")
	role := query.Get("role")

	employees, err := ReadCSV("get")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Errorf("\nCould not read data from DB during controllers.SearchEmployees function call")
		return
	}

	var fetchedEmployees []models.Employee

	empFound := false

	// Iterate through the employees and add matching employees to the results slice
	for _, emp := range employees {
		if (firstName == "" || emp.FirstName == firstName) &&
			(lastName == "" || emp.LastName == lastName) &&
			(email == "" || emp.Email == email) &&
			(role == "" || emp.Role == role) {
			fetchedEmployees = append(fetchedEmployees, emp)
			empFound = true
		}
	}

	if !empFound {
		http.Error(w, "Employee not found!", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Serialize the results to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fetchedEmployees)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	vars := mux.Vars(r)

	empId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid Employee ID", http.StatusBadRequest)
		logger.Errorf("\nInvalid id was employeeId was provided on DeleteEmployee function call")
		return
	}

	employees, err := ReadCSV("delete")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Errorf("\nCould not read data from DB during controllers.DeleteEmployee function call")
		return
	}

	empFound := false

	// Finding employee from employees slice created with ReadCSV() to update it...
	for i, emp := range employees {
		if empId == emp.ID {
			empFound = true
			employees = append(employees[:i], employees[i+1:]...)
			break
		}
	}

	if !empFound {
		http.Error(w, "Employee not found!", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Truncating DB file to avoid duplicate entries
	_, err = os.Create("ems.csv")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Errorf("\nError occured while to truncating DB (ems.csv) during controllers.DeleteEmployee function call")
		return
	}

	// Writing each rmployee in []employees to DB (ems.csv)
	for _, emp := range employees {
		_, err := WriteToCSV(emp)

		if err != nil {
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			logger.Errorf("\nCould not write data to DB during controllers.DeleteEmployee function call")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
