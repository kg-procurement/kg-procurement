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

	r.GET(cfg.GetCurrentUser, func(ctx *gin.Context) {
		utils.Logger.Info("Received getCurrentUser request")

		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			utils.Logger.Error("Authorization header is required")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}
	
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
	
		account, err := accountSvc.GetCurrentUser(ctx, tokenString)
		if err != nil {
			utils.Logger.Error(err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}
	
		ctx.JSON(http.StatusOK, gin.H{
			"id":        account.ID,
			"email":     account.Email,
			"createdAt": account.CreatedAt,
			"modifiedAt": account.ModifiedDate,
		})
	})
}
