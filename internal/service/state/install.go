package state

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Install describes the current installation state of a launchd service.
type Install struct {
	enum  installState
	err   error
	color color.Attribute
}

// Installable is the interface we depend on to determine install state.
type Installable interface {
	DefinitionPath() (string, error)
	Print() ([]byte, error)
}

// NewInstall computes and returns the installation state for a service.
func NewInstall(svc Installable) *Install {
	plist, err := svc.DefinitionPath()
	if err != nil {
		return &Install{Unknown, err, color.FgYellow}
	}

	if _, err = os.Stat(plist); errors.Is(err, os.ErrNotExist) {
		return &Install{NotInstalled, nil, color.FgRed}
	} else if err != nil {
		return &Install{Unknown, err, color.FgYellow}
	}

	output, err := svc.Print()
	if err != nil {
		if strings.Contains("Could not find service", string(output)) {
			return &Install{PlistPresentButNotLoaded, nil, color.FgRed}
		} else {
			return &Install{Unknown, err, color.FgYellow}
		}
	}

	return &Install{Installed, nil, color.FgGreen}
}

// Is compares using the underlying enum value
func (s *Install) Is(enum installState) bool {
	return enum == s.enum
}

// Pretty is a description of the state formatted nicely for display.
func (s *Install) Pretty() string {
	return fmt.Sprintf("Install state: [%s]%s", s.Color(), s.Err())
}

// Color is the state name rendered with ANSI color escape sequences.
func (s *Install) Color() string {
	return color.New(s.color).SprintFunc()(s.String())
}

// String is the state name.
func (s *Install) String() string {
	return s.enum.String()
}

// Err is a description of the error encountered while computing the state, if any.
func (s *Install) Err() string {
	if s.err != nil {
		return fmt.Sprintf(" (%v)", s.err.Error())
	} else {
		return ""
	}
}

// installState is an enumeration of the possible installation states a
// service can be in.
//go:generate stringer -type=installState
type installState int

const (
	Unknown installState = iota
	NotInstalled
	PlistPresentButNotLoaded
	Installed
)
