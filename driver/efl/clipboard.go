package efl

import (
	"log"
)

// clipboard represents the system clipboard
type clipboard struct {
}

// Content returns the clipboard content
func (c *clipboard) Content() string {
	log.Println("Not implemented")
	return ""
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	log.Println("Not implemented")
}
