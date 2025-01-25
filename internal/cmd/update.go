package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/files"
)

func init() {
	rootCmd.AddCommand(&updateCmd)
}

var updateCmd = cobra.Command{
	Use:   "update",
	Short: "Update dictionaries",
	Run: func(cmd *cobra.Command, args []string) {
		files.UpdateDictionaries(
			config.Shared().Dictionaries,
			config.Shared().DictionaryDirectory,
		)
	},
}
