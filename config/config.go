package config

import (
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetDefault("CustomScheme", "loc-auth")
	viper.SetDefault("Authserver", "adal")
	viper.SetDefault("AuthType", "adal")
	viper.SetDefault("TenantID", "tenant id")
	viper.SetDefault("ClientID", "client id")
	viper.SetDefault("ClientSecret", "client secret")
	viper.SetDefault("ResourceID", "resource id")
}

func AuthServer() string {
	return viper.GetString("Authserver")
}
func TenantID() string {
	return viper.GetString("TenantID")
}
func ClientID() string {
	return viper.GetString("ClientID")
}
func ClientSecret() string {
	return viper.GetString("ClientSecret")
}
func ResourceID() string {
	return viper.GetString("ResourceID")
}
func CustomScheme() string {
	return viper.GetString("CustomScheme")
}
func AuthServerType() string {
	return viper.GetString("AuthType")
}
