package glfw

//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
//GLFWmonitor* GetMonitorAtIndex(GLFWmonitor **monitors, int index);
//GLFWvidmode GetVidmodeAtIndex(GLFWvidmode *vidmodes, int index);
//void glfwSetMonitorCallbackCB();
//unsigned int GetGammaAtIndex(unsigned short *color, int i);
//void SetGammaAtIndex(unsigned short *color, int i, unsigned short value);
import "C"

import (
	"unsafe"
)

// Monitor represents a monitor.
type Monitor struct {
	data *C.GLFWmonitor
}

// PeripheralEvent corresponds to a peripheral(Monitor or Joystick)
// configuration event.
type PeripheralEvent int

// GammaRamp describes the gamma ramp for a monitor.
type GammaRamp struct {
	Red   []uint16 // A slice of value describing the response of the red channel.
	Green []uint16 // A slice of value describing the response of the green channel.
	Blue  []uint16 // A slice of value describing the response of the blue channel.
}

// PeripheralEvent events.
const (
	Connected    PeripheralEvent = C.GLFW_CONNECTED
	Disconnected PeripheralEvent = C.GLFW_DISCONNECTED
)

// VidMode describes a single video mode.
type VidMode struct {
	Width       int // The width, in pixels, of the video mode.
	Height      int // The height, in pixels, of the video mode.
	RedBits     int // The bit depth of the red channel of the video mode.
	GreenBits   int // The bit depth of the green channel of the video mode.
	BlueBits    int // The bit depth of the blue channel of the video mode.
	RefreshRate int // The refresh rate, in Hz, of the video mode.
}

var fMonitorHolder func(monitor *Monitor, event PeripheralEvent)

//export goMonitorCB
func goMonitorCB(monitor unsafe.Pointer, event C.int) {
	fMonitorHolder(&Monitor{(*C.GLFWmonitor)(monitor)}, PeripheralEvent(event))
}

// GetMonitors returns a slice of handles for all currently connected monitors.
func GetMonitors() []*Monitor {
	var length int

	mC := C.glfwGetMonitors((*C.int)(unsafe.Pointer(&length)))
	panicError()
	if mC == nil {
		return nil
	}

	m := make([]*Monitor, length)

	for i := 0; i < length; i++ {
		m[i] = &Monitor{C.GetMonitorAtIndex(mC, C.int(i))}
	}

	return m
}

// GetPrimaryMonitor returns the primary monitor. This is usually the monitor
// where elements like the Windows task bar or the OS X menu bar is located.
func GetPrimaryMonitor() *Monitor {
	m := C.glfwGetPrimaryMonitor()
	panicError()
	if m == nil {
		return nil
	}
	return &Monitor{m}
}

// GetPos returns the position, in screen coordinates, of the upper-left
// corner of the monitor.
func (m *Monitor) GetPos() (x, y int) {
	var xpos, ypos C.int
	C.glfwGetMonitorPos(m.data, &xpos, &ypos)
	panicError()
	return int(xpos), int(ypos)
}

// GetWorkarea returns the position, in screen coordinates, of the upper-left
// corner of the work area of the specified monitor along with the work area
// size in screen coordinates. The work area is defined as the area of the
// monitor not occluded by the operating system task bar where present. If no
// task bar exists then the work area is the monitor resolution in screen
// coordinates.
//
// This function must only be called from the main thread.
func (m *Monitor) GetWorkarea() (x, y, width, height int) {
	var cX, cY, cWidth, cHeight C.int
	C.glfwGetMonitorWorkarea(m.data, &cX, &cY, &cWidth, &cHeight)
	x, y, width, height = int(cX), int(cY), int(cWidth), int(cHeight)
	return
}

// GetContentScale function retrieves the content scale for the specified monitor.
// The content scale is the ratio between the current DPI and the platform's
// default DPI. If you scale all pixel dimensions by this scale then your content
// should appear at an appropriate size. This is especially important for text
// and any UI elements.
//
// This function must only be called from the main thread.
func (m *Monitor) GetContentScale() (float32, float32) {
	var x, y C.float
	C.glfwGetMonitorContentScale(m.data, &x, &y)
	return float32(x), float32(y)
}

// SetUserPointer sets the user-defined pointer of the monitor. The current value
// is retained until the monitor is disconnected. The initial value is nil.
//
// This function may be called from the monitor callback, even for a monitor
// that is being disconnected.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) SetUserPointer(pointer unsafe.Pointer) {
	C.glfwSetMonitorUserPointer(m.data, pointer)
}

