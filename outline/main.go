package main

import (
	"bytes"
	"common"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const waitDelay = time.Second * 5

var endpoint string
var outline string
var waitTimeout time.Duration
var serviceAccountKey string
var parent string
var driveId string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Exporting all Outline collections to json")

	readConfig()

	operationId, err := startExportOperation()
	if err != nil {
		log.Fatalln("Failed to start export operation", err)
	}

	log.Println("Started export operation:", operationId)
	defer func(operationId string) {
		err := deleteExportOperation(operationId)
		if err != nil {
			log.Println("Cleanup failed", err)
		}
	}(operationId)
	err = waitForFileOperationToFinish(operationId)
	if err != nil {
		log.Fatalln("Export operation failed", err)
		return
	}

	data, err := retrieveExportedFile(operationId)
	log.Println("retrieved data with size", len(data))

	today := time.Now().Format(time.DateOnly)
	fileName := fmt.Sprintf("%s-outline-exported-collections-json.zip", today)
	log.Println("Saving export to Google Drive:", fileName)

	err = common.UploadToGoogleDrive(fileName, bytes.NewReader(data), parent, driveId, serviceAccountKey)
	if err != nil {
		log.Fatalln("Failed to upload file to Google Drive:", err)
	}
}

func waitForFileOperationToFinish(operationId string) error {
	resultChan := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	go func() {
		for {
			log.Println("Waiting for operation to finish:", operationId)
			select {
			case <-ctx.Done():
				resultChan <- ctx.Err()
				return
			case <-time.After(waitDelay):
			}

			finished, err := checkIfFileOperationFinished(ctx, operationId)
			if err != nil || finished {
				resultChan <- err
				return
			}
		}
	}()

	return <-resultChan
}

func readConfig() {
	var exists bool
	endpoint, exists = os.LookupEnv("OUTLINE_ENDPOINT")
	if !exists {
		log.Fatalln("OUTLINE_ENDPOINT not set")
	}

	outline, exists = os.LookupEnv("OUTLINE_SECRET")
	if !exists {
		log.Fatalln("OUTLINE_SECRET not set")
	}

	waitTimeoutString, exists := os.LookupEnv("WAIT_TIMEOUT_SECONDS")
	if exists {
		timeout, err := strconv.Atoi(waitTimeoutString)
		if err != nil {
			log.Fatalln("WAIT_TIMEOUT_SECONDS must be a number")

		}
		waitTimeout = time.Second * time.Duration(timeout)
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
