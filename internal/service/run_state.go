package service

import (
	"fmt"
	"regexp"
	"strings"
)

// runState is the current runtime state of a launchctl service.
type runState int

//go:generate stringer -type=runState
const (
	unknownRS runState = iota
	NoSuchService
	Running
	Starting
	NotRunning
)

// RunState returns the current runtime status of the service.
func (s *Service) RunState() (rs runState, err error) {
	out, err := s.launchctlPrint()
	if err != nil {
		if strings.Contains(err.Error(), "Could not find service") {
			return NoSuchService, nil
		} else {
			return unknownRS, err
		}
	}
	matches := parser.FindSubmatch(out)
	if len(matches) == 0 {
		return rs, fmt.Errorf("parse error: no status:\n%s", string(out))
	}
	match := string(matches[1])
	status, ok := outputToStatus[match]
	if !ok {
		return rs, fmt.Errorf("unrecognized status value %s in command output", match)
	}
	return status, nil
}

// output looks like:
//
//     gui/501/com.blakewilliams.rdm = {
//             active count = 0
//             path = /Users/carl/Library/LaunchAgents/com.blakewilliams.rdm.plist
//             state = not running
//
var parser = regexp.MustCompile(`(?m)^\s+state = (.*)$`)

var outputToStatus = map[string]runState{
	"running":         Running,
	"xpcproxy":        Starting,
	"spawn scheduled": Starting,
	"not running":     NotRunning,
}
