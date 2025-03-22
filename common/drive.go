package common

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
)

func UploadToGoogleDrive(fileName string, data io.Reader, parent, driveId, serviceAccountKey string) error {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(serviceAccountKey)))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	driveFile := &drive.File{Name: fileName, Parents: []string{parent}, DriveId: driveId}
	uploadedFile, err := srv.Files.Create(driveFile).Media(data).SupportsAllDrives(true).Do()
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	log.Println("File uploaded successfully! File ID:", uploadedFile)
	return nil
}
