package state

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// Run describes the current runtime state of a launchd service.
type Run struct {
	enum  RunState
	err   error
	color color.Attribute
}

// Runnable is the interface we depend on to determine a service's run state.
type Runnable interface {
	Print() ([]byte, error)
}

// NewRun computes and returns the current runtime status of a service.
func NewRun(s Runnable) *Run {
	out, err := s.Print()
	if err != nil {
		if strings.Contains(err.Error(), "Could not find service") {
			return &Run{NoSuchService, nil, red}
		} else {
			return &Run{unknownRS, err, yellow}
		}
	}
	matches := parser.FindSubmatch(out)
	if len(matches) == 0 {
		return &Run{unknownRS, fmt.Errorf("parse error: no status:\n%s", string(out)), red}
	}
	match := string(matches[1])
	parsed, ok := outputToStatus[match]
	if !ok {
		return &Run{unknownRS, fmt.Errorf("unrecognized status value %s in command output", match), red}
	}
	return &Run{parsed.state, nil, parsed.color}
}

// Is compares using the underlying enum value
func (s *Run) Is(enum RunState) bool {
	return enum == s.enum
}

// Pretty is a description of the state formatted nicely for display.
func (s *Run) Pretty() string {
	return fmt.Sprintf("Run state:     [%s]%s", s.Color(), s.Err())
}

// Color is the state name rendered with ANSI color escape sequences.
func (s *Run) Color() string {
	return color.New(s.color).SprintFunc()(s.String())
}

// String is the state name.
func (s *Run) String() string {
	return s.enum.String()
}

// Err is a description of the error encountered while computing the state, if any.
func (s *Run) Err() string {
	if s.err != nil {
		return fmt.Sprintf(" (%v)", s.err.Error())
	} else {
		return ""
	}
}

// A simple regexp parser for the `state` line of `launchctl print` output,
// which looks like:
//
//     gui/501/com.blakewilliams.rdm = {
//             active count = 0
//             path = /Users/carl/Library/LaunchAgents/com.blakewilliams.rdm.plist
//             state = not running
//
// For this purpose we only care about a single line, so extract it with a
// multiline-enabled regex.
var parser = regexp.MustCompile(`(?m)^\s+state = (.*)$`)

type coloredState struct {
	state RunState
	color color.Attribute
}

var outputToStatus = map[string]*coloredState{
	"running":         {Running, green},
	"xpcproxy":        {Starting, yellow},
	"spawn scheduled": {Starting, yellow},
	"not running":     {NotRunning, yellow},
}

// RunState enumerates the runtime states that a launchd service may be in.
type RunState int

//go:generate stringer -type=RunState
const (
	unknownRS RunState = iota
	NoSuchService
	Running
	Starting
	NotRunning
)
