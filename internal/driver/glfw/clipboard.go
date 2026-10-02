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
func (clipboard) Content() string {
	// This retry logic is to work around the "Access Denied" error often thrown in windows PR#1679
	if runtime.GOOS != goos.Windows {
		return glfw.GetClipboardString()
	}
	for i := 3; i > 0; i-- {
		cb := glfw.GetClipboardString()
		if cb != "" {
			return cb
		}
		time.Sleep(50 * time.Millisecond)
	}
	// can't log retry as it would also log errors for an empty clipboard
	return ""
}

// SetContent sets the clipboard content
func (clipboard) SetContent(content string) {
	// This retry logic is to work around the "Access Denied" error often thrown in windows PR#1679
	if runtime.GOOS != goos.Windows {
		glfw.SetClipboardString(content)
		return
	}
	for i := 3; i > 0; i-- {
		glfw.SetClipboardString(content)
		if glfw.GetClipboardString() == content {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	fyne.LogError("GLFW clipboard set failed", nil)
}
