package driver

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

const (
	cfgFile      = "/etc/minfs/config.json"
	cfgDir       = "/etc/minfs/"
	vers         = "1"
	volumePrefix = "miniovol-"
	bucketPrefix = "miniobucket-"
	location     = "us-east-1"
)

type minfsCfg struct {
	Version   string `json:"version"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// EnvVarError implements the error interface and uses it to return a generic
// validation error for environment variables that are mandatory.
type EnvVarError struct {
	envVar string
}

func (e EnvVarError) Error() string {
	return fmt.Sprintf("environment variable %s cannot be empty", e.envVar)
}

func newErrEmptyEnvVar(v string) error {
	return EnvVarError{
		envVar: v,
	}
}

// VolumeError implements the error interface to eliminate message duplication
// when the driver checks for a specific volume
type VolumeError struct {
	volumeName string
}

func (e VolumeError) Error() string {
	return fmt.Sprintf("volume %s not found", e.volumeName)
}

func newErrVolNotFound(v string) error {
	return VolumeError{
		volumeName: v,
	}
}

// ProvisionConfig updates the minfs config with the Minio instance details
// (accessKeyID, secretAccessKey, serverURI)
// This is necessary for minfs to autheticate with the Minio instance.
// NOTE: move this to the driver to streamline testing?
// NOTE: if the API is correct, it should be possible to do this via env vars.
func provisionConfig(m *MinioDriver) error {
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cfgDir, 0755); err != nil {
			glog.V(1).Infof("Error while creating MinFS config dir: %s", err)
			return err
		}
	} else if err != nil {
		glog.V(1).Infof("Error while writing MinFS config: %s", err)
		return err
	}

	details := fmt.Sprintf(`{"version":"%s","accessKey":"%s","secretKey":"%s"}`,
		vers,
		m.accessKey,
		m.secretKey,
	)

	fh, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	defer fh.Close()

	fh.WriteString(details)
	if err != nil {
		return err
	}

	return nil
}

func newCfg(accessKey, secretKey, version string) *minfsCfg {
	return &minfsCfg{
		Version:   version,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func createName(prefix string) string {
	return fmt.Sprintf("%s%08x", prefix, rand.Uint32())
}

func checkParam(param string, opts map[string]string) (string, error) {
	stringParam, exists := opts[param]
	if stringParam == "" || !exists {
		return "", fmt.Errorf("%s option is required", param)
	}
	return stringParam, nil
}

func volumeResp(mountPoint, rName string, volumes []*volume.Volume, capabilities volume.Capability, err string) volume.Response {
	return volume.Response{
		Err: err,
		Volume: &volume.Volume{
			Mountpoint: mountPoint,
			Name:       rName,
		},
		Volumes:      volumes,
		Capabilities: capabilities,
	}
}
