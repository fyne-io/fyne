//go:build !wasm && !test_web_driver

package glfw

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"

	"github.com/go-gl/glfw/v3.4/glfw"
)

// platform values returned by forcePlatform to override GLFW's auto-detection.
const (
	platformAuto    = ""
	platformX11     = "x11"
	platformWayland = "wayland"
)

func (*gLDriver) initGLFW() {
	switch forcePlatform() {
	case platformX11:
		glfw.InitHint(glfw.PlatformHint, int(glfw.PlatformX11))
	case platformWayland:
		glfw.InitHint(glfw.PlatformHint, int(glfw.PlatformWayland))
	}
	// glfw.Terminate resets InitHints. Re-apply compositor-side decorations
	// when the app asked for them (FYNE_DISABLE_LIBDECOR=1).
	if os.Getenv("FYNE_DISABLE_LIBDECOR") == "1" {
		glfw.InitHint(glfw.WaylandLibdecor, glfw.WaylandDisableLibdecor)
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
	d.hideIfDRMOutputLost()
	if d.shouldSkipPoll() {
		return
	}
	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func (*gLDriver) Terminate() {
	glfw.Terminate()
}
