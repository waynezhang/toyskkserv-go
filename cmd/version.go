package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/waynezhang/toyskkserv/internal/defs"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(defs.VersionString())
	},
}
