package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/blakewilliams/remote-development-manager/internal/client"
	"github.com/blakewilliams/remote-development-manager/internal/clipboard"
	"github.com/stretchr/testify/require"
)

var path string = client.UnixSocketPath() + ".test"

var httpClient http.Client = http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		DialContext: func(_ctx context.Context, _network string, _address string) (net.Conn, error) {
			return net.Dial("unix", path)
		},
	},
}

func TestServer_Copy(t *testing.T) {
	nullLogger := log.New(io.Discard, "", log.LstdFlags)

	testClipboard := clipboard.NewTestClipboard()
	server := New(path, testClipboard, nullLogger)

	listener, err := net.Listen("unix", server.path)
	defer os.Remove(server.path)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := server.Serve(ctx, listener)
		require.ErrorIs(t, err, context.Canceled)
	}()

	copyCommand := client.Command{
		Name:      "copy",
		Arguments: []string{"test 1 2 3"},
	}

	data, err := json.Marshal(copyCommand)
	require.NoError(t, err)

	_, err = httpClient.Post("http://unix://"+path, "application/json", bytes.NewReader(data))
	require.NoError(t, err)

	require.Equal(t, "test 1 2 3", testClipboard.Buffer)

	pasteCommand := client.Command{
		Name:      "paste",
		Arguments: []string{},
	}

	data, err = json.Marshal(pasteCommand)
	require.NoError(t, err)

	result, err := httpClient.Post("http://unix://"+path, "application/json", bytes.NewReader(data))
	require.NoError(t, err)

	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	require.Equal(t, "test 1 2 3", string(body))
}

func TestServer_Ping(t *testing.T) {
	nullLogger := log.New(io.Discard, "", log.LstdFlags)

	testClipboard := clipboard.NewTestClipboard()
	server := New(path, testClipboard, nullLogger)

	listener, err := net.Listen("unix", server.path)
	defer os.Remove(server.path)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := server.Serve(ctx, listener)
		require.ErrorIs(t, err, context.Canceled)
	}()

	statusCommand := client.Command{
		Name:      "status",
		Arguments: []string{},
	}

	data, err := json.Marshal(statusCommand)
	require.NoError(t, err)

	result, err := httpClient.Post("http://unix://"+path, "application/json", bytes.NewReader(data))
	require.NoError(t, err)

	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	require.Equal(t, `{ "status": "running" }`, string(body))
}

func TestServer_ExistingSocket(t *testing.T) {
	nullLogger := log.New(io.Discard, "", log.LstdFlags)

	testClipboard := clipboard.NewTestClipboard()
	server := New(path, testClipboard, nullLogger)

	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	file, err := os.Create(path)
	require.NoError(t, err)
	file.Close()
	defer os.Remove(server.path)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Listen(ctx)
		require.ErrorIs(t, err, context.Canceled)
	}()
}
