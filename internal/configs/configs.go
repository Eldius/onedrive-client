package configs

import (
	"github.com/spf13/viper"
)

var (
	DefaultAuthScopes   = []string{"profile", "email", "openid", "offline_access", "User.Read", "Files.ReadWrite.AppFolder"}
	DefaultRedirectURL  = "http://localhost:9999/authentication"
	DefaultBaseEndpoint = "https://graph.microsoft.com/v1.0"
)

func GetSecretID() string {
	return viper.GetString("auth.secret_id")
}

func GetRedirectURL() string {
	return viper.GetString("auth.redirect_url")
}

func GetAuthScopes() []string {
	return viper.GetStringSlice("auth.scopes")
}

func GetAppName() string {
	return "onedrive-client"
}
