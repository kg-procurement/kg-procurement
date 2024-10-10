//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=account
package account

import (
	"context"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type accountDBAccessor interface {
	RegisterAccount(ctx context.Context, account Account) error
}

type AccountService struct {
	accountDBAccessor
}

func (a *AccountService) RegisterAccount(ctx context.Context, spec RegisterAccountSpec) error {
	// Validate email
	if _, err := mail.ParseAddress(spec.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(spec.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create account
	account := Account{
		Email:    spec.Email,
		Password: string(hashedPassword),
	}

	return a.accountDBAccessor.RegisterAccount(ctx, account)
}
