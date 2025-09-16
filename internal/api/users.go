package api

import (
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/internal/deps"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
	"net/http"
)

func RegisterUsersRoutes(mux *mux.Router, c *deps.Container) {
	protected := mux.PathPrefix("/users").Subrouter()
	protected.Use(jwt.Middleware(c.Env))
	protected.
		Handle("/{id}", http.HandlerFunc(c.UsersHttpAPI.GetUserByID)).
		Methods("GET")
	protected.
		Handle("/{id}", http.HandlerFunc(c.UsersHttpAPI.DeleteUserByID)).
		Methods("DELETE")
	protected.
		Handle("", http.HandlerFunc(c.UsersHttpAPI.ListUsers)).
		Methods("GET")
	protected.
		Handle("/{id}", http.HandlerFunc(c.UsersHttpAPI.UpdateUserByID)).
		Methods("PUT", "PATCH")

	public := mux.PathPrefix("/users").Subrouter()
	public.
		Handle("", http.HandlerFunc(c.UsersHttpAPI.CreateUser)).
		Methods("POST")
	public.
		Handle("/signin", http.HandlerFunc(c.UsersHttpAPI.SignIn)).
		Methods("POST")
}
