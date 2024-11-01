package product

import (
	"context"
	"errors"
	"kg/procurement/internal/common/database"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewProductService(t *testing.T) {
	_ = NewProductService(nil, nil)
}

func TestProductService_GetProductsByVendor(t *testing.T) {
	t.Parallel()

	var (
		now      = time.Now()
		vendorID = "1234"
		products = []Product{
			{
				ID:           "1111",
				Name:         "Mixer",
				ModifiedDate: now,
			},
			{
				ID:           "2222",
				Name:         "Rice Cooker",
				ModifiedDate: now,
			},
		}
		paginationSpec = database.PaginationSpec{Limit: 10, Order: "DESC", Page: 1}
	)

	t.Run("success", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
			spec                = GetProductsByVendorSpec{PaginationSpec: paginationSpec}
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		expect := &AccessorGetProductsByVendorPaginationData{
			Products: products,
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 2,
			},
		}

		mockProductAccessor.EXPECT().GetProductsByVendor(ctx, vendorID, spec).
			Return(expect, nil)

		res, err := svc.GetProductsByVendor(ctx, vendorID, spec)
		g.Expect(res).Should(gomega.BeComparableTo(expect))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("success with order by and filter by name", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
			spec                = GetProductsByVendorSpec{
				Name:           "Rice Cooker",
				PaginationSpec: database.PaginationSpec{OrderBy: "name"},
			}
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		expect := &AccessorGetProductsByVendorPaginationData{
			Products: products[1:],
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 2,
			},
		}

		mockProductAccessor.EXPECT().GetProductsByVendor(ctx, vendorID, spec).
			Return(expect, nil)

		res, err := svc.GetProductsByVendor(ctx, vendorID, spec)
		g.Expect(res).Should(gomega.BeComparableTo(expect))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns err on accessor error", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
			spec                = GetProductsByVendorSpec{}
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		mockProductAccessor.EXPECT().GetProductsByVendor(ctx, vendorID, spec).
			Return(nil, errors.New("error"))

		res, err := svc.GetProductsByVendor(ctx, vendorID, spec)
		g.Expect(res).To(gomega.BeNil())
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func TestProductService_GetProductVendors(t *testing.T) {
	t.Parallel()

	var (
		productVendors = []GetProductVendorsDBResponse{
			{
				ID:                  "1",
				ProductID:           "1",
				Code:                "",
				Name:                "Buku",
				QuantityMin:         1,
				QuantityMax:         300,
				CurrencyName:        "Rupiah",
				CurrencyCode:        "IDR",
				Price:               23000,
				PriceQuantity:       1,
				VendorName:          "Multi Kharisma Solusindo, PT",
				VendorRating:        -100,
				IncomeTaxID:         "0",
				IncomeTaxName:       "",
				IncomeTaxPercentage: "0",
				Description:         "Buku",
				UOMID:               "26",
				SAPCode:             "",
				ModifiedDate:        time.Date(2020, 11, 11, 13, 22, 16, 0, time.UTC),
				ModifiedBy:          "151",
			},
			{
				ID:                  "2",
				ProductID:           "2",
				Code:                "",
				Name:                "Koran",
				QuantityMin:         1,
				QuantityMax:         4,
				CurrencyName:        "Rupiah",
				CurrencyCode:        "IDR",
				Price:               290000,
				PriceQuantity:       1,
				VendorName:          "Toko Amazon",
				VendorRating:        0,
				IncomeTaxID:         "0",
				IncomeTaxName:       "",
				IncomeTaxPercentage: "0",
				Description:         "Koran",
				UOMID:               "26",
				SAPCode:             "",
				ModifiedDate:        time.Date(2020, 11, 2, 14, 49, 6, 0, time.UTC),
				ModifiedBy:          "0",
			},
		}
		paginationSpec = database.PaginationSpec{Limit: 10, Order: "DESC", Page: 1}
	)

	t.Run("success", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
			spec                = GetProductVendorsSpec{PaginationSpec: paginationSpec}
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		expect := &AccessorGetProductVendorsPaginationData{
			ProductVendors: productVendors,
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 2,
			},
		}

		mockProductAccessor.EXPECT().GetAllProductVendors(ctx, spec).
			Return(expect, nil)

		res, err := svc.GetProductVendors(ctx, spec)
		g.Expect(res).Should(gomega.BeComparableTo(expect))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("success with order by name", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
			spec                = GetProductVendorsSpec{
				PaginationSpec: database.PaginationSpec{OrderBy: "name"},
			}
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		expect := &AccessorGetProductVendorsPaginationData{
			ProductVendors: productVendors,
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 2,
			},
		}

		mockProductAccessor.EXPECT().GetAllProductVendors(ctx, spec).
			Return(expect, nil)

		res, err := svc.GetProductVendors(ctx, spec)
		g.Expect(res).Should(gomega.BeComparableTo(expect))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns err on accessor error", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
			spec                = GetProductVendorsSpec{}
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		mockProductAccessor.EXPECT().GetAllProductVendors(ctx, spec).
			Return(nil, errors.New("error"))

		res, err := svc.GetProductVendors(ctx, spec)
		g.Expect(res).To(gomega.BeNil())
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	t.Parallel()
	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)
	var (
		updatedTime = time.Now()
		updateSpec  = Product{
			ID:                "some_id",
			ProductCategoryID: "category_id_updated",
			UOMID:             "uom_id_updated",
			IncomeTaxID:       "income_tax_id_updated",
			ProductTypeID:     "product_type_id_updated",
			Name:              "Updated Product",
			Description:       "Updated description",
			ModifiedDate:      fixedTime,
			ModifiedBy:        "modified_by_updated",
		}
		updatedProduct = Product{
			ID:                "some_id",
			ProductCategoryID: "category_id_updated",
			UOMID:             "uom_id_updated",
			IncomeTaxID:       "income_tax_id_updated",
			ProductTypeID:     "product_type_id_updated",
			Name:              "Updated Product",
			Description:       "Updated description",
			ModifiedDate:      updatedTime,
			ModifiedBy:        "modified_by_updated",
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
		)

		svc := &ProductService{
			productDBAccessor: mockProductAccessor,
		}

		mockProductAccessor.EXPECT().UpdateProduct(ctx, updateSpec).Return(updatedProduct, nil)

		res, err := svc.UpdateProduct(ctx, updateSpec)
		g.Expect(res).Should(gomega.BeComparableTo(updatedProduct))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns error on accessor failure", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
		)

		svc := &ProductService{
			productDBAccessor: mockProductAccessor,
		}

		mockProductAccessor.EXPECT().UpdateProduct(ctx, updateSpec).Return(Product{}, errors.New("update error"))

		res, err := svc.UpdateProduct(ctx, updateSpec)
		g.Expect(res).To(gomega.Equal(Product{}))
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func TestPriceService_UpdatePrice(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 1, time.UTC)
	var (
		priceID     = "ID"
		updatedTime = time.Now()
		updateSpec  = Price{
			ID:              priceID,
			PurchasingOrgID: "org_id_updated",
			VendorID:        "vendor_id_updated",
			Price:           99.99,
			QuantityMin:     10,
			QuantityMax:     99,
			ModifiedDate:    fixedTime,
			ModifiedBy:      "modified_by_updated",
		}
		updatedPrice = Price{
			ID:              priceID,
			PurchasingOrgID: "org_id_updated",
			VendorID:        "vendor_id_updated",
			Price:           99.99,
			QuantityMin:     10,
			QuantityMax:     100,
			ModifiedDate:    updatedTime,
			ModifiedBy:      "modified_by_updated",
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
		)

		svc := &ProductService{
			productDBAccessor: mockProductAccessor,
		}

		mockProductAccessor.EXPECT().UpdatePrice(ctx, updateSpec).Return(updatedPrice, nil)

		res, err := svc.UpdatePrice(ctx, updateSpec)
		g.Expect(res).Should(gomega.BeComparableTo(updatedPrice))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns error on accessor failure", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
		)

		svc := &ProductService{
			productDBAccessor: mockProductAccessor,
		}

		mockProductAccessor.EXPECT().UpdatePrice(ctx, updateSpec).Return(Price{}, errors.New("update error"))

		res, err := svc.UpdatePrice(ctx, updateSpec)
		g.Expect(res).To(gomega.Equal(Price{}))
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}
