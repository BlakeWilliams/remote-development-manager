package clipboard

// Clipboard interacts with the system clipboard.
type Clipboard interface {
	// Copy a string to the clipboard.
	Copy(string) error
	// Retrieve the current contents of the clipboard.
	Paste() ([]byte, error)
}
