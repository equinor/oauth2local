package oauth2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/browser"

	"github.com/zaker/oauth2local/config"
)

type Client struct {
	cfg *config.Config
	net *http.Client
}

func NewClient(cfg *config.Config) (*Client, error) {
	cli := &Client{cfg: cfg, net: new(http.Client)}

	return cli, nil
}

func tokenURL(tenant string) string {

	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenant)
}

func (cli *Client) OpenLoginProvider() error {
	params := url.Values{}

	params.Set("redirect_uri", cli.cfg.AppRedirect)
	params.Set("client_id", cli.cfg.ClientID)
	params.Set("response_type", "code")
	params.Set("state", "none")
	loginURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/authorize?%s", cli.cfg.TenantID, params.Encode())
	browser.OpenURL(loginURL)
	return nil
}

func (cli *Client) CodeFromCallback(callbackURL string) (string, error) {
	u, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}

	if u.Scheme != cli.cfg.HandleScheme {
		log.Println("We dont handle", u.Scheme)
		time.Sleep(time.Second * 3)
		return "", err
	}
	params := u.Query()
	code := params.Get("code")

	return code, nil
}

func (cli *Client) GetToken(code string) (string, error) {

	params := url.Values{}

	params.Set("redirect_uri", cli.cfg.AppRedirect)
	params.Set("client_id", cli.cfg.ClientID)
	params.Set("client_secret", cli.cfg.ClientSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("resource", cli.cfg.ClientID)
	body := bytes.NewBufferString(params.Encode())

	tokenURL := tokenURL(cli.cfg.TenantID)

	resp, err := cli.net.Post(tokenURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", fmt.Errorf("Error posting to token url %s: %s ", tokenURL, err)
	}
	log.Println("Got token", resp.Body)
	decoder := json.NewDecoder(resp.Body)
	var dat map[string]interface{}
	err = decoder.Decode(&dat)
	if err != nil {
		return "", err
	}

	accessToken := dat["access_token"].(string)
	return accessToken, nil
}
