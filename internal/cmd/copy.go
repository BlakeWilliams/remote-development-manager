package cmd

import (
	"bufio"
	"context"
	"io"
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

			content, err := readBuffer(bufio.NewReader(os.Stdin))

			if err != nil {
				log.Printf("Can not get input to copy: %v", err)
				return
			}

			_, err = c.SendCommand(ctx, "copy", content)

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}
		},
	}
}

func readBuffer(r *bufio.Reader) (string, error) {
	var content strings.Builder

	for {
		line, err := r.ReadBytes('\n')

		switch err {
		case io.EOF:
			content.Write(line)
			return content.String(), nil
		case nil:
			content.Write(line)
		default:
			return "", err
		}
	}
}
