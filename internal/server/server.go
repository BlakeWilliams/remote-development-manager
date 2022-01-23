package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Server struct {
	path       string
	logger     *log.Logger
	httpServer *http.Server
}

type Command struct {
	Name      string
	Arguments []string
}

func UnixSocketPath() string {
	return os.TempDir() + "rdm.sock"
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %v", err)
	}
	r.Body.Close()

	var command Command
	json.Unmarshal(body, &command)

	switch command.Name {
	case "copy":
		err := copy(command.Arguments[0])
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
		contents, err := paste()
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

func New(path string, logger *log.Logger) *Server {
	server := &Server{path: UnixSocketPath(), logger: logger}
	server.httpServer = &http.Server{
		Handler:      server,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return server
}

/// TODO extract into commands package and support more than macOS
func copy(input string) error {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return fmt.Errorf("could not create pbcopy stdin: %w", err)
	}

	_, err = stdin.Write([]byte(input))

	if err != nil {
		return fmt.Errorf("could not create write to pbcopy: %w", err)
	}

	err = stdin.Close()

	if err != nil {
		return fmt.Errorf("could not create close pbcopy stdin: %w", err)
	}

	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("could not run pbcopy command: %w", err)
	}

	return nil
}

func open(target string) error {
	cmd := exec.Command("open", target)

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("could not run open command: %w", err)
	}

	return nil
}

func paste() ([]byte, error) {
	cmd := exec.Command("pbpaste")

	contents, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not run pbpaste: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("could not read stdout for pbpaste: %w", err)
	}

	return contents, nil
}
