package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/clipboard"
)

type Server struct {
	path       string
	logger     *log.Logger
	clipboard  clipboard.Clipboard
	httpServer *http.Server
	cancel     context.CancelFunc
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %v", err)
	}
	r.Body.Close()

	var command client.Command
	json.Unmarshal(body, &command)

	switch command.Name {
	case "status":
		rw.Write([]byte(`{ "status": "running" }`))
	case "copy":
		err := s.clipboard.Copy(command.Arguments[0])
		if err != nil {
			log.Printf("error running copy command: %v", err)
		}
	case "open":
		err := open(command.Arguments[0])
		if err != nil {
			log.Printf("error running open command: %v", err)
		}
	case "stop":
		log.Printf("received stop command, shutting down")
		s.cancel()
	case "paste":
		contents, err := s.clipboard.Paste()
		if err != nil {
			s.logger.Printf("error running paste command: %v", err)
		} else {
			_, err := rw.Write(contents)
			if err != nil {
				s.logger.Printf("could not write paste message: %v", err)
			}
		}
	default:
		s.logger.Printf("command not found: %s", command.Name)
	}
}

func (s *Server) Serve(ctx context.Context, listener net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	go func() {
		err := s.httpServer.Serve(listener)
		if err != nil {
			cancel()
		}
	}()

	<-ctx.Done()

	return ctx.Err()
}

func (s *Server) Listen(ctx context.Context) error {
	sock, err := net.Listen("unix", s.path)
	if err != nil {
		var errNo syscall.Errno

		if errors.As(err, &errNo) && errNo == syscall.EADDRINUSE {
			c := client.New()
			_, err := c.SendCommand(ctx, "status")

			if err != nil {
				os.Remove(s.path)
				sock, err = net.Listen("unix", s.path)
			} else {
				return fmt.Errorf("could not listen to unix socket: %w", errNo)
			}
		}
	}
	defer os.Remove(s.path)

	return s.Serve(ctx, sock)
}

func New(path string, clipboard clipboard.Clipboard, logger *log.Logger) *Server {
	server := &Server{path: path, clipboard: clipboard, logger: logger}
	server.httpServer = &http.Server{
		Handler:      server,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return server
}

func open(target string) error {
	cmd := exec.Command("open", target)

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("could not run open command: %w", err)
	}

	return nil
}
