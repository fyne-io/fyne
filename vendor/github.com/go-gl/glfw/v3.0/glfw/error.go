package glfw

//#include <GLFW/glfw3.h>
//void glfwSetErrorCallbackCB();
import "C"

//ErrorCode corresponds to an error code.
type ErrorCode int

//Error codes.
const (
	NotInitialized     ErrorCode = C.GLFW_NOT_INITIALIZED     //GLFW has not been initialized.
	NoCurrentContext   ErrorCode = C.GLFW_NO_CURRENT_CONTEXT  //No context is current.
	InvalidEnum        ErrorCode = C.GLFW_INVALID_ENUM        //One of the enum parameters for the function was given an invalid enum.
	InvalidValue       ErrorCode = C.GLFW_INVALID_VALUE       //One of the parameters for the function was given an invalid value.
	OutOfMemory        ErrorCode = C.GLFW_OUT_OF_MEMORY       //A memory allocation failed.
	ApiUnavailable     ErrorCode = C.GLFW_API_UNAVAILABLE     //GLFW could not find support for the requested client API on the system.
	VersionUnavailable ErrorCode = C.GLFW_VERSION_UNAVAILABLE //The requested client API version is not available.
	PlatformError      ErrorCode = C.GLFW_PLATFORM_ERROR      //A platform-specific error occurred that does not match any of the more specific categories.
	FormatUnavailable  ErrorCode = C.GLFW_FORMAT_UNAVAILABLE  //The clipboard did not contain data in the requested format.
)

var fErrorHolder func(code ErrorCode, desc string)

//export goErrorCB
func goErrorCB(code C.int, desc *C.char) {
	fErrorHolder(ErrorCode(code), C.GoString(desc))
}

//SetErrorCallback sets the error callback, which is called with an error code
//and a human-readable description each time a GLFW error occurs.
//
//This function may be called before Init.
func SetErrorCallback(cbfun func(code ErrorCode, desc string)) {
	if cbfun == nil {
		C.glfwSetErrorCallback(nil)
	} else {
		fErrorHolder = cbfun
		C.glfwSetErrorCallbackCB()
	}
}
