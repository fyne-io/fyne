package glfw

//#include <stdlib.h>
//#include <GLFW/glfw3.h>
//void glfwSetWindowPosCallbackCB(GLFWwindow *window);
//void glfwSetWindowSizeCallbackCB(GLFWwindow *window);
//void glfwSetFramebufferSizeCallbackCB(GLFWwindow *window);
//void glfwSetWindowCloseCallbackCB(GLFWwindow *window);
//void glfwSetWindowRefreshCallbackCB(GLFWwindow *window);
//void glfwSetWindowFocusCallbackCB(GLFWwindow *window);
//void glfwSetWindowIconifyCallbackCB(GLFWwindow *window);
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

// Internal window list stuff
type windowList struct {
	l sync.Mutex
	m map[*C.GLFWwindow]*Window
}

var windows = windowList{m: map[*C.GLFWwindow]*Window{}}

func (w *windowList) put(wnd *Window) {
	w.l.Lock()
	defer w.l.Unlock()
	w.m[wnd.data] = wnd
}

func (w *windowList) remove(wnd *C.GLFWwindow) {
	w.l.Lock()
	defer w.l.Unlock()
	delete(w.m, wnd)
}

func (w *windowList) get(wnd *C.GLFWwindow) *Window {
	w.l.Lock()
	defer w.l.Unlock()
	return w.m[wnd]
}

//Hint corresponds to hints that can be set before creating a window.
//
//Hint also corresponds to the attributes of the window that can be get after
//its creation.
type Hint int

//Window related hints.
const (
	Focused   Hint = C.GLFW_FOCUSED   //Specifies whether the window will be focused.
	Iconified Hint = C.GLFW_ICONIFIED //Specifies whether the window will be minimized.
	Visible   Hint = C.GLFW_VISIBLE   //Specifies whether the window will be initially visible.
	Resizable Hint = C.GLFW_RESIZABLE //Specifies whether the window will be resizable by the user.
	Decorated Hint = C.GLFW_DECORATED //Specifies whether the window will have window decorations such as a border, a close widget, etc.
)

//Context related hints.
const (
	ClientApi               Hint = C.GLFW_CLIENT_API            //Specifies which client API to create the context for. Hard constraint.
	ContextVersionMajor     Hint = C.GLFW_CONTEXT_VERSION_MAJOR //Specifies the client API version that the created context must be compatible with.
	ContextVersionMinor     Hint = C.GLFW_CONTEXT_VERSION_MINOR //Specifies the client API version that the created context must be compatible with.
	ContextRobustness       Hint = C.GLFW_CONTEXT_ROBUSTNESS    //Specifies the robustness strategy to be used by the context.
	OpenglForwardCompatible Hint = C.GLFW_OPENGL_FORWARD_COMPAT //Specifies whether the OpenGL context should be forward-compatible. Hard constraint.
	OpenglDebugContext      Hint = C.GLFW_OPENGL_DEBUG_CONTEXT
	OpenglProfile           Hint = C.GLFW_OPENGL_PROFILE //Specifies which OpenGL profile to create the context for. Hard constraint.
)

