package driver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

const (
	cfgFile      = "/etc/minfs/config.json"
	vers         = "1"
	volumePrefix = "miniovol-"
	bucketPrefix = "miniobucket-"
	location     = "us-east-1"
)

// Error implements the error interface to eliminate message duplication when
// the driver checks for a specific volume
type Error struct {
	volumeName string
}

func (e Error) Error() string {
	return fmt.Sprintf("volume %s not found", e.volumeName)
}

func newErrVolNotFound(v string) error {
	return Error{
		volumeName: v,
	}
}

type minfsCfg struct {
	Version   string `json:"version"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// ProvisionConfig updates the minfs config with the Minio instance details
// (accessKeyID, secretAccessKey, serverURI)
// This is necessary for minfs to autheticate with the Minio instance.
// NOTE: move this to the driver to streamline testing?
func provisionConfig(m *MinioDriver) error {
	cfg := newCfg(m.c.AccesKeyID, m.c.SecretAccessKey, vers)

	details, err := json.Marshal(cfg)
	if err != nil {
		return nil
	}
	fh, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer fh.Close()
	fh.Write(details)
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

func volumeResp(mountPoint,
	rName string,
	volumes []*volume.Volume,
	capabilities volume.Capability,
	err error,
) volume.Response {

	return volume.Response{
		Err: err.Error(),
		Volume: &volume.Volume{
			Mountpoint: mountPoint,
			Name:       rName,
		},
		Volumes:      volumes,
		Capabilities: capabilities,
	}
}
