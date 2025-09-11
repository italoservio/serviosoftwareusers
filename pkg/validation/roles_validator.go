package validation

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

var AllowedRoles = []string{"ads:admin", "ads:user"}

func validateRole(fl validator.FieldLevel) bool {
	for _, role := range AllowedRoles {
		if strings.EqualFold(fl.Field().String(), role) { // Case-insensitive comparison
			return true
		}
	}

	return false
}
