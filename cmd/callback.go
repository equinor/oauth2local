package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/equinor/oauth2local/ipc"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var breakB bool

// callbackCmd represents the callback command
var callbackCmd = &cobra.Command{
	Use:   "callback",
	Short: "Send callback url to sovereign",
	Long:  `Send callback url to sovereign`,
	Run: func(cmd *cobra.Command, args []string) {
		err := sendCallback(args)
		if err != nil {

			jww.ERROR.Println("Error...", err)
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		if breakB {
			jww.INFO.Println("sent calback", args)
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	},
}

func sendCallback(args []string) error {
	cli, err := ipc.NewClient()
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf("Only one arg supported")
	}

	err = cli.SendCallback(args[0])
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(callbackCmd)

	callbackCmd.Flags().BoolVarP(&breakB, "break", "b", false, "Break before exit")
}
