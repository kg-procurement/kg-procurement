//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=mailer
package mailer

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

type emailStatusDBAccessor interface {
	WriteEmailStatus(ctx context.Context, payload EmailStatus) error
	GetAll(ctx context.Context, spec GetAllEmailStatusSpec) (*AccessorGetAllPaginationData, error)
	UpdateEmailStatus(ctx context.Context, payload EmailStatus) (*EmailStatus, error)
}

type EmailStatusService struct {
	emailStatusDBAccessor
}

func (p *EmailStatusService) WriteEmailStatus(ctx context.Context, payload EmailStatus) error {
	return p.emailStatusDBAccessor.WriteEmailStatus(ctx, payload)
}

func (p *EmailStatusService) GetAllEmailStatus(ctx context.Context, spec GetAllEmailStatusSpec) (*AccessorGetAllPaginationData, error) {
	return p.emailStatusDBAccessor.GetAll(ctx, spec)
}

func (p *EmailStatusService) UpdateEmailStatus(ctx context.Context, payload EmailStatus) (*EmailStatus, error) {
	return p.emailStatusDBAccessor.UpdateEmailStatus(ctx, payload)
}

func NewEmailStatusService(
	conn database.DBConnector,
	clock clock.Clock,
) *EmailStatusService {
	return &EmailStatusService{
		emailStatusDBAccessor: newPostgresEmailStatusAccessor(conn, clock),
	}
}
