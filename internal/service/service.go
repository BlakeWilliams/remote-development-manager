package service

import (
	"fmt"
	"os"
	"path/filepath"
)

// Everything we deal with is limited to the current user's gui domain.
var domain = fmt.Sprintf("gui/%d", os.Getuid())

// Service is a LaunchAgent service.
type Service struct {
	Name string
}

// New returns a Service with the given name for the current user
func New(name string) *Service {
	return &Service{name}
}

// UserSpecifier unambiguously identifies the service in subcommands.
// e.g. gui/501/com.blakewilliams.rdm
// See launchctl(1).
func (s *Service) UserSpecifier() string {
	return fmt.Sprintf("%s/%s", domain, s.Name)
}

// The absolute fs path where the service's plist config lives
func (s *Service) plistPath() (string, error) {
	dir, err := launchAgentsDir()
	if err != nil {
		return "", err
	}
	plistFileName := s.Name + ".plist"
	return filepath.Join(dir, plistFileName), nil
}
