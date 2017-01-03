package driver

import (
	"encoding/json"
	"os"
)

const (
	cfgFile = "/etc/minfs/config.json"
	vers    = "1"
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
	cfg := newCfg(m.AccesKeyID, m.SecretAccessKey, vers)

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
