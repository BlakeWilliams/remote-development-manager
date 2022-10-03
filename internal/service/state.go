package service

import "github.com/blakewilliams/remote-development-manager/internal/service/state"

// InstallState returns information about whether or not a service is installed.
func (s *Service) InstallState() *state.Install {
	return state.NewInstall(s)
}

// RunState returns information about the runtime status of a service.
func (s *Service) RunState() *state.Run {
	return state.NewRun(s)
}
