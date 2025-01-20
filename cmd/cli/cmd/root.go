package cmd

import (
	cfg "github.com/eldius/initial-config-go/configs"
	"github.com/eldius/initial-config-go/setup"
	"github.com/eldius/onedrive-client/internal/configs"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "onedrive",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return setup.InitSetup(
			configs.GetAppName(),
			setup.WithConfigFileToBeUsed(cfgFile),
			setup.WithDefaultValues(map[string]any{
				cfg.LogFormatKey:           cfg.LogFormatJSON,
				cfg.LogLevelKey:            cfg.LogLevelDEBUG,
				cfg.LogOutputFileKey:       "execution.log",
				cfg.LogOutputToStdoutKey:   false,
				cfg.LogKeysToRedactKey:     configs.RedactedKeyList,
				configs.AuthRedirectURLKey: configs.DefaultRedirectURL,
				configs.AuthScopesKey:      configs.DefaultAuthScopes,
				configs.DBFileKey:          ".db",
			}),
		)
	},
	Short: "A simple command line interface to manage onedrive files",
	Long:  `A simple command line interface to manage onedrive files.`,
}

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.onedrive-client.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
