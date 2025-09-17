package rbac

import (
	"errors"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"github.com/italoservio/serviosoftwareusers/pkg/jwt"
	"net/http"
	"strings"
)

func Middleware(allowed []Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := getContextSession(r)
			if err != nil {
				exception.NewUnauthorizedException(err.Error()).WriteJSON(w)
				return
			}

			allowedMap := make(map[string]bool)
			for _, role := range allowed {
				allowedMap[role.String()] = true
			}

			hasRole := false
			for _, role := range session.Roles {
				hasRole = allowedMap[role]

				if hasRole && strings.Contains(role, "admin") {
					session.IsAdmin = true
				}
			}

			if !hasRole {
				exception.NewForbiddenException().WriteJSON(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getContextSession(r *http.Request) (*jwt.Session, error) {
	ctx := r.Context()
	ctxSession := ctx.Value("session")

	if ctxSession == nil {
		return nil, errors.New("Sessao nao encontrada no contexto")
	}

	session, ok := ctxSession.(*jwt.Session)

	if !ok {
		return nil, errors.New("Sessao invalida no contexto")
	}

	return session, nil
}