// GetUserPointer returns the current value of the user-defined pointer of the
// monitor. The initial value is nil.
//
// This function may be called from the monitor callback, even for a monitor
// that is being disconnected.
//
// This function may be called from any thread. Access is not synchronized.
func (m *Monitor) GetUserPointer() unsafe.Pointer {
	return C.glfwGetMonitorUserPointer(m.data)
}

// GetPhysicalSize returns the size, in millimetres, of the display area of the
// monitor.
//
// Note: Some operating systems do not provide accurate information, either
// because the monitor's EDID data is incorrect, or because the driver does not
// report it accurately.
func (m *Monitor) GetPhysicalSize() (width, height int) {
	var wi, h C.int
	C.glfwGetMonitorPhysicalSize(m.data, &wi, &h)
	panicError()
	return int(wi), int(h)
}

// GetName returns a human-readable name of the monitor, encoded as UTF-8.
func (m *Monitor) GetName() string {
	mn := C.glfwGetMonitorName(m.data)
	panicError()
	if mn == nil {
		return ""
	}
	return C.GoString(mn)
}

// MonitorCallback is the signature for monitor configuration callback
// functions.
type MonitorCallback func(monitor *Monitor, event PeripheralEvent)

// SetMonitorCallback sets the monitor configuration callback, or removes the
// currently set callback. This is called when a monitor is connected to or
// disconnected from the system.
//
// This function must only be called from the main thread.
func SetMonitorCallback(cbfun MonitorCallback) MonitorCallback {
	previous := fMonitorHolder
	fMonitorHolder = cbfun
	if cbfun == nil {
		C.glfwSetMonitorCallback(nil)
	} else {
		C.glfwSetMonitorCallbackCB()
	}
	return previous
}

// GetVideoModes returns an array of all video modes supported by the monitor.
// The returned array is sorted in ascending order, first by color bit depth
// (the sum of all channel depths) and then by resolution area (the product of
// width and height).
func (m *Monitor) GetVideoModes() []*VidMode {
	var length int

	vC := C.glfwGetVideoModes(m.data, (*C.int)(unsafe.Pointer(&length)))
	panicError()
	if vC == nil {
		return nil
	}

	v := make([]*VidMode, length)

	for i := 0; i < length; i++ {
		t := C.GetVidmodeAtIndex(vC, C.int(i))
		v[i] = &VidMode{int(t.width), int(t.height), int(t.redBits), int(t.greenBits), int(t.blueBits), int(t.refreshRate)}
	}

	return v
}

// GetVideoMode returns the current video mode of the monitor. If you
// are using a full screen window, the return value will therefore depend on
// whether it is focused.
func (m *Monitor) GetVideoMode() *VidMode {
	t := C.glfwGetVideoMode(m.data)
	if t == nil {
		return nil
	}
	panicError()
	return &VidMode{int(t.width), int(t.height), int(t.redBits), int(t.greenBits), int(t.blueBits), int(t.refreshRate)}
}

// SetGamma generates a 256-element gamma ramp from the specified exponent and then calls
// SetGamma with it.
func (m *Monitor) SetGamma(gamma float32) {
	C.glfwSetGamma(m.data, C.float(gamma))
	panicError()
}

// GetGammaRamp retrieves the current gamma ramp of the monitor.
func (m *Monitor) GetGammaRamp() *GammaRamp {
	var ramp GammaRamp

	rampC := C.glfwGetGammaRamp(m.data)
	panicError()
	if rampC == nil {
		return nil
	}

	length := int(rampC.size)
	ramp.Red = make([]uint16, length)
	ramp.Green = make([]uint16, length)
	ramp.Blue = make([]uint16, length)

	for i := 0; i < length; i++ {
		ramp.Red[i] = uint16(C.GetGammaAtIndex(rampC.red, C.int(i)))
		ramp.Green[i] = uint16(C.GetGammaAtIndex(rampC.green, C.int(i)))
		ramp.Blue[i] = uint16(C.GetGammaAtIndex(rampC.blue, C.int(i)))
	}

	return &ramp
}

// SetGammaRamp sets the current gamma ramp for the monitor.
func (m *Monitor) SetGammaRamp(ramp *GammaRamp) {
	var rampC C.GLFWgammaramp

	length := len(ramp.Red)

	for i := 0; i < length; i++ {
		C.SetGammaAtIndex(rampC.red, C.int(i), C.ushort(ramp.Red[i]))
		C.SetGammaAtIndex(rampC.green, C.int(i), C.ushort(ramp.Green[i]))
		C.SetGammaAtIndex(rampC.blue, C.int(i), C.ushort(ramp.Blue[i]))
	}

	C.glfwSetGammaRamp(m.data, &rampC)
	panicError()
}
