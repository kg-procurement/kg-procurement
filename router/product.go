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

		ctx.JSON(http.StatusOK, res)
	})
}
