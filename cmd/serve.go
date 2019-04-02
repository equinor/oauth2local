package cmd

import (
	"fmt"
	"log"

	"github.com/equinor/oauth2local/ipc"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve a local auth provider",
	Long:  `serve a local auth provider.`,
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {

	if viper.ConfigFileUsed() == "" {
		log.Fatal("No config file loaded")
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
	if ipc.HasSovereign() {
		log.Println("A server is already running")
		return
	}

	cli, err := oauth2.NewAdalHandler(storage.Memory())
	if err != nil {
		log.Printf("Error with oauth client: %v", err)
		return
	}

	fmt.Println("starting browser...")
	cli.OpenLoginProvider()
	s := ipc.NewServer(cli)

	log.Fatalf("Cannot serve: %v", s.Serve())
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
