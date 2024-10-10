package token

import (
	"errors"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestTokenService_GenerateToken(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
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

	t.Run("returns error when tokenManager fails", func(t *testing.T) {
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

	t.Run("returns error when ClaimSpec is invalid", func(t *testing.T) {
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
