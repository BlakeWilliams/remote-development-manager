package clipboard

type Clipboard interface {
	Copy(string) error
	Paste() ([]byte, error)
}
