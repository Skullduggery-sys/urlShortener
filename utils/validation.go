package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateErrors(errs validator.ValidationErrors) error {
	errorMsgs := make([]string, 0)

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s was not filled", err.Field()))
		case "url":
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s url is wrong", err.Field()))
		case "numeric":
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s url is wrong", err.Field()))
		default:
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s unknown error", err.Field()))
		}
	}

	return errors.New(strings.Join(errorMsgs, ", "))
}
