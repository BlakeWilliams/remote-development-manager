package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blakewilliams/remote-development-manager/internal/cmd"
)

func main() {
	userMessages := log.New(os.Stderr, "", 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		for sig := range c {
			userMessages.Printf("received signal %v", sig)
			cancel()
		}
	}()

	err := cmd.Execute(ctx, userMessages)

	if err != nil {
		userMessages.Printf("error executing command: %v", err)
		os.Exit(1)
	}
}
