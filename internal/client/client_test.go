package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_path(t *testing.T) {
	client := New()
	require.Regexp(t, regexp.MustCompile("http://unix://"), client.path)

	os.Setenv("SSH_TTY", "fake")
	defer os.Unsetenv("SSH_TTY")

	client = New()
	require.Equal(t, "http://localhost:7391", client.path)
}

func TestClient_SendCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var command Command
		err = json.Unmarshal(content, &command)

		require.NoError(t, err)

		require.Equal(t, "copy", command.Name)
		require.Equal(t, "test 1 2 3", command.Arguments[0])

		rw.Write([]byte("test result"))
	}))
	defer server.Close()

	client := &Client{
		path:       server.URL,
		httpClient: *http.DefaultClient,
	}

	responseContent, err := client.SendCommand(context.Background(), "copy", "test 1 2 3")

	require.NoError(t, err)
	require.Equal(t, "test result", string(responseContent))
}

func TestClient_DoesNotHang(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 200)
		rw.Write([]byte("test result"))
	}))
	defer server.Close()

	client := &Client{
		path:       server.URL,
		httpClient: *http.DefaultClient,
	}

	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()

	_, err := client.SendCommand(ctx, "copy", "test 1 2 3")
	require.Error(t, err)

	require.Less(t, time.Now().Sub(start), time.Millisecond*20)
}
