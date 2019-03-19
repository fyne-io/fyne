package fyne

// Clipboard represents the system clipboard interface
type Clipboard interface {
	// Content returns the clipboard content
	Content() string
	// SetContent sets the clipboard content
	SetContent(content string)
}
