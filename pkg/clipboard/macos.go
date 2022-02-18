package clipboard

import (
	"fmt"
	"os/exec"
)

type macosClipboard struct{}

func (m macosClipboard) Copy(input string) error {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return fmt.Errorf("could not create pbcopy stdin: %w", err)
	}

	_, err = stdin.Write([]byte(input))

	if err != nil {
		return fmt.Errorf("could not create write to pbcopy: %w", err)
	}

	err = stdin.Close()

	if err != nil {
		return fmt.Errorf("could not create close pbcopy stdin: %w", err)
	}

	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("could not run pbcopy command: %w", err)
	}

	return nil
}

func (m macosClipboard) Paste() ([]byte, error) {
	cmd := exec.Command("pbpaste")

	contents, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not run pbpaste: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("could not read stdout for pbpaste: %w", err)
	}

	return contents, nil
}

var MacosClipboard Clipboard = macosClipboard{}
