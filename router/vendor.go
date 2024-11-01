package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/internal/vendors"
	"net/http"

	u "kg/procurement/cmd/utils"

	"github.com/gin-gonic/gin"
)

func NewVendorEngine(
	r *gin.Engine,
	cfg config.VendorRoutes,
	vendorSvc *vendors.VendorService,
) {
	r.GET(cfg.GetAll, func(ctx *gin.Context) {
		u.GeneralLogger.Println("Received getAllVendor request")

		paginationSpec := GetPaginationSpec(ctx.Request)
		spec := vendors.GetAllVendorSpec{
			Location:       ctx.Query("location"),
			Product:        ctx.Query("product"),
			PaginationSpec: paginationSpec,
		}

		res, err := vendorSvc.GetAll(ctx, spec)
		if err != nil {
			u.ErrorLogger.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		u.GeneralLogger.Println("Completed getAllVendor request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdateDetail, func(ctx *gin.Context) {
		u.GeneralLogger.Println("Received updateVendorDetail request")

		id := ctx.Param("id")

		spec := &vendors.PutVendorSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			u.ErrorLogger.Println(err.Error())
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
			u.ErrorLogger.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		u.GeneralLogger.Println("Completed updateDetailVendor request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetById, func(ctx *gin.Context) {
		u.GeneralLogger.Println("Received getVendorById request")

		id := ctx.Param("id")

		res, err := vendorSvc.GetById(ctx, id)
		if err != nil {
			u.ErrorLogger.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		u.GeneralLogger.Println("Completed getVendorById request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetLocations, func(ctx *gin.Context) {
		u.GeneralLogger.Println("Received getVendorByLocations request")

		res, err := vendorSvc.GetLocations(ctx)
		if err != nil {
			u.ErrorLogger.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		u.GeneralLogger.Println("Completed getVendorByLocations request process")

		ctx.JSON(http.StatusOK, gin.H{
			"locations": res,
		})
	})

	r.POST(cfg.EmailBlast, func(ctx *gin.Context) {
		u.GeneralLogger.Println("Received emailBlast request")

		payload := vendors.EmailBlastContract{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			u.ErrorLogger.Println(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		errList, err := vendorSvc.BlastEmail(ctx, payload.VendorIDs, payload.EmailTemplate)
		if err != nil {
			if len(errList) > 0 {
				u.ErrorLogger.Println(err.Error())
				ctx.JSON(http.StatusMultiStatus, gin.H{
					"error": errList,
				})
				return
			}
			u.ErrorLogger.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		u.GeneralLogger.Println("Completed emailBlast request process")

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Emails successfully sent",
		})
	})
}
