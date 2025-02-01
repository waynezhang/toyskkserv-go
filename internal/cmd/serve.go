package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/dictionary"
	"github.com/waynezhang/toyskkserv/internal/scheduler"
	"github.com/waynezhang/toyskkserv/internal/server"
	"github.com/waynezhang/toyskkserv/internal/tcp"
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

		dm := dictionary.NewDictManager(dictionary.Config{
			Dictionaires:     cfg.Dictionaries,
			Directory:        cfg.DictionaryDirectory,
			FallbackToGoogle: cfg.FallbackToGoogle,
			UseDiskCache:     cfg.UseDiskCache,
		})

		server.
			New(cfg.ListenAddr, dm).
			Start()
	},
}
