package main

import (
	"common"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

var dumpName string
var driveFileName string
var serviceAccountKey string
var parent string
var driveId string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	readConfig()

	reader, err := os.Open(dumpName)
	if err != nil {
		log.Fatalln("Failed to open dump:", err)
	}
	defer reader.Close()

	log.Println("Uploading MySQL dump to Google Drive")
	today := time.Now().Format(time.DateOnly)
	fileName := fmt.Sprintf("%s-%s-mysql.sql", today, driveFileName)
	log.Println("Saving export to Google Drive:", fileName)

	driveService, err := common.InitDriveService(context.Background(), serviceAccountKey)
	if err != nil {
		log.Fatalln("Failed to init Google Drive client:", err)
	}
	err = common.UploadToGoogleDrive(fileName, reader, parent, driveId, driveService)
	if err != nil {
		log.Fatalln("Failed to upload file to Google Drive:", err)
	}
}

func readConfig() {
	var exists bool
	driveFileName, exists = os.LookupEnv("DRIVE_FILE_NAME")
	if !exists {
		log.Fatalln("DRIVE_FILE_NAME not set")
	}

	dumpName, exists = os.LookupEnv("DUMP_NAME")
	if !exists {
		log.Fatalln("DUMP_NAME not set")
	}

	serviceAccountKey, exists = os.LookupEnv("GOOGLE_SERVICE_ACCOUNT_KEY")
	if !exists {
		log.Fatalln("GOOGLE_SERVICE_ACCOUNT_KEY not set")
	}

	parent, exists = os.LookupEnv("DRIVE_PARENT_ID")
	if !exists {
		log.Fatalln("DRIVE_PARENT_ID not set")
	}

	driveId, exists = os.LookupEnv("DRIVE_ID")
	if !exists {
		log.Fatalln("DRIVE_ID not set")
	}
}
