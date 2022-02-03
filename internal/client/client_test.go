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

	"github.com/blakewilliams/remote-development-manager/internal/server"
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

		var command server.Command
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
