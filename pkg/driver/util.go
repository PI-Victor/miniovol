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
)

type minfsCfg struct {
	Version   string `json:"version"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// ProvisionConfig updates the minfs config with the Minio instance details
// (accessKeyID, secretAccessKey, serverURI)
// This is necessary for minfs to autheticate with the Minio instance.
func ProvisionConfig(m *MinioDriver) error {
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

func checkValidParameter(param string, opts map[string]string) (string, error) {
	stringParam, exists := opts[param]
	if stringParam == "" || !exists {
		return "", fmt.Errorf("%s option is required", param)
	}
	return stringParam, nil
}

func respError(errMsg string) volume.Response {
	return volume.Response{
		Err: errMsg,
	}
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
