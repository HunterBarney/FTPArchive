package gcp

import (
	"FTPArchive/internal/config"
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
)

func CreateGcpClient(profile *config.Profile) (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(profile.UploadProfile.CredentialsFile))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func UploadArchiveGcp() string {
	return "not implemented"
}
