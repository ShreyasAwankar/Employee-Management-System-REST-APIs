package models

type Employee struct {
	ID        int     `json:"employee_id"`
	FirstName string  `json:"first name" validate:"required,alpha_space"`
	LastName  string  `json:"last name" validate:"required,alpha_space"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password,omitempty" validate:"required,min=6"`
	PhoneNo   string  `json:"phone no" validate:"required,phone"`
	Role      string  `json:"role" validate:"required,role"`
	Salary    float64 `json:"salary" validate:"required,numeric,gt=0"`
	Birthdate string  `json:"birthdate" validate:"required,date"`
}
