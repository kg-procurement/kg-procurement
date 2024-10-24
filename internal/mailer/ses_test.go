package mailer

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewSESProvider(t *testing.T) {
	_ = NewSESProvider(aws.Config{})
}

func Test_ProviderSendEmail(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)

	mockSES := NewMockSESSendEmailAPI(ctrl)
	provider := sesProvider{
		sesClient: mockSES,
	}

	mockSES.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&ses.SendEmailOutput{}, nil)

	err := provider.SendEmail(Email{})
	g.Expect(err).To(gomega.BeNil())
}

func TestSendEmail(t *testing.T) {
	t.Parallel()

	var (
		email = Email{
			From:    "from@example.com",
			To:      []string{"to@example.com"},
			Subject: "Test Subject",
			Body:    "Test Body",
		}
	)

	t.Run("success", func(t *testing.T) {
		g := gomega.NewWithT(t)
		ctrl := gomock.NewController(t)
		ctx := context.TODO()

		mockSES := NewMockSESSendEmailAPI(ctrl)
		provider := sesProvider{mockSES}
		mockSES.EXPECT().SendEmail(ctx, gomock.Any(), gomock.Any()).
			Return(&ses.SendEmailOutput{}, nil)

		err := provider.sendEmail(ctx, mockSES, email)
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		g := gomega.NewWithT(t)
		ctrl := gomock.NewController(t)
		ctx := context.TODO()

		mockSES := NewMockSESSendEmailAPI(ctrl)
		provider := sesProvider{mockSES}
		mockSES.EXPECT().SendEmail(ctx, gomock.Any(), gomock.Any()).
			Return(nil, errors.New("oh noo"))

		err := provider.sendEmail(ctx, mockSES, email)
		g.Expect(err).ToNot(gomega.BeNil())
	})
}
