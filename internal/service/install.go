package service

import (
	"fmt"
	"os"
	"path/filepath"
)

// Install sets up a new service by writing a plist file and telling launchd about it.
func (s *Service) Install(plistFileContent []byte) (err error) {
	path, err := s.plistPath()
	if err != nil {
		return
	}
	if err = s.writePlist(path, plistFileContent); err != nil {
		return
	}
	_, err = s.launchctlBootstrap(path)
	return
}

func (s *Service) writePlist(path string, content []byte) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", path, err)
	}
	defer f.Close()
	_, err = f.Write(content)
	return
}

// The full path of the executable that is currently running.
func executablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}
