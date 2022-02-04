package cmd

import (
	"context"
	"log"
	"os"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/spf13/cobra"
)

func newOpenCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "open url",
		Short: "Sends given url to the open command",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			c := client.New()

			_, err := c.SendCommand(ctx, "open", args[0])

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}
		},
	}
}
