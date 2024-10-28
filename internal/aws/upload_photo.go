package aws

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"mime/multipart"
)

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
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("uploadPhotoToS3: failed to upload object: %v", err)
	}

	url := fmt.Sprintf("https://s3-%s.amazonaws.com/%s", bucketName, key)

	return url, nil
}
