package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Uploader encapsulates the AWS S3 uploader
type S3Uploader struct {
	S3 *s3.S3
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

	// Create a new S3 service client
	s3Svc := s3.New(sess)

	// Return an S3Uploader which wraps the S3 service client
	return &S3Uploader{
		S3: s3Svc,
	}, nil
}

// UploadPart handles the uploading of a single part to S3
func (u *S3Uploader) UploadPart(file multipart.File, bucketName, key, uploadID string, partNumber int64) (*s3.UploadPartOutput, error) {
	// Read the file part into a buffer
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file part: %v", err)
	}

	// Create the upload part request
	uploadPartInput := &s3.UploadPartInput{
		Bucket:     aws.String(bucketName),
		Key:        aws.String(key),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int64(partNumber),
		Body:       bytes.NewReader(buffer),
	}

	// Upload the part and return the response
	resp, err := u.S3.UploadPart(uploadPartInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload part to S3: %v", err)
	}

	return resp, nil
}
