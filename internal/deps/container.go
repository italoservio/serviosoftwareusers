package deps

import (
	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users"
	userscmds "github.com/italoservio/serviosoftwareusers/internal/modules/users/commands"
	usersmodels "github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	usersrepos "github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/db"
	"github.com/italoservio/serviosoftwareusers/pkg/env"
	"github.com/italoservio/serviosoftwareusers/pkg/validation"
)

type Container struct {
	DB                db.DB
	Env               env.Env
	Validator         validator.Validate
	UsersRepo         usersrepos.UsersRepo
	UsersHttpAPI      users.UsersHttpAPI
	CreateUserCmd     userscmds.Cmd[userscmds.CreateUserInput, usersmodels.User]
	GetUserByIDCmd    userscmds.Cmd[userscmds.GetUserByIDCmdInput, usersmodels.User]
	UpdateUserByIdCmd userscmds.Cmd[userscmds.UpdateUserByIDCmdInput, usersmodels.User]
	DeleteUserByIdCmd userscmds.CmdNoOutput[userscmds.DeleteUserByIDCmdInput]
	ListUsersCmd      userscmds.Cmd[userscmds.ListUserCmdInput, userscmds.ListUserCmdOutput]
	SignInCmd         userscmds.Cmd[userscmds.SignInCmdInput, userscmds.SignInCmdOutput]
}

func NewContainer(envVars *env.Env) *Container {
	dbConn, err := db.NewDB((*envVars).MONGODB_URI)
	if err != nil {
		panic(err)
	}

	valdtr, err := validation.NewValidator()

	usersRepository := usersrepos.NewMongoUsersRepo(dbConn)
	createUserCmd := userscmds.NewCreateUserCmd(*envVars, usersRepository)
	getUserByIDCmd := userscmds.NewGetUserByIDCmd(usersRepository)
	updateUserByIdCmd := userscmds.NewUpdateUserByIDCmd(*envVars, usersRepository)
	deleteUserByIdCmd := userscmds.NewDeleteUserByIDCmd(usersRepository)
	listUserCmd := userscmds.NewListUserCmd(usersRepository)
	signInCmd := userscmds.NewSignInCmd(*envVars, usersRepository)

	usersHttpAPI := users.NewUsersHttpAPI(
		valdtr,
		createUserCmd,
		getUserByIDCmd,
		updateUserByIdCmd,
		deleteUserByIdCmd,
		listUserCmd,
		signInCmd,
	)

	if err != nil {
		panic(err)
	}

	return &Container{
		DB:                *dbConn,
		Env:               *envVars,
		Validator:         *valdtr,
		UsersRepo:         usersRepository,
		CreateUserCmd:     createUserCmd,
		GetUserByIDCmd:    getUserByIDCmd,
		UpdateUserByIdCmd: updateUserByIdCmd,
		DeleteUserByIdCmd: deleteUserByIdCmd,
		ListUsersCmd:      listUserCmd,
		SignInCmd:         signInCmd,
		UsersHttpAPI:      *usersHttpAPI,
	}
}
