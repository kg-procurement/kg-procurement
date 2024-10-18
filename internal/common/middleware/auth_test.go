package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"kg/procurement/internal/token"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_newAuthMiddleware(t *testing.T) {
	_ = NewAuthMiddleware(nil)
}

func TestAuthMiddleware_MustAuthenticated(t *testing.T) {
	var (
		g              *gomega.WithT
		mockCtrl       *gomock.Controller
		mockTokenMgr   *token.MocktokenManager
		tokenSvc       *token.TokenService
		authMiddleware *AuthMiddleware
	)

	setup := func(t *testing.T) {
		g = gomega.NewWithT(t)
		mockCtrl = gomock.NewController(t)
		mockTokenMgr = token.NewMocktokenManager(mockCtrl)
		tokenSvc = &token.TokenService{
			TokenManager: mockTokenMgr,
		}
		authMiddleware = NewAuthMiddleware(tokenSvc)
	}

	teardown := func() {
		mockCtrl.Finish()
	}

	t.Run("NoAuthorizationHeaderReturnsError", func(t *testing.T) {
		setup(t)
		defer teardown()

		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler := authMiddleware.MustAuthenticated()
		handler(c)

		g.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
		g.Expect(w.Body.String()).To(gomega.ContainSubstring(ErrorAuthHeader))
	})
}
