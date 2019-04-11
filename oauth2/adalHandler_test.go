package oauth2

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/equinor/oauth2local/storage"
)

var testSettings = Oauth2Settings{
	AuthServer:   "https://example.com/",
	TenantID:     "comon",
	ClientID:     "clientid",
	ClientSecret: "secret",
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestAdalHandler_UpdateFromRedirect(t *testing.T) {
	type args struct {
		redirect *url.URL
	}

	testTokenCli := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		wantURL := "https://example.com/comon/oauth2/token"
		gotURL := req.URL.String()
		if gotURL != wantURL {
			t.Errorf("Not creating the correct token endpoint: got = %v , want = %v", gotURL, wantURL)
		}
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"access_token":"newAccessToken","id_token":"newIdToken","refresh_token":"newRefreshToken"}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	testStore := storage.Memory()
	h, err := NewAdal(
		WithOauth2Settings(testSettings),
		WithClient(testTokenCli),
		WithState("none"),
		WithStore(testStore),
	)

	if err != nil {
		t.Errorf("Failed creating handler %v", err)
	}

	redir, err := url.Parse("loc-auth://callback?state=none")
	failScheme, err := url.Parse("loki-auth://callback?state=none")
	noState, err := url.Parse("loki-auth://callback")
	if err != nil {
		t.Errorf("Couldn't parse url %v", err)
	}
	tests := []struct {
		name    string
		h       Handler
		args    args
		wantErr bool
	}{
		{name: "Update tokens", h: h, args: args{redir}, wantErr: false},
		{"Fail update ", h, args{failScheme}, true},
		{"No state in return", h, args{noState}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.UpdateFromRedirect(tt.args.redirect); (err != nil) != tt.wantErr {
				t.Errorf("AdalHandler.UpdateFromRedirect() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}

	gotToken, err := testStore.GetToken(storage.AccessToken)
	if err != nil {
		t.Errorf("Couldn't retreive token from store %v", err)
	}
	wantToken := "newAccessToken"
	if gotToken != wantToken {
		t.Errorf("Token wasn't set, got = %v, want %v", gotToken, wantToken)
	}
}

func TestAdalHandler_GetAccessToken(t *testing.T) {

	errCli := NewTestClient(func(req *http.Request) *http.Response {
		t.Error("Shouldn't do a http request")

		return &http.Response{
			StatusCode: 404,
		}
	})
	testStore := storage.Memory()
	testStore.SetToken(storage.AccessToken, "accessToken")
	h, err := NewAdal(
		WithOauth2Settings(testSettings),
		WithClient(errCli),
		WithState("none"),
		WithStore(testStore),
	)

	if err != nil {
		t.Errorf("Failed creating handler %v", err)
	}

	tests := []struct {
		name    string
		h       Handler
		want    string
		wantErr bool
	}{
		{"Fetch token from internal store", h, "accessToken", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.GetAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdalHandler.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AdalHandler.GetAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdalHandler_UpdateFromCode(t *testing.T) {
	errCli := NewTestClient(func(req *http.Request) *http.Response {
		t.Error("Shouldn't do a http request")

		return &http.Response{
			StatusCode: 404,
		}
	})
	h, err := NewAdal(
		WithOauth2Settings(testSettings),
		WithClient(errCli),
		WithState("test"),
	)

	if err != nil {
		t.Errorf("Failed creating handler %v", err)
	}

	type args struct {
		code string
	}
	tests := []struct {
		name    string
		h       Handler
		args    args
		wantErr bool
	}{
		{"Not implemented", h, args{""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.UpdateFromCode(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("AdalHandler.UpdateFromCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
