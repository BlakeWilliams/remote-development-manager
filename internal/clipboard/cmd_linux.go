//go:build linux
// +build linux

package clipboard

// New returns a Clipboard using xclip for linux
func New() Clipboard {
	return &commandClipboard{
		copy:  &command{"xclip", []string{"-in", "-selection", "clipboard"}},
		paste: &command{"xclip", []string{"-out", "-selection", "clipboard"}},
	}
}
