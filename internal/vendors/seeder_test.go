package vendors

import (
	"context"
	"errors"
	"testing"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewSeeder(t *testing.T) {
	_ = NewSeeder(nil)
}

func Test_NewDBSeederWriter(t *testing.T) {
	_ = NewDBSeederWriter(nil, nil)
}

func Test_Seeder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var (
		mockWriter *MockseederDataWriter
		subject    *Seeder
	)

	setup := func(t *testing.T) *gomega.GomegaWithT {
		ctrl := gomock.NewController(t)

		mockWriter = NewMockseederDataWriter(ctrl)
		subject = NewSeeder(mockWriter)
		return gomega.NewWithT(t)
	}
	t.Run("setup product", func(t *testing.T) {
		vendors := []Vendor{{ID: "1"}, {ID: "2"}}

		t.Run("success", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeVendor(ctx, Vendor{ID: "1"})
			mockWriter.EXPECT().writeVendor(ctx, Vendor{ID: "2"})

			err := subject.SetupVendors(ctx, vendors)
			g.Expect(err).ShouldNot(gomega.HaveOccurred())
		})

		t.Run("error", func(t *testing.T) {
			g := setup(t)

			mockWriter.EXPECT().writeVendor(ctx, Vendor{ID: "1"})
			mockWriter.EXPECT().writeVendor(ctx, Vendor{ID: "2"}).Return(errors.New("error"))

			err := subject.SetupVendors(ctx, vendors)
			g.Expect(err).Should(gomega.HaveOccurred())
		})
	})

}
