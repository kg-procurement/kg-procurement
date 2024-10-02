package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/internal/vendors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewVendorEngine(
	r *gin.Engine,
	cfg config.VendorRoutes,
	vendorSvc *vendors.VendorService,
) {
	r.GET(cfg.GetAll, func(ctx *gin.Context) {

		paginationSpec := GetPaginationSpec(ctx.Request)
		spec := vendors.GetAllVendorSpec{
			Location:       ctx.Query("location"),
			Product:        ctx.Query("product"),
			PaginationSpec: paginationSpec,
		}

		res, err := vendorSvc.GetAll(ctx, spec)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetByLocation, func(ctx *gin.Context) {
		location := ctx.Query("location")

		res, err := vendorSvc.GetByLocation(ctx, location)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetByProduct, func(ctx *gin.Context) {
		product := ctx.Query("product")

		res, err := vendorSvc.GetByProduct(ctx, product)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, res)
	})
}