//Framebuffer related hints.
const (
	ContextRevision Hint = C.GLFW_CONTEXT_REVISION
	RedBits         Hint = C.GLFW_RED_BITS         //Specifies the desired bit depth of the default framebuffer.
	GreenBits       Hint = C.GLFW_GREEN_BITS       //Specifies the desired bit depth of the default framebuffer.
	BlueBits        Hint = C.GLFW_BLUE_BITS        //Specifies the desired bit depth of the default framebuffer.
	AlphaBits       Hint = C.GLFW_ALPHA_BITS       //Specifies the desired bit depth of the default framebuffer.
	DepthBits       Hint = C.GLFW_DEPTH_BITS       //Specifies the desired bit depth of the default framebuffer.
	StencilBits     Hint = C.GLFW_STENCIL_BITS     //Specifies the desired bit depth of the default framebuffer.
	AccumRedBits    Hint = C.GLFW_ACCUM_RED_BITS   //Specifies the desired bit depth of the accumulation buffer.
	AccumGreenBits  Hint = C.GLFW_ACCUM_GREEN_BITS //Specifies the desired bit depth of the accumulation buffer.
	AccumBlueBits   Hint = C.GLFW_ACCUM_BLUE_BITS  //Specifies the desired bit depth of the accumulation buffer.
	AccumAlphaBits  Hint = C.GLFW_ACCUM_ALPHA_BITS //Specifies the desired bit depth of the accumulation buffer.
	AuxBuffers      Hint = C.GLFW_AUX_BUFFERS      //Specifies the desired number of auxiliary buffers.
	Stereo          Hint = C.GLFW_STEREO           //Specifies whether to use stereoscopic rendering. Hard constraint.
	Samples         Hint = C.GLFW_SAMPLES          //Specifies the desired number of samples to use for multisampling. Zero disables multisampling.
	SrgbCapable     Hint = C.GLFW_SRGB_CAPABLE     //Specifies whether the framebuffer should be sRGB capable.
	RefreshRate     Hint = C.GLFW_REFRESH_RATE     //specifies the desired refresh rate for full screen windows. If set to zero, the highest available refresh rate will be used. This hint is ignored for windowed mode windows.
)

//Values for the ClientApi hint.
const (
	OpenglApi   int = C.GLFW_OPENGL_API
	OpenglEsApi int = C.GLFW_OPENGL_ES_API
)

//Values for the ContextRobustness hint.
const (
	NoRobustness        int = C.GLFW_NO_ROBUSTNESS
	NoResetNotification int = C.GLFW_NO_RESET_NOTIFICATION
	LoseContextOnReset  int = C.GLFW_LOSE_CONTEXT_ON_RESET
)

//Values for the OpenglProfile hint.
const (
	OpenglAnyProfile    int = C.GLFW_OPENGL_ANY_PROFILE
	OpenglCoreProfile   int = C.GLFW_OPENGL_CORE_PROFILE
	OpenglCompatProfile int = C.GLFW_OPENGL_COMPAT_PROFILE
)

//TRUE and FALSE values to use with hints.
const (
	True  int = C.GL_TRUE
	False int = C.GL_FALSE
)

type Window struct {
	data *C.GLFWwindow

	// Window
	fPosHolder             func(w *Window, xpos int, ypos int)
	fSizeHolder            func(w *Window, width int, height int)
	fFramebufferSizeHolder func(w *Window, width int, height int)
	fCloseHolder           func(w *Window)
	fRefreshHolder         func(w *Window)
	fFocusHolder           func(w *Window, focused bool)
	fIconifyHolder         func(w *Window, iconified bool)

	// Input
	fMouseButtonHolder func(w *Window, button MouseButton, action Action, mod ModifierKey)
	fCursorPosHolder   func(w *Window, xpos float64, ypos float64)
	fCursorEnterHolder func(w *Window, entered bool)
	fScrollHolder      func(w *Window, xoff float64, yoff float64)
	fKeyHolder         func(w *Window, key Key, scancode int, action Action, mods ModifierKey)
	fCharHolder        func(w *Window, char uint)
}

//export goWindowPosCB
func goWindowPosCB(window unsafe.Pointer, xpos, ypos C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fPosHolder(w, int(xpos), int(ypos))
}

//export goWindowSizeCB
func goWindowSizeCB(window unsafe.Pointer, width, height C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fSizeHolder(w, int(width), int(height))
}

//export goFramebufferSizeCB
func goFramebufferSizeCB(window unsafe.Pointer, width, height C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fFramebufferSizeHolder(w, int(width), int(height))
}

//export goWindowCloseCB
func goWindowCloseCB(window unsafe.Pointer) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fCloseHolder(w)
}

//export goWindowRefreshCB
func goWindowRefreshCB(window unsafe.Pointer) {
	w := windows.get((*C.GLFWwindow)(window))
	w.fRefreshHolder(w)
}

