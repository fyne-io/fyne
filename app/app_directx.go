//go:build directx && windows && !ci && !android && !ios && !mobile && !tamago && !noos && !tinygo

package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/directx"
)

// NewWithID returns a new app instance using the native Win32 and Direct3D 11
// driver instead of GLFW and OpenGL. Selected with the `directx` build tag.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(directx.NewDriver(), directx.NewClipboard(), id)
}
