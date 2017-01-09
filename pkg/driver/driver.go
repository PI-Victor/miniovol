package driver

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	server    string
	accessKey string
	secretKey string
	secure    bool
	volumes   map[string]*minioVolume
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
		bucketName: bucketname,
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
			fmt.Errorf("error creating client: %s", err),
		)
	}
	log.Printf("Got here!: %#v", d)
	volumePath := createName(volumePrefix)
	volumeMount := filepath.Join("/mnt/miniomnt_%s", volumePath)
	if err := d.createVolumeMount(volumeMount); err != nil {
		return volumeResp("", "", nil, capability, err)
	}

	volumeName := createName(volumePrefix)
	v := newVolume(volumeName, volumeMount, d.c.BucketName)
	d.volumes[r.Name] = v

	return volumeResp("", "", nil, capability, nil)
}

// List lists all currently available volumes.
func (d MinioDriver) List(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	var vols []*volume.Volume
	for name, v := range d.volumes {
		vols = append(vols,
			&volume.Volume{
				Name:       name,
				Mountpoint: v.mountpoint,
			})
	}
	return volumeResp("", "", vols, capability, nil)
}

// Get retrieves information about a current volume.
func (d MinioDriver) Get(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name))
	}

	return volumeResp(v.mountpoint, r.Name, nil, capability, nil)
}

// Remove attempts to remove a volume if it's not currently in use.
func (d MinioDriver) Remove(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Lock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name))
	}

	if v.connections == 0 {
		if err := os.RemoveAll(v.mountpoint); err != nil {
			return volumeResp("", "", nil, capability, err)
		}
		delete(d.volumes, r.Name)
		return volumeResp("", "", nil, capability, nil)
	}

	return volumeResp("",
		"",
		nil,
		capability,
		fmt.Errorf("volume %s currently in use by container", r.Name),
	)
}

// Path returns the mount path of the current volume.
func (d MinioDriver) Path(r volume.Request) volume.Response {
	d.m.RLock()
	defer d.m.RUnlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name))
	}
	return volumeResp(v.mountpoint, "", nil, capability, nil)
}

// Mount tries to mount a path inside the docker volume to a minio bucket
// instance with a bucket defined.
func (d MinioDriver) Mount(r volume.MountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name))
	}

	if err := d.mountVolume(v); err != nil {
		return volumeResp("", "", nil, capability, err)
	}
	// if the mount was successful, then increment the number of connections we
	// have to the mount.
	v.connections++

	return volumeResp("", "", nil, capability, nil)
}

// Unmount will unmount a specified volume.
func (d MinioDriver) Unmount(r volume.UnmountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	v, exists := d.volumes[r.Name]
	if !exists {
		return volumeResp("", "", nil, capability, newErrVolNotFound(r.Name))
	}

	if v.connections <= 1 {
		if err := d.unmountVolume(v); err != nil {
			return volumeResp("", "", nil, capability, err)
		}
		v.connections = 0
		return volumeResp("", "", nil, capability, nil)
	}
	v.connections--
	return volumeResp("", "", nil, capability, nil)
}

// Capabilities returns the capabilities needed for this plugin.
func (d MinioDriver) Capabilities(r volume.Request) volume.Response {
	localCapability := volume.Capability{
		Scope: "local",
	}
	return volumeResp("", "", nil, localCapability, nil)
}

// mountVolume is a helper function for the docker interface that mounts the
// filesystem with the minfs driver.
func (d MinioDriver) mountVolume(volume *minioVolume) error {
	minioPath := fmt.Sprintf("%s/%s", d.server, volume.bucketName)
	cmd := fmt.Sprintf("mount -t minfs %s %s", volume.mountpoint, minioPath)

	return exec.Command("sh", "-c", cmd).Run()
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
	d.server = server

	accessKey, err := checkParam("accessKey", options)
	if err != nil {
		return err
	}
	d.accessKey = accessKey

	secretKey, err := checkParam("secretKey", options)
	if err != nil {
		return err
	}
	d.secretKey = secretKey

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

	// TODO: implement reusability of a bucket by passing its name as a parameter.
	bucketName, err := checkParam("bucket", options)
	if err != nil || bucketName == "" {
		if err = d.createBucket(); err != nil {
			return err
		}
	}
	d.c.BucketName = bucketName
	return nil
}

// createBucket is a helper function that creates a bucket on minio to be used
// by the volume plugin to mount a minio bucket locally.
func (d MinioDriver) createBucket() error {
	bucket := createName(bucketPrefix)
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

func (d MinioDriver) createVolumeMount(volumeName string) error {
	if _, err := os.Stat(volumeName); os.IsNotExist(err) {
		if err = os.MkdirAll(volumeName, 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
