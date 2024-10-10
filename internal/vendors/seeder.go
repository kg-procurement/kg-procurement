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
	for _, vendor := range vendors {
		if err := s.seederDataWriter.writeVendor(ctx, vendor); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) Close() error {
	return s.seederDataWriter.Close()
}

func NewSeeder(
	seederDataWriter seederDataWriter,
) *Seeder {
	return &Seeder{seederDataWriter}
}

func NewDBSeederWriter(
	dbClient database.DBConnector,
	clock clock.Clock,
) seederDataWriter {
	return newPostgresVendorAccessor(dbClient, clock)
}
