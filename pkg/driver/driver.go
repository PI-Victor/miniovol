package driver

import (
	"fmt"
	"os"
	"os/exec"
	_ "path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/client"
)

var capability volume.Capability

type minioVolume struct {
	name        string
	mountpoint  string
	connections int

	// NOTE: check to see if buckets would really collide if we specify them only
	// in the driver, instead of attaching them individually to each volume.
	bucketName string
}

// MinioDriver is the driver used by docker.
type MinioDriver struct {
	m *sync.RWMutex
	c *client.MinioClient

	volumes map[string]*minioVolume
}

// NewMinioDriver creates a new driver for the docker plugin.
func NewMinioDriver(c *client.MinioClient) MinioDriver {
	return MinioDriver{
		m: &sync.RWMutex{},
	}
}

func newVolume(name, mountPoint, bucketname string) *minioVolume {
	return &minioVolume{
		name:       name,
		mountpoint: mountPoint,
	}
}

// Create creates a new volume with the appropiate data.
func (d MinioDriver) Create(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	if err := d.createClient(r.Options); err != nil {
		return volumeResp("",
			"",
			nil,
			capability,
			fmt.Sprintf("error creating client: %s", err),
		)
	}

	v := newVolume("", "", d.c.BucketName)
	d.volumes[r.Name] = v
	volumePath := createName(volumePrefix)
	v.mountpoint = fmt.Sprintf("/mnt/miniomnt_%s", volumePath)

	return volumeResp("", "", nil, capability, "")
}

// List lists all currently available volumes.
func (d MinioDriver) List(r volume.Request) volume.Response {
	return volume.Response{}
}

// Get retrieves information about a current volume.
func (d MinioDriver) Get(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name).Error())
	}

	return volumeResp(v.mountpoint, r.Name, nil, capability, "")
}

// Remove attempts to remove a volume if it's not currently in use.
func (d MinioDriver) Remove(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Lock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name).Error())
	}

	if v.connections == 0 {
		if err := os.RemoveAll(v.mountpoint); err != nil {
			return volumeResp("", "", nil, capability, err.Error())
		}
		delete(d.volumes, r.Name)
		return volumeResp("", "", nil, capability, "")
	}

	return volumeResp("",
		"",
		nil,
		capability,
		fmt.Errorf("volume %s currently in use by container", r.Name).Error(),
	)
}

// Path returns the mount path of the current volume.
func (d MinioDriver) Path(r volume.Request) volume.Response {
	d.m.RLock()
	defer d.m.RUnlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name).Error())
	}
	return volumeResp(v.mountpoint, "", nil, capability, "")
}

// Mount tries to mount a path inside the docker volume to a minio bucket
// instance with a bucket defined.
func (d MinioDriver) Mount(r volume.MountRequest) volume.Response {
	return volume.Response{}
}

// Unmount will unmount a specified volume.
func (d MinioDriver) Unmount(r volume.UnmountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name).Error())
	}

	err := d.unmountVolume(v)
	if err != nil {
		return volumeResp("", "", nil, capability, err.Error())
	}

	if v.connections <= 1 {
		if err := d.unmountVolume(v); err != nil {
			return volumeResp("", "", nil, capability, err.Error())
		}
		v.connections = 0
		return volumeResp("", "", nil, capability, "")
	}
	v.connections--
	return volumeResp("", "", nil, capability, "")
}

// Capabilities
func (d MinioDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}

// mountVolume is a helper function for the docker interface that mounts the
// filesystem with the minfs driver.
func (d MinioDriver) mountVolume(volume *minioVolume) error {
	return nil
}

// unmountVolume is a helper function for the docker interface that unmounts
// the mounted minio bucket from the local fs.
func (d MinioDriver) unmountVolume(volume *minioVolume) error {
	return exec.Command("umount", volume.mountpoint).Run()
}

// createClient is a helper function that uses minio go bindings to instantiate
// a new session with minio's API.
func (d MinioDriver) createClient(options map[string]string) error {
	var secure bool

	server, err := checkParam("server", options)
	if err != nil {
		return err
	}
	accessKey, err := checkParam("accessKey", options)
	if err != nil {
		return err
	}
	secretKey, err := checkParam("secretKey", options)
	if err != nil {
		return err
	}
	// TODO: remember to fix this, since the user could pass false and it would
	// become true.
	_, err = checkParam("secure", options)
	if err == nil {
		secure = true
	}

	if d.c == nil {
		d.c, err = client.NewMinioClient(server, accessKey, secretKey, secure)
		if err != nil {
			return err
		}
	}

	bucketName, err := checkParam("bucket", options)
	if err != nil || bucketName == "" {
		if err = d.createBucket(bucketName); err != nil {
			return err
		}
	}

	return nil
}

// createBucket is a helper function that creates a bucket on minio to be used
// by the volume plugin to mount a minio bucket locally.
func (d MinioDriver) createBucket(bucket string) error {
	exists, err := d.c.Client.BucketExists(bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := d.c.Client.MakeBucket(bucket, location); err != nil {
			return err
		}
	}
	d.c.BucketName = bucket
	return nil
}
