package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/client"
	"github.com/cloudflavor/miniovol/pkg/driver"
)

const (
	socketAddress = "/run/docker/plugins/miniovol.sock"
	rootID        = 0
)

// Error implements the error interface and uses it to return a generic
// validation error for environment variables that are mandatory.
type Error struct {
	envVar string
}

func (e Error) Error() string {
	return fmt.Sprintf("environment variable %s cannot be empty", e.envVar)
}

func newErrEmptyEnvVar(v string) error {
	return Error{
		envVar: v,
	}
}

func getEnvDetails() (serverURI, accessKeyID, secretAccessKey, bucket string, secure bool, err error) {
	secure, err = strconv.ParseBool(os.Getenv("SECURE"))
	if err != nil {
		return
	}
	serverURI = os.Getenv("MINIO_SERVER")
	if serverURI == "" {
		err = newErrEmptyEnvVar("MINIO_SERVER")
		return
	}
	accessKeyID = os.Getenv("MINIO_ACCESSKEY")
	if accessKeyID == "" {
		err = newErrEmptyEnvVar("MINIO_ACCESSKEY")
		return
	}
	secretAccessKey = os.Getenv("MINIO_SECRETKEY")
	if secretAccessKey == "" {
		err = newErrEmptyEnvVar("MINIO_SECRETKEY")
		return
	}
	// bucket can be empty, since we're gonna generate a new bucket name if there
	// wasn't one specified.
	bucket = os.Getenv("MINIO_BUCKET")

	return
}

func main() {
	serverURI, accessKeyID, secretAccessKey, bucket, secure, err := getEnvDetails()
	if err != nil {
		log.Fatalf("An error occured while fetching environment settings: %s", err)
	}
	c, err := client.NewMinioClient(serverURI, accessKeyID, secretAccessKey, bucket, secure)
	if err != nil {
		log.Fatalf("An error occured while creating a new Minio client: %s", err)
	}
	d := driver.NewMinioDriver(c, secure)
	h := volume.NewHandler(d)
	h.ServeUnix(socketAddress, rootID)
}
