package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	// Determines if command should connect locally via unix socket or if port
	// should be forwarded via ssh
	path       string
	httpClient http.Client
}

type Command struct {
	Name      string
	Arguments []string
}

func UnixSocketPath() string {
	return strings.TrimRight(os.TempDir(), "/") + "/rdm.sock"
}

func (c *Client) SendCommand(ctx context.Context, commandName string, arguments ...string) ([]byte, error) {
	command := Command{
		Name:      commandName,
		Arguments: arguments,
	}

	result, err := json.Marshal(command)
	reader := bytes.NewReader(result)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.path, reader)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %w", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("could not send command: %w", err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read response from server: %w", err)
	}

	return contents, nil
}

const (
	RunLocal  = "unix"
	RunRemote = "tcp"
)

func New() *Client {
	return NewWithSocketPath(UnixSocketPath())
}

func NewWithSocketPath(socketPath string) *Client {
	runType := RunLocal

	if os.Getenv("SSH_TTY") != "" || os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_CONNECTION") != "" {
		runType = RunRemote
	}

	client := &Client{
		httpClient: http.Client{
			Timeout: time.Second * 10,
		},
	}

	if runType == RunLocal {
		client.path = "http://unix://" + socketPath
	} else {
		client.path = "http://localhost:7391"
	}

	if runType == RunLocal {
		client.httpClient.Transport = &http.Transport{
			DialContext: func(_ctx context.Context, _network string, _address string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		}
	}

	return client
}
