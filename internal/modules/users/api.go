package users

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/commands"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"net/http"
)

type UsersHttpAPI struct {
	validate      validator.Validate
	CreateUserCmd commands.Cmd[models.User, models.User]
}

func NewUsersHttpAPI(
	validate *validator.Validate,
	createUserCmd *commands.CreateUserCmd,
) *UsersHttpAPI {
	return &UsersHttpAPI{
		validate:      *validate,
		CreateUserCmd: createUserCmd,
	}
}

func (u *UsersHttpAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	var payload models.User
	err := json.NewDecoder(body).Decode(&payload)

	if err != nil {
		exception.NewPayloadException().WriteJSON(w)
		return
	}

	err = u.validate.Struct(payload)
	if err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	user, err := u.CreateUserCmd.Exec(&payload)

	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, _ := json.Marshal(user)
	w.Write(response)

	return
}
