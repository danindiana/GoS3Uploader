package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "mime/multipart"
    "net/http"
)

// S3Uploader encapsulates the AWS S3 uploader
type S3Uploader struct {
    Uploader *s3manager.Uploader
}

// NewS3Uploader creates a new S3Uploader instance
func NewS3Uploader() (*S3Uploader, error) {
    // Start a new AWS session
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2"), // Replace with your AWS region
    })
    if err != nil {
        return nil, fmt.Errorf("unable to create AWS session: %v", err)
    }

    // Create a new Uploader with the session
    uploader := s3manager.NewUploader(sess)

    // Return an S3Uploader which wraps the Uploader
    return &S3Uploader{
        Uploader: uploader,
    }, nil
}

// UploadFile handles the uploading of a file to S3
func (u *S3Uploader) UploadFile(file multipart.File, fileName string, bucketName string) (*s3manager.UploadOutput, error) {
    // Prepare the file to be uploaded to S3
    uploadInput := &s3manager.UploadInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(fileName),
        Body:   file,
    }

    // Perform the upload to S3
    result, err := u.Uploader.Upload(uploadInput)
    if err != nil {
        return nil, fmt.Errorf("failed to upload file to S3: %v", err)
    }

    return result, nil
}

// UploadHandler is an HTTP handler that processes file upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // Only accept POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
        return
    }

    // Parse the multipart form
    const maxUploadSize = 10 << 20 // 10 MB
    if err := r.ParseMultipartForm(maxUploadSize); err != nil {
        http.Error(w, fmt.Sprintf("Could not parse multipart form: %v", err), http.StatusBadRequest)
        return
    }

    // The form field "file" is used for the file upload
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, fmt.Sprintf("Could not get uploaded file: %v", err), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create an S3 uploader
    s3Uploader, err := NewS3Uploader()
    if err != nil {
        http.Error(w, fmt.Sprintf("Could not create S3 uploader: %v", err), http.StatusInternalServerError)
        return
    }

    // Define the bucket to which we're uploading the file
    bucketName := "your-s3-bucket-name" // Replace with your actual bucket name

    // Upload the file to S3
    uploadResult, err := s3Uploader.UploadFile(file, header.Filename, bucketName)
    if err != nil {
        http.Error(w, fmt.Sprintf("Could not upload file to S3: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with the location of the file in S3
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "File uploaded successfully: %s\n", uploadResult.Location)
}
