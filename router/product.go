package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
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
		utils.Logger.Info("Received getProductsByVendor request")

		vendorID := ctx.Param("vendor_id")

		if vendorID == "" {
			utils.Logger.Error("vendor_id is required")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "vendor_id is required",
			})
			return
		}

		paginationSpec := GetPaginationSpec(ctx.Request)

		spec := product.GetProductVendorByVendorSpec{
			Name:           ctx.Query("name"),
			PaginationSpec: paginationSpec,
		}

		res, err := productSvc.GetProductVendorsByVendor(ctx, vendorID, spec)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.Logger.Info("Completed getproductsByVendor request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetProductVendors, func(ctx *gin.Context) {
		paginationSpec := GetPaginationSpec(ctx.Request)

		spec := product.GetProductVendorsSpec{
			PaginationSpec: paginationSpec,
		}

		res, err := productSvc.GetProductVendors(ctx, spec)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdateProduct, func(ctx *gin.Context) {
		utils.Logger.Info("Received updateProductDetail request")

		id := ctx.Param("id")

		spec := product.PutProductSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
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
			return
		}

		utils.Logger.Info("Completed updateProductDetail request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdatePrice, func(ctx *gin.Context) {
		utils.Logger.Info("Received updateProductPrice request")

		id := ctx.Param("id")
		spec := product.PutPriceSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			utils.Logger.Error(err.Error())
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
			return
		}

		utils.Logger.Info("Completed updateProductPrice request process")

		ctx.JSON(http.StatusOK, res)
	})
}
