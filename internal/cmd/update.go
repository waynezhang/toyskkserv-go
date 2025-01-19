package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/tskks/internal/config"
	"github.com/waynezhang/tskks/internal/dictionary"
)

func init() {
	rootCmd.AddCommand(&updateCmd)
}

var updateCmd = cobra.Command{
	Use:   "update",
	Short: "Update dictionaries",
	Run: func(cmd *cobra.Command, args []string) {
		dictionary.UpdateDictionaries(
			config.Shared().Dictionaries,
			config.Shared().DictionaryDirectory,
			config.Shared().CacheDirectory,
		)
	},
}
