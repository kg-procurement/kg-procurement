//go:generate mockgen -typed -source=seeder.go -destination=seeder_mock.go -package=product
package product

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

type seederDataWriter interface {
	writeProduct(ctx context.Context, product Product) error
	writeProductCategory(ctx context.Context, productCategory ProductCategory) error
	writeProductType(ctx context.Context, productType ProductType) error
	writeUOM(ctx context.Context, uom UOM) error
	Close() error
}

type Seeder struct {
	seederDataWriter
}

func (s *Seeder) SetupProducts(ctx context.Context, products []Product) error {
	for _, product := range products {
		if err := s.seederDataWriter.writeProduct(ctx, product); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) SetupProductCategory(ctx context.Context, productCategory []ProductCategory) error {
	for _, category := range productCategory {
		if err := s.seederDataWriter.writeProductCategory(ctx, category); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) SetupProductType(ctx context.Context, productType []ProductType) error {
	for _, pType := range productType {
		if err := s.seederDataWriter.writeProductType(ctx, pType); err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) SetupUOM(ctx context.Context, uoms []UOM) error {
	for _, uoms := range uoms {
		if err := s.seederDataWriter.writeUOM(ctx, uoms); err != nil {
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
	return newPostgresProductAccessor(dbClient, clock)
}
