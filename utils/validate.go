package utils

import (
	"fmt"
	"unicode"

	"githb.com/demo-employee-api/internal/entity"
	"github.com/go-playground/validator/v10"
)

type ValidationResult struct {
	ErrorMessage  string
	FieldErrors   []map[string]string
	OriginalError error
	Success       bool
}

func (res *ValidationResult) HasFieldErrors() bool {
	return len(res.FieldErrors) > 0

}

func validatePassword(fl validator.FieldLevel) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	pass := fl.Field().String()

	if len(pass) >= 8 {
		hasMinLen = true
	}
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func validateRole(fl validator.FieldLevel) bool {
	if fl.Field().String() == entity.RoleAdmin || fl.Field().String() == entity.RoleEmployee {
		return true
	}
	return false
}

func ValidateStruct(val interface{}) ValidationResult {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("role", validateRole)
	// validate.RegisterAlias("upwd", "min=8,max=20,alphanum")

	if err := validate.Struct(val); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {

			return ValidationResult{
				OriginalError: err,
				ErrorMessage:  "failed to perform validation",
				Success:       false,
			}

		}
		var FieldErrors []map[string]string

		for _, err := range err.(validator.ValidationErrors) {
			message := func() string {
				switch err.Tag() {
				case "required":
					return "value is required"
				case "email":
					return "invalid email format"
				case "password":
					return "invalid password format"
				case "role":
					return "invalid role"
				default:
					return fmt.Sprintf("%s validation failed", err.Tag())

				}
			}()
			FieldErrors = append(FieldErrors, map[string]string{"field": err.Field(), "message": message, "tag": err.Tag()})
		}
		return ValidationResult{
			OriginalError: err,
			FieldErrors:   FieldErrors,
			ErrorMessage:  "request validation failed",
			Success:       false,
		}
	}
	return ValidationResult{
		Success: true,
	}
}
