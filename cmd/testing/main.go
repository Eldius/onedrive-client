package main

import (
	"fmt"
	cfg "github.com/eldius/initial-config-go/configs"
	"github.com/eldius/initial-config-go/setup"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/configs"
	"github.com/spf13/pflag"
	"log/slog"
)

func init() {
}

func main() {
	configFile := pflag.String("config", "", "config file path")

	if err := setup.InitSetup(
		configs.GetAppName(),
		setup.WithConfigFileToBeUsed(*configFile),
		setup.WithDefaultValues(map[string]any{
			cfg.LogFormatKey:         cfg.LogFormatJSON,
			cfg.LogLevelKey:          cfg.LogLevelDEBUG,
			cfg.LogOutputFileKey:     "execution.log",
			cfg.LogOutputToStdoutKey: false,
			"auth.redirect_url":      configs.DefaultRedirectURL,
			"auth.scope":             configs.DefaultAuthScopes,
		}),
	); err != nil {
		panic(err)
	}

	c := client.New(
		client.WithSecretID(configs.GetSecretID()),
		client.WithScopes(configs.GetAuthScopes()...),
		client.WithRedirectURL(configs.GetRedirectURL()),
	)

	token, err := c.Authenticate()
	if err != nil {
		panic(err)
	}

	fmt.Println("token:", token)
	slog.With("token", token).Info("token received")

	u, err := c.AuthenticatedUser()
	if err != nil {
		panic(err)
	}
	fmt.Println("user:", u)
}
