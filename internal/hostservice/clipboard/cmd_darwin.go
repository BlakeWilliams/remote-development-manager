//go:build darwin
// +build darwin

package clipboard

// New returns a Clipboard using pbpaste for macos
func New() Clipboard {
	return &commandClipboard{
		copy:  &command{"pbcopy", []string{}},
		paste: &command{"pbpaste", []string{}},
	}
}
