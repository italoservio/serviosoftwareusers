package users

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/commands"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/repos"
	"github.com/italoservio/serviosoftwareusers/pkg/cast"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"net/http"
)

type UsersHttpAPI struct {
	validate          validator.Validate
	CreateUserCmd     commands.Cmd[commands.CreateUserInput, models.User]
	GetUserByIDCmd    commands.Cmd[commands.GetUserByIDCmdInput, models.User]
	UpdateUserByIDCmd commands.Cmd[commands.UpdateUserByIDCmdInput, models.User]
	DeleteUserByIDCmd commands.CmdNoOutput[commands.DeleteUserByIDCmdInput]
	ListUserCmd       commands.Cmd[commands.ListUserCmdInput, commands.ListUserCmdOutput]
}

func NewUsersHttpAPI(
	validate *validator.Validate,
	createUserCmd *commands.CreateUserCmd,
	getUserByIdCmd *commands.GetUserByIDCmd,
	updateUserByIdCmd *commands.UpdateUserByIDCmd,
	deleteUserByIdCmd *commands.DeleteUserByIDCmd,
	listUserCmd *commands.ListUserCmd,
) *UsersHttpAPI {
	return &UsersHttpAPI{
		validate:          *validate,
		CreateUserCmd:     createUserCmd,
		GetUserByIDCmd:    getUserByIdCmd,
		UpdateUserByIDCmd: updateUserByIdCmd,
		DeleteUserByIDCmd: deleteUserByIdCmd,
		ListUserCmd:       listUserCmd,
	}
}

func (u *UsersHttpAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	var payload commands.CreateUserInput
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

func (u *UsersHttpAPI) GetUserByID(w http.ResponseWriter, r *http.Request) {
	pathParameters := mux.Vars(r)
	userID := pathParameters["id"]

	payload := commands.GetUserByIDCmdInput{ID: userID}

	err := u.validate.Struct(payload)
	if err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	user, err := u.GetUserByIDCmd.Exec(&payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(user)

	w.Write(response)
}

func (u *UsersHttpAPI) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	pathParameters := mux.Vars(r)
	userID := pathParameters["id"]

	body := r.Body
	defer body.Close()

	var payload commands.UpdateUserByIDCmdInput
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		exception.NewPayloadException().WriteJSON(w)
		return
	}

	payload.ID = userID
	err = u.validate.Struct(payload)
	if err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	user, err := u.UpdateUserByIDCmd.Exec(&payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(user)

	w.Write(response)
}

func (u *UsersHttpAPI) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	pathParameters := mux.Vars(r)
	userID := pathParameters["id"]

	payload := commands.DeleteUserByIDCmdInput{ID: userID}

	err := u.validate.Struct(payload)
	if err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	err = u.DeleteUserByIDCmd.Exec(&payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *UsersHttpAPI) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	payload := commands.ListUserCmdInput{
		&repos.ListInput{
			Limit:    cast.StrToInt64(query.Get("limit")),
			Page:     cast.StrToInt64(query.Get("page")),
			SortBy:   cast.StrToStringPtr(query.Get("sortBy")),
			Order:    cast.StrToStringPtr(query.Get("order")),
			Email:    cast.StrSliceToPtr(query["email"]),
			FullName: cast.StrSliceToPtr(query["fullName"]),
			Roles:    cast.StrSliceToPtr(query["roles"]),
		},
	}

	err := u.validate.Struct(payload)
	if err != nil {
		exception.NewValidatorException(err).WriteJSON(w)
		return
	}

	result, err := u.ListUserCmd.Exec(&payload)
	if err != nil {
		exception.ToAppException(err).WriteJSON(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)

	w.Write(response)
}
