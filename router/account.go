package router

import (
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/account"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewAccountEngine(
	r *gin.Engine,
	cfg config.AccountRoutes,
	accountSvc *account.AccountService,
) {
	r.POST(cfg.Register, func(ctx *gin.Context) {
		utils.Logger.Info("Received accountRegister request")

		payload := account.RegisterContract{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		err := accountSvc.RegisterAccount(ctx, payload)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Account registered successfully",
		})
	})

	r.POST(cfg.Login, func(ctx *gin.Context) {
		utils.Logger.Info("Received accountLogin request")

		payload := account.LoginContract{}
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			utils.Logger.Error(err.Error())

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		token, err := accountSvc.Login(ctx, payload)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		utils.Logger.Info("Completed accountLogin request process")

		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})
}
