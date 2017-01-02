package driver

import (
	"sync"

	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/client"
)

// MinioDriver is the driver used by docker.
type MinioDriver struct {
	*sync.Mutex
	*client.MinioClient
}

func NewMinioDriver(c *client.MinioClient) MinioDriver {
	return MinioDriver{
		&sync.Mutex{},
		c,
	}
}

func (d MinioDriver) Create(volume.Request) volume.Response {

	return volume.Response{}
}
func (d MinioDriver) List(volume.Request) volume.Response {
	return volume.Response{}
}
func (d MinioDriver) Get(volume.Request) volume.Response {
	return volume.Response{}
}
func (d MinioDriver) Remove(volume.Request) volume.Response {
	return volume.Response{}
}
func (d MinioDriver) Path(volume.Request) volume.Response {
	return volume.Response{}
}
func (d MinioDriver) Mount(volume.MountRequest) volume.Response {
	return volume.Response{}
}
func (d MinioDriver) Unmount(volume.UnmountRequest) volume.Response {
	return volume.Response{}
}
func (d MinioDriver) Capabilities(volume.Request) volume.Response {
	return volume.Response{}
}
