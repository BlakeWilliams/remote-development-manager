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
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/config"
	"github.com/blakewilliams/remote-development-manager/pkg/clipboard"
	"github.com/blakewilliams/remote-development-manager/pkg/process"
)

type Server struct {
	path           string
	rdmConfig      *config.RdmConfig
	logger         *log.Logger
	clipboard      clipboard.Clipboard
	httpServer     *http.Server
	cancel         context.CancelFunc
	processManager *process.Manager
	context        context.Context
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
	case "kill":
		pid, err := strconv.Atoi(command.Arguments[0])

		if err != nil {
			rw.Write([]byte("Invalid pid provided\n"))
			return
		}

		err = s.processManager.Kill(pid)

		if err != nil {
			rw.Write([]byte(fmt.Sprintf("Could not kill process: %s\n", err)))
			return
		}

		rw.Write([]byte("Process killed"))
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
	case "ps":
		commands := s.processManager.RunningProcesses()

		if len(commands) == 0 {
			rw.Write([]byte("No processes running\n"))
			return
		}

		rw.Write([]byte("Processes:\n"))

		writer := tabwriter.NewWriter(rw, 20, 2, 4, ' ', 0)
		fmt.Fprintln(writer, "pid\tcommand")
		for _, command := range commands {
			fmt.Fprintf(writer, "%d\t%s\n", command.Process.Pid, command.String())
		}

		writer.Flush()
	case "runbg":
		userCommandName := command.Arguments[0]
		userCommandArgs := command.Arguments[1:]

		if userCommand, ok := s.rdmConfig.Commands[userCommandName]; ok {
			err := s.processManager.RunInBackground(s.context, userCommandName, userCommand.ExecutablePath, userCommandArgs...)

			if err != nil {
				rw.Write([]byte(fmt.Sprintf("Could not run command: %v", err)))
				return
			}

			rw.Write([]byte(fmt.Sprintf("Started command: %s\n", userCommandName)))
		} else {
			rw.Write([]byte("Command not found"))
		}

	case "run":
		userCommandName := command.Arguments[0]
		userCommandArgs := command.Arguments[1:]

		if userCommand, ok := s.rdmConfig.Commands[userCommandName]; ok {
			ctx, cancel := context.WithTimeout(s.context, time.Second*30)
			defer cancel()

			output, err := s.processManager.RunNow(ctx, userCommandName, userCommand.ExecutablePath, userCommandArgs...)

			if err != nil {
				rw.Write([]byte(fmt.Sprintf("Could not run command: %v", err)))
				return
			}

			rw.Write([]byte(output))
		} else {
			rw.Write([]byte("Command not found\n"))
		}

	case "commands":
		var out strings.Builder
		for name := range s.rdmConfig.Commands {
			out.WriteString(fmt.Sprintf("%s\n", name))
		}

		rw.Write([]byte(out.String()))
	default:
		s.logger.Printf("command not found: %s", command.Name)
	}
}

func (s *Server) Serve(ctx context.Context, listener net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	s.context = ctx
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

func New(path string, clipboard clipboard.Clipboard, logger *log.Logger, rdmConfig *config.RdmConfig) *Server {
	server := &Server{
		path:           path,
		clipboard:      clipboard,
		logger:         logger,
		rdmConfig:      rdmConfig,
		processManager: process.NewManager(),
	}
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
