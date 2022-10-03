package service

import (
	"fmt"
	"os"
)

// Install sets up a new service by writing a plist file and telling launchd about it.
func (s *Service) Install(plistFileContent []byte) (err error) {
	if err = s.writePlist(plistFileContent); err != nil {
		return
	}
	_, err = s.Bootstrap()
	return
}

func (s *Service) writePlist(content []byte) error {
	path, err := s.DefinitionPath()
	if err != nil {
		return fmt.Errorf("could not find definition path: %w", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", path, err)
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}
