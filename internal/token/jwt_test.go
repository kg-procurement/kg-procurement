package token

import (
	"errors"
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

	t.Run("GenerateTokenWithInvalidSecretReturnsError", func(t *testing.T) {
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

		// sign the token with an invalid secret
		token, err := tokenObject.SignedString("haha")

		g.Expect(token).Should(gomega.Equal(""))
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
		tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   spec.UserID,
				ExpiresAt: jwt.NewNumericDate(mockClock.Now().UTC().Add(30 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(mockClock.Now().UTC()),
				ID:        uuid.NewString(),
			},
		})

		// sign the token
		tokenString, err := tokenObject.SignedString([]byte(jwtMgr.cfg.Secret))
		g.Expect(err).ShouldNot(gomega.HaveOccurred())

		// validate the token
		claims, err := jwtMgr.ValidateToken(tokenString)

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

		// create a valid token
		tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   spec.UserID,
				ExpiresAt: jwt.NewNumericDate(mockClock.Now().UTC().Add(30 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(mockClock.Now().UTC()),
				ID:        uuid.NewString(),
			},
		})

		// sign the token with different secret
		tokenString, err := tokenObject.SignedString([]byte("differentSecret"))
		g.Expect(err).ShouldNot(gomega.HaveOccurred())

		// validate the token
		_, err = jwtMgr.ValidateToken(tokenString)

		// assertions
		g.Expect(err).Should(gomega.HaveOccurred())
		g.Expect(errors.Is(err, jwt.ErrSignatureInvalid)).Should(gomega.BeTrue())
	})
}
