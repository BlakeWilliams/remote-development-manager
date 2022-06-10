package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/spf13/cobra"
)

func newPsCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "ps",
		Short: "Lists running processes started by the run command",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			c := client.New()
			result, err := c.SendCommand(ctx, "ps")

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}

			fmt.Print(string(result))
		},
	}
}
