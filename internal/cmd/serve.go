package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/tskks/internal/config"
	"github.com/waynezhang/tskks/internal/dictionary"
	"github.com/waynezhang/tskks/internal/scheduler"
	"github.com/waynezhang/tskks/internal/server"
	"github.com/waynezhang/tskks/internal/tcp"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start SKK server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Shared()
		cfg.OnConfigChange(func() {
			tcp.SendReloadCommand(cfg.ListenAddr)
		})

		scheduler.StartUpdateWatcher(cfg)
		dm := dictionary.Shared()

		server.
			New(cfg.ListenAddr, dm).
			Start()
	},
}
