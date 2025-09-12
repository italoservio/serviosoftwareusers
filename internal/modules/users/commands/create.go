package commands

import (
	"fmt"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
)

type CreateUserCmd struct {
	repo repos.UsersRepo
}

func NewCreateUserCmd(repo repos.UsersRepo) *CreateUserCmd {
	return &CreateUserCmd{repo: repo}
}

type CreateUserInput struct {
	FirstName string   `bson:"firstName" json:"firstName" validate:"required,min=2,max=100"`
	LastName  string   `bson:"lastName" json:"lastName" validate:"required,min=2,max=100"`
	Email     string   `bson:"email" json:"email" validate:"required,email,min=5,max=200"`
	Roles     []string `bson:"roles" json:"roles" validate:"required,min=1,dive,required,oneofrole"`
	Password  string   `bson:"password,omitempty" json:"password,omitempty" validate:"required,min=8,max=100"`
}

func (c *CreateUserCmd) Exec(input *CreateUserInput) (*models.User, error) {
	if existing, err := c.repo.GetByEmail(input.Email); err != nil || existing != nil {
		return nil, exception.NewResourceExistsException("Usuario ja existente")
	}

	user := &models.User{}
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email
	user.Roles = input.Roles

	user.FullName = fmt.Sprintf("%s %s", input.FirstName, input.LastName)

	hashedPass, err := HashPass(input.Password)
	if err != nil {
		return nil, exception.NewInternalException(err.Error())
	}

	user.Password = hashedPass
	created, err := c.repo.Create(user)

	if err != nil {
		err := exception.NewRepoException(
			"Nao foi possivel criar o usuario, tente novamente mais tarde",
			err.Error(),
		)
		return nil, err
	}

	return created, nil
}
