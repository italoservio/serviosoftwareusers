package exception

import (
	"net/http"
)

type RepositoryException struct {
	Messages   []string `json:"messages"`
	StatusText string   `json:"statusText"`
	StatusCode int      `json:"statusCode"`
}

func NewRepoException(messages ...string) *AppException {
	if len(messages) == 0 {
		messages = []string{"Erro inesperado ao processar a requisicao ao repositorio"}
	}

	return NewAppException(messages, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
