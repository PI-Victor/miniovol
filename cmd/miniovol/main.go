package main

import (
	"log"
	
	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/driver"
)

const (
	socketAddress = "/run/docker/plugins/miniovol.sock"
	rootID        = 0
)

func main() {
	d := driver.NewMinioDriver(nil)
	h := volume.NewHandler(d)
	log.Printf("Trying to listen on socket %s", socketAddress)
	err := h.ServeUnix(socketAddress, rootID)
	if err != nil {
		log.Fatalf("An error occured while trying to serve through socket: %v", err)
	}
}
