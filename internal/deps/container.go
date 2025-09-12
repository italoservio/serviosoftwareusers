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
	Validator         validator.Validate
	UsersRepo         usersrepos.UsersRepo
	UsersHttpAPI      users.UsersHttpAPI
	CreateUserCmd     userscmds.Cmd[userscmds.CreateUserInput, usersmodels.User]
	GetUserByIDCmd    userscmds.Cmd[userscmds.GetUserByIDCmdInput, usersmodels.User]
	DeleteUserByIdCmd userscmds.CmdNoOutput[userscmds.DeleteUserByIDCmdInput]
}

func NewContainer(db *db.DB) *Container {
	v, err := validation.NewValidator()

	usersRepository := usersrepos.NewMongoUsersRepo(db)
	createUserCmd := userscmds.NewCreateUserCmd(usersRepository)
	getUserByIDCmd := userscmds.NewGetUserByIDCmd(usersRepository)
	deleteUserByIdCmd := userscmds.NewDeleteUserByIDCmd(usersRepository)
	listUserCmd := userscmds.NewListUserCmd(usersRepository)

	usersHttpAPI := users.NewUsersHttpAPI(
		v,
		createUserCmd,
		getUserByIDCmd,
		deleteUserByIdCmd,
		listUserCmd,
	)

	if err != nil {
		panic(err)
	}

	return &Container{
		Validator:         *v,
		UsersRepo:         usersRepository,
		CreateUserCmd:     createUserCmd,
		GetUserByIDCmd:    getUserByIDCmd,
		UsersHttpAPI:      *usersHttpAPI,
		DeleteUserByIdCmd: deleteUserByIdCmd,
	}
}
