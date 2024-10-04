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
		vendorID = "1234"
		products = []Product{
			{
				ID:           "1111",
				Name:         "Mixer",
				ModifiedDate: time.Now(),
			},
			{
				ID:           "2222",
				Name:         "Rice Cooker",
				ModifiedDate: time.Now(),
			},
		}
	)

	t.Run("success", func(t *testing.T) {
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
			Return(products, nil)

		res, err := svc.GetProductsByVendor(ctx, vendorID, spec)
		g.Expect(res).Should(gomega.BeComparableTo(products))
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

		mockProductAccessor.EXPECT().GetProductsByVendor(ctx, vendorID, spec).
			Return(products[1:], nil)

		res, err := svc.GetProductsByVendor(ctx, vendorID, spec)
		g.Expect(res).Should(gomega.BeComparableTo(products[1:]))
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

func TestProductService_UpdateProduct(t *testing.T) {
    t.Parallel()
	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)
    var (
        updatedTime = time.Now()
        updateSpec  = Product{
            ID:                "some_id",
            ProductCategoryID: "category_id_updated",
            UOMID:            "uom_id_updated",
            IncomeTaxID:      "income_tax_id_updated",
            ProductTypeID:    "product_type_id_updated",
            Name:             "Updated Product",
            Description:      "Updated description",
            ModifiedDate:     fixedTime,
            ModifiedBy:       "modified_by_updated",
        }
        updatedProduct = &Product{
            ID:                "some_id",
            ProductCategoryID: "category_id_updated",
            UOMID:            "uom_id_updated",
            IncomeTaxID:      "income_tax_id_updated",
            ProductTypeID:    "product_type_id_updated",
            Name:             "Updated Product",
            Description:      "Updated description",
            ModifiedDate:     updatedTime,
            ModifiedBy:       "modified_by_updated",
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

        mockProductAccessor.EXPECT().UpdateProduct(ctx, updateSpec).Return(nil, errors.New("update error"))

        res, err := svc.UpdateProduct(ctx, updateSpec)
        g.Expect(res).To(gomega.BeNil())
        g.Expect(err).ShouldNot(gomega.BeNil())
    })
}


func TestPriceService_UpdatePrice(t *testing.T) {
    t.Parallel()

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 1, time.UTC)
    var (
        priceID    = "ID"
		updatedTime = time.Now()
        updateSpec = Price{
            ID:              priceID,
            PurchasingOrgID: "org_id_updated",
            VendorID:        "vendor_id_updated",
            Price:           99.99,
            QuantityMin:     10,
            QuantityMax:     99,
            ModifiedDate:    fixedTime,
            ModifiedBy:      "modified_by_updated",
        }
        updatedPrice = &Price{
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
            mockProductAccessor   = NewMockproductDBAccessor(mockCtrl)
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
            mockProductAccessor   = NewMockproductDBAccessor(mockCtrl)
        )

        svc := &ProductService{
            productDBAccessor: mockProductAccessor,
        }

        mockProductAccessor.EXPECT().UpdatePrice(ctx, updateSpec).Return(nil, errors.New("update error"))

        res, err := svc.UpdatePrice(ctx, updateSpec)
        g.Expect(res).To(gomega.BeNil())
        g.Expect(err).ShouldNot(gomega.BeNil())
    })
}
