package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"glk-web-app/config"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// GDriveUploadResponse is the structure returned by our Google Apps Script.
type GDriveUploadResponse struct {
	Status  string `json:"status"`
	FileID  string `json:"fileId"`
	FileURL string `json:"fileUrl"`
	Message string `json:"message"`
}

// UploadToGDrive reads a multipart.FileHeader and uploads it to Google Drive via Apps Script.
func UploadToGDrive(fileHeader *multipart.FileHeader) (string, error) {
	// Ambil URL Apps Script dari .env
	gasURL := config.GetEnv("GAS_UPLOAD_URL", "")
	if gasURL == "" {
		return "", errors.New("GAS_UPLOAD_URL is not set in .env")
	}

	// Buka file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Baca seluruh isi file ke memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Ubah file menjadi base64 string
	base64Data := base64.StdEncoding.EncodeToString(fileBytes)

	// Tentukan Mime Type
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Dapatkan extension untuk penamaan
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".file" // Fallback
	}

	// Buat payload JSON
	payload := map[string]string{
		"filename": fileHeader.Filename,
		"fileData": base64Data,
		"mimeType": mimeType,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Kirim HTTP POST request ke Google Apps Script
	resp, err := http.Post(gasURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Baca response dari GAS
	var gasResp GDriveUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&gasResp); err != nil {
		return "", fmt.Errorf("failed to decode response from GAS: %v", err)
	}

	if gasResp.Status != "success" {
		return "", fmt.Errorf("GAS Error: %s", gasResp.Message)
	}

	// Berhasil, return File URL
	return gasResp.FileURL, nil
}
