package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
)

type DeleteUserByIDCmd struct {
	repo repos.UsersRepo
}

type DeleteUserByIDCmdInput struct {
	ID string `validate:"required,mongodb"`
}

func NewDeleteUserByIDCmd(repo repos.UsersRepo) *DeleteUserByIDCmd {
	return &DeleteUserByIDCmd{repo}
}

func (c *DeleteUserByIDCmd) Exec(input *DeleteUserByIDCmdInput) error {
	err := c.repo.DeleteByID(input.ID)
	if err != nil {
		return exception.NewRepoException(
			"Nao foi possivel remover o usuario, tente novamente mais tarde",
			err.Error(),
		)
	}

	return nil
}
