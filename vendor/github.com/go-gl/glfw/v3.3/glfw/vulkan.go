package glfw

//#define GLFW_INCLUDE_NONE
//#include "glfw/include/GLFW/glfw3.h"
import "C"

// VulkanSupported reports whether the Vulkan loader has been found. This check is performed by Init.
//
// The availability of a Vulkan loader does not by itself guarantee that window surface creation or
// even device creation is possible. Call GetRequiredInstanceExtensions to check whether the
// extensions necessary for Vulkan surface creation are available and GetPhysicalDevicePresentationSupport
// to check whether a queue family of a physical device supports image presentation.
func VulkanSupported() bool {
	return glfwbool(C.glfwVulkanSupported())
}
