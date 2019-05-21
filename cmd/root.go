package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "oauth2local",
	Short: "oauth2local is providing oauth2 authenticated tokens to local processes",
	Long:  "oauth2local is providing oauth2 authenticated tokens to local processes",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oauth2local.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "log to console to console")

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".oauth2local")
	}

	if verbose {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		jww.ERROR.Println(err)
	}

}
