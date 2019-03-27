package main
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zaker/oauth2local/config"
	"github.com/zaker/oauth2local/ipc"
	"github.com/zaker/oauth2local/oauth2"
	"github.com/zaker/oauth2local/register"
)

var redirectCallback = flag.String("r", "", "Handles redirect from azure ad")

func main() {
	flag.Parse()
	var cfg *config.Config
	isClient, err := ipc.HasSovereign()
	isRedirect := len(*redirectCallback) > 0
	if err != nil {
		log.Println("Sovereign check failed", err)

	}

	if !isClient {

		cfg, err = config.Load()
		if err != nil {
			log.Println("Couldn't load config", err)
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			log.Fatal(err)
		}
		// cli, _ := oauth2.NewClient(cfg)
		// cli.OpenLoginProvider()
		go ipc.StartServer()
		return

	} else {
		fmt.Printf("Press 'Enter' to continue redirect %v", isRedirect)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		if isRedirect {
			code, err := oauth2.CodeFromURL(*redirectCallback, cfg.HandleScheme)
			if err != nil {

				log.Println("Couldn't get code from url", err)
				fmt.Print("Press 'Enter' to continue...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				return
			}
			err = ipc.SendCode(code)
			if err != nil {

				log.Println("Couldn't send code to sovereign", err)
				fmt.Print("Press 'Enter' to continue...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				return
			}
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}

		return
	}
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cli, err := oauth2.NewClient(cfg)
	if err != nil {
		log.Println("Couldn't start client", err)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		log.Fatal(err)
	}

	if isRedirect {
		// log.Println("Handle redirect", *redirectCallback)
		code, err := oauth2.CodeFromURL(*redirectCallback, cfg.HandleScheme)
		if err != nil {
			log.Println("Couldn't retreive code from url", err)
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		accessToken, err := cli.GetToken(code)
		if err != nil {
			log.Println("Error parsing url", err)
			fmt.Print("Press 'Enter' to continue...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		fmt.Println("Access Token", accessToken)
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		return
	}

	register.RegMe(cfg.HandleScheme, os.Args[0])
	fmt.Println("starting browser...")
	cli.OpenLoginProvider()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return
}
