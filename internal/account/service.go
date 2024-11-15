//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=account
package account

import (
	"context"
	"errors"
	"fmt"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/token"
	"net/mail"

	"github.com/benbjohnson/clock"
	"golang.org/x/crypto/bcrypt"
)

var ErrLoginFailed = errors.New("login failed")

type accountDBAccessor interface {
	RegisterAccount(ctx context.Context, account Account) error
	FindAccountByEmail(ctx context.Context, email string) (*Account, error)
	FindAccountByID(ctx context.Context, id string) (*Account, error)
}

type tokenService interface {
	GenerateToken(spec token.ClaimSpec) (string, error)
	ValidateToken(tokenString string) (*token.Claims, error)
}

type AccountService struct {
	accountDBAccessor
	tokenService
}

func (a *AccountService) RegisterAccount(ctx context.Context, spec RegisterContract) error {
	// Validate email
	if _, err := mail.ParseAddress(spec.Email); err != nil {
		utils.Logger.Errorf("invalid email: %v", err)
		return fmt.Errorf("invalid email: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(spec.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Errorf("failed to hash password: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate ID
	id, err := generateRandomID()
	if err != nil {
		utils.Logger.Errorf("failed to generate random ID: %v", err)
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

func (a *AccountService) Login(ctx context.Context, spec LoginContract) (string, error) {

	// Find the account by email
	account, err := a.accountDBAccessor.FindAccountByEmail(ctx, spec.Email)
	if err != nil {
		utils.Logger.Errorf("account not found: %s", spec.Email)
		return "", ErrLoginFailed
	}

	// Verify the password
	if err := account.VerifyPassword(spec.Password); err != nil {
		utils.Logger.Errorf("invalid password for email: %s", spec.Email)
		return "", ErrLoginFailed
	}

	// Generate a JWT token
	token, err := a.tokenService.GenerateToken(token.ClaimSpec{UserID: account.ID})
	if err != nil {
		utils.Logger.Errorf("failed to generate token for email: %s, error: %v", spec.Email, err)
		return "", ErrLoginFailed
	}

	return token, nil
}

func (a *AccountService) GetCurrentUser(ctx context.Context, tokenString string) (*Account, error) {
	// Parse and validate the JWT token
	claims, err := a.tokenService.ValidateToken(tokenString)
	if err != nil {
		utils.Logger.Errorf("failed to parse token: %v", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Find the account associated with the user ID in the token claims
	userID := claims.Subject
	account, err := a.accountDBAccessor.FindAccountByID(ctx, userID)
	if err != nil {
		utils.Logger.Errorf("account not found for user ID: %v", userID)
		return nil, fmt.Errorf("account not found: %w", err)
	}

	return account, nil
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
