package controllers

import (
	"Task2/models"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

// Helper functions
func ReadCSV(methodType string) ([]models.Employee, error) {
	file, err := os.Open("ems.csv")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var employees []models.Employee

	for _, record := range records {

		if err != nil {
			return nil, err
		}

		id, _ := strconv.Atoi(record[0])
		salaryFloat, _ := strconv.ParseFloat(record[7], 64)

		var password string
		if methodType == "get" {
			password = ""
		} else {
			password = record[4]
		}

		employee := models.Employee{
			ID:        id,
			FirstName: record[1],
			LastName:  record[2],
			Email:     record[3],
			Password:  password,
			PhoneNo:   record[5],
			Role:      record[6],
			Salary:    salaryFloat,
			Birthdate: record[8],
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func WriteToCSV(employee models.Employee) (emp models.Employee, err error) {
	file, err := os.OpenFile("ems.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Printf("%v :occured while opening DB file ", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	employee.FirstName = strings.TrimSpace(employee.FirstName)
	employee.FirstName = strings.ToUpper(employee.FirstName[0:1]) + strings.ToLower(employee.FirstName[1:])

	employee.LastName = strings.TrimSpace(employee.LastName)
	employee.LastName = strings.ToUpper(employee.LastName[0:1]) + strings.ToLower(employee.LastName[1:])

	err1 := writer.Write([]string{
		strconv.Itoa(employee.ID),
		employee.FirstName,
		employee.LastName,
		employee.Email,
		employee.Password,
		employee.PhoneNo,
		employee.Role,
		strconv.FormatFloat(employee.Salary, 'f', -1, 64),
		employee.Birthdate,
	})

	if err1 != nil {
		log.Println("Error occured while writting record to DB (ems.csv)")
		return employee, err1
	}

	return employee, nil
}

func GenerateEmployeeId() int {
	file, err := os.Open("ems.csv")

	if err != nil {
		log.Println("Error opening csv")
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	max := 0
	var eId int
	for _, emp := range records {
		eId, _ = strconv.Atoi(emp[0])
		if eId >= max {
			max = eId
		}
	}
	return max + 1
}
