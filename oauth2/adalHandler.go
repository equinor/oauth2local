package oauth2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/equinor/oauth2local/storage"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

type AdalHandler struct {
	net           *http.Client
	tenantID      string
	appRedirect   string
	clientID      string
	clientSecret  string
	handleScheme  string
	codeChallenge string
	store         storage.Storage
	jwtParser     *jwt.Parser
	ticker        *time.Ticker
}

const (
	authGrant    = "authorization_code"
	refreshGrant = "refresh_token"
)

func NewAdalHandler(store storage.Storage) (AdalHandler, error) {

	h := AdalHandler{
		net:           new(http.Client),
		tenantID:      viper.GetString("TenantID"),
		appRedirect:   viper.GetString("CustomScheme") + "://callback",
		clientID:      viper.GetString("ClientID"),
		clientSecret:  viper.GetString("ClientSecret"),
		handleScheme:  viper.GetString("CustomScheme"),
		codeChallenge: uuid.New().String() + "-" + uuid.New().String(),
		store:         store,
		jwtParser:     new(jwt.Parser),
		ticker:        time.NewTicker(1 * time.Minute)}

	go func() {
		for range h.ticker.C {
			err := h.renewTokens()
			if err != nil {
				log.Println("Couldn't renew token", err)
			}
		}
	}()
	return h, nil
}

func (h AdalHandler) renewTokens() error {

	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return err
	}

	token, _, err := h.jwtParser.ParseUnverified(a, &jwt.StandardClaims{})
	//Reissue to authorize if old
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {

		tokenPeriod := claims.ExpiresAt - claims.IssuedAt
		currentPeriod := claims.ExpiresAt - time.Now().Unix()

		if currentPeriod > tokenPeriod/5 {
			log.Println("Token still in grace period")
			return nil
		}

	}

	r, err := h.store.GetToken(storage.RefreshToken)
	if err != nil {
		return err
	}
	err = h.updateTokens(r, refreshGrant)
	if err != nil {
		return err
	}
	return nil
}

func tokenURL(tenant string) string {

	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenant)
}

func (h AdalHandler) OpenLoginProvider() error {
	params := url.Values{}

	params.Set("redirect_uri", h.appRedirect)
	params.Set("client_id", h.clientID)
	params.Set("response_type", "code")
	params.Set("state", "none")
	params.Set("code_challenge", h.codeChallenge)
	loginURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/authorize?%s", h.tenantID, params.Encode())
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

func (h AdalHandler) CodeFromURL(callbackURL string) (string, error) {
	return CodeFromURL(callbackURL, h.handleScheme)
}

func (h AdalHandler) updateTokens(code, grant string) error {

	params := url.Values{}
	params.Set("client_id", h.clientID)
	params.Set("client_secret", h.clientSecret)
	params.Set("grant_type", grant)

	if grant == authGrant {
		params.Set("code_verifier", h.codeChallenge)
		params.Set("code", code)
		params.Set("redirect_uri", h.appRedirect)
	} else if grant == refreshGrant {
		params.Set("refresh_token", code)
	}
	params.Set("resource", h.clientID)
	body := bytes.NewBufferString(params.Encode())

	tokenURL := tokenURL(h.tenantID)
	resp, err := h.net.Post(tokenURL, "application/x-www-form-urlencoded", body)
	if err != nil {
		return fmt.Errorf("Error posting to token url %s: %s ", tokenURL, err)
	}
	if resp.StatusCode != 200 {

		return fmt.Errorf("Did not receive token: %v", resp.Status)

	}

	decoder := json.NewDecoder(resp.Body)
	var dat map[string]interface{}
	err = decoder.Decode(&dat)
	if err != nil {
		return err
	}

	if t, ok := dat["access_token"]; ok {

		err = h.store.SetToken(storage.AccessToken, t.(string))
		if err != nil {
			return err
		}
	}
	if t, ok := dat["id_token"]; ok {
		err = h.store.SetToken(storage.IDToken, t.(string))
		if err != nil {
			return err
		}
	}
	if t, ok := dat["refresh_token"]; ok {
		err = h.store.SetToken(storage.RefreshToken, t.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h AdalHandler) getValidAccessToken() (string, error) {
	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return "", err
	}

	token, _, err := h.jwtParser.ParseUnverified(a, &jwt.StandardClaims{})
	//Reissue to authorize if old
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return a, nil
	}

	return "", err
}

func (h AdalHandler) GetAccessToken() (string, error) {

	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return "", err
	}

	return a, nil
}
func (h AdalHandler) UpdateFromRedirect(redirect *url.URL) error {

	// TODO: Validate state/nonce
	// Decode to authorize code
	c, err := h.CodeFromURL(redirect.String())
	if err != nil {
		return err
	}

	err = h.updateTokens(c, authGrant)
	if err != nil {
		return err
	}

	return nil
}

func (h AdalHandler) UpdateFromCode(code string) error {

	return nil
}
