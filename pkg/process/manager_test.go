package process

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestManager_RunInBackground(t *testing.T) {
	manager := NewManager()

	tmpScript, err := ioutil.TempFile("", "manager_sleep.sh")
	require.NoError(t, err)
	tmpScript.WriteString("#!/usr/bin/env bash\nsleep 5")
	tmpScript.Chmod(0700)
	tmpScript.Close()
	defer os.Remove(tmpScript.Name())

	err = manager.RunInBackground(context.TODO(), "test", tmpScript.Name())
	require.NoError(t, err)
	commands := manager.RunningProcesses()
	require.Equal(t, 1, len(commands))

	command := commands[0]
	pid := command.Process.Pid
	manager.Kill(pid)
	err = command.Wait()
	time.Sleep(time.Millisecond * 10)

	commands = manager.RunningProcesses()
	require.Equal(t, 0, len(commands))
}

func TestManager_RunNow(t *testing.T) {
	manager := NewManager()

	tmpScript, err := ioutil.TempFile("", "manager_hello.sh")
	require.NoError(t, err)
	tmpScript.WriteString("#!/usr/bin/env bash\necho 'hi'")
	tmpScript.Chmod(0700)
	tmpScript.Close()
	defer os.Remove(tmpScript.Name())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	output, err := manager.RunNow(ctx, "test", tmpScript.Name())
	require.NoError(t, err)

	require.Equal(t, "hi\n", string(output))

	commands := manager.RunningProcesses()
	require.Equal(t, 0, len(commands))
}
