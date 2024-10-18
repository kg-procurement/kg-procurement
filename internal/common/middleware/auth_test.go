package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	t.Run("InvalidAuthorizationHeaderReturnsError", func(t *testing.T) {
		setup(t)
		defer teardown()

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(AuthorizationHeader, "InvalidHeader")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler := authMiddleware.MustAuthenticated()
		handler(c)

		g.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
		g.Expect(w.Body.String()).To(gomega.ContainSubstring(ErrorAuthInvalid))
	})

	t.Run("InvalidAuthorizationTypeReturnsError", func(t *testing.T) {
		setup(t)
		defer teardown()

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(AuthorizationHeader, "Basic token123")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler := authMiddleware.MustAuthenticated()
		handler(c)

		g.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
		g.Expect(w.Body.String()).To(gomega.ContainSubstring(ErrorAuthType))
	})

	t.Run("TokenValidationErrorReturnsError", func(t *testing.T) {
		setup(t)
		defer teardown()

		mockTokenMgr.EXPECT().
			ValidateToken("invalidtoken").
			Return(nil, errors.New("token validation error"))

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(AuthorizationHeader, "Bearer invalidtoken")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler := authMiddleware.MustAuthenticated()
		handler(c)

		g.Expect(w.Code).To(gomega.Equal(http.StatusInternalServerError))
		g.Expect(w.Body.String()).To(gomega.ContainSubstring("token validation error"))
	})

	t.Run("ValidTokenReturnsStatusOKAndSetsAuthPayload", func(t *testing.T) {
		setup(t)
		defer teardown()

		expectedClaims := &token.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "user123",
			},
		}

		// mock
		mockTokenMgr.EXPECT().ValidateToken("validtoken").Return(expectedClaims, nil)

		// create request with valid Bearer token
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(AuthorizationHeader, "Bearer validtoken")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// call the middleware
		handler := authMiddleware.MustAuthenticated()
		handler(c)

		g.Expect(w.Code).To(gomega.Equal(http.StatusOK))

		// Check if "auth_payload" is set in context
		value, exists := c.Get("auth_payload")
		g.Expect(exists).To(gomega.BeTrue())
		g.Expect(value).To(gomega.Equal(token.ClaimSpec{
			UserID: "user123",
		}))
	})
}
