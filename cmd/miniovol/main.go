package main

import (
	"flag"
	"log"
	"os"

	"github.com/golang/glog"

	"github.com/docker/go-plugins-helpers/volume"

	"github.com/cloudflavor/miniovol/pkg/driver"
)

const (
	socketAddress = "/run/docker/plugins/miniovol.sock"
	rootID        = 0
)

func main() {
	// set logging capabilities. Flush logs and set stderr as default..
	defer glog.Flush()
	flag.Set("logtostderr", "true")
	if os.Getenv("MINIOVOL_LOG_LEVEL") != "" {
		flag.Set("v", os.Getenv("MINIOVOL_LOG_LEVEL"))
	}
	flag.Parse()

	d := driver.NewMinioDriver(nil, false)
	h := volume.NewHandler(d)
	glog.V(0).Infof("Trying to serve on %s", socketAddress)
	if err := h.ServeUnix(socketAddress, rootID); err != nil {
		log.Fatalf("An error occured while trying to serve: %s", err)
	}
}
