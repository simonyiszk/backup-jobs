package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type FileOperationInfoResult struct {
	Data struct {
		State string `json:"state"`
	} `json:"data"`
}

type ExportAllResponse struct {
	Data struct {
		FileOperation struct {
			Id string `json:"id"`
		} `json:"fileOperation"`
	} `json:"data"`
}

func postApi(ctx context.Context, path, body string) ([]byte, error) {
	apiUrl, err := url.JoinPath(endpoint, path)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(body))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+outline)
	req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if res.StatusCode != 200 || err != nil {
		return []byte{}, fmt.Errorf("go status [%v] %v %v", res.StatusCode, string(responseBody), err)
	}

	return responseBody, nil
}

func startExportOperation() (string, error) {
	exportResult, err := postApi(context.Background(), "/api/collections.export_all", `{"format":"json"}`)
	if err != nil {
		return "", err
	}

	log.Println("Got response for collections.export_all:", string(exportResult))
	var exportAllResponse ExportAllResponse
	err = json.Unmarshal(exportResult, &exportAllResponse)
	if err != nil {
		return "", err
	}

	return exportAllResponse.Data.FileOperation.Id, nil
}

func checkIfFileOperationFinished(ctx context.Context, operationId string) (bool, error) {
	operationResult, err := postApi(ctx, "/api/fileOperations.info", fmt.Sprintf(`{"id":"%s"}`, operationId))
	if err != nil {
		return false, err
	}

	log.Println("Got response for collections.export_all:", string(operationResult))
	var operationState FileOperationInfoResult
	err = json.Unmarshal(operationResult, &operationState)
	if err != nil {
		return false, err
	}

	if operationState.Data.State == "error" || operationState.Data.State == "expired" {
		return false, errors.New("file operation failed")
	}

	return operationState.Data.State == "complete", nil
}

func retrieveExportedFile(operationId string) ([]byte, error) {
	log.Println("Retrieving exported file")
	operationResult, err := postApi(context.Background(), "/api/fileOperations.redirect", fmt.Sprintf(`{"id":"%s"}`, operationId))
	if err != nil {
		return []byte{}, err
	}

	return operationResult, nil
}

func deleteExportOperation(operationId string) error {
	log.Println("Cleaning up file operation")
	operationResult, err := postApi(context.Background(), "/api/fileOperations.delete", fmt.Sprintf(`{"id":"%s"}`, operationId))
	if err != nil {
		return err
	}

	log.Println("Got response for collections.delete:", string(operationResult))
	return nil
}
