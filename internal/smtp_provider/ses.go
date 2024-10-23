//go:generate mockgen -typed -source=ses.go -destination=ses_mock.go -package=smtp_provider
package smtp_provider

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESSendEmailAPI interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

type sesProvider struct {
	sesClient SESSendEmailAPI
}

// SendEmail is an abstraction for SES' send email function
func (s sesProvider) SendEmail(email Email) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.sendEmail(ctx, s.sesClient, email)
}

// sendEmail calls the SES API
func (s sesProvider) sendEmail(ctx context.Context, client SESSendEmailAPI, email Email) error {
	inputPayload := s.buildInputPayload(email)

	result, err := client.SendEmail(ctx, inputPayload)
	if err != nil {
		log.Printf("failed executing SendEmail : %v", err)
		return err
	}

	log.Printf("email sent: %v", result)
	return nil
}

func (sesProvider) buildInputPayload(email Email) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: email.To,
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(email.Body),
				},
			},
			Subject: &types.Content{
				Data: aws.String(email.Subject),
			},
		},
		Source: aws.String(email.From),
	}
}

func NewSESProvider(cfg aws.Config) *sesProvider {
	return &sesProvider{
		sesClient: ses.NewFromConfig(cfg),
	}
}
