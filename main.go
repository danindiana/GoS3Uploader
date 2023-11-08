package main

import (
    "log"
    "net/http"
)

func main() {
    // Set up routes
    http.HandleFunc("/upload", UploadHandler) // Use the UploadHandler from s3uploader.go

    // Start the server
    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}

// Note: The UploadHandler is now defined in s3uploader.go
