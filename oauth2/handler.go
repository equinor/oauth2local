package oauth2

import "net/url"

type Handler interface {
	GetAccessToken() (string, error)
	UpdateFromRedirect(*url.URL) error
	UpdateFromCode(string) error
}
