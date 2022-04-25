package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "remote-development-manager",
	Short: "A server and client for better remote development integration.",
	Long: `Embetter your remove development experience!
	Complete documentation is available at https://github.com/BlakeWilliams/remote-development-manager`,
}

var LogPath string = os.TempDir() + "rdm.log"

func Execute(ctx context.Context) error {
	logger := log.Default()

	rootCmd.AddCommand(newServerCmd(ctx, logger))
	rootCmd.AddCommand(newCopyCmd(ctx, logger))
	rootCmd.AddCommand(newPasteCmd(ctx, logger))
	rootCmd.AddCommand(newOpenCmd(ctx, logger))
	rootCmd.AddCommand(newSocketCmd(ctx))
	rootCmd.AddCommand(newStopCmd(ctx, logger))
	rootCmd.AddCommand(newLogpathCmd(ctx))

	return rootCmd.Execute()
}
