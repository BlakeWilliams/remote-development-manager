package cmd

import (
	"context"
	"fmt"

	"github.com/blakewilliams/remote-development-manager/internal/server"
	"github.com/spf13/cobra"
)

func newSocketCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "socket",
		Short: "Prints the location of the unix socket",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(server.UnixSocketPath())
		},
	}
}
