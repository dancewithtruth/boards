package validator

import "github.com/go-playground/validator/v10"

type Validate struct {
	*validator.Validate
}

func New() Validate {
	validate := validator.New()
	return Validate{validate}
}
