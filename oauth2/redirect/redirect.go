package redirect

import "net/url"

type Params struct {
	Scheme string
	Code   string
	State  string
}

func DecodeRedirect(u *url.URL) *Params {

	params := u.Query()
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return &Params{
		Scheme: u.Scheme,
		Code:   params.Get("code"),
		State:  params.Get("state"),
	}

}
