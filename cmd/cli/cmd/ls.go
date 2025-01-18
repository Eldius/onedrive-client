package cmd

import (
	"context"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/configs"
	"github.com/eldius/onedrive-client/internal/usecase"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists files",
	Long:  `Lists files.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := client.New(
			client.WithSecretID(configs.GetSecretID()),
		)
		uc := usecase.NewListFilesUseCase(c)
		if err := uc.ListFilesFromDrive(ctx, lsArgs.accountName); err != nil {
			panic(err)
		}
	},
}

var (
	lsArgs struct {
		accountName string
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringVarP(&lsArgs.accountName, "account-name", "a", "", "account name")
}
