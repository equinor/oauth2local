package oauth2

import "strings"

type Oauth2Settings struct {
	TenantID     string
	AuthServer   string
	ClientID     string
	ClientSecret string
	ResourceID   string
}

func (o2o Oauth2Settings) Valid() bool {

	if o2o.AuthServer == "" {
		return false
	}
	strings.TrimRight(o2o.AuthServer, "/")

	if o2o.ClientID == "" {
		return false
	}

	if o2o.ClientSecret == "" {
		return false
	}

	if o2o.ResourceID == "" {
		return false
	}
	return true
}
