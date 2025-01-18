package configs

import (
	"github.com/spf13/viper"
)

var (
	DefaultAuthScopes  = []string{"profile", "email", "openid", "offline_access", "User.Read", "Files.ReadWrite.AppFolder"}
	DefaultRedirectURL = "http://localhost:9999/authentication"
	AppName            = "onedrive-client"

	DBFileKey          = "db.filepath"
	AuthSecretIDKey    = "auth.secret_id"
	AuthRedirectURLKey = "auth.redirect_url"
	AuthScopesKey      = "auth.scopes"
)

func GetSecretID() string {
	return viper.GetString(AuthSecretIDKey)
}

func GetRedirectURL() string {
	return viper.GetString(AuthRedirectURLKey)
}

func GetAuthScopes() []string {
	return viper.GetStringSlice(AuthScopesKey)
}

func GetAppName() string {
	return AppName
}

func GetDBFilePath() string {
	return viper.GetString(DBFileKey)
}
