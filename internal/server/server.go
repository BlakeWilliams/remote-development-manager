package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

type Server struct {
	path   string
	logger *log.Logger
}

type Command struct {
	Name      string
	Arguments []string
}

func UnixSocketPath() string {
	return os.TempDir() + "rdm.sock"
}

func (s *Server) Listen(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sock, err := net.Listen("unix", s.path)

	if err != nil {
		return err
	}
	defer sock.Close()
	defer os.Remove(s.path)

	connections := make(chan net.Conn)

	go func() {
		for {
			conn, err := sock.Accept()
			if err != nil {
				s.logger.Fatalf("could not accept connection: %v", err)
				continue
			}

			connections <- conn
		}
	}()

	go func() {
		for {
			conn := <-connections

			reader := bufio.NewReader(conn)
			byteSize, err := reader.ReadBytes('\n')

			if err != nil {
				s.logger.Printf("Could not read size from conn")
				continue
			}

			size, err := strconv.Atoi(string(byteSize[:len(byteSize)-1]))
			if err != nil {
				s.logger.Printf("Could not decode size from conn: %v", err)
			}

			content := make([]byte, size)
			io.ReadFull(reader, content)

			if err != nil {
				s.logger.Printf("could not accept connection: %v", err)
				continue
			}

			var command Command
			json.Unmarshal(content, &command)

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
				cancel()
			case "paste":
				contents, err := paste()
				if err != nil {
					s.logger.Printf("error running paste command: %v", err)
				} else {
					sizeMessage := fmt.Sprintf("%d\n", len(contents))
					_, err = conn.Write([]byte(sizeMessage))
					if err != nil {
						s.logger.Printf("could not write paste size message: %v", err)
					}

					_, err := conn.Write(contents)
					if err != nil {
						s.logger.Printf("could not write paste message: %v", err)
					}
				}
			default:
				s.logger.Printf("command not found: %s", command.Name)
			}

			conn.Close()
		}
	}()

	select {
	case err := <-ctx.Done():
		s.logger.Printf("Context cancelled: %v", err)
	}

	return nil
}

func New(path string, logger *log.Logger) *Server {
	return &Server{path: UnixSocketPath(), logger: logger}
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
