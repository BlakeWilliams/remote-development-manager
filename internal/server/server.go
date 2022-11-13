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
	"syscall"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/hostservice"
)

type Server struct {
	host       hostservice.Runner
	path       string
	logger     *log.Logger
	httpServer *http.Server
	cancel     context.CancelFunc
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("could not read request body: %v", err)
	}
	r.Body.Close()

	var command client.Command
	json.Unmarshal(body, &command)

	switch command.Name {
	case "status":
		rw.Write([]byte(`{ "status": "running" }`))
	case "copy":
		err := s.host.Copy(command.Arguments[0])
		if err != nil {
			s.logger.Printf("error running copy command: %v", err)
		}
	case "open":
		err := s.host.Open(command.Arguments[0])
		if err != nil {
			s.logger.Printf("error running open command: %v", err)
		}
	case "stop":
		s.logger.Printf("received stop command")
		s.cancel()
	case "paste":
		contents, err := s.host.Paste()
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
		s.logger.Printf("HTTP server listening on %s", s.path)
		err := s.httpServer.Serve(listener)
		if err != nil {
			cancel()
		}
	}()

	<-ctx.Done()
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Printf("HTTP server shutdown (err=%v)", err)
	} else {
		s.logger.Println("HTTP server shutdown (clean)")
	}
	return err
}

func (s *Server) Listen(ctx context.Context) error {
	sock, err := net.Listen("unix", s.path)
	if err != nil {
		var errNo syscall.Errno

		if errors.As(err, &errNo) && errNo == syscall.EADDRINUSE {
			c := client.NewWithSocketPath(s.path)

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

func New(path string, service hostservice.Runner, logger *log.Logger) *Server {
	server := &Server{
		host:   service,
		path:   path,
		logger: logger,
	}
	server.httpServer = &http.Server{
		Handler:      server,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorLog:     logger,
	}

	return server
}
