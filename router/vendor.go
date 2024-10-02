package router

import (
	"encoding/json"
	"io"
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

		spec := GetPaginationSpec(ctx.Request)

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

	r.PUT(cfg.Put, func(ctx *gin.Context) {
		id := ctx.Query("id")

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		spec := PutVendorSpec{}
		if err = json.Unmarshal(body, &spec); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		newVendor := vendors.Vendor{
			ID:            id,
			Name:          spec.Name,
			Description:   spec.Description,
			BpID:          spec.BpID,
			BpName:        spec.BpName,
			Rating:        spec.Rating,
			AreaGroupID:   spec.AreaGroupID,
			AreaGroupName: spec.AreaGroupName,
			SapCode:       spec.SapCode,
		}

		res, err := vendorSvc.Put(ctx, newVendor)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, res)

	})
}
