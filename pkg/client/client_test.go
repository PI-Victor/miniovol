package client

import (
	"testing"
)

func TestNewMinioClient(t *testing.T) {
	var (
		serverURI       = "testlocal:9000"
		accessKeyID     = "abc123"
		secretAccessKey = "secretKey"
		bucket          = "testBucket"
		secure          = false
	)

	// NOTE: remember to do something better than this.
	c, err := NewMinioClient(serverURI, accessKeyID, secretAccessKey, bucket, secure)
	if err != nil {
		t.Fatalf("An error occured while creating a new client: %#v", err)
	}
	if c.AccesKeyID != accessKeyID {
		t.Errorf("Expected %s got %s", accessKeyID, c.AccesKeyID)
	}
	if c.SecretAccessKey != secretAccessKey {
		t.Errorf("Expected %s got %s", secretAccessKey, c.SecretAccessKey)
	}
	if c.BucketName != bucket {
		t.Errorf("Expected %s got %s", bucket, c.BucketName)
	}
}
