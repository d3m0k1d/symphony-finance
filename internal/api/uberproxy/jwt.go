package uberproxy

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s server) parseJwt(j string) (uid string, err error) {
	tok, err := jwt.ParseWithClaims(j, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{s.jwtMethod.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return "", err
	}
	return tok.Claims.GetSubject()
}
func (s server) newJWTString(user_id string) (string, error) {
	j, e := s.newJwt(user_id)
	if e != nil {
		return "", e
	}
	return j.SignedString(s.jwtSecret)
}

func (s server) newJwt(uid string) (*jwt.Token, error) {
	claims :=
		jwt.RegisteredClaims{
			Issuer:    "WS AG",
			Subject:   uid,
			Audience:  jwt.ClaimStrings{"symphony banking Frontend"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		}
	return jwt.NewWithClaims(s.jwtMethod, claims), nil

}
