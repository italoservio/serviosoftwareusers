package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftwareusers/pkg/rbac"
	"strings"
)

func validateRole(fl validator.FieldLevel) bool {
	for _, role := range rbac.GetAllRoles() {
		if strings.EqualFold(fl.Field().String(), role.String()) {
			return true
		}
	}

	return false
}
