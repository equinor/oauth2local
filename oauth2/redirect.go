package oauth2

import "net/url"

type RedirectParams struct {
	Scheme string
	Code   string
	State  string
}

func DecodeRedirect(u *url.URL) *RedirectParams {

	params := u.Query()
	return &RedirectParams{
		Scheme: u.Scheme,
		Code:   params.Get("code"),
		State:  params.Get("state"),
	}

}
