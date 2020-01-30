package config

import (
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetEnvPrefix("O2L")

	viper.RegisterAlias("CustomScheme", "CUSTOM_SCHEME")
	viper.SetDefault("CUSTOM_SCHEME", "loc-auth")

	viper.RegisterAlias("AuthServer", "AUTH_SERVER")
	viper.SetDefault("AUTH_SERVER", "https://login.example.com/")

	viper.RegisterAlias("AuthType", "AUTH_TYPE")
	viper.SetDefault("AUTH_TYPE", "adal")

	viper.RegisterAlias("TenantID", "TENANT_ID")
	viper.SetDefault("TENANT_ID", "")

	viper.RegisterAlias("ClientID", "CLIENT_ID")
	viper.SetDefault("CLIENT_ID", "")

	viper.RegisterAlias("ClientSecret", "CLIENT_SECRET")
	viper.SetDefault("CLIENT_SECRET", "")

	viper.RegisterAlias("ResourceID", "RESOURCE_ID")
	viper.SetDefault("RESOURCE_ID", "")
}

func AuthServer() string {
	return viper.GetString("AUTH_SERVER")
}
func TenantID() string {
	return viper.GetString("TENANT_ID")
}
func ClientID() string {
	return viper.GetString("CLIENT_ID")
}
func ClientSecret() string {
	return viper.GetString("CLIENT_SECRET")
}
func ResourceID() string {
	return viper.GetString("RESOURCE_ID")
}

func Scopes() string {
	return viper.GetString("SCOPES")
}
func CustomScheme() string {
	return viper.GetString("CUSTOM_SCHEME")
}
func AuthServerType() string {
	return viper.GetString("AUTH_TYPE")
}