//export goWindowFocusCB
func goWindowFocusCB(window unsafe.Pointer, focused C.int) {
	w := windows.get((*C.GLFWwindow)(window))
	isFocused := glfwbool(focused)
	w.fFocusHolder(w, isFocused)
}

//export goWindowIconifyCB
func goWindowIconifyCB(window unsafe.Pointer, iconified C.int) {
	isIconified := glfwbool(iconified)
	w := windows.get((*C.GLFWwindow)(window))
	w.fIconifyHolder(w, isIconified)
}

//DefaultHints resets all window hints to their default values.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func DefaultWindowHints() {
	C.glfwDefaultWindowHints()
}

//Hint function sets hints for the next call to CreateWindow. The hints,
//once set, retain their values until changed by a call to Hint or
//DefaultHints, or until the library is terminated with Terminate.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func WindowHint(target Hint, hint int) {
	C.glfwWindowHint(C.int(target), C.int(hint))
}

//CreateWindow creates a window and its associated context. Most of the options
//controlling how the window and its context should be created are specified
//through Hint.
//
//Successful creation does not change which context is current. Before you can
//use the newly created context, you need to make it current using
//MakeContextCurrent.
//
//Note that the created window and context may differ from what you requested,
//as not all parameters and hints are hard constraints. This includes the size
//of the window, especially for full screen windows. To retrieve the actual
//attributes of the created window and context, use queries like
//GetWindowAttrib and GetWindowSize.
//
//To create the window at a specific position, make it initially invisible using
//the Visible window hint, set its position and then show it.
//
//If a fullscreen window is active, the screensaver is prohibited from starting.
//
//Windows: If the executable has an icon resource named GLFW_ICON, it will be
//set as the icon for the window. If no such icon is present, the IDI_WINLOGO
//icon will be used instead.
//
//Mac OS X: The GLFW window has no icon, as it is not a document window, but the
//dock icon will be the same as the application bundle's icon. Also, the first
//time a window is opened the menu bar is populated with common commands like
//Hide, Quit and About. The (minimal) about dialog uses information from the
//application's bundle. For more information on bundles, see the Bundle
//Programming Guide provided by Apple.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var (
		m *C.GLFWmonitor
		s *C.GLFWwindow
	)

	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))

	if monitor != nil {
		m = monitor.data
	}

	if share != nil {
		s = share.data
	}

	w := C.glfwCreateWindow(C.int(width), C.int(height), t, m, s)

	if w == nil {
		return nil, errors.New("Can't create window.")
	}
	wnd := &Window{data: w}
	windows.put(wnd)
	return wnd, nil
}

//Destroy destroys the specified window and its context. On calling this
//function, no further callbacks will be called for that window.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) Destroy() {
	windows.remove(w.data)
	C.glfwDestroyWindow(w.data)
}

//ShouldClose returns the value of the close flag of the specified window.
func (w *Window) ShouldClose() bool {
	return glfwbool(C.glfwWindowShouldClose(w.data))
}

//SetShouldClose sets the value of the close flag of the window. This can be
//used to override the user's attempt to close the window, or to signal that it
//should be closed.
func (w *Window) SetShouldClose(value bool) {
	if !value {
		C.glfwSetWindowShouldClose(w.data, C.GL_FALSE)
	} else {
		C.glfwSetWindowShouldClose(w.data, C.GL_TRUE)
	}
}

//SetTitle sets the window title, encoded as UTF-8, of the window.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) SetTitle(title string) {
	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))
	C.glfwSetWindowTitle(w.data, t)
}

//GetPosition returns the position, in screen coordinates, of the upper-left
//corner of the client area of the window.
func (w *Window) GetPosition() (x, y int) {
	var xpos, ypos C.int

	C.glfwGetWindowPos(w.data, &xpos, &ypos)
	return int(xpos), int(ypos)
}

