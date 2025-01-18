package cmd

import (
	"context"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/configs"
	"github.com/eldius/onedrive-client/internal/usecase"

	"github.com/spf13/cobra"
)

// driveAddCmd represents the add command
var driveAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a drive configuration",
	Long:  `Adds a drive configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := client.New(
			client.WithSecretID(configs.GetSecretID()),
		)
		uc := usecase.NewDriveAddUseUseCase(c)
		if err := uc.DriveAdd(ctx, driveName); err != nil {
			panic(err)
		}
	},
}

var (
	driveName string
)

func init() {
	driveCmd.AddCommand(driveAddCmd)
	driveAddCmd.Flags().StringVarP(&driveName, "name", "n", "", "name of the drive")
}
