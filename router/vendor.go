package router

import (
	"encoding/json"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/mailer"
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
		utils.Logger.Info("Received getAllVendor request")

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

		utils.Logger.Info("Completed getAllVendor request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdateDetail, func(ctx *gin.Context) {
		utils.Logger.Info("Received updateVendorDetail request")

		id := ctx.Param("id")

		spec := &vendors.PutVendorSpec{}
		if err := ctx.ShouldBindJSON(&spec); err != nil {
			utils.Logger.Error(err.Error())
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

		utils.Logger.Info("Completed updateDetailVendor request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetById, func(ctx *gin.Context) {
		utils.Logger.Info("Received getVendorById request")

		id := ctx.Param("id")

		res, err := vendorSvc.GetById(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.Logger.Info("Completed getVendorById request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.GET(cfg.GetLocations, func(ctx *gin.Context) {
		utils.Logger.Info("Received getVendorByLocations request")

		res, err := vendorSvc.GetLocations(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.Logger.Info("Completed getVendorByLocations request process")

		ctx.JSON(http.StatusOK, gin.H{
			"locations": res,
		})
	})

	r.POST(cfg.EmailBlast, func(ctx *gin.Context) {
		utils.Logger.Info("Received emailBlast request")

		// parse vendor ids
		var vendorIDs []string
		if err := json.Unmarshal([]byte(ctx.PostForm("vendor_ids")), &vendorIDs); err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid vendor_ids JSON",
			})
			return
		}

		maxMemory := int64(16 << 20) // 16 MB
		if err := ctx.Request.ParseMultipartForm(maxMemory); err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse multipart form",
			})
			return
		}

		payload := vendors.EmailBlastContract{}
		if err := ctx.ShouldBind(&payload); err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		form := ctx.Request.MultipartForm
		files := form.File["attachments"]
		attachments, err := mailer.BulkFromMultipartForm(files)
		if err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		email := mailer.Email{
			Subject:     payload.Subject,
			Body:        payload.Body,
			Attachments: attachments,
		}

		errList, err := vendorSvc.BlastEmail(ctx, vendorIDs, email)
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

		utils.Logger.Info("Completed emailBlast request process")

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Emails successfully sent",
		})
	})

	r.POST(cfg.AutomatedEmailBlast, func(ctx *gin.Context) {
		productName := ctx.Param("product_name")
		if productName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Product name is required",
			})
			return
		}

		errList, err := vendorSvc.AutomatedEmailBlast(ctx, productName)
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
}
