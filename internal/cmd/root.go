package cmd

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
	"github.com/spf13/cobra"
)

var rootCmd = func() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "tskks",
		Short: "Toy SKK Server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := slog.LevelWarn
			if verbose {
				level = slog.LevelInfo
			}
			slog.SetDefault(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: level})))
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	return cmd
}()

func Execute() {
	_ = rootCmd.Execute()
}
