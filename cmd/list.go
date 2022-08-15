package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list command",
	Long:  "indicate remote storage files.",
	Run: func(cmd *cobra.Command, args []string) {
		storage.List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
