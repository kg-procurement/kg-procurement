package account

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"testing"

	"bou.ke/monkey"
	"github.com/onsi/gomega"
)

func Test_GenerateRandomID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		g := gomega.NewWithT(t)

		id, err := GenerateRandomID()

		g.Expect(err).To(gomega.BeNil())
		g.Expect(len(id)).To(gomega.Equal(15))
	})

	t.Run("success - unique IDs", func(t *testing.T) {
		g := gomega.NewWithT(t)

		id1, err1 := GenerateRandomID()
		id2, err2 := GenerateRandomID()

		g.Expect(err1).To(gomega.BeNil())
		g.Expect(err2).To(gomega.BeNil())
		g.Expect(id1).ToNot(gomega.Equal(id2))
	})

	t.Run("error - random generation failure", func(t *testing.T) {
		g := gomega.NewWithT(t)

		monkey.Patch(rand.Int, func(_ io.Reader, _ *big.Int) (*big.Int, error) {
			return nil, errors.New("failed to randomize int")
		})
		defer monkey.Unpatch(rand.Int)

		_, err := GenerateRandomID()

		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}
