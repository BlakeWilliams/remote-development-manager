package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/clipboard"
)

type Server struct {
	path       string
	logger     *log.Logger
	clipboard  clipboard.Clipboard
	httpServer *http.Server
}

type Command struct {
	Name      string
	Arguments []string
}

func UnixSocketPath() string {
	return strings.TrimRight(os.TempDir(), "/") + "/rdm.sock"
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %v", err)
	}
	r.Body.Close()

	var command Command
	json.Unmarshal(body, &command)

	switch command.Name {
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

func (s *Server) Listen(ctx context.Context) error {
	sock, err := net.Listen("unix", s.path)
	if err != nil {
		return fmt.Errorf("could not listen to unix socket: %w", err)
	}
	defer os.Remove(s.path)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		err = s.httpServer.Serve(sock)
		if err != nil {
			cancel()
		}
	}()

	<-ctx.Done()

	return ctx.Err()
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
