package cmd

import (
	"context"
	"log"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/spf13/cobra"
)

func newStopCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stops the server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			c := client.New()
			_, err := c.SendCommand(ctx, "stop")

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}
		},
	}
}
