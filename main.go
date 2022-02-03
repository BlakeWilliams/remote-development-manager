package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blakewilliams/remote-development-manager/internal/cmd"
)

var LogPath string = os.TempDir() + "rdm.log"

func main() {
	logger := log.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range c {
			logger.Printf("received signal %v", sig)
			cancel()
		}
	}()

	_ = cmd.Execute(ctx)
}
