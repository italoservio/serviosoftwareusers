package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/italoservio/serviosoftwareusers/internal/modules/users/models"
	"time"
)

func ParseToken(secret string, tokenStr string) (*models.Session, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	session := &models.Session{}

	if userID, ok := claims["userId"].(string); ok {
		session.UserID = userID
	}

	if roles, ok := claims["roles"].([]any); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				session.Roles = append(session.Roles, roleStr)
			}
		}
	}

	if startedAt, err := time.Parse(time.RFC3339, claims["startedAt"].(string)); err == nil {
		session.StartedAt = startedAt
	}

	return session, nil
}
