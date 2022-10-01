package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Set up a service.
func (s *Service) launchctlBootstrap(plistPath string) ([]byte, error) {
	return s.launchctl("bootstrap", domain, plistPath)
}

// Print service state.
func (s *Service) launchctlPrint() ([]byte, error) {
	return s.launchctl("print", s.UserSpecifier())
}

// Run a launchctl(1) subcommand for the service and return the output or an error
func (s *Service) launchctl(args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("launchctl", args...)
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("(%w) running `launchctl %v` for %s:\n%s", err, args, s.Name, stderr.String())
	}
	return stdout.Bytes(), nil
}
