package cmd

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadBuffer(t *testing.T) {
	testCases := map[string]struct {
		value string
	}{
		"single line": {value: "foo"},
		"multi-line":  {value: "foo\nbar\nbaz"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			handle := strings.NewReader(tc.value)
			reader := bufio.NewReader(handle)
			content, err := readBuffer(reader)

			require.NoError(t, err)
			require.Equal(t, tc.value, content)
		})
	}
}
