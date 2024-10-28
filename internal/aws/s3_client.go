package aws

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
	"log"
	"os"
)

var s3Client *s3.Client

func InitS3Client() {
	log.Println("AWS_ACCESS_KEY_ID:", os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Println("AWS_SECRET_ACCESS_KEY:", os.Getenv("AWS_SECRET_ACCESS_KEY"))

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal("unable to load SDK config, " + err.Error())
	}
	log.Println("AWS_ACCESS_KEY_ID:", os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Println("AWS_SECRET_ACCESS_KEY:", os.Getenv("AWS_SECRET_ACCESS_KEY"))

	s3Client = s3.NewFromConfig(cfg)
}
