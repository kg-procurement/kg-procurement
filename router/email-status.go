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
}
