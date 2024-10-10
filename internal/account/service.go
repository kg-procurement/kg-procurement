//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=account
package account

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"net/mail"

	"github.com/benbjohnson/clock"
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

	// Generate ID
	id, err := generateRandomID()
	if err != nil {
		return fmt.Errorf("failed to generate random ID: %w", err)
	}

	// Create account
	account := Account{
		ID:       id,
		Email:    spec.Email,
		Password: string(hashedPassword),
	}

	return a.accountDBAccessor.RegisterAccount(ctx, account)
}

func NewAccountService(
	conn database.DBConnector,
	clock clock.Clock,
) *AccountService {
	return &AccountService{
		accountDBAccessor: newPostgresAccountAccessor(conn, clock),
	}
}
