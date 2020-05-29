package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/equinor/oauth2local/oauth2/redirect"
	"github.com/equinor/oauth2local/storage"
	jww "github.com/spf13/jwalterweatherman"
)

type AdalHandler struct {
	client        *http.Client
	o2o           Oauth2Settings
	redirectURL   string
	scheme        string
	sessionState  string
	codeChallenge string
	store         storage.Storage
	jwtParser     *jwt.Parser
	ticker        *time.Ticker
	mut           *sync.Mutex
	renewer       func()
}

func (h *AdalHandler) defaultRenewer() {
	for range h.ticker.C {
		err := h.renewTokens()
		if err != nil {
			jww.INFO.Println("Couldn't renew token", err)
		}
	}
}

func NewAdal(opts ...Option) (*AdalHandler, error) {

	dopts := &AdalHandler{
		client:        new(http.Client),
		o2o:           Oauth2Settings{},
		redirectURL:   "loc-auth://callback",
		scheme:        "loc-auth",
		sessionState:  generateSessionState(),
		codeChallenge: generateCodeChallenge(),
		store:         storage.Memory(),
		jwtParser:     new(jwt.Parser),
		ticker:        time.NewTicker(1 * time.Minute),
		mut:           new(sync.Mutex)}

	dopts.renewer = dopts.defaultRenewer

	for _, opt := range opts {
		err := opt.apply(dopts)
		if err != nil {
			return nil, fmt.Errorf("Oauth2 Settings loading %w", err)
		}
	}

	if !dopts.o2o.Valid() {
		return nil, fmt.Errorf("Oauth2 Settings is not valid")
	}

	if dopts.renewer == nil {
		return nil, fmt.Errorf("No acces token renewer defined")
	}

	go dopts.renewer()
	return dopts, nil
}

func (h *AdalHandler) renewTokens() error {

	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return err
	}

	token, _, err := h.jwtParser.ParseUnverified(a, &jwt.StandardClaims{})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {

		tokenPeriod := claims.ExpiresAt - claims.IssuedAt
		currentPeriod := claims.ExpiresAt - time.Now().Unix()

		if currentPeriod > tokenPeriod/5 {
			jww.INFO.Println("Token still in grace period")
			return nil
		}
		jww.INFO.Println("Token is out grace period")
	}
	jww.DEBUG.Println("Fetching refresh token from store")
	r, err := h.store.GetToken(storage.RefreshToken)
	if err != nil {
		return err
	}
	jww.DEBUG.Println("Updating refresh grant in store")
	err = h.updateTokens(r, refreshGrant)
	if err != nil {
		return err
	}
	return nil
}

func (h *AdalHandler) tokenURL() string {

	return fmt.Sprintf("%s/oauth2/token", h.getAuthEndpoint())
}

func (h *AdalHandler) getAuthEndpoint() string {
	if h.o2o.TenantID == "" {
		return h.o2o.AuthServer
	}
	return fmt.Sprintf("%s/%s", h.o2o.AuthServer, h.o2o.TenantID)
}

func (h *AdalHandler) LoginProviderURL() (string, error) {
	u, err := url.Parse(fmt.Sprintf("%s/oauth2/authorize", h.getAuthEndpoint()))
	if err != nil {
		return "", err
	}
	jww.DEBUG.Println("LoginProvider at:", u)
	params := u.Query()

	params.Set("redirect_uri", h.redirectURL)
	params.Set("client_id", h.o2o.ClientID)
	params.Set("response_type", "code")
	params.Set("state", h.sessionState)
	params.Set("code_challenge", h.codeChallenge)

	u.RawQuery = params.Encode()
	return u.String(), nil

}

func (h *AdalHandler) updateTokens(code, grant string) error {
	defer h.client.CloseIdleConnections()
	params := url.Values{}
	params.Set("client_id", h.o2o.ClientID)
	params.Set("client_secret", h.o2o.ClientSecret)
	params.Set("grant_type", grant)

	if grant == authGrant {
		params.Set("code_verifier", h.codeChallenge)
		params.Set("code", code)
		params.Set("redirect_uri", h.redirectURL)
	} else if grant == refreshGrant {
		params.Set("refresh_token", code)
	}
	params.Set("resource", h.o2o.ResourceID)

	tokenURL := h.tokenURL()
	jww.DEBUG.Println("Getting token from:", tokenURL)
	resp, err := h.client.PostForm(tokenURL, params)
	if err != nil {
		return fmt.Errorf("Error posting to token url %s: %s ", tokenURL, err)
	}
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("Did not receive token: %v - No body", resp.Status)
		}
		bodyString := string(body)
		return fmt.Errorf("Did not receive token: %v - %s", resp.Status, bodyString)

	}

	decoder := json.NewDecoder(resp.Body)
	var dat map[string]interface{}
	err = decoder.Decode(&dat)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
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

func (h *AdalHandler) GetAccessToken() (string, error) {

	a, err := h.store.GetToken(storage.AccessToken)
	if err != nil {
		return "", err
	}

	return a, nil
}
func (h *AdalHandler) UpdateFromRedirect(rp *redirect.Params) error {

	if rp.State != h.sessionState {
		return errors.New("Invalid state in redirect")
	}

	if rp.Scheme != h.scheme {
		return errors.New("Invalid scheme in redirect")
	}

	h.mut.Lock()
	defer h.mut.Unlock()

	err := h.updateTokens(rp.Code, authGrant)
	if err != nil {
		return err
	}

	return nil
}

func (h *AdalHandler) UpdateFromCode(code string) error {
	h.mut.Lock()
	defer h.mut.Unlock()
	return fmt.Errorf("Not implemented")
}
