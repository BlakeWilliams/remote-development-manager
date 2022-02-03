package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/clipboard"
	"github.com/blakewilliams/remote-development-manager/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rdm",
	Short: "A server and client for better remote development integration.",
}

var LogPath string = os.TempDir() + "rdm.log"

func Execute(ctx context.Context) error {
	logger := log.Default()

	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Starts a server on the local machine.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			logFile, err := os.OpenFile(LogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			defer logFile.Close()
			log.SetOutput(logFile)

			s := server.New(server.UnixSocketPath(), clipboard.MacosClipboard, logger)
			err = s.Listen(ctx)

			if err != nil {
				log.Printf("Server could not be started: %v", err)
				cancel()
				return
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
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
	})

	rootCmd.AddCommand(&cobra.Command{
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
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "open url",
		Short: "Sends given url to the open command",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			if len(args) == 1 {
				log.Println("please provide a url to open")
				cancel()
				return
			}

			c := client.New()

			_, err := c.SendCommand(ctx, "open", args[0])

			if err != nil {
				log.Printf("Can not send command: %v", err)
				cancel()
				return
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "socket",
		Short: "Prints the location of the unix socket",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(server.UnixSocketPath())
		},
	})

	rootCmd.AddCommand(&cobra.Command{
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
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "logpath",
		Short: "Prints the location of the log file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(LogPath)
		},
	})

	return rootCmd.Execute()
}
