package glfw

//#include "glfw/include/GLFW/glfw3.h"
import "C"

// GetTime returns the value of the GLFW timer. Unless the timer has been set
// using SetTime, the timer measures time elapsed since GLFW was initialized.
//
// The resolution of the timer is system dependent, but is usually on the order
// of a few micro- or nanoseconds. It uses the highest-resolution monotonic time
// source on each supported platform.
func GetTime() float64 {
	ret := float64(C.glfwGetTime())
	panicError()
	return ret
}

// SetTime sets the value of the GLFW timer. It then continues to count up from
// that value.
//
// The resolution of the timer is system dependent, but is usually on the order
// of a few micro- or nanoseconds. It uses the highest-resolution monotonic time
// source on each supported platform.
func SetTime(time float64) {
	C.glfwSetTime(C.double(time))
	panicError()
}

// GetTimerFrequency returns frequency of the timer, in Hz, or zero if an error occurred.
func GetTimerFrequency() uint64 {
	ret := uint64(C.glfwGetTimerFrequency())
	panicError()
	return ret
}

// GetTimerValue returns the current value of the raw timer, measured in 1 / frequency seconds.
func GetTimerValue() uint64 {
	ret := uint64(C.glfwGetTimerValue())
	panicError()
	return ret
}
