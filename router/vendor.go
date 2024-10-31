package router

import (
	"fmt"
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
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdateDetail, func(ctx *gin.Context) {
		id := ctx.Param("id")

		spec := &vendors.PutVendorSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
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

		res, err := vendorSvc.UpdateDetail(ctx, newVendor)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetById, func(ctx *gin.Context) {
		id := ctx.Param("id")

		res, err := vendorSvc.GetById(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetLocations, func(ctx *gin.Context) {
		res, err := vendorSvc.GetLocations(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"locations": res,
		})
	})

	r.POST(cfg.EmailBlast, func(ctx *gin.Context) {
		payload := vendors.EmailBlastContract{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		errList, err := vendorSvc.BlastEmail(ctx, payload.VendorIDs, payload.EmailTemplate)
		if err != nil {
			if len(errList) > 0 {
				ctx.JSON(http.StatusMultiStatus, gin.H{
					"error": errList,
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Emails successfully sent",
		})
	})
	r.POST(cfg.AutomatedEmailBlast, func(ctx *gin.Context) {
		product_name := ctx.Param("product_name")
		if product_name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Product name is required",
			})
			return
		}
		vendorIDs := []string{}
		vendorIDs, err := vendorSvc.AutomatedEmailBlast(ctx, product_name)
		fmt.Println("vendorIds: ", vendorIDs)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"vendorIDs": vendorIDs,
		})
	})
}
