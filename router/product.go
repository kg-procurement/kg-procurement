package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/internal/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewProductEngine(
	r *gin.Engine,
	cfg config.ProductRoutes,
	productSvc *product.ProductService,
) {
	r.GET(cfg.GetProductsByVendor, func(ctx *gin.Context) {
		vendorID := ctx.Param("vendor_id")

		if vendorID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "vendor_id is required",
			})
			return
		}

		paginationSpec := GetPaginationSpec(ctx.Request)

		spec := product.GetProductsByVendorSpec{
			Name:           ctx.Query("name"),
			PaginationSpec: paginationSpec,
		}

		res, err := productSvc.GetProductsByVendor(ctx, vendorID, spec)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"products": res,
		})
	})

	r.PUT(cfg.UpdateProduct, func(ctx *gin.Context) {
		id := ctx.Param("id")

		spec := product.PutProductSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		newProduct := product.Product{
			ID:                product.ProductID(id),
			ProductCategoryID: spec.ProductCategoryID,
			UOMID:             spec.UOMID,
			IncomeTaxID:       spec.IncomeTaxID,
			ProductTypeID:     spec.ProductTypeID,
			Name:              spec.Name,
			Description:       spec.Description,
		}
		res, err := productSvc.UpdateProduct(ctx, newProduct)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdatePrice, func(ctx *gin.Context) {
		id := ctx.Param("id")
		spec := product.PutPriceSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		newPrice := product.Price{
			ID:              id,
			PurchasingOrgID: spec.PurchasingOrgID,
			VendorID:        spec.VendorID,
			ProductVendorID: spec.ProductVendorID,
			QuantityMin:     spec.QuantityMin,
			QuantityMax:     spec.QuantityMax,
			QuantityUOMID:   spec.QuantityUOMID,
			LeadTimeMin:     spec.LeadTimeMin,
			LeadTimeMax:     spec.LeadTimeMax,
			CurrencyID:      spec.CurrencyID,
			Price:           spec.Price,
			PriceQuantity:   spec.PriceQuantity,
			PriceUOMID:      spec.PriceUOMID,
			ValidFrom:       spec.ValidFrom,
			ValidTo:         spec.ValidTo,
			ValidPatternID:  spec.ValidPatternID,
			AreaGroupID:     spec.AreaGroupID,
			ReferenceNumber: spec.ReferenceNumber,
			ReferenceDate:   spec.ReferenceDate,
			DocumentTypeID:  spec.DocumentTypeID,
			DocumentID:      spec.DocumentID,
			ItemID:          spec.ItemID,
			TermOfPaymentID: spec.TermOfPaymentID,
			InvocationOrder: spec.InvocationOrder,
		}

		res, err := productSvc.UpdatePrice(ctx, newPrice)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, res)
	})
}
