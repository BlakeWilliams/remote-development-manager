package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/spf13/cobra"
)

func newPasteCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "paste",
		Short: "Prints the contents of host host machines clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			c := client.New()

			result, err := c.SendCommand(ctx, "paste")

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}

			fmt.Print(string(result))
		},
	}
}
