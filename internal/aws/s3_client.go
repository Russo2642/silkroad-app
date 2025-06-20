package aws

import (
	"fmt"
	appConfig "silkroad/m/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
)

var (
	s3Client   *s3.Client
	awsRegion  string
	bucketName string
)

func InitS3Client(awsConfig appConfig.AWSConfig) error {
	awsRegion = awsConfig.Region
	bucketName = awsConfig.BucketName
	var cfg aws.Config
	var err error

	if awsConfig.AccessKeyID != "" && awsConfig.SecretAccessKey != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(awsConfig.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				awsConfig.AccessKeyID,
				awsConfig.SecretAccessKey,
				"",
			)),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsConfig.Region))
	}

	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	return nil
}
