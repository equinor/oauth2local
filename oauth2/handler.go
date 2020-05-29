package oauth2

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/equinor/oauth2local/oauth2/redirect"

	"github.com/equinor/oauth2local/storage"
	"github.com/google/uuid"
)

type Handler interface {
	GetAccessToken() (string, error)
	UpdateFromRedirect(*redirect.Params) error
	UpdateFromCode(string) error
	LoginProviderURL() (string, error)
}
type Tokens struct {
	AccessToken  string
	IDToken      string
	RefreshToken string
}

const (
	authGrant    = "authorization_code"
	refreshGrant = "refresh_token"
)

type Option interface {
	apply(interface{}) error
}

func generateCodeChallenge() string {
	return uuid.New().String() + "-" + uuid.New().String()
}

func generateSessionState() string {
	return uuid.New().String() + "-" + uuid.New().String()
}

func WithState(state string) Option {
	return newFuncOption(func(h interface{}) error {
		switch h := h.(type) {
		case *AdalHandler:
			h.sessionState = state
		case *MsalHandler:
			h.sessionState = state
		default:
			return fmt.Errorf("Not implemented for this type %v", reflect.TypeOf(h))
		}

		return nil
	})
}

func WithRenewer(renewer func()) Option {
	return newFuncOption(func(h interface{}) error {
		switch h := h.(type) {
		case *AdalHandler:
			h.renewer = renewer
		case *MsalHandler:
			h.renewer = renewer
		default:
			return fmt.Errorf("Not implemented for this type %v", reflect.TypeOf(h))
		}

		return nil
	})
}

func WithClient(client *http.Client) Option {
	return newFuncOption(func(h interface{}) error {
		switch h := h.(type) {
		case *AdalHandler:
			h.client = client
		case *MsalHandler:
			h.client = client
		default:
			return fmt.Errorf("Not implemented for this type %v", reflect.TypeOf(h))
		}

		return nil
	})
}

func WithStore(store storage.Storage) Option {
	return newFuncOption(func(h interface{}) error {
		switch h := h.(type) {
		case *AdalHandler:
			h.store = store
		case *MsalHandler:
			h.store = store
		default:
			return fmt.Errorf("Not implemented for this type %v", reflect.TypeOf(h))
		}

		return nil
	})
}

func WithOauth2Settings(o2o Oauth2Settings) Option {
	return newFuncOption(func(h interface{}) error {

		switch h := h.(type) {
		case *AdalHandler:
			h.o2o = o2o
			h.o2o.AuthServer = strings.TrimRight(h.o2o.AuthServer, "/")
			h.o2o.AuthServer = strings.TrimSpace(h.o2o.AuthServer)
			if h.o2o.ResourceID == "" {
				h.o2o.ResourceID = h.o2o.ClientID
			}
		case *MsalHandler:
			h.o2o = o2o
			h.o2o.AuthServer = strings.TrimRight(h.o2o.AuthServer, "/")
			h.o2o.AuthServer = strings.TrimSpace(h.o2o.AuthServer)
			if h.o2o.ResourceID == "" {
				h.o2o.ResourceID = h.o2o.ClientID
			}
		default:
			return fmt.Errorf("Not implemented for this type %v", reflect.TypeOf(h))
		}

		return nil
	})
}

func WithLocalhostHttpServer(port uint) Option {
	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)
	return newFuncOption(func(h interface{}) error {
		var s *redirect.Server
		switch h := h.(type) {
		case *AdalHandler:
			h.redirectURL = redirectURL
			h.scheme = "http"
			s = redirect.Init(port, h.UpdateFromRedirect)
		case *MsalHandler:
			h.redirectURL = redirectURL
			h.scheme = "http"
			s = redirect.Init(port, h.UpdateFromRedirect)
		default:
			return fmt.Errorf("Not implemented for this type %v", reflect.TypeOf(h))
		}

		if s != nil {
			go s.Serve()
		}
		return nil
	})
}
