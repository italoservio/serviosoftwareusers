package rbac

import "strings"

type Role string

const (
	ADS_ADMIN_ROLE = "ads:admin"
	ADS_USER_ROLE  = "ads:user"
)

func (r Role) IsAdmin() bool {
	return strings.Contains(string(r), "admin")
}

func (r Role) String() string {
	return string(r)
}

func GetAllRoles() []Role {
	return []Role{
		ADS_ADMIN_ROLE,
		ADS_USER_ROLE,
	}
}

func GetAdminRoles() []Role {
	return []Role{
		ADS_ADMIN_ROLE,
	}
}
