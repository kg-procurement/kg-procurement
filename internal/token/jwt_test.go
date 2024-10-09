package token

import (
	"github.com/benbjohnson/clock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"kg/procurement/cmd/config"
	"testing"
	"time"
)

func Test_newJWTManager(t *testing.T) {
	_ = newJWTManager(config.Token{}, nil)
}

func Test_GenerateToken(t *testing.T) {
	t.Parallel()

	var (
		mockClock *clock.Mock
		jwtMgr    *jwtManager
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		mockClock = clock.NewMock()

		jwtMgr = &jwtManager{
			cfg:   config.Token{Secret: "secret"},
			clock: mockClock,
		}

		return gomega.NewWithT(t)
	}

	t.Run("GenerateTokenWithValidClaimsReturnsTokenSuccessfully", func(t *testing.T) {
		g := setup(t)
		spec := ClaimSpec{UserID: "123"}

		// create a valid token
		tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   spec.UserID,
				ExpiresAt: jwt.NewNumericDate(mockClock.Now().UTC().Add(30 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(mockClock.Now().UTC()),
				ID:        uuid.NewString(),
			},
		})

		// sign the token
		token, err := tokenObject.SignedString([]byte(jwtMgr.cfg.Secret))

		g.Expect(token).ShouldNot(gomega.BeEmpty())
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}
