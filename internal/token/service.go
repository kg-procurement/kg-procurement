//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=token
package token

import (
	"github.com/benbjohnson/clock"
	"kg/procurement/cmd/config"
)

type TokenManager interface {
	GenerateToken(spec ClaimSpec) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type TokenService struct {
	TokenManager
}

func (s *TokenService) GenerateToken(spec ClaimSpec) (string, error) {
	return s.TokenManager.GenerateToken(spec)
}

func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	return s.TokenManager.ValidateToken(tokenString)
}

func NewTokenService(cfg config.Token, clock clock.Clock) *TokenService {
	return &TokenService{
		TokenManager: newJWTManager(cfg, clock),
	}
}
