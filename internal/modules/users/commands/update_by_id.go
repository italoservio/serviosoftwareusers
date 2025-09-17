package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/env"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
)

type UpdateUserByIDCmd struct {
	envVars env.Env
	repo    repos.UsersRepo
}

func NewUpdateUserByIDCmd(envVars env.Env, repo repos.UsersRepo) *UpdateUserByIDCmd {
	return &UpdateUserByIDCmd{envVars, repo}
}

type UpdateUserByIDCmdInput struct {
	ID        string       `bson:"_id,omitempty" json:"id,omitempty,string" validate:"omitempty,omitnil,mongodb"`
	Session   *jwt.Session `validate:"-"`
	FirstName string       `bson:"firstName" json:"firstName" validate:"omitempty,omitnil,min=2,max=100"`
	LastName  string       `bson:"lastName" json:"lastName" validate:"omitempty,omitnil,min=2,max=100"`
	FullName  string       `bson:"fullName" json:"fullName" validate:"omitempty,omitnil,required,min=2,max=200"`
	Email     string       `bson:"email" json:"email" validate:"omitempty,omitnil,email,min=5,max=200"`
	Roles     []string     `bson:"roles" json:"roles" validate:"omitempty,omitnil,min=1,dive,required,oneofrole"`
	Password  string       `bson:"password,omitempty" json:"password,omitempty" validate:"omitempty,min=8,max=100"`
}

func (c *UpdateUserByIDCmd) Exec(input *UpdateUserByIDCmdInput) (*models.User, error) {
	sess := input.Session
	if sess == nil || (!sess.IsAdmin && input.ID != sess.UserID) {
		return nil, exception.NewForbiddenException()
	}

	existing, err := c.repo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, exception.NewNotFoundException("Usuario nao encontrado")
	}

	if input.FirstName != "" && input.LastName != "" {
		input.FullName = input.FirstName + " " + input.LastName
	}

	if input.FirstName != "" && input.LastName == "" {
		input.FullName = input.FirstName + " " + existing.LastName
	}

	if input.FirstName == "" && input.LastName != "" {
		input.FullName = existing.FirstName + " " + input.LastName
	}

	if input.Password != "" {
		hashedPass, err := HashPass(c.envVars.PASS_SECRET, input.Password)
		if err != nil {
			return nil, exception.NewInternalException(err.Error())
		}

		input.Password = hashedPass
	}

	updated, err := c.repo.UpdateByID(input.ID, &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		FullName:  input.FullName,
		Email:     input.Email,
		Roles:     input.Roles,
		Password:  input.Password,
	})
	if err != nil {
		err := exception.NewRepoException(
			"Nao foi possivel atualizar o usuario, tente novamente mais tarde",
			err.Error(),
		)
		return nil, err
	}

	return updated, nil
}
