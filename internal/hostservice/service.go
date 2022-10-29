package hostservice

import (
	"github.com/blakewilliams/remote-development-manager/internal/hostservice/clipboard"
	"github.com/blakewilliams/remote-development-manager/internal/hostservice/open"
)

// Service is the set of functionalities made available by the server to the client.
type Service interface {
	clipboard.Clipboard
	open.Opener
}

// HostService is a default implementation of Service which offers copy paste
// and open functionality.
type HostService struct {
	clipboard.Clipboard
	open.Opener
}

// New returns a HostService.
func New() *HostService {
	return &HostService{
		Clipboard: clipboard.New(),
	}
}

// Copy a string to the host system's clipboard.
func (svc *HostService) Copy(s string) error {
	return svc.Copy(s)
}

// Paste a string from the host system's clipboard.
func (svc *HostService) Paste() ([]byte, error) {
	return svc.Paste()
}

// Open the target on the host system, most likely by opening a browser.
func (svc *HostService) Open(target string) error {
	return open.Open(target)
}
