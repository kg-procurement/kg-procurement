package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/mailer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewEmailStatusEngine(
	r *gin.Engine,
	cfg config.EmailStatusRoutes,
	emailStatusSvc *mailer.EmailStatusService,
) {
	r.GET(cfg.GetAll, func(ctx *gin.Context) {
		utils.Logger.Info("Received getAllEmailStatus request")

		paginationSpec := GetPaginationSpec(ctx.Request)
		spec := mailer.GetAllEmailStatusSpec{
			PaginationSpec: paginationSpec,
		}

		res, err := emailStatusSvc.GetAll(ctx, spec)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.Logger.Info("Completed getAllEmailStatus request process")

		ctx.JSON(http.StatusOK, res)
	})

	r.PUT(cfg.UpdateEmailStatus, func(ctx *gin.Context) {
		utils.Logger.Info("Received updateEmailStatus request")

		id := ctx.Param("id")
		var payload mailer.EmailStatus
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if payload.Status == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "status field cannot be empty",
			})
			return
		}
		if _, err := mailer.ParseEmailStatusEnum(payload.Status); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		payload.ID = id

		updatedStatus, err := emailStatusSvc.UpdateEmailStatus(ctx, payload)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.Logger.Info("Completed updateEmailStatus request process")

		ctx.JSON(http.StatusOK, updatedStatus)
	})
}
