package cmd

import (
	"fmt"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/configs"

	"github.com/spf13/cobra"
)

// driveAddCmd represents the add command
var driveAddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(
			client.WithSecretID(configs.GetSecretID()),
			client.WithScopes(configs.GetAuthScopes()...),
			client.WithRedirectURL(configs.GetRedirectURL()),
		)

		auth, err := c.Authenticate()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Auth Token: %s\n", auth.AccessToken)

		user, err := c.AuthenticatedUser()
		if err != nil {
			panic(err)
		}
		fmt.Printf("user: %+v\n", user)

		if err := c.GetAppDriveInfo(); err != nil {
			panic(err)
		}
	},
}

func init() {
	driveCmd.AddCommand(driveAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// driveAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// driveAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
