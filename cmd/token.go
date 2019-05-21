package cmd

import (
	"fmt"
	"os"

	"github.com/equinor/oauth2local/ipc"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Gets access token from the local server instance",
	Long:  `Gets access token from the local server instance"`,
	Run: func(cmd *cobra.Command, args []string) {

		cli, err := ipc.NewClient()
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}

		a, err := cli.GetAccessToken()
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}

		fmt.Println(a)
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
