package validations

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

var V = validator.New()

func init() {
	// Registering custom validation functions...
	V.RegisterValidation("phone", validatePhoneNumber)
	V.RegisterValidation("alpha_space", validateName)
	V.RegisterValidation("role", validateRole)
	V.RegisterValidation("date", validateDate)
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	// This regex pattern will cover almost all international numbers.
	regexPattern := regexp.MustCompile(`\+?\d{1,4}(?:[-|\s])?\d{1,15}`)
	return regexPattern.MatchString(fl.Field().String())
}

func validateName(fl validator.FieldLevel) bool {
	regexPattern := regexp.MustCompile(`^[a-zA-Z ]+$`)
	return regexPattern.MatchString(fl.Field().String())
}

func validateRole(fl validator.FieldLevel) bool {
	allowedValues := []string{"admin", "manager", "developer", "tester"}
	role := fl.Field().String()
	for _, allowedValue := range allowedValues {
		if role == allowedValue {
			return true
		}
	}
	return false
}

func validateDate(fl validator.FieldLevel) bool {
	// Defining a regular expression pattern for a valid date format (e.g., YYYY-MM-DD)
	datePattern := `^\d{4}-\d{2}-\d{2}$`
	date := fl.Field().String()
	match, _ := regexp.MatchString(datePattern, date)
	return match
}
