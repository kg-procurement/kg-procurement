package product

import (
	"context"
	"errors"
	"testing"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewSeeder(t *testing.T) {
	_ = NewSeeder(nil)
}

func Test_NewDBSeederWriter(t *testing.T) {
	_ = NewDBSeederWriter(nil, nil)
}

func Test_Seeder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var (
		mockWriter *MockseederDataWriter
		subject    *Seeder
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)

		mockWriter = NewMockseederDataWriter(ctrl)
		subject = NewSeeder(mockWriter)

		return gomega.NewWithT(t)
	}

	t.Run("setup product", func(t *testing.T) {
		product := []Product{{ID: "123"}, {ID: "321"}}

		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProduct(ctx, Product{ID: "123"})
			mockWriter.EXPECT().writeProduct(ctx, Product{ID: "321"})

			err := subject.SetupProducts(ctx, product)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProduct(ctx, Product{ID: "123"})
			mockWriter.EXPECT().writeProduct(ctx, Product{ID: "321"}).Return(errors.New("error"))

			err := subject.SetupProducts(ctx, product)
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	t.Run("setup product category", func(t *testing.T) {
		category := []ProductCategory{{ID: "123"}, {ID: "321"}}

		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProductCategory(ctx, ProductCategory{ID: "123"})
			mockWriter.EXPECT().writeProductCategory(ctx, ProductCategory{ID: "321"})

			err := subject.SetupProductCategory(ctx, category)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProductCategory(ctx, ProductCategory{ID: "123"})
			mockWriter.EXPECT().writeProductCategory(ctx, ProductCategory{ID: "321"}).Return(errors.New("error"))

			err := subject.SetupProductCategory(ctx, category)
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	t.Run("setup product type", func(t *testing.T) {
		pType := []ProductType{{ID: "123"}, {ID: "321"}}

		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProductType(ctx, ProductType{ID: "123"})
			mockWriter.EXPECT().writeProductType(ctx, ProductType{ID: "321"})

			err := subject.SetupProductType(ctx, pType)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProductType(ctx, ProductType{ID: "123"})
			mockWriter.EXPECT().writeProductType(ctx, ProductType{ID: "321"}).Return(errors.New("error"))

			err := subject.SetupProductType(ctx, pType)
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	t.Run("setup uom", func(t *testing.T) {
		uom := []UOM{{ID: "123"}, {ID: "321"}}

		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeUOM(ctx, UOM{ID: "123"})
			mockWriter.EXPECT().writeUOM(ctx, UOM{ID: "321"})

			err := subject.SetupUOM(ctx, uom)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeUOM(ctx, UOM{ID: "123"})
			mockWriter.EXPECT().writeUOM(ctx, UOM{ID: "321"}).Return(errors.New("error"))

			err := subject.SetupUOM(ctx, uom)
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	t.Run("setup product-vendor", func(t *testing.T) {
		pv := []ProductVendor{{ProductID: "123", VendorID: "321"}, {ProductID: "321", VendorID: "123"}}

		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProductVendor(ctx, ProductVendor{ProductID: "123", VendorID: "321"})
			mockWriter.EXPECT().writeProductVendor(ctx, ProductVendor{ProductID: "321", VendorID: "123"})

			err := subject.SetupProductVendor(ctx, pv)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeProductVendor(ctx, ProductVendor{ProductID: "123", VendorID: "321"})
			mockWriter.EXPECT().writeProductVendor(ctx, ProductVendor{ProductID: "321", VendorID: "123"}).Return(errors.New("error"))

			err := subject.SetupProductVendor(ctx, pv)
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})

	t.Run("close", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().Close()

			err := subject.Close()
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().Close().Return(errors.New("error"))

			err := subject.Close()
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})
}
