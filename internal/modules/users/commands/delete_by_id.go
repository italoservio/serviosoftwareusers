package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
)

type DeleteUserByIDCmd struct {
	repo repos.UsersRepo
}

type DeleteUserByIDCmdInput struct {
	ID      string       `validate:"required,mongodb"`
	Session *jwt.Session `validate:"-"`
}

func NewDeleteUserByIDCmd(repo repos.UsersRepo) *DeleteUserByIDCmd {
	return &DeleteUserByIDCmd{repo}
}

func (c *DeleteUserByIDCmd) Exec(input *DeleteUserByIDCmdInput) error {
	sess := input.Session
	if sess == nil || (!sess.IsAdmin && input.ID != sess.UserID) {
		return exception.NewForbiddenException()
	}

	err := c.repo.DeleteByID(input.ID)
	if err != nil {
		return exception.NewRepoException(
			"Nao foi possivel remover o usuario, tente novamente mais tarde",
			err.Error(),
		)
	}

	return nil
}
