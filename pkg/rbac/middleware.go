package rbac

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
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
			hasAdminRole := false
			for _, role := range session.Roles {
				hasRole = allowedMap[role]

				if hasRole && strings.Contains(role, "admin") {
					hasAdminRole = true
				}
			}

			if !hasRole || !isResourceOwnerOrAdmin(r, session, hasAdminRole) {
				exception.NewForbiddenException().WriteJSON(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getContextSession(r *http.Request) (*models.Session, error) {
	ctx := r.Context()
	ctxSession := ctx.Value("session")

	if ctxSession == nil {
		return nil, errors.New("Sessao nao encontrada no contexto")
	}

	session, ok := ctxSession.(*models.Session)

	if !ok {
		return nil, errors.New("Sessao invalida no contexto")
	}

	return session, nil
}

func isResourceOwnerOrAdmin(r *http.Request, session *models.Session, hasAdminRole bool) bool {
	pathParameters := mux.Vars(r)
	userId, hasUserIdInPathParams := pathParameters["userId"]

	isResourceOwner := userId == (*session).UserID

	if hasUserIdInPathParams && (!isResourceOwner && !hasAdminRole) {
		return false
	}

	return true
}
