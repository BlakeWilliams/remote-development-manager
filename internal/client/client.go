package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/blakewilliams/remote-development-manager/internal/server"
)

type Client struct {
	// Determines if command should connect locally via unix socket or if port
	// should be forwarded via ssh
	runType    string
	conn       net.Conn
	httpClient http.Client
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
	reader := bytes.NewReader(result)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://unix://"+c.path(), reader)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %w", err)
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("could not send command: %w", err)
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read response from server: %w", err)
	}

	return contents, nil
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

		httpClient: http.Client{
			Transport: &http.Transport{
				DialContext: func(_ctx context.Context, _network string, _address string) (net.Conn, error) {
					return net.Dial("unix", server.UnixSocketPath())
				},
			},
		},
	}
}
