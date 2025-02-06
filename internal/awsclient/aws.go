package awsclient

import (
	config2 "FTPArchive/internal/config"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
	"path/filepath"
)

func CreateAwsClient() (*s3.Client, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func UploadFileAWS(profile *config2.Profile) error {
	log.Printf("Uploading file %s to bucket %s", profile.ArchivePath, profile.BucketName)
	client, err := CreateAwsClient()
	if err != nil {
		return fmt.Errorf("could not create aws client: %w", err)
	}

	localFile, err := os.Open(profile.ArchivePath)
	if err != nil {
		return fmt.Errorf("could not open local file: %w", err)
	}
	defer localFile.Close()

	fileInfo, err := localFile.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %w", err)
	}

	baseName := filepath.Base(profile.ArchivePath)
	_, err = client.PutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket:        aws.String(profile.BucketName),
			Key:           aws.String(baseName),
			Body:          localFile,
			ContentLength: aws.Int64(fileInfo.Size()),
			ContentType:   aws.String("application/zip"),
			ACL:           types.ObjectCannedACLPrivate,
		})
	if err != nil {
		return fmt.Errorf("could not upload file: %w", err)
	}

	log.Printf("successfully uploaded file: %s", profile.ArchivePath)
	return nil
}
