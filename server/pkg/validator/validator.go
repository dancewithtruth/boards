package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validate struct {
	*validator.Validate
}

func New() Validate {
	validate := validator.New()
	return Validate{validate}
}

func GetValidationErrMsg(err error) string {
	errMsg := ""
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErr := fieldErrors[0]
		switch fieldErr.Tag() {
		case "required":
			errMsg = fmt.Sprintf("%s is a required field", fieldErr.Field())
		default:
			errMsg = fmt.Sprintf("Invalid input on %s", fieldErr.Field())
		}
	}
	return errMsg
}