//SetPosition sets the position, in screen coordinates, of the upper-left corner
//of the client area of the window.
//
//If it is a full screen window, this function does nothing.
//
//If you wish to set an initial window position you should create a hidden
//window (using Hint and Visible), set its position and then show it.
//
//It is very rarely a good idea to move an already visible window, as it will
//confuse and annoy the user.
//
//The window manager may put limits on what positions are allowed.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) SetPosition(xpos, ypos int) {
	C.glfwSetWindowPos(w.data, C.int(xpos), C.int(ypos))
}

//GetSize returns the size, in screen coordinates, of the client area of the
//specified window.
func (w *Window) GetSize() (width, height int) {
	var wi, h C.int
	C.glfwGetWindowSize(w.data, &wi, &h)
	return int(wi), int(h)
}

//SetSize sets the size, in screen coordinates, of the client area of the
//window.
//
//For full screen windows, this function selects and switches to the resolution
//closest to the specified size, without affecting the window's context. As the
//context is unaffected, the bit depths of the framebuffer remain unchanged.
//
//The window manager may put limits on what window sizes are allowed.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) SetSize(width, height int) {
	C.glfwSetWindowSize(w.data, C.int(width), C.int(height))
}

//GetFramebufferSize retrieves the size, in pixels, of the framebuffer of the
//specified window.
func (w *Window) GetFramebufferSize() (width, height int) {
	var wi, h C.int
	C.glfwGetFramebufferSize(w.data, &wi, &h)
	return int(wi), int(h)
}

//Iconfiy iconifies/minimizes the window, if it was previously restored. If it
//is a full screen window, the original monitor resolution is restored until the
//window is restored. If the window is already iconified, this function does
//nothing.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) Iconify() {
	C.glfwIconifyWindow(w.data)
}

//Restore restores the window, if it was previously iconified/minimized. If it
//is a full screen window, the resolution chosen for the window is restored on
//the selected monitor. If the window is already restored, this function does
//nothing.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) Restore() {
	C.glfwRestoreWindow(w.data)
}

//Show makes the window visible, if it was previously hidden. If the window is
//already visible or is in full screen mode, this function does nothing.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) Show() {
	C.glfwShowWindow(w.data)
}

//Hide hides the window, if it was previously visible. If the window is already
//hidden or is in full screen mode, this function does nothing.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func (w *Window) Hide() {
	C.glfwHideWindow(w.data)
}

//GetMonitor returns the handle of the monitor that the window is in
//fullscreen on.
func (w *Window) GetMonitor() (*Monitor, error) {
	m := C.glfwGetWindowMonitor(w.data)

	if m == nil {
		return nil, errors.New("Can't get the monitor.")
	}
	return &Monitor{m}, nil
}

//GetAttribute returns an attribute of the window. There are many attributes,
//some related to the window and others to its context.
func (w *Window) GetAttribute(attrib Hint) int {
	return int(C.glfwGetWindowAttrib(w.data, C.int(attrib)))
}

//SetUserPointer sets the user-defined pointer of the window. The current value
//is retained until the window is destroyed. The initial value is nil.
func (w *Window) SetUserPointer(pointer unsafe.Pointer) {
	C.glfwSetWindowUserPointer(w.data, pointer)
}

//GetUserPointer returns the current value of the user-defined pointer of the
//window. The initial value is nil.
func (w *Window) GetUserPointer() unsafe.Pointer {
	return C.glfwGetWindowUserPointer(w.data)
}

//SetPositionCallback sets the position callback of the window, which is called
//when the window is moved. The callback is provided with the screen position
//of the upper-left corner of the client area of the window.
func (w *Window) SetPositionCallback(cbfun func(w *Window, xpos int, ypos int)) {
	if cbfun == nil {
		C.glfwSetWindowPosCallback(w.data, nil)
	} else {
		w.fPosHolder = cbfun
		C.glfwSetWindowPosCallbackCB(w.data)
	}
}

