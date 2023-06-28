package validator

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Validate is a struct that holds the validator implementation.
type Validate struct {
	*validator.Validate
}

// New initializes a new validator.
func New() Validate {
	validate := validator.New()
	return Validate{validate}
}

// GetValidationErrMsg checks to see if the provided err is a validation error and
// returns the first validation error message.
func GetValidationErrMsg(s interface{}, err error) string {
	errMsg := ""
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErr := fieldErrors[0]
		fieldName := getStructTag(s, fieldErr.Field(), "json")
		switch fieldErr.Tag() {
		case "required":
			errMsg = fmt.Sprintf("%s is a required field", fieldName)
		default:
			errMsg = fmt.Sprintf("Invalid input on %s", fieldName)
		}
	}
	return errMsg
}

func getStructTag(s interface{}, fieldName string, tagKey string) string {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return fieldName
	}
	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName
	}

	return field.Tag.Get(tagKey)
}

// IsValidationError checks to see if error is of type validator.ValidationErrors
func IsValidationError(err error) bool {
	if _, ok := err.(validator.ValidationErrors); ok {
		return true
	}
	return false
}
