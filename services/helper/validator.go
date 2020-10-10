package helper

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}
