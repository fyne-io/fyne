//go:build !wasm && !test_web_driver

package glfw

import (
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/goos"

	"github.com/go-gl/glfw/v3.4/glfw"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = clipboard{}

func NewClipboard() fyne.Clipboard {
	return clipboard{}
}

// clipboard represents the system clipboard
type clipboard struct{}

// Content returns the clipboard content
func (c clipboard) Content() string {
	// This retry logic is to work around the "Access Denied" error often thrown in windows PR#1679
	if runtime.GOOS != goos.Windows {
		return c.content()
	}
	for i := 3; i > 0; i-- {
		cb := c.content()
		if cb != "" {
			return cb
		}
		time.Sleep(50 * time.Millisecond)
	}
	// can't log retry as it would also log errors for an empty clipboard
	return ""
}

func (clipboard) content() string {
	return glfw.GetClipboardString()
}

// SetContent sets the clipboard content
func (c clipboard) SetContent(content string) {
	// This retry logic is to work around the "Access Denied" error often thrown in windows PR#1679
	if runtime.GOOS != goos.Windows {
		c.setContent(content)
		return
	}
	for i := 3; i > 0; i-- {
		c.setContent(content)
		if c.content() == content {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	fyne.LogError("GLFW clipboard set failed", nil)
}

func (clipboard) setContent(content string) {
	glfw.SetClipboardString(content)
}
