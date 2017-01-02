package client

import (
	minio "github.com/minio/minio-go"
)

type MinioClient struct {
	Client     *minio.Client
	BucketInfo *minio.BucketInfo
}

// NewMinioClient returns a new minio client based on passed access specs.
func NewMinioClient(serverURI, accessKeyID, secretAccessKey string, secure bool) (*minio.Client, error) {
	c, err := minio.New(serverURI, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return nil, err
	}
	return c, nil
}
