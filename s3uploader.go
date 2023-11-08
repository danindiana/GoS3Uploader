package main

import (
    "fmt"
    "io"
    "mime/multipart"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Uploader encapsulates the AWS S3 uploader
type S3Uploader struct {
    Uploader *s3manager.Uploader
    Bucket   string
}

// NewS3Uploader creates a new S3Uploader instance
func NewS3Uploader(bucket string) (*S3Uploader, error) {
    // Initialize a session that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials
    // and region from the shared configuration file ~/.aws/config.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")}, // for example, "us-west-2"
    )
    if err != nil {
        return nil, err
    }

    // Create an uploader with the session and default options
    uploader := s3manager.NewUploader(sess)

    // Return an S3Uploader which wraps the uploader
    return &S3Uploader{
        Uploader: uploader,
        Bucket:   bucket,
    }, nil
}

// UploadFile handles the uploading of a file to S3
func (u *S3Uploader) UploadFile(file multipart.File, fileName string) (*s3manager.UploadOutput, error) {
    // Upload input parameters
    upParams := &s3manager.UploadInput{
        Bucket: &u.Bucket,
        Key:    aws.String(fileName),
        Body:   file,
    }

    // Perform an upload.
    result, err := u.Uploader.Upload(upParams)
    if err != nil {
        return nil, fmt.Errorf("failed to upload file, %v", err)
    }
    return result, nil
}

// UploadHandler will be used in the HTTP server to handle file upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the multipart form
    err := r.ParseMultipartForm(10 << 20) // Max upload size ~10MB
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Retrieve the file from posted form-data
    file, handler, err := r.FormFile("myfile")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create an S3 uploader
    s3Uploader, err := NewS3Uploader("my-s3-bucket") // Replace with your actual bucket name
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Upload the file to S3
    uploadResult, err := s3Uploader.UploadFile(file, handler.Filename)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with the location of the file in S3
    w.Write([]byte(fmt.Sprintf("Successfully uploaded to %s", uploadResult.Location)))
}
