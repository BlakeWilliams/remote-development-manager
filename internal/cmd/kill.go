package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/spf13/cobra"
)

func newKillCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "kill pid",
		Short: "Kills a process started by rdm.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			c := client.New()
			result, err := c.SendCommand(ctx, "kill", args[0])

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}

			fmt.Println(string(result))
		},
	}
}
