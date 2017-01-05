package driver

import (
	"fmt"
	"log"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/client"
)

var capability volume.Capability

type minioVolume struct {
	name       string
	mountpoint string
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

func newVolume(name string, mountPoint string) *minioVolume {
	return &minioVolume{
		name:       name,
		mountpoint: mountPoint,
	}
}

// Create creates a new volume with the appropiate date.
func (d MinioDriver) Create(r volume.Request) volume.Response {
	d.m.Lock()

	defer d.m.Unlock()
	v := minioVolume{}
	log.Print(v)
	return volume.Response{}
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
		err := fmt.Sprintf("requested volume is not found: %s", r.Name)
		volumeResp("", "", nil, capability, err)
	}

	return volumeResp(v.mountpoint, r.Name, nil, capability, "")
}

// Remove deletes a volume.
func (d MinioDriver) Remove(r volume.Request) volume.Response {
	return volume.Response{}
}

// Path returns the mount path of the current volume.
func (d MinioDriver) Path(r volume.Request) volume.Response {
	d.m.RLock()
	defer d.m.RUnlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volume.Response{
			Err: fmt.Sprintf("requested volume is not found: %s", r.Name),
		}
	}
	return volumeResp(v.mountpoint, "", nil, capability, "")
}

// Mount tries to mount a path inside the docker volume to a minio bucket
// instance with a bucket defined.
func (d MinioDriver) Mount(r volume.MountRequest) volume.Response {
	return volume.Response{}
}

// Unmount will unmount a specified
func (d MinioDriver) Unmount(r volume.UnmountRequest) volume.Response {
	return volume.Response{}
}

func (d MinioDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{}
}

func (d MinioDriver) mountVolume(volume *minioVolume) error {
	return nil
}

func (d MinioDriver) unmountVolume(volume *minioVolume) error {
	return nil
}

func (d MinioDriver) createClient(r volume.Request) error {
	var secure bool
	server, err := checkValidParameter("server", r.Options)
	if err != nil {
		return err
	}

	accessKey, err := checkValidParameter("accessKey", r.Options)
	if err != nil {
		return err
	}

	secretKey, err := checkValidParameter("secretKey", r.Options)
	if err != nil {
		return err
	}

	// TODO: remember to fix this, since the user could pass false and it would
	// become true.
	_, err = checkValidParameter("secure", r.Options)
	if err == nil {
		secure = true
	}

	if d.c == nil {
		d.c, err = client.NewMinioClient(server, accessKey, secretKey, secure)
		if err != nil {
			return err
		}
	}

	bucketName, err := checkValidParameter("bucket", r.Options)
	if err != nil {
		if err = d.createBucket(bucketName); err != nil {
			return err
		}
	}
	return nil
}

func (d MinioDriver) createBucket(bucket string) error {
	exists, err := d.c.Client.BucketExists(bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := client.CreateBucket(bucket, d.c.Client); err != nil {
			return err
		}
	}
	d.c.BucketName = bucket
	return nil
}
