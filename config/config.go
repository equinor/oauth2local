package config

import (
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetDefault("CustomScheme", "loc-auth")
	viper.SetDefault("Authserver", "https://login.example.com/")
	viper.SetDefault("AuthType", "adal")
	viper.SetDefault("TenantID", "")
	viper.SetDefault("ClientID", "")
	viper.SetDefault("ClientSecret", "")
	viper.SetDefault("ResourceID", "")
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
