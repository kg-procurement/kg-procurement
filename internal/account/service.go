//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=account
package account

import (
	"context"
)

type accountDBAccessor interface {
	RegisterAccount(ctx context.Context, account Account) error
}

type AccountService struct {
	accountDBAccessor
}

func (a *AccountService) RegisterAccount(ctx context.Context, spec RegisterAccountSpec) error {
	return nil
}
