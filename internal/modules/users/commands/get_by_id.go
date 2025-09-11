package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
)

type GetUserByIDCmd struct {
	repo repos.UsersRepo
}

type GetUserByIDCmdInput struct {
	ID string `validate:"required,mongodb"`
}

func NewGetUserByIDCmd(repo repos.UsersRepo) *GetUserByIDCmd {
	return &GetUserByIDCmd{repo: repo}
}

func (c *GetUserByIDCmd) Exec(input *GetUserByIDCmdInput) (*models.User, error) {
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
