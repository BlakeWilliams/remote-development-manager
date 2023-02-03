package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var LogPath string = filepath.Join(os.TempDir(), "rdm.log")

func newLogpathCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "logpath",
		Short: "Prints the location of the log file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(LogPath)
		},
	}
}