//SetSizeCallback sets the size callback of the window, which is called when
//the window is resized. The callback is provided with the size, in screen
//coordinates, of the client area of the window.
func (w *Window) SetSizeCallback(cbfun func(w *Window, width int, height int)) {
	if cbfun == nil {
		C.glfwSetWindowSizeCallback(w.data, nil)
	} else {
		w.fSizeHolder = cbfun
		C.glfwSetWindowSizeCallbackCB(w.data)
	}
}

//SetFramebufferSizeCallback sets the framebuffer resize callback of the specified
//window, which is called when the framebuffer of the specified window is resized.
func (w *Window) SetFramebufferSizeCallback(cbfun func(w *Window, width int, height int)) {
	if cbfun == nil {
		C.glfwSetFramebufferSizeCallback(w.data, nil)
	} else {
		w.fFramebufferSizeHolder = cbfun
		C.glfwSetFramebufferSizeCallbackCB(w.data)
	}
}

//SetCloseCallback sets the close callback of the window, which is called when
//the user attempts to close the window, for example by clicking the close
//widget in the title bar.
//
//The close flag is set before this callback is called, but you can modify it at
//any time with SetShouldClose.
//
//Mac OS X: Selecting Quit from the application menu will trigger the close
//callback for all windows.
func (w *Window) SetCloseCallback(cbfun func(w *Window)) {
	if cbfun == nil {
		C.glfwSetWindowCloseCallback(w.data, nil)
	} else {
		w.fCloseHolder = cbfun
		C.glfwSetWindowCloseCallbackCB(w.data)
	}
}

//SetRefreshCallback sets the refresh callback of the window, which
//is called when the client area of the window needs to be redrawn, for example
//if the window has been exposed after having been covered by another window.
//
//On compositing window systems such as Aero, Compiz or Aqua, where the window
//contents are saved off-screen, this callback may be called only very
//infrequently or never at all.
func (w *Window) SetRefreshCallback(cbfun func(w *Window)) {
	if cbfun == nil {
		C.glfwSetWindowRefreshCallback(w.data, nil)
	} else {
		w.fRefreshHolder = cbfun
		C.glfwSetWindowRefreshCallbackCB(w.data)
	}
}

//SetFocusCallback sets the focus callback of the window, which is called when
//the window gains or loses focus.
//
//After the focus callback is called for a window that lost focus, synthetic key
//and mouse button release events will be generated for all such that had been
//pressed. For more information, see SetKeyCallback and SetMouseButtonCallback.
func (w *Window) SetFocusCallback(cbfun func(w *Window, focused bool)) {
	if cbfun == nil {
		C.glfwSetWindowFocusCallback(w.data, nil)
	} else {
		w.fFocusHolder = cbfun
		C.glfwSetWindowFocusCallbackCB(w.data)
	}
}

//SetIconifyCallback sets the iconification callback of the window, which is
//called when the window is iconified or restored.
func (w *Window) SetIconifyCallback(cbfun func(w *Window, iconified bool)) {
	if cbfun == nil {
		C.glfwSetWindowIconifyCallback(w.data, nil)
	} else {
		w.fIconifyHolder = cbfun
		C.glfwSetWindowIconifyCallbackCB(w.data)
	}
}

//PollEvents processes only those events that have already been received and
//then returns immediately. Processing events will cause the window and input
//callbacks associated with those events to be called.
//
//This function is not required for joystick input to work.
//
//This function may not be called from a callback.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func PollEvents() {
	C.glfwPollEvents()
}

//WaitEvents puts the calling thread to sleep until at least one event has been
//received. Once one or more events have been recevied, it behaves as if
//PollEvents was called, i.e. the events are processed and the function then
//returns immediately. Processing events will cause the window and input
//callbacks associated with those events to be called.
//
//Since not all events are associated with callbacks, this function may return
//without a callback having been called even if you are monitoring all
//callbacks.
//
//This function may not be called from a callback.
//
//This function may only be called from the main thread. See
//https://code.google.com/p/go-wiki/wiki/LockOSThread
func WaitEvents() {
	C.glfwWaitEvents()
}
