//go:build !wasm && !test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"

	"github.com/go-gl/glfw/v3.4/glfw"
)

func (d *gLDriver) initGLFW() {
	if shouldForceX11() {
		// A Wayland compositor that forces client-side decorations would push
		// GLFW onto libdecor, which is unstable. Prefer X11 through XWayland.
		glfw.InitHint(glfw.PlatformHint, int(glfw.PlatformX11))
	}

	err := glfw.Init()
	if err != nil {
		fyne.LogError("failed to initialise GLFW", err)
		return
	}

	initCursors()
	if glfw.GetPlatform() == glfw.PlatformWayland {
		build.IsWayland = true
	}
}

func (d *gLDriver) pollEvents() {
	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func (d *gLDriver) Terminate() {
	glfw.Terminate()
}
