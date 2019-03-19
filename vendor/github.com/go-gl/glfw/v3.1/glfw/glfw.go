package glfw

//#include "glfw/include/GLFW/glfw3.h"
import "C"

const (
	VersionMajor    = C.GLFW_VERSION_MAJOR    // This is incremented when the API is changed in non-compatible ways.
	VersionMinor    = C.GLFW_VERSION_MINOR    // This is incremented when features are added to the API but it remains backward-compatible.
	VersionRevision = C.GLFW_VERSION_REVISION // This is incremented when a bug fix release is made that does not contain any API changes.
)

// Init initializes the GLFW library. Before most GLFW functions can be used,
// GLFW must be initialized, and before a program terminates GLFW should be
// terminated in order to free any resources allocated during or after
// initialization.
//
// If this function fails, it calls Terminate before returning. If it succeeds,
// you should call Terminate before the program exits.
//
// Additional calls to this function after successful initialization but before
// termination will succeed but will do nothing.
//
// This function may take several seconds to complete on some systems, while on
// other systems it may take only a fraction of a second to complete.
//
// On Mac OS X, this function will change the current directory of the
// application to the Contents/Resources subdirectory of the application's
// bundle, if present.
//
// This function may only be called from the main thread.
func Init() error {
	C.glfwInit()
	return acceptError(APIUnavailable)
}

// Terminate destroys all remaining windows, frees any allocated resources and
// sets the library to an uninitialized state. Once this is called, you must
// again call Init successfully before you will be able to use most GLFW
// functions.
//
// If GLFW has been successfully initialized, this function should be called
// before the program exits. If initialization fails, there is no need to call
// this function, as it is called by Init before it returns failure.
//
// This function may only be called from the main thread.
func Terminate() {
	flushErrors()
	C.glfwTerminate()
}

// GetVersion retrieves the major, minor and revision numbers of the GLFW
// library. It is intended for when you are using GLFW as a shared library and
// want to ensure that you are using the minimum required version.
//
// This function may be called before Init.
func GetVersion() (major, minor, revision int) {
	var (
		maj C.int
		min C.int
		rev C.int
	)

	C.glfwGetVersion(&maj, &min, &rev)
	return int(maj), int(min), int(rev)
}

// GetVersionString returns a static string generated at compile-time according
// to which configuration macros were defined. This is intended for use when
// submitting bug reports, to allow developers to see which code paths are
// enabled in a binary.
//
// This function may be called before Init.
func GetVersionString() string {
	return C.GoString(C.glfwGetVersionString())
}
