package cmd

import (
	"log"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/spf13/cobra"

	"github.com/cloudflavor/miniovol/pkg/client"
	"github.com/cloudflavor/miniovol/pkg/driver"
)

const (
	socketAddress = "/run/docker/plugins/miniovol.sock"
	rootID        = 0
)

var (
	serverURI       string
	bucket          string
	accessKeyID     string
	secretAccessKey string
	secure          bool
)

// RunCmd is responsible for starting the plugin and connecting to a Minio
// instance.
var RunCmd = &cobra.Command{
	Use:   "miniovol",
	Short: "Minio Volume plugin for Docker",
	Long: `
You can use environment variables to connect to your Minio instance.
e.g.:
export MINIO_SERVER_URI="localhost:9000"
export MINIO_ACCESS_KEY="myAccessKey"
export MINIO_SECRET_KEY="mySecretKey"
`,
	Run: func(cmd *cobra.Command, args []string) {

		// if os.Getenv("MINIO_SERVER_URI") == "" && serverURI == "" ||
		// 	os.Getenv("MINIO_ACCESS_KEY") == "" && accessKeyID == "" ||
		// 	os.Getenv("MINIO_SECRET_KEY") == "" && secretAccessKey == "" {
		//
		// 	log.Printf("server, accesskeyID and secretKey are mandatory parameters.\n")
		// 	cmd.Help()
		//
		// 	// Use os.Exit() once, only here so that we can print the help when there's
		// 	// an error.
		// 	//os.Exit(1)
		// }
		// ugly af.
		if serverURI == "" {
			serverURI = os.Getenv("MINIO_SERVER_URI")
		}
		if accessKeyID == "" {
			accessKeyID = os.Getenv("MINIO_ACCESS_KEY")
		}
		if secretAccessKey == "" {
			secretAccessKey = os.Getenv("MINIO_SECRET_KEY")
		}

		c, err := client.NewMinioClient(serverURI, accessKeyID, secretAccessKey, bucket, secure)
		log.Printf("%s, %s, %s", serverURI, accessKeyID, secretAccessKey)
		if err != nil {
			log.Fatalf("An error occured while connecting to the Minio instance: %v", err)
		}
		log.Printf("this is the client %#v", c)
		d := driver.NewMinioDriver(c, secure)
		h := volume.NewHandler(d)
		log.Printf("Trying to listen on socket %s", socketAddress)
		err = h.ServeUnix(socketAddress, rootID)
		if err != nil {
			log.Fatalf("An error occured while trying to serve through socket: %v", err)
		}
	},
}

func init() {
	RunCmd.PersistentFlags().StringVar(&serverURI, "server", serverURI, "Specify the Minio server URI")
	RunCmd.PersistentFlags().StringVar(&accessKeyID, "accessKeyID", accessKeyID, "Specify your Minio Access Key")
	RunCmd.PersistentFlags().StringVar(&bucket, "bucket", "docker-volumes", "Specify the name of the Minio Bucket to be created")
	RunCmd.PersistentFlags().StringVar(&secretAccessKey, "secretKey", secretAccessKey, "Specify your Minio secret key")
	RunCmd.PersistentFlags().BoolVar(&secure, "secure", secure, "Set to true to use a secure connection")
}
