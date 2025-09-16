package api

import (
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/internal/deps"
	"net/http"
)

func RegisterUsersRoutes(mux *mux.Router, c *deps.Container) {
	usersRouter := mux.PathPrefix("/users").Subrouter()

	usersRouter.
		Handle("", http.HandlerFunc(c.UsersHttpAPI.CreateUser)).
		Methods("POST")
	usersRouter.
		Handle("/{id}", http.HandlerFunc(c.UsersHttpAPI.GetUserByID)).
		Methods("GET")
	usersRouter.
		Handle("/{id}", http.HandlerFunc(c.UsersHttpAPI.DeleteUserByID)).
		Methods("DELETE")
	usersRouter.
		Handle("", http.HandlerFunc(c.UsersHttpAPI.ListUsers)).
		Methods("GET")
	usersRouter.
		Handle("/{id}", http.HandlerFunc(c.UsersHttpAPI.UpdateUserByID)).
		Methods("PUT", "PATCH")

	usersRouter.
		Handle("/signin", http.HandlerFunc(c.UsersHttpAPI.SignIn)).
		Methods("POST")
}
