package dependency

import (
	"context"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func NewAWSConfig(cfg config.AWS) *aws.Config {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		utils.Logger.Fatalf("unable to load SES configuration: %v", err)
	}

	return &awsCfg
}
