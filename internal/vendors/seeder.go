//go:generate mockgen -typed -source=seeder.go -destination=seeder_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

type seederDataWriter interface {
	writeVendor(ctx context.Context, vendor Vendor) error
	Close() error
}

type Seeder struct {
	seederDataWriter
}

func (s *Seeder) SetupVendors(ctx context.Context, vendors []Vendor) error {
	return nil
}

func (s *Seeder) Close() error {
	return nil
}

func NewSeeder(
	seederDataWriter seederDataWriter,
) *Seeder {
	return nil
}

func NewDBSeederWriter(
	dbClient database.DBConnector,
	clock clock.Clock,
) seederDataWriter {
	return nil
}
