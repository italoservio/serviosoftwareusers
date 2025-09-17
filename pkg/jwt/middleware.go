package jwt

import (
	"context"
	"github.com/italoservio/serviosoftwareusers/pkg/env"
	"github.com/italoservio/serviosoftwareusers/pkg/exception"
	"net/http"
)

func Middleware(e env.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header

			if headers.Get("Authorization") == "" {
				exception.NewUnauthorizedException("Nao ha cabecalho de autorizacao").WriteJSON(w)
				return
			}

			authorization := headers.Get("Authorization")

			if len(authorization) < 7 || authorization[:7] != "Bearer " {
				exception.NewUnauthorizedException("Formato de token invalido").WriteJSON(w)
				return
			}

			token := authorization[7:]

			if token == "" {
				exception.NewUnauthorizedException("Formato de token invalido").WriteJSON(w)
				return
			}

			session, err := ParseToken(e.AUTH_SECRET, token)
			if err != nil {
				exception.NewUnauthorizedException("Token invalido: " + err.Error()).WriteJSON(w)
				return
			}

			if session.Expired {
				exception.NewUnauthorizedException("Token expirado").WriteJSON(w)
				return
			}

			ctx := r.Context()
			ctxWithValue := context.WithValue(ctx, "session", session)
			r = r.WithContext(ctxWithValue)

			next.ServeHTTP(w, r)
		})
	}
}
