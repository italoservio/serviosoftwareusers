package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"net/http"
)

type Health struct {
	Status string `json:"status"`
}

func RegisterInfraRoutes(mux *mux.Router) {
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(Health{Status: "ok"})

		w.Write(res)
	})).Methods("GET")
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	res, _ := json.Marshal(exception.AppException{
		Messages:   []string{"Method Not Allowed"},
		StatusText: http.StatusText(http.StatusMethodNotAllowed),
		StatusCode: http.StatusMethodNotAllowed,
	})

	w.Write(res)
}
