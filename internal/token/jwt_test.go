package token

import (
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/golang-jwt/jwt/v5"
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
		token, err := jwtMgr.GenerateToken(spec)

		g.Expect(token).ShouldNot(gomega.BeEmpty())
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
	})

	t.Run("GenerateTokenWithInvalidSecretReturnsError", func(t *testing.T) {
		g := setup(t)
		spec := ClaimSpec{UserID: "123"}

		// invalid secret
		jwtMgr.cfg.Secret = ""
		token, err := jwtMgr.GenerateToken(spec)

		// assertions
		g.Expect(token).Should(gomega.BeEmpty())
		g.Expect(err).Should(gomega.HaveOccurred())
	})
}

func Test_ValidateToken(t *testing.T) {
	t.Parallel()

	var (
		mockClock *clock.Mock
		jwtMgr    *jwtManager
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		mockClock = clock.NewMock()
		mockClock.Set(time.Now())

		jwtMgr = &jwtManager{
			cfg:   config.Token{Secret: "secret"},
			clock: mockClock,
		}

		return gomega.NewWithT(t)
	}

	t.Run("ValidateTokenWithValidTokenReturnsClaimsSuccessfully", func(t *testing.T) {
		g := setup(t)
		spec := ClaimSpec{UserID: "123"}

		// create a valid token
		token, err := jwtMgr.GenerateToken(spec)

		// validate the token
		claims, err := jwtMgr.ValidateToken(token)

		// assertions
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(claims.Subject).Should(gomega.Equal("123"))
	})

	t.Run("ValidateTokenWithMalformedTokenReturnsError", func(t *testing.T) {
		g := setup(t)

		// malformed token
		invalidToken := "invalid.token"

		// validate token
		_, err := jwtMgr.ValidateToken(invalidToken)

		// assertions
		g.Expect(err).Should(gomega.HaveOccurred())
		g.Expect(errors.Is(err, jwt.ErrTokenMalformed)).Should(gomega.BeTrue())
	})

	t.Run("ValidateTokenWithInvalidSignatureReturnsError", func(t *testing.T) {
		g := setup(t)
		spec := ClaimSpec{UserID: "123"}

		token, err := jwtMgr.GenerateToken(spec)

		// validate the token with different secret
		jwtMgr.cfg.Secret = "differentSecret"
		_, err = jwtMgr.ValidateToken(token)

		// assertions
		g.Expect(err).Should(gomega.HaveOccurred())
		g.Expect(errors.Is(err, jwt.ErrSignatureInvalid)).Should(gomega.BeTrue())
	})

	t.Run("ValidateTokenWithInvalidConfigShouldReturnError", func(t *testing.T) {
		g := setup(t)
		spec := ClaimSpec{UserID: "123"}

		// create a valid token
		token, err := jwtMgr.GenerateToken(spec)

		// validate the token
		jwtMgr.cfg.Secret = ""
		_, err = jwtMgr.ValidateToken(token)

		// assertions
		g.Expect(err).Should(gomega.HaveOccurred())
		g.Expect(err.Error()).Should(gomega.ContainSubstring("secret key is empty"))
	})
}
