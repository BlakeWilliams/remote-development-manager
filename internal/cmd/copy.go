package cmd

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/spf13/cobra"
)

func newCopyCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "copy",
		Short: "Copies stdin to clipboard on the host machine.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			c := client.New()
			// TODO handle argument in addition to STDIN
			// TODO check if stdin has data first, otherwise exit
			scanner := bufio.NewScanner(os.Stdin)
			var content strings.Builder

			for scanner.Scan() {
				content.Write(scanner.Bytes())
				content.WriteByte('\n')
			}

			if scanner.Err() != nil {
				log.Printf("Can not get input to copy: %v", scanner.Err())
				cancel()
				return
			}

			_, err := c.SendCommand(ctx, "copy", content.String())

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}
		},
	}
}
