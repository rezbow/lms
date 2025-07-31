package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	CodeInvalidJsonErr = "INVALID_JSON"
	CodeInvalidDataErr = "INVALID_DATA"
	CodeNotFoundErr    = "NOT_FOUND"
	CodeInternalErr    = "INTERNAL_ERR"
	CodeActiveLoan     = "ACTIVE_LOAN"
)

func extractValidationErrs(err error) []gin.H {
	errors := []gin.H{}
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, gin.H{"field": e.Field(), "error": fieldErrorMessage(e)})
	}
	return errors
}

func fieldErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return "Minimum length is " + e.Param()
	case "max":
		return "Maximum length is " + e.Param()
	default:
		return "invalid value"
	}
}
