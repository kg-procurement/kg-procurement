package mailer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewEmailStatusService(t *testing.T) {
	_ = NewEmailStatusService(nil, nil)
}

func TestEmailStatusService_WriteEmailStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		now       = time.Now()
		writeSpec = EmailStatus{
			ID:           "some_id",
			EmailTo:      "email@email.com",
			Status:       "sent",
			ModifiedDate: now,
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			g                       = gomega.NewWithT(t)
			ctx                     = context.Background()
			mockCtrl                = gomock.NewController(t)
			mockEmailStatusAccessor = NewMockemailStatusDBAccessor(mockCtrl)
		)

		svc := &EmailStatusService{
			emailStatusDBAccessor: mockEmailStatusAccessor,
		}

		mockEmailStatusAccessor.EXPECT().WriteEmailStatus(ctx, writeSpec).Return(nil)

		err := svc.WriteEmailStatus(ctx, writeSpec)
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("returns error on accessor failure", func(t *testing.T) {
		var (
			g                       = gomega.NewWithT(t)
			ctx                     = context.Background()
			mockCtrl                = gomock.NewController(t)
			mockEmailStatusAccessor = NewMockemailStatusDBAccessor(mockCtrl)
		)

		svc := &EmailStatusService{
			emailStatusDBAccessor: mockEmailStatusAccessor,
		}

		mockEmailStatusAccessor.EXPECT().WriteEmailStatus(ctx, writeSpec).Return(errors.New("create error"))

		err := svc.WriteEmailStatus(ctx, writeSpec)
		g.Expect(err).ShouldNot(gomega.BeNil())
	})
}
