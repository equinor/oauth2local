package cmd

import (
	"os"

	"github.com/equinor/oauth2local/register"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register app ass url handler for custom url",
	Long:  `Register app ass url handler for custom url`,
	Run: func(cmd *cobra.Command, args []string) {
		register.RegMe(viper.GetString("CustomScheme"), os.Args[0])
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
