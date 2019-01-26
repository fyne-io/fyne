package efl

import (
	"log"
)

// Clipboard represents the system clipboard
type Clipboard struct {
}

// Content returns the clipboard content
func (c *Clipboard) Content() string {
	log.Println("Not implemented")
	return ""
}

// SetContent sets the clipboard content
func (c *Clipboard) SetContent(content string) {
	log.Println("Not implemented")
}
