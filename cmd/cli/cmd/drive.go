package cmd

import (
	"github.com/spf13/cobra"
)

// driveCmd represents the drive command
var driveCmd = &cobra.Command{
	Use:   "drive",
	Short: "Drives related commands",
	Long:  `Drives related commands.`,
}

func init() {
	rootCmd.AddCommand(driveCmd)
}
