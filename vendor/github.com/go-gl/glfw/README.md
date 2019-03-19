# GLFW 3.2 for Go [![Build Status](https://travis-ci.org/go-gl/glfw.svg?branch=master)](https://travis-ci.org/go-gl/glfw) [![GoDoc](https://godoc.org/github.com/go-gl/glfw/v3.2/glfw?status.svg)](https://godoc.org/github.com/go-gl/glfw/v3.2/glfw)

## Installation

* GLFW C library source is included and built automatically as part of the Go package. But you need to make sure you have dependencies of GLFW:
	* On macOS, you need Xcode or Command Line Tools for Xcode (`xcode-select --install`) for required headers and libraries.
	* On Ubuntu/Debian-like Linux distributions, you need `libgl1-mesa-dev` and `xorg-dev` packages.
	* On CentOS/Fedora-like Linux distributions, you need `libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel` packages.
	* See [here](http://www.glfw.org/docs/latest/compile.html#compile_deps) for full details.
* Go 1.4+ is required on Windows (otherwise you must use MinGW v4.8.1 exactly, see [Go issue 8811](https://github.com/golang/go/issues/8811)).

```
go get -u github.com/go-gl/glfw/v3.2/glfw
```

## Usage

```Go
package main

import (
	"runtime"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
```

## Changelog

* Easy `go get` installation. GLFW source code is now included in-repo and compiled in so you don't have to build GLFW on your own and distribute shared libraries. The revision of GLFW C library used is listed in [GLFW_C_REVISION.txt](https://github.com/go-gl/glfw/blob/master/v3.2/glfw/GLFW_C_REVISION.txt) file.
* The error callback is now set internally. Functions return an error with corresponding code and description (do a type assertion to glfw3.Error for accessing the variables) if the error is recoverable. If not a panic will occur.

### GLFW 3.2 Specfic Changes
* Added function `Window.SetSizeLimits`.
* Added function `Window.SetAspectRatio`.
* Added function `Window.SetMonitor`.
* Added function `Window.Maximize`.
* Added function `Window.SetIcon`.
* Added function `Window.Focus`.
* Added function `GetKeyName`.
* Added function `VulkanSupported`.
* Added function `GetTimerValue`.
* Added function `GetTimerFrequency`.
* Added function `WaitEventsTimeout`.
* Added function `SetJoystickCallback`.
* Added window hint `Maximized`.
* Added hint `NoAPI`.
* Added hint `NativeContextAPI`.
* Added hint `EGLContextAPI`.

### GLFW 3.1 Specfic Changes
* Added type `Cursor`.
* Added function `Window.SetDropCallback`.
* Added function `Window.SetCharModsCallback`.
* Added function `PostEmptyEvent`.
* Added function `CreateCursor`.
* Added function `CreateStandardCursor`.
* Added function `Cursor.Destroy`.
* Added function `Window.SetCursor`.
* Added function `Window.GetFrameSize`.
* Added window hint `Floating`.
* Added window hint `AutoIconify`.
* Added window hint `ContextReleaseBehavior`.
* Added window hint `DoubleBuffer`.
* Added hint value `AnyReleaseBehavior`.
* Added hint value `ReleaseBehaviorFlush`.
* Added hint value `ReleaseBehaviorNone`.
* Added hint value `DontCare`.

### API changes
* `Window.Iconify` Returns an error.
* `Window.Restore` Returns an error.
* `Init` Returns an error instead of `bool`.
* `GetJoystickAxes` No longer returns an error.
* `GetJoystickButtons` No longer returns an error.
* `GetJoystickName` No longer returns an error.
* `GetMonitors` No longer returns an error.
* `GetPrimaryMonitor` No longer returns an error.
* `Monitor.GetGammaRamp` No longer returns an error.
* `Monitor.GetVideoMode` No longer returns an error.
* `Monitor.GetVideoModes` No longer returns an error.
* `GetCurrentContext` No longer returns an error.
* `Window.SetCharCallback` Accepts `rune` instead of `uint`.
* Added type `Error`.
* Removed `SetErrorCallback`.
* Removed error code `NotInitialized`.
* Removed error code `NoCurrentContext`.
* Removed error code `InvalidEnum`.
* Removed error code `InvalidValue`.
* Removed error code `OutOfMemory`.
* Removed error code `PlatformError`.
* Removed `KeyBracket`.
* Renamed `Window.SetCharacterCallback` to `Window.SetCharCallback`.
* Renamed `Window.GetCursorPosition` to `GetCursorPos`.
* Renamed `Window.SetCursorPosition` to `SetCursorPos`.
* Renamed `CursorPositionCallback` to `CursorPosCallback`.
* Renamed `Window.SetCursorPositionCallback` to `SetCursorPosCallback`.
* Renamed `VideoMode` to `VidMode`.
* Renamed `Monitor.GetPosition` to `Monitor.GetPos`.
* Renamed `Window.GetPosition` to `Window.GetPos`.
* Renamed `Window.SetPosition` to `Window.SetPos`.
* Renamed `Window.GetAttribute` to `Window.GetAttrib`.
* Renamed `Window.SetPositionCallback` to `Window.SetPosCallback`.
* Renamed `PositionCallback` to `PosCallback`.
* Ranamed `Cursor` to `CursorMode`.
* Renamed `StickyKeys` to `StickyKeysMode`.
* Renamed `StickyMouseButtons` to `StickyMouseButtonsMode`.
* Renamed `ApiUnavailable` to `APIUnavailable`.
* Renamed `ClientApi` to `ClientAPI`.
* Renamed `OpenglForwardCompatible` to `OpenGLForwardCompatible`.
* Renamed `OpenglDebugContext` to `OpenGLDebugContext`.
* Renamed `OpenglProfile` to `OpenGLProfile`.
* Renamed `SrgbCapable` to `SRGBCapable`.
* Renamed `OpenglApi` to `OpenGLAPI`.
* Renamed `OpenglEsApi` to `OpenGLESAPI`.
* Renamed `OpenglAnyProfile` to `OpenGLAnyProfile`.
* Renamed `OpenglCoreProfile` to `OpenGLCoreProfile`.
* Renamed `OpenglCompatProfile` to `OpenGLCompatProfile`.
* Renamed `KeyKp...` to `KeyKP...`.
