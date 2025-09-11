package deps

import (
	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users"
	userscmds "github.com/italoservio/serviosoftwareusers/internal/modules/users/commands"
	usersmodels "github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	usersrepos "github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/db"
	"github.com/italoservio/serviosoftwareusers/pkg/validation"
)

type Container struct {
	Validator     validator.Validate
	UsersRepo     usersrepos.UsersRepo
	UsersHttpAPI  users.UsersHttpAPI
	CreateUserCmd userscmds.Cmd[usersmodels.User, usersmodels.User]
}

func NewContainer(db *db.DB) *Container {
	val, err := validation.NewValidator()

	ur := usersrepos.NewMongoUsersRepo(db)
	cuc := userscmds.NewCreateUserCmd(ur)

	ua := users.NewUsersHttpAPI(val, cuc)

	if err != nil {
		panic(err)
	}

	return &Container{
		Validator:     *val,
		UsersRepo:     ur,
		CreateUserCmd: cuc,
		UsersHttpAPI:  *ua,
	}
}
