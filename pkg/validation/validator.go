package validation

import "github.com/go-playground/validator/v10"

func NewValidator() (*validator.Validate, error) {
	validate := validator.New()
	err := validate.RegisterValidation("oneofrole", validateRole)

	if err != nil {
		return nil, err
	}

	return validate, nil
}
