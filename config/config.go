package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config for oauth local service
type Config struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	AppRedirect  string
	HandleScheme string
	IsServer     bool
}

const handleURLScheme = "loc-auth"

func redirectURL() string {
	return handleURLScheme + "://callback"
}

// Load config for app
func Load() (*Config, error) {
	err := godotenv.Load("C:\\bin\\.env")
	if err != nil {
		return nil, err
	}
	c := &Config{HandleScheme: handleURLScheme, AppRedirect: redirectURL(), IsServer: true}
	tenant := os.Getenv("AAD_TENANT_ID")
	if len(tenant) == 0 {
		return nil, fmt.Errorf("No tenant id => AAD_TENANT_ID")
	}
	c.TenantID = tenant

	clientID := os.Getenv("AAD_CLIENT_ID")
	if len(clientID) == 0 {
		return nil, fmt.Errorf("No client id => AAD_CLIENT_ID")
	}
	c.ClientID = clientID

	clientSecret := os.Getenv("AAD_CLIENT_SECRET")
	if len(clientSecret) == 0 {
		return nil, fmt.Errorf("No client secret => AAD_CLIENT_SECRET")
	}
	c.ClientSecret = clientSecret
	return c, nil
}
