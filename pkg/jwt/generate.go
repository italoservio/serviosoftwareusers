package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"time"
)

func GenerateToken(secret string, session models.Session) (string, error) {
	if secret == "" {
		return "", errors.New("Nao ha chave de criptografia de token")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"userId":    session.UserID,
		"roles":     session.Roles,
		"startedAt": session.StartedAt.UTC().Format(time.RFC3339),
		"exp":       expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
