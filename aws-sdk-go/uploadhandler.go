package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// UploadHandler is an HTTP handler that processes file upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a max memory of 32MB
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse multipart form: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve the 'file' from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get uploaded file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Retrieve the 'partNumber' from form data
	partNumberStr := r.FormValue("partNumber")
	partNumber, err := strconv.Atoi(partNumberStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid part number: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve the 'uploadID' from form data
	uploadID := r.FormValue("uploadID")
	if uploadID == "" {
		http.Error(w, "uploadID is required", http.StatusBadRequest)
		return
	}

	// Create an S3 uploader
	s3Uploader, err := NewS3Uploader()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create S3 uploader: %v", err), http.StatusInternalServerError)
		return
	}

	// Upload the file chunk to S3
	uploadResult, err := s3Uploader.UploadPart(file, header.Filename, uploadID, int64(partNumber))
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not upload file part to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the ETag of the uploaded part
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Part uploaded successfully: %s\n", *uploadResult.ETag)
}
