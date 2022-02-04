package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newLogpathCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "logpath",
		Short: "Prints the location of the log file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(LogPath)
		},
	}
}
