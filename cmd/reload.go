package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/toyskkserv/internal/config"
	"github.com/waynezhang/toyskkserv/internal/tcp"
)

func init() {
	rootCmd.AddCommand(&reloadCmd)
}

var reloadCmd = cobra.Command{
	Use:   "reload",
	Short: "Reloadf dictionaries",
	Run: func(cmd *cobra.Command, args []string) {
		addr := config.Shared().ListenAddr
		tcp.SendReloadCommand(addr)
	},
}
