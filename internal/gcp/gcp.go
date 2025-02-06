package gcp

import (
	"FTPArchive/internal/config"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CreateGcpClient() (*storage.Client, error) {
	ctx := context.Background()
	log.Println("Connecting to Google Cloud Platform")
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Cloud Platform client: %v", err)
	}
	return client, nil
}

func UploadArchiveGcp(profile *config.Profile) error {
	client, err := CreateGcpClient()
	if err != nil {
		return err
	}

	bucket := client.Bucket(profile.BucketName)

	baseName := filepath.Base(profile.ArchivePath)
	obj := bucket.Object(baseName)
	writer := obj.NewWriter(context.Background())
	log.Printf("Uploading Archive as %s to bucket %s", baseName, profile.BucketName)

	archive, err := os.Open(profile.ArchivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %v", err)
	}
	defer archive.Close()

	_, err = io.Copy(writer, archive)
	if err != nil {
		return fmt.Errorf("failed to upload archive: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to upload archive: %v", err)
	}

	fmt.Println("Uploaded archive successfully")
	return nil
}
