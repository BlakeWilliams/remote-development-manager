package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/blakewilliams/remote-development-manager/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rdm",
	Short: "A server and client for better remote development integration.",
	Long: `Embetter your remote development experience!
	Complete documentation is available at https://github.com/BlakeWilliams/remote-development-manager`,
}

var LogPath string = os.TempDir() + "rdm.log"

func Execute(ctx context.Context) error {
	logger := log.Default()

	rdmConfig := config.New()
	home, err := os.UserHomeDir()

	if err == nil {
		path := filepath.Join(home, ".config/rdm/rdm.json")
		err := rdmConfig.Load(path)

		if err != nil && err != config.ErrConfigDoesNotExist {
			panic(err)
		}
	}

	rootCmd.AddCommand(newServerCmd(ctx, logger, rdmConfig))
	rootCmd.AddCommand(newCopyCmd(ctx, logger))
	rootCmd.AddCommand(newPasteCmd(ctx, logger))
	rootCmd.AddCommand(newOpenCmd(ctx, logger))
	rootCmd.AddCommand(newSocketCmd(ctx))
	rootCmd.AddCommand(newStopCmd(ctx, logger))
	rootCmd.AddCommand(newLogpathCmd(ctx))
	rootCmd.AddCommand(newRunCmd(ctx, logger, rdmConfig))
	rootCmd.AddCommand(newPsCmd(ctx, logger))
	rootCmd.AddCommand(newKillCmd(ctx, logger))

	if rdmConfig != nil {
		rootCmd.AddCommand(newRunCmd(ctx, logger, rdmConfig))
	}

	return rootCmd.Execute()
}
