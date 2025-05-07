package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StatusType string

const (
	StatusOK    StatusType = "OK"
	StatusError StatusType = "Error"
)

type Response struct {
	Status StatusType `json:"status"`
	Error  string     `json:"error,omitempty"`
}

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

func ValidationError(errs validator.ValidationErrors) Response {
	return Response{
		Status: StatusError,
		Error:  composeValidationMessage(errs),
	}
}

type ValidationRule struct {
	Tag     string
	Message string
}

var DefaultValidationRules = []ValidationRule{
	{"required", "field %s is a required field"},
	{"url", "field %s is not a valid URL"},
	{"email", "field %s is not a valid email"},
	{"min", "field %s must be at least %s characters"},
	{"max", "field %s must be at most %s characters"},
}

func composeValidationMessage(errs validator.ValidationErrors) string {
	var messages []string

	for _, err := range errs {
		message := getValidationMessage(err)
		messages = append(messages, message)
	}

	return strings.Join(messages, "; ")
}

func getValidationMessage(err validator.FieldError) string {
	for _, rule := range DefaultValidationRules {
		if rule.Tag == err.ActualTag() {
			switch err.ActualTag() {
			case "min", "max":
				return fmt.Sprintf(rule.Message, err.Field(), err.Param())
			default:
				return fmt.Sprintf(rule.Message, err.Field())
			}
		}
	}

	return fmt.Sprintf("field %s is not valid (validation: %s)", err.Field(), err.ActualTag())
}
