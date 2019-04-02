package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/equinor/oauth2local/ipc"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/storage"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve a local auth provider",
	Long:  `serve a local auth provider.`,
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
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
	s := ipc.NewServer(*cli)

	log.Fatalf("Cannot serve: %v", s.Serve())
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
