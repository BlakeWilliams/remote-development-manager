package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/config"
	"github.com/spf13/cobra"
)

func newRunCmd(ctx context.Context, logger *log.Logger, config *config.RdmConfig) *cobra.Command {
	// TODO this needs to diverge, server should hold all the commands and
	// client should query for available commands
	cmd := &cobra.Command{
		Use:   "run command [ARGS]",
		Short: "Runs a custom command defined in the rdm config",
		Run: func(cmd *cobra.Command, args []string) {
			c := client.New()
			needsBackground, err := cmd.Flags().GetBool("background")
			if err != nil {
				panic(err)
			}

			serverCmd := "run"
			if needsBackground {
				serverCmd = "runbg"
			}

			content, err := c.SendCommand(ctx, serverCmd, args...)

			if err != nil {
				fmt.Printf("Could not run command: %v", err)
				return
			}

			fmt.Print(string(content))
		},
	}

	var background bool
	cmd.Flags().BoolVarP(&background, "background", "b", false, "Runs process in background instead of foreground")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		fmt.Println(longRunDescription(config))
	})

	return cmd
}

func longRunDescription(config *config.RdmConfig) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var out strings.Builder

	out.WriteString("\n")

	c := client.New()
	result, err := c.SendCommand(ctx, "commands")

	if err != nil {
		cancel()
		out.WriteString("Could not communicate with server to get commands")
		return out.String()
	}

	result = bytes.TrimRight(result, "\n")

	out.WriteString("Available commands:\n")
	commands := bytes.Split(result, []byte("\n"))

	for i, command := range commands {
		out.WriteString(fmt.Sprintf("  %s", command))
		if i != len(commands)-1 {
			out.WriteByte('\n')
		}
	}

	return out.String()
}
