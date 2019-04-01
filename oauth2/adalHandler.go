package oauth2
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/equinor/oauth2local/storage"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

type AdalHandler struct {
	net          *http.Client
	tenantID     string
	appRedirect  string
	clientID     string
	clientSecret string
	handleScheme string
	store        storage.Storage
}

func NewAdalHandler(store storage.Storage) (*AdalHandler, error) {

	cli := &AdalHandler{
		net:          new(http.Client),
		tenantID:     viper.GetString("TenantID"),
		appRedirect:  viper.GetString("CustomScheme") + "://callback",
		clientID:     viper.GetString("ClientID"),
		clientSecret: viper.GetString("ClientSecret"),
		handleScheme: viper.GetString("CustomScheme"),
		store:        store}

	return cli, nil
}

func tokenURL(tenant string) string {

	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenant)
}

func (cli *AdalHandler) OpenLoginProvider() error {
	params := url.Values{}

	params.Set("redirect_uri", cli.appRedirect)
	params.Set("client_id", cli.clientID)
	params.Set("response_type", "code")
	params.Set("state", "none")
	loginURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/authorize?%s", cli.tenantID, params.Encode())
	browser.OpenURL(loginURL)
	return nil
}

func CodeFromURL(callbackURL, scheme string) (string, error) {
	u, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}

	if u.Scheme != scheme {
		return "", fmt.Errorf("App doesn't handle scheme: %s", u.Scheme)

	}
	params := u.Query()
	code := params.Get("code")

	return code, nil
}

func (cli *AdalHandler) CodeFromURL(callbackURL string) (string, error) {
	return CodeFromURL(callbackURL, cli.handleScheme)
}

func (cli *AdalHandler) GetToken(code string) (string, error) {

	params := url.Values{}
	params.Set("redirect_uri", cli.appRedirect)
	params.Set("client_id", cli.clientID)
	params.Set("client_secret", cli.clientSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("resource", cli.clientID)
	body := bytes.NewBufferString(params.Encode())

	tokenURL := tokenURL(cli.tenantID)
	resp, err := cli.net.Post(tokenURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", fmt.Errorf("Error posting to token url %s: %s ", tokenURL, err)
	}

	decoder := json.NewDecoder(resp.Body)
	var dat map[string]interface{}
	err = decoder.Decode(&dat)
	if err != nil {
		return "", err
	}

	if accessToken, ok := dat["access_token"]; ok {
		return accessToken.(string), nil
	}
	if clientError, ok := dat["error"]; ok {
		return "", fmt.Errorf("Did not receive token: %v", clientError)
	}

	return "", fmt.Errorf("Token response not valid: %v", dat)
}

func (h AdalHandler) GetAccessToken() (string, error) {

	//Check storage
	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return "", err
	}
	//Reissue to authorize if old

	return a, nil
}
func (h AdalHandler) UpdateFromRedirect(redirect *url.URL) error {

	// Decode to authorize code

	// Validate state/nonce

	//
	return nil
}

func (h AdalHandler) UpdateFromCode(code string) error {

	return nil
}
