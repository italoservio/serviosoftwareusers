package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
)

type GetUserByIDCmd struct {
	repo repos.UsersRepo
}

type GetUserByIDCmdInput struct {
	ID      string       `validate:"required,mongodb"`
	Session *jwt.Session `validate:"-"`
}

func NewGetUserByIDCmd(repo repos.UsersRepo) *GetUserByIDCmd {
	return &GetUserByIDCmd{repo}
}

func (c *GetUserByIDCmd) Exec(input *GetUserByIDCmdInput) (*models.User, error) {
	sess := input.Session
	if sess == nil || (!sess.IsAdmin && input.ID != sess.UserID) {
		return nil, exception.NewForbiddenException()
	}

	user, err := c.repo.GetByID(input.ID)
	if err != nil {
		return nil, exception.NewRepoException(
			"Nao foi possivel buscar o usuario, tente novamente mais tarde",
			err.Error(),
		)
	}

	if user == nil {
		return nil, exception.NewNotFoundException("Usuario nao encontrado")
	}

	return user, nil
}
