package api

import (
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/internal/deps"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
	"github.com/italoservio/serviosoftwareusers/pkg/rbac"
	"net/http"
)

func RegisterUsersRoutes(mux *mux.Router, c *deps.Container) {
	protected := mux.PathPrefix("/users").Subrouter()
	protected.Use(jwt.Middleware(c.Env))

	protectedAdminOnly := protected.PathPrefix("").Subrouter()
	protectedAdminOnly.Use(rbac.Middleware(rbac.GetAdminRoles()))

	protectedAdminAndOwners := protected.PathPrefix("").Subrouter()
	protectedAdminAndOwners.Use(rbac.Middleware(rbac.GetAllRoles()))

	protectedAdminOnly.
		Handle("/{userId}", http.HandlerFunc(c.UsersHttpAPI.DeleteUserByID)).
		Methods("DELETE")
	protectedAdminOnly.
		Handle("", http.HandlerFunc(c.UsersHttpAPI.ListUsers)).
		Methods("GET")
	protectedAdminAndOwners.
		Handle("/{userId}", http.HandlerFunc(c.UsersHttpAPI.GetUserByID)).
		Methods("GET")
	protectedAdminAndOwners.
		Handle("/{userId}", http.HandlerFunc(c.UsersHttpAPI.UpdateUserByID)).
		Methods("PUT", "PATCH")

	public := mux.PathPrefix("/users").Subrouter()
	public.
		Handle("", http.HandlerFunc(c.UsersHttpAPI.CreateUser)).
		Methods("POST")
	public.
		Handle("/signin", http.HandlerFunc(c.UsersHttpAPI.SignIn)).
		Methods("POST")
}
