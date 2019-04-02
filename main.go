package main

import (
	"github.com/equinor/oauth2local/cmd"
)

// var redirectCallback = flag.String("r", "", "Handles redirect from azure ad")

// func runClient(cfg *config.Config) error {
// 	ipcClient, err := ipc.NewClient()
// 	if err != nil {
// 		return err
// 	}
// 	defer ipcClient.Close()
// 	log.Println("Running client", cfg.ClientType())
// 	switch cfg.ClientType() {
// 	case config.Redirect:
// 		err := ipcClient.SendCallback(cfg.RedirectURL)
// 		if err != nil {

// 			return err
// 		}
// 	case config.AccessToken:
// 		a, err := ipcClient.GetAccessToken()
// 		if err != nil {

// 			fmt.Print(a)
// 		}
// 	}
// 	return nil
// }

func main() {
	cmd.Execute()
	// return
	// isClient := ipc.HasSovereign()

	// flag.Parse()

	// cfg := config.Init(isClient)
	// cfg.RedirectURL = *redirectCallback
	// if cfg.AsClient() {
	// 	err := runClient(cfg)
	// 	if err != nil {

	// 		fmt.Println("Error...", err)
	// 		bufio.NewReader(os.Stdin).ReadBytes('\n')
	// 	}
	// 	fmt.Println("Success...", err)
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// } else {

	// }
}
