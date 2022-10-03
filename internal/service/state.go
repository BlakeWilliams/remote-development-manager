package service

import "github.com/blakewilliams/remote-development-manager/internal/service/state"

func (s *Service) IsHealthy() bool {
	return s.InstallState().Is(state.Installed) && s.RunState().Is(state.Running)
}

// InstallState returns information about whether or not a service is installed.
func (s *Service) InstallState() *state.Install {
	return state.NewInstall(s)
}

// RunState returns information about the runtime status of a service.
func (s *Service) RunState() *state.Run {
	return state.NewRun(s)
}
