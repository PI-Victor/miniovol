package driver

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

func TestError(t *testing.T) {
	volName := "test"
	err := newErrVolNotFound("test")
	if err.Error() != fmt.Sprintf("volume %s not found", volName) {
		t.Errorf("Expected error to match \"volume %s not found\" got \"%s\"", volName, err.Error())
	}
}

func TestNewCfg(t *testing.T) {
	accesKey := "testKey"
	secretKey := "secretKey"
	version := "1"

	fakeCfg := &minfsCfg{
		Version:   version,
		AccessKey: accesKey,
		SecretKey: secretKey,
	}

	newCfg := newCfg(accesKey, secretKey, vers)
	if !reflect.DeepEqual(fakeCfg, newCfg) {
		t.Errorf("Expected struct to match %#v, got %#v", newCfg, fakeCfg)
	}
}

func TestCreateName(t *testing.T) {
	testPrefix := "testPrefix"
	newName := createName(testPrefix)
	if !strings.Contains(newName, testPrefix) {
		t.Errorf("Expected %s to contain \"%s\"", newName, testPrefix)
	}
}

func TestVolumeResp(t *testing.T) {
	var (
		mountPoint   = "/mnt/test"
		rName        = "test"
		err          = errors.New("new test error")
		capabilities = capability
	)

	fakeVolumesResponse := volume.Response{
		Err: err.Error(),
		Volume: &volume.Volume{
			Mountpoint: mountPoint,
			Name:       rName,
		},
	}

	volumeResponse := volumeResp(mountPoint, rName, nil, capabilities, err.Error())
	if !reflect.DeepEqual(fakeVolumesResponse, volumeResponse) {
		t.Errorf("Expected %#v to match %#v", volumeResponse, fakeVolumesResponse)
	}
}
