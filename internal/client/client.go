package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/server"
)

type Client struct {
	// Determines if command should connect locally via unix socket or if port
	// should be forwarded via ssh
	runType string
	conn    net.Conn
}

func (c *Client) path() string {
	if c.runType == RunLocal {
		return server.UnixSocketPath()
	} else {
		return "localhost:7391"
	}
}

func (c *Client) Connect(ctx context.Context) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, c.runType, c.path())
	c.conn = conn

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendCommand(ctx context.Context, commandName string, arguments ...string) ([]byte, error) {
	command := server.Command{
		Name:      commandName,
		Arguments: arguments,
	}

	result, err := json.Marshal(command)

	if err != nil {
		return nil, err
	}

	sizeMessage := fmt.Sprintf("%d\n", len(result))
	_, err = c.conn.Write([]byte(sizeMessage))

	if err != nil {
		return nil, fmt.Errorf("Could not write size: %w", err)
	}
	_, err = c.conn.Write(result)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*10))
	defer cancel()

	contentCh := make(chan []byte)
	errCh := make(chan error)

	go func() {
		reader := bufio.NewReader(c.conn)
		byteSize, err := reader.ReadBytes('\n')

		if errors.Is(err, io.EOF) {
			contentCh <- []byte{}
			return
		} else if err != nil {
			errCh <- fmt.Errorf("could not read size message: %w", err)
		}

		size, err := strconv.Atoi(string(byteSize[:len(byteSize)-1]))
		if err != nil {
			errCh <- fmt.Errorf("could not read message: %w", err)
		}

		content := make([]byte, size)
		io.ReadFull(reader, content)

		contentCh <- content
	}()

	select {
	case contents := <-contentCh:
		return contents, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

const (
	RunLocal  = "unix"
	RunRemote = "tcp"
)

func New() *Client {
	runType := RunLocal

	if os.Getenv("SSH_TTY") != "" {
		runType = RunRemote
	}

	return &Client{
		runType: runType,
	}
}
