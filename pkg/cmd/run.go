package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	_ "github.com/cloudflavor/miniovol/pkg/client"
	_ "github.com/cloudflavor/miniovol/pkg/driver"
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
	Run: func(cmd *cobra.Command, args []string) {

		if os.Getenv("MINIO_SERVER_URI") == "" && serverURI == "" ||
			os.Getenv("MINIO_ACCESS_KEY") == "" && accessKeyID == "" ||
			os.Getenv("MINIO_SECRET_KEY") == "" && secretAccessKey == "" {
			cmd.Help()
			log.Fatal("server, accesskeyID and secretKey are mandatory parameters.")
		}
		// c := client.NewMinioClient(serverURI, accessKeyID, secretAccessKey, secure)
	},
}

func init() {
	RunCmd.PersistentFlags().StringVar(&serverURI, "server", "", "Specify the Minio server URI.")
	RunCmd.PersistentFlags().StringVar(&accessKeyID, "accessKeyID", "", "Specify your Minio Access Key.")
	RunCmd.PersistentFlags().StringVar(&bucket, "bucket", "docker-volumes", "Specify the name of the bucket to be created.")
	RunCmd.PersistentFlags().StringVar(&secretAccessKey, "secretKey", "", "Specify your Minio secret key.")
	RunCmd.PersistentFlags().BoolVar(&secure, "secure", false, "Specify true to use a secure connection.")
}
