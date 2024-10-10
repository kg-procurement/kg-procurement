package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"kg/procurement/cmd/config"
	"testing"
)

func Test_NewTokenService(t *testing.T) {
	_ = NewTokenService(config.Token{}, nil)
}

func TestTokenService_GenerateToken(t *testing.T) {
	t.Parallel()

	t.Run("GenerateTokenWithValidClaimsShouldReturnExpectedToken", func(t *testing.T) {
		var (
			g             = gomega.NewWithT(t)
			mockCtrl      = gomock.NewController(t)
			mockTokenMgr  = NewMocktokenManager(mockCtrl)
			claimSpec     = ClaimSpec{UserID: "123"}
			expectedToken = "valid_token_string"
		)

		// mock expectation
		mockTokenMgr.EXPECT().GenerateToken(claimSpec).Return(expectedToken, nil)

		svc := &TokenService{
			tokenManager: mockTokenMgr,
		}
		tokenString, err := svc.GenerateToken(claimSpec)

		// assertions
		g.Expect(err).To(gomega.BeNil())
		g.Expect(tokenString).To(gomega.Equal(expectedToken))
	})

	t.Run("FailedToGenerateTokenShouldReturnExpectedError", func(t *testing.T) {
		var (
			g            = gomega.NewWithT(t)
			mockCtrl     = gomock.NewController(t)
			mockTokenMgr = NewMocktokenManager(mockCtrl)
			claimSpec    = ClaimSpec{UserID: "123"}
			expectedErr  = errors.New("failed to generate token")
		)

		// mock expectation
		mockTokenMgr.EXPECT().GenerateToken(claimSpec).Return("", expectedErr)

		svc := &TokenService{
			tokenManager: mockTokenMgr,
		}
		tokenString, err := svc.GenerateToken(claimSpec)

		// assertions
		g.Expect(err).To(gomega.Equal(expectedErr))
		g.Expect(tokenString).To(gomega.BeEmpty())
	})

	t.Run("GenerateTokenWithInvalidClaimsShouldReturnExpectedToken", func(t *testing.T) {
		var (
			g            = gomega.NewWithT(t)
			mockCtrl     = gomock.NewController(t)
			mockTokenMgr = NewMocktokenManager(mockCtrl)
			claimSpec    = ClaimSpec{UserID: ""}
			expectedErr  = errors.New("invalid ClaimSpec")
		)

		// mock expectation
		mockTokenMgr.EXPECT().GenerateToken(claimSpec).Return("", expectedErr)

		svc := &TokenService{
			tokenManager: mockTokenMgr,
		}
		tokenString, err := svc.GenerateToken(claimSpec)

		// assertions
		g.Expect(err).To(gomega.Equal(expectedErr))
		g.Expect(tokenString).To(gomega.BeEmpty())
	})
}

func TestTokenService_ValidateToken(t *testing.T) {
	t.Parallel()

	t.Run("ValidateTokenWithShouldReturnExpectedClaims", func(t *testing.T) {
		var (
			g              = gomega.NewWithT(t)
			mockCtrl       = gomock.NewController(t)
			mockTokenMgr   = NewMocktokenManager(mockCtrl)
			tokenString    = "valid_token_string"
			expectedClaims = &Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "123",
				},
			}
		)

		// mock expectation
		mockTokenMgr.EXPECT().ValidateToken(tokenString).Return(expectedClaims, nil)

		svc := &TokenService{
			tokenManager: mockTokenMgr,
		}
		claims, err := svc.ValidateToken(tokenString)

		// assertions
		g.Expect(err).To(gomega.BeNil())
		g.Expect(claims).To(gomega.Equal(expectedClaims))
	})

	t.Run("ErrorValidateTokenWithShouldReturnExpectedError", func(t *testing.T) {
		var (
			g            = gomega.NewWithT(t)
			mockCtrl     = gomock.NewController(t)
			mockTokenMgr = NewMocktokenManager(mockCtrl)
			tokenString  = "invalid_token_string"
			expectedErr  = errors.New("failed to validate token")
		)

		mockTokenMgr.EXPECT().ValidateToken(tokenString).Return(nil, expectedErr)

		svc := &TokenService{
			tokenManager: mockTokenMgr,
		}
		claims, err := svc.ValidateToken(tokenString)

		// assertions
		g.Expect(err).To(gomega.Equal(expectedErr))
		g.Expect(claims).To(gomega.BeNil())
	})

	t.Run("ValidateTokenWithInvalidTokenShouldReturnExpectedError", func(t *testing.T) {
		var (
			g            = gomega.NewWithT(t)
			mockCtrl     = gomock.NewController(t)
			mockTokenMgr = NewMocktokenManager(mockCtrl)
			tokenString  = ""
			expectedErr  = errors.New("token string is empty")
		)

		mockTokenMgr.EXPECT().
			ValidateToken(tokenString).
			Return(nil, expectedErr)

		svc := &TokenService{
			tokenManager: mockTokenMgr,
		}
		claims, err := svc.ValidateToken(tokenString)

		// assertions
		g.Expect(err).To(gomega.Equal(expectedErr))
		g.Expect(claims).To(gomega.BeNil())
	})
}
