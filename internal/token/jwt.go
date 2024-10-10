package token

import (
	"github.com/benbjohnson/clock"
	"github.com/golang-jwt/jwt/v5"
	"kg/procurement/cmd/config"
	"time"

	"github.com/google/uuid"
)

type jwtManager struct {
	cfg   config.Token
	clock clock.Clock
}

func (s *jwtManager) GenerateToken(spec ClaimSpec) (string, error) {
	tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   spec.UserID,
			ExpiresAt: jwt.NewNumericDate(s.clock.Now().UTC().Add(30 * 24 * time.Hour)), // one month from now
			IssuedAt:  jwt.NewNumericDate(s.clock.Now().UTC()),
			ID:        uuid.NewString(),
		},
	})

	token, err := tokenObject.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			_, isAccSigningMethod := t.Method.(*jwt.SigningMethodHMAC)
			if !isAccSigningMethod {
				return nil, jwt.ErrInvalidKeyType
			}
			return []byte(s.cfg.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKeyType
	}

	return claims, nil
}

// newJWTManager is only accessible by the Token package
// entrypoint for other verticals should refer to the interface declared on service
func newJWTManager(cfg config.Token, clock clock.Clock) *jwtManager {
	return &jwtManager{
		cfg:   cfg,
		clock: clock,
	}
}
