package aws

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
	"log"
)

var s3Client *s3.Client

func InitS3Client() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal("unable to load SDK config, " + err.Error())
	}

	s3Client = s3.NewFromConfig(cfg)
}
