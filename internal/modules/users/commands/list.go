package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
)

type ListUserCmd struct {
	repo repos.UsersRepo
}

func NewListUserCmd(repo repos.UsersRepo) *ListUserCmd {
	return &ListUserCmd{repo}
}

type ListUserCmdInput struct {
	*repos.ListInput
}

type ListUserCmdOutput struct {
	*repos.ListOutput
}

func (c *ListUserCmd) Exec(input *ListUserCmdInput) (*ListUserCmdOutput, error) {
	if input.SortBy == nil {
		createdAt := "createdAt"
		input.SortBy = &createdAt
	}

	if input.Order == nil {
		desc := "desc"
		input.Order = &desc
	}

	output, err := c.repo.List(input.ListInput)
	if err != nil {
		return nil, err
	}

	return &ListUserCmdOutput{output}, nil
}
