package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClaimSpec struct {
	UserID string
}

type Claims struct {
	jwt.RegisteredClaims
}
