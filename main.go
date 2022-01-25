package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/clipboard"
	"github.com/blakewilliams/remote-development-manager/internal/server"
)

var LogPath string = os.TempDir() + "rdm.log"

func main() {
	logger := log.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := os.Args[1:]

	if len(args) == 0 {
		// TODO print help message
		fmt.Println("No arguments passed")
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range c {
			logger.Printf("received signal %v", sig)
			cancel()
		}
	}()

	switch args[0] {
	case "server":
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
	case "copy":
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
	case "paste":
		c := client.New()

		result, err := c.SendCommand(ctx, "paste")

		if err != nil {
			log.Printf("Can not send command: %v", err)
			cancel()
			return
		}

		fmt.Print(string(result))
	case "open":
		if len(args) == 1 {
			log.Println("please provide a url to open")
			cancel()
			return
		}

		c := client.New()

		_, err := c.SendCommand(ctx, "open", args[1])

		if err != nil {
			log.Printf("Can not send command: %v", err)
			cancel()
			return
		}
	case "socket":
		fmt.Println(server.UnixSocketPath())
	case "stop":
		c := client.New()
		_, err := c.SendCommand(ctx, "stop")

		if err != nil {
			log.Printf("Can not send command: %v", err)
			cancel()
			return
		}
	case "logpath":
		fmt.Println(LogPath)
	}
}
