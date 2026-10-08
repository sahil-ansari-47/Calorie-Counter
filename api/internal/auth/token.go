package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	Secret     []byte
	Expiration time.Duration
}

func (s TokenService) Issue(userID uuid.UUID) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userID.String(), "iat": now.Unix(), "exp": now.Add(s.Expiration).Unix()}).SignedString(s.Secret)
}
func (s TokenService) Parse(token string) (uuid.UUID, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.Secret, nil
	})
	if err != nil || !parsed.Valid {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}
	subject, err := parsed.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(subject)
}
