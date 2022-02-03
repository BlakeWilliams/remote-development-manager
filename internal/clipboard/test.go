package clipboard

type TestClipboard struct {
	Buffer string
}

func (tc *TestClipboard) Copy(input string) error {
	tc.Buffer = input

	return nil
}

func (tc *TestClipboard) Paste() ([]byte, error) {
	return []byte(tc.Buffer), nil
}

func NewTestClipboard() *TestClipboard {
	return &TestClipboard{}
}

var _ Clipboard = (*TestClipboard)(nil)
