package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/hostservice"
	"github.com/blakewilliams/remote-development-manager/internal/server"
	"github.com/spf13/cobra"
)

func newServerCmd(ctx context.Context, logger *log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Starts a server on the local machine.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			logFile := updateLoggerForServer(logger)
			defer logFile.Close()

			s := server.New(client.UnixSocketPath(), hostservice.New(), logger)
			err := s.Listen(ctx)

			if err != nil {
				logger.Printf("Server could not be started: %v\n", err)
				cancel()
				return
			}
		},
	}
}

// We have an existing logger attached to stderr for human consumption. Rather
// than creating a new one for the server we are about to launch, reconfigure
// the existing one with more appropriate settings for a server.
func updateLoggerForServer(logger *log.Logger) io.Closer {
	logFile, err := os.OpenFile(LogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Errorf("unable to open logfile %s: %w", LogPath, err))
	}

	// Since our server process runs in the foreground as well as writes to a
	// log, we will log to stdout and the log file at the same time.
	logSink := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(logSink)

	// In a server context we want timestamps and file locations in our logs.
	logger.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile | log.Lmicroseconds)

	// Return this file so it can be closed when the server exits.
	return logFile
}
