package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/browser"

	"github.com/equinor/oauth2local/config"
	"github.com/equinor/oauth2local/ipc"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
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
		jww.ERROR.Println("No config file loaded")
		os.Exit(1)
	}

	if ipc.HasSovereign() {
		jww.INFO.Println("A server is already running")
		os.Exit(9)
	}

	opts := []oauth2.Option{oauth2.WithOauth2Settings(
		oauth2.Oauth2Settings{
			AuthServer:   config.AuthServer(),
			TenantID:     config.TenantID(),
			ClientID:     config.ClientID(),
			ClientSecret: config.ClientSecret(),
			ResourceID:   config.ResourceID(),
		})}
	if config.IgnoreStateCheck() {
		opts = append(opts, oauth2.WithState("none"))
	}

	var port = config.LocalHttpServerPort()
	if port > 0 && port < 65536 {
		jww.DEBUG.Println("Using local http server to receive callback")
		opts = append(opts, oauth2.WithLocalhostHttpServer(port))

	}
	var oauthHandler oauth2.Handler
	var err error
	switch config.AuthServerType() {
	case "adal":
		oauthHandler, err = oauth2.NewMsal(opts...)
	case "msal":
		oauthHandler, err = oauth2.NewAdal(opts...)
	default:
		err = fmt.Errorf("No authserver selected")
	}

	if err != nil {
		jww.ERROR.Printf("Error with oauth client: %v", err)
		os.Exit(1)
	}

	jww.INFO.Println("starting browser...")

	lpu, err := oauthHandler.LoginProviderURL()
	if err != nil {
		jww.ERROR.Printf("Login provider url isn't an url: %v", err)
		os.Exit(1)
	}

	if !testing {
		err = browser.OpenURL(lpu)
		if err != nil {
			jww.ERROR.Println(err)
			os.Exit(1)
		}
	}
	s := ipc.NewServer(oauthHandler)

	jww.ERROR.Println("Cannot serve:", s.Serve())
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(serveCmd)

}
