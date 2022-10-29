package open

import (
	"fmt"
	"os/exec"
)

// Opener causes the side effect of opening a referenced target in an
// appropriate way on the host system, probably by launching a browser.
type Opener interface {
	Open(target string) error
}

// Open opens the target on the host system using a platform-specific command.
func Open(target string) error {
	cmd := exec.Command(openCommand, target)

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("could not run open command: %w", err)
	}

	return nil
}
