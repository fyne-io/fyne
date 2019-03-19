package glfw

//#include "glfw/include/GLFW/glfw3.h"
//void glfwSetErrorCallbackCB();
import "C"

import (
	"fmt"
	"log"
)

// ErrorCode corresponds to an error code.
type ErrorCode int

// Error codes that are translated to panics and the programmer should not
// expect to handle.
const (
	notInitialized   ErrorCode = C.GLFW_NOT_INITIALIZED    // GLFW has not been initialized.
	noCurrentContext ErrorCode = C.GLFW_NO_CURRENT_CONTEXT // No context is current.
	invalidEnum      ErrorCode = C.GLFW_INVALID_ENUM       // One of the enum parameters for the function was given an invalid enum.
	invalidValue     ErrorCode = C.GLFW_INVALID_VALUE      // One of the parameters for the function was given an invalid value.
	outOfMemory      ErrorCode = C.GLFW_OUT_OF_MEMORY      // A memory allocation failed.
	platformError    ErrorCode = C.GLFW_PLATFORM_ERROR     // A platform-specific error occurred that does not match any of the more specific categories.
)

const (
	// APIUnavailable is the error code used when GLFW could not find support
	// for the requested client API on the system.
	//
	// The installed graphics driver does not support the requested client API,
	// or does not support it via the chosen context creation backend. Below
	// are a few examples.
	//
	// Some pre-installed Windows graphics drivers do not support OpenGL. AMD
	// only supports OpenGL ES via EGL, while Nvidia and Intel only supports it
	// via a WGL or GLX extension. OS X does not provide OpenGL ES at all. The
	// Mesa EGL, OpenGL and OpenGL ES libraries do not interface with the
	// Nvidia binary driver.
	APIUnavailable ErrorCode = C.GLFW_API_UNAVAILABLE

	// VersionUnavailable is the error code used when the requested OpenGL or
	// OpenGL ES (including any requested profile or context option) is not
	// available on this machine.
	//
	// The machine does not support your requirements. If your application is
	// sufficiently flexible, downgrade your requirements and try again.
	// Otherwise, inform the user that their machine does not match your
	// requirements.
	//
	// Future invalid OpenGL and OpenGL ES versions, for example OpenGL 4.8 if
	// 5.0 comes out before the 4.x series gets that far, also fail with this
	// error and not GLFW_INVALID_VALUE, because GLFW cannot know what future
	// versions will exist.
	VersionUnavailable ErrorCode = C.GLFW_VERSION_UNAVAILABLE

	// FormatUnavailable is the error code used for both window creation and
	// clipboard querying format errors.
	//
	// If emitted during window creation, the requested pixel format is not
	// supported. This means one or more hard constraints did not match any of
	// the available pixel formats. If your application is sufficiently
	// flexible, downgrade your requirements and try again. Otherwise, inform
	// the user that their machine does not match your requirements.
	//
	// If emitted when querying the clipboard, the contents of the clipboard
	// could not be converted to the requested format. You should ignore the
	// error or report it to the user, as appropriate.
	FormatUnavailable ErrorCode = C.GLFW_FORMAT_UNAVAILABLE
)

func (e ErrorCode) String() string {
	switch e {
	case notInitialized:
		return "NotInitialized"
	case noCurrentContext:
		return "NoCurrentContext"
	case invalidEnum:
		return "InvalidEnum"
	case invalidValue:
		return "InvalidValue"
	case outOfMemory:
		return "OutOfMemory"
	case platformError:
		return "PlatformError"
	case APIUnavailable:
		return "APIUnavailable"
	case VersionUnavailable:
		return "VersionUnavailable"
	case FormatUnavailable:
		return "FormatUnavailable"
	default:
		return fmt.Sprintf("ErrorCode(%d)", e)
	}
}

// Error holds error code and description.
type Error struct {
	Code ErrorCode
	Desc string
}

// Error prints the error code and description in a readable format.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code.String(), e.Desc)
}

// Note: There are many cryptic caveats to proper error handling here.
// See: https://github.com/go-gl/glfw3/pull/86

// Holds the value of the last error.
var lastError = make(chan *Error, 1)

//export goErrorCB
func goErrorCB(code C.int, desc *C.char) {
	flushErrors()
	err := &Error{ErrorCode(code), C.GoString(desc)}
	select {
	case lastError <- err:
	default:
		fmt.Println("GLFW: An uncaught error has occurred:", err)
		fmt.Println("GLFW: Please report this bug in the Go package immediately.")
	}
}

// Set the glfw callback internally
func init() {
	C.glfwSetErrorCallbackCB()
}

// flushErrors is called by Terminate before it actually calls C.glfwTerminate,
// this ensures that any uncaught errors buffered in lastError are printed
// before the program exits.
func flushErrors() {
	err := fetchError()
	if err != nil {
		fmt.Println("GLFW: An uncaught error has occurred:", err)
		fmt.Println("GLFW: Please report this bug in the Go package immediately.")
	}
}

// acceptError fetches the next error from the error channel, it accepts only
// errors with one of the given error codes. If any other error is encountered,
// a panic will occur.
//
// Platform errors are always printed, for information why please see:
//
//  https://github.com/go-gl/glfw/issues/127
//
func acceptError(codes ...ErrorCode) error {
	// Grab the next error, if there is one.
	err := fetchError()
	if err == nil {
		return nil
	}

	// Only if the error has the specific error code accepted by the caller, do
	// we return the error.
	for _, code := range codes {
		if err.Code == code {
			return err
		}
	}

	// The error isn't accepted by the caller. If the error code is not a code
	// defined in the GLFW C documentation as a programmer error, then the
	// caller should have accepted it. This is effectively a bug in this
	// package.
	switch err.Code {
	case platformError:
		log.Println(err)
		return nil
	case notInitialized, noCurrentContext, invalidEnum, invalidValue, outOfMemory:
		panic(err)
	default:
		fmt.Println("GLFW: An invalid error was not accepted by the caller:", err)
		fmt.Println("GLFW: Please report this bug in the Go package immediately.")
		panic(err)
	}
}

// panicError is a helper used by functions which expect no errors (except
// programmer errors) to occur. It will panic if it finds any such error.
func panicError() {
	err := acceptError()
	if err != nil {
		panic(err)
	}
}

// fetchError fetches the next error from the error channel, it does not block
// and returns nil if there is no error present.
func fetchError() *Error {
	select {
	case err := <-lastError:
		return err
	default:
		return nil
	}
}
