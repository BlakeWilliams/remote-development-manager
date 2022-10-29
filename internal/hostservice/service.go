package hostservice

import (
	"github.com/blakewilliams/remote-development-manager/internal/hostservice/clipboard"
	"github.com/blakewilliams/remote-development-manager/internal/hostservice/open"
)

// Runner is the set of functionalities made available by the server to the client.
type Runner interface {
	clipboard.Clipboard
	open.Opener
}

// HostService is a default implementation of Runner which exposes
// system capabilities.
type HostService struct {
	clipboard clipboard.Clipboard
}

// New returns a HostService.
func New() *HostService {
	return &HostService{
		clipboard: clipboard.New(),
	}
}

// Copy a string to the host system's clipboard.
func (svc *HostService) Copy(s string) error {
	return svc.clipboard.Copy(s)
}

// Paste a string from the host system's clipboard.
func (svc *HostService) Paste() ([]byte, error) {
	return svc.clipboard.Paste()
}

// Open the target on the host system, most likely by opening a browser.
func (svc *HostService) Open(target string) error {
	return open.Open(target)
}

// Compile-time assertion that HostService implements Runner.
var _ Runner = (*HostService)(nil)
