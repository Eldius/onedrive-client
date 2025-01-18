package cmd

import (
	"context"
	"github.com/eldius/onedrive-client/client"
	"github.com/eldius/onedrive-client/internal/configs"
	"github.com/eldius/onedrive-client/internal/usecase"

	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to the OneDrive",
	Long:  `Upload files to the OneDrive.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := client.New(
			client.WithSecretID(configs.GetSecretID()),
		)
		uc := usecase.NewFileUpload(c)
		if err := uc.Upload(ctx, uploadOpts.accountName, uploadOpts.inputFile, uploadOpts.outputFile); err != nil {
			panic(err)
		}

	},
}

var (
	uploadOpts struct {
		accountName string
		inputFile   string
		outputFile  string
	}
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uploadCmd.Flags().StringVarP(&uploadOpts.accountName, "account", "a", "", "Account name")
	uploadCmd.Flags().StringVarP(&uploadOpts.inputFile, "input-file", "i", "", "File to upload")
	uploadCmd.Flags().StringVarP(&uploadOpts.outputFile, "output-file", "o", "", "Remote path")
}
