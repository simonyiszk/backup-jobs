package common

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
)

func InitDriveService(ctx context.Context, serviceAccountKey string) (*drive.Service, error) {
	return drive.NewService(ctx, option.WithCredentialsJSON([]byte(serviceAccountKey)))
}

func CreateDriveFolder(name, parent, driveId string, service *drive.Service) (string, error) {
	driveFolder := &drive.File{Name: name, MimeType: "application/vnd.google-apps.folder", Parents: []string{parent}, DriveId: driveId}
	folder, err := service.Files.Create(driveFolder).SupportsAllDrives(true).Do()
	if err != nil {
		return "", err
	}
	return folder.Id, nil
}

func UploadToGoogleDrive(fileName string, data io.Reader, parent, driveId string, service *drive.Service) error {
	driveFile := &drive.File{Name: fileName, Parents: []string{parent}, DriveId: driveId}
	uploadedFile, err := service.Files.Create(driveFile).Media(data).SupportsAllDrives(true).Do()
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	log.Println("File uploaded successfully! File ID:", uploadedFile)
	return nil
}
