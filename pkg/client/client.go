package client

import (
	minio "github.com/minio/minio-go"
)

const (
	location = "us-east-1"
)

// MinioClient is a wrapper around minio.Client that also holds the bucket name
// where we want to copy files.
type MinioClient struct {
	Client          *minio.Client
	BucketName      string
	ServerURI       string
	AccesKeyID      string
	SecretAccessKey string
}

// NewMinioClient returns a new minio client based on passed access specs and
// creates a new bucket if it doesn't exist.
func NewMinioClient(serverURI, accessKeyID, secretAccessKey, bucket string, secure bool) (*MinioClient, error) {
	c, err := minio.New(serverURI, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return nil, err
	}

	exists, err := c.BucketExists(bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := c.MakeBucket(bucket, location); err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		Client:          c,
		BucketName:      bucket,
		ServerURI:       serverURI,
		AccesKeyID:      accessKeyID,
		SecretAccessKey: secretAccessKey,
	}, nil
}
