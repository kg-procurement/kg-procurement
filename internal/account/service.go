//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=account
package account

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/token"
	"net/mail"

	"github.com/benbjohnson/clock"
	"golang.org/x/crypto/bcrypt"
)

type accountDBAccessor interface {
	RegisterAccount(ctx context.Context, account Account) error
	FindAccountByEmail(ctx context.Context, email string) (*Account, error)
}

type tokenService interface {
	GenerateToken(spec token.ClaimSpec) (string, error)
}

type AccountService struct {
	accountDBAccessor
	tokenService
}

func (a *AccountService) RegisterAccount(ctx context.Context, spec AccountCredentialSpec) error {
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

func (a *AccountService) Login(ctx context.Context, spec AccountCredentialSpec) (string, error) {
	// Find the account by email
	account, err := a.accountDBAccessor.FindAccountByEmail(ctx, spec.Email)
	if err != nil {
		return "", fmt.Errorf("account not found: %w", err)
	}

	// Verify the password
	if err := account.VerifyPassword(spec.Password); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	// Generate a JWT token
	token, err := a.tokenService.GenerateToken(token.ClaimSpec{UserID: account.ID})
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}


func NewAccountService(
	conn database.DBConnector,
	clock clock.Clock,
	tokenSvc tokenService,
) *AccountService {
	return &AccountService{
		accountDBAccessor: newPostgresAccountAccessor(conn, clock),
		tokenService:      tokenSvc,
	}
}
