package commands

import (
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
	"os"
	"time"
)

type SignInCmd struct {
	repo repos.UsersRepo
}

func NewSignInCmd(repo repos.UsersRepo) *SignInCmd {
	return &SignInCmd{repo: repo}
}

type SignInCmdInput struct {
	Email    string `json:"email" validate:"required,email,min=5,max=200"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=100"`
}

type SignInCmdOutput struct {
	Token string `json:"token"`
}

func (c *SignInCmd) Exec(input *SignInCmdInput) (*SignInCmdOutput, error) {
	user, err := c.repo.GetByEmail(input.Email)
	if err != nil {
		return nil, exception.NewPayloadException("Usuario ou senha invalidos")
	}

	if !c.comparePass(user.Password, input.Password) {
		return nil, exception.NewPayloadException("Usuario ou senha invalidos")
	}

	session := models.Session{
		UserID:    user.ID.Hex(),
		Roles:     user.Roles,
		StartedAt: time.Now().UTC(),
	}

	token, err := jwt.GenerateToken(os.Getenv("AUTH_SECRET"), session)
	if err != nil {
		return nil, exception.NewInternalException("Nao foi possivel gerar o token de autenticacao")
	}

	go c.updateSignedInAt(user, &session)

	return &SignInCmdOutput{Token: token}, nil
}

func (c *SignInCmd) comparePass(userPass, plainPass string) bool {
	hashedPass, err := HashPass(plainPass)
	if err != nil {
		return false
	}

	return userPass == hashedPass
}

func (c *SignInCmd) updateSignedInAt(user *models.User, session *models.Session) {
	_, err := c.repo.UpdateByID(user.StringID(), &models.User{SignedInAt: &session.StartedAt})
	if err != nil {
		println("failed to update user signedInAt:", err.Error())
	}
}
