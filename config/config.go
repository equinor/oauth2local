package config

import (
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetEnvPrefix("O2L")

	viper.SetDefault("CUSTOM_SCHEME", "loc-auth")
	viper.SetDefault("AUTH_SERVER", "https://login.example.com/")
	viper.SetDefault("AUTH_TYPE", "adal")
	viper.SetDefault("TENANT_ID", "")
	viper.SetDefault("CLIENT_ID", "")
	viper.SetDefault("CLIENT_SECRET", "")
	viper.SetDefault("RESOURCE_ID", "")
	viper.SetDefault("SCOPES", "")
}

func AuthServer() string {
	viper.RegisterAlias("AuthServer", "AUTH_SERVER")
	return viper.GetString("AUTH_SERVER")
}

func AuthServerType() string {
	viper.RegisterAlias("AuthType", "AUTH_TYPE")
	return viper.GetString("AUTH_TYPE")
}

func TenantID() string {
	viper.RegisterAlias("TenantID", "TENANT_ID")
	return viper.GetString("TENANT_ID")
}

func ClientID() string {
	viper.RegisterAlias("ClientID", "CLIENT_ID")
	return viper.GetString("CLIENT_ID")
}

func ClientSecret() string {
	viper.RegisterAlias("ClientSecret", "CLIENT_SECRET")
	return viper.GetString("CLIENT_SECRET")
}

func ResourceID() string {
	viper.RegisterAlias("ResourceID", "RESOURCE_ID")
	return viper.GetString("RESOURCE_ID")
}

func Scopes() []string {
	return viper.GetStringSlice("SCOPES")
}

func CustomScheme() string {
	viper.RegisterAlias("CustomScheme", "CUSTOM_SCHEME")
	return viper.GetString("CUSTOM_SCHEME")
}
