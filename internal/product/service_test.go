package product

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewProductService(t *testing.T) {
	_ = NewProductService(nil)
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
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		mockProductAccessor.EXPECT().GetProductsByVendor(ctx, vendorID).
			Return(products, nil)

		res, err := svc.GetProductsByVendor(ctx, vendorID)
		g.Expect(res).Should(gomega.BeComparableTo(products))
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns err on accessor error", func(t *testing.T) {
		var (
			g                   = gomega.NewWithT(t)
			ctx                 = context.Background()
			mockCtrl            = gomock.NewController(t)
			mockProductAccessor = NewMockproductDBAccessor(mockCtrl)
		)

		svc := &ProductService{
			mockProductAccessor,
		}

		mockProductAccessor.EXPECT().GetProductsByVendor(ctx, vendorID).
			Return(nil, errors.New("error"))

		res, err := svc.GetProductsByVendor(ctx, vendorID)
		g.Expect(res).To(gomega.BeNil())
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func TestPriceService_UpdatePrice(t *testing.T) {
    t.Parallel()

    var (
        priceID    = "ID"
        updatedTime = time.Date(2024, time.September, 27, 12, 30, 0, 1, time.UTC)
        updateSpec = Price{
            ID:              priceID,
            PurchasingOrgID: "org_id_updated",
            VendorID:        "vendor_id_updated",
            Price:           99.99,
            QuantityMin:     10,
            QuantityMax:     99,
            ModifiedDate:    time.Time{},
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
