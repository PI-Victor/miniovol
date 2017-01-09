package main

import (
	"log"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/driver"
)

const (
	socketAddress = "/run/docker/plugins/miniovol.sock"
	rootID        = 0
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	d := driver.NewMinioDriver(nil)
	h := volume.NewHandler(d)
	log.Printf("Listening on socket %s...", socketAddress)
	err := h.ServeUnix(socketAddress, rootID)
	if err != nil {
		log.Fatalf("Error while trying to serve through socket: %v", err)
	}
}
