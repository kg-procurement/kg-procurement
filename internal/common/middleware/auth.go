package middleware

import (
	"github.com/gin-gonic/gin"
	"kg/procurement/internal/token"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	BearerType          = "Bearer"
	ErrorAuthHeader     = "authorization header not provided"
	ErrorAuthInvalid    = "authorization header not valid"
	ErrorAuthType       = "authorization type not valid"
)

type AuthMiddleware struct {
	tokenManager token.TokenManager
}

func NewAuthMiddleware(manager token.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: manager,
	}
}

func (m *AuthMiddleware) MustAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check authorization header
		authorizationHeader := ctx.GetHeader(AuthorizationHeader)
		if len(authorizationHeader) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": ErrorAuthHeader,
			})
			return
		}

		// check validity of authorization header
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": ErrorAuthInvalid,
			})
			return
		}

		// check authorization type
		authorizationType := fields[0]
		if authorizationType != BearerType {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": ErrorAuthType,
			})
			return
		}

		// claim claims
		tokenStr := fields[1]
		claims, err := m.tokenManager.ValidateToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Set("auth_payload", token.ClaimSpec{
			UserID: claims.Subject,
		})

		ctx.Next()
	}
}
