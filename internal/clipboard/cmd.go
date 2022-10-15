package clipboard

import (
	"fmt"
	"os/exec"
)

type command struct {
	name string
	argv []string
}

type commandClipboard struct {
	copy  *command
	paste *command
}

func (m *commandClipboard) Copy(input string) error {
	cmd := exec.Command(m.copy.name, m.copy.argv...)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return fmt.Errorf("could not create %v stdin: %w", m.copy.name, err)
	}

	_, err = stdin.Write([]byte(input))

	if err != nil {
		return fmt.Errorf("could not create write to %v: %w", m.copy.name, err)
	}

	err = stdin.Close()

	if err != nil {
		return fmt.Errorf("could not create close %v stdin: %w", m.copy.name, err)
	}

	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("could not run %v command: %w", m.copy.name, err)
	}

	return nil
}

func (m *commandClipboard) Paste() ([]byte, error) {
	cmd := exec.Command(m.paste.name, m.paste.argv...)

	contents, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not run %v: %w", m.paste.name, err)
	}

	if err != nil {
		return nil, fmt.Errorf("could not read stdout for %v: %w", m.paste.name, err)
	}

	return contents, nil
}
