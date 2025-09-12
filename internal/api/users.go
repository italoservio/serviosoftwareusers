package api

import (
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/internal/deps"
	"net/http"
)

func RegisterUsersRoutes(mux *mux.Router, c *deps.Container) {
	mux.
		Handle("/users", http.HandlerFunc(c.UsersHttpAPI.CreateUser)).
		Methods("POST")
	mux.
		Handle("/users/{id}", http.HandlerFunc(c.UsersHttpAPI.GetUserByID)).
		Methods("GET")
	mux.
		Handle("/users/{id}", http.HandlerFunc(c.UsersHttpAPI.DeleteUserByID)).
		Methods("DELETE")
	mux.
		Handle("/users", http.HandlerFunc(c.UsersHttpAPI.ListUsers)).
		Methods("GET")
	mux.
		Handle("/users/{id}", http.HandlerFunc(c.UsersHttpAPI.UpdateUserByID)).
		Methods("PUT", "PATCH")
}
