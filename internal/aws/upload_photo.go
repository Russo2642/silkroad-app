package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetBucketName возвращает имя bucket'а из конфигурации
func GetBucketName() string {
	return bucketName
}

func processFile(file multipart.File) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("processFile: failed to copy file to buffer: %v", err)
	}
	return buf, nil
}

func UploadPhotoToS3(bucketName, key string, file multipart.File) (string, error) {
	buf, err := processFile(file)
	if err != nil {
		return "", fmt.Errorf("uploadPhotoToS3: failed to copy file to buffer: %v", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("uploadPhotoToS3: failed to close file: %v", err)
		}
	}(file)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(http.DetectContentType(buf.Bytes())),
	})
	if err != nil {
		return "", fmt.Errorf("uploadPhotoToS3: failed to upload object: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, key)

	return url, nil
}
