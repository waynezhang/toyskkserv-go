package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/tskks/internal/server"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start SKK server",
	Run: func(cmd *cobra.Command, args []string) {
		server.New().Start()
	},
}
