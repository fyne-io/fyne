package glfw

/*
#include "glfw/src/internal.h"

GLFWAPI VkResult glfwCreateWindowSurface(VkInstance instance, GLFWwindow* window, const VkAllocationCallbacks* allocator, VkSurfaceKHR* surface);
GLFWAPI GLFWvkproc glfwGetInstanceProcAddress(VkInstance instance, const char* procname);

// Helper function for doing raw pointer arithmetic
static inline const char* getArrayIndex(const char** array, unsigned int index) {
	return array[index];
}

void* getVulkanProcAddr() {
	return glfwGetInstanceProcAddress;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// VulkanSupported reports whether the Vulkan loader has been found. This check is performed by Init.
//
// The availability of a Vulkan loader does not by itself guarantee that window surface creation or
// even device creation is possible. Call GetRequiredInstanceExtensions to check whether the
// extensions necessary for Vulkan surface creation are available and GetPhysicalDevicePresentationSupport
// to check whether a queue family of a physical device supports image presentation.
func VulkanSupported() bool {
	return glfwbool(C.glfwVulkanSupported())
}

// GetVulkanGetInstanceProcAddress returns the function pointer used to find Vulkan core or
// extension functions. The return value of this function can be passed to the Vulkan library.
//
// Note that this function does not work the same way as the glfwGetInstanceProcAddress.
func GetVulkanGetInstanceProcAddress() unsafe.Pointer {
	return C.getVulkanProcAddr()
}

// GetRequiredInstanceExtensions returns a slice of Vulkan instance extension names required
// by GLFW for creating Vulkan surfaces for GLFW windows. If successful, the list will always
// contain VK_KHR_surface, so if you don't require any additional extensions you can pass this list
// directly to the VkInstanceCreateInfo struct.
//
// If Vulkan is not available on the machine, this function returns nil. Call
// VulkanSupported to check whether Vulkan is available.
//
// If Vulkan is available but no set of extensions allowing window surface creation was found, this
// function returns nil. You may still use Vulkan for off-screen rendering and compute work.
func (window *Window) GetRequiredInstanceExtensions() []string {
	var count C.uint32_t
	strarr := C.glfwGetRequiredInstanceExtensions(&count)
	if count == 0 {
		return nil
	}

	extensions := make([]string, count)
	for i := uint(0); i < uint(count); i++ {
		extensions[i] = C.GoString(C.getArrayIndex(strarr, C.uint(i)))
	}
	return extensions
}

// CreateWindowSurface creates a Vulkan surface for this window.
func (window *Window) CreateWindowSurface(instance interface{}, allocCallbacks unsafe.Pointer) (surface uintptr, err error) {
	if instance == nil {
		return 0, errors.New("vulkan: instance is nil")
	}
	val := reflect.ValueOf(instance)
	if val.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("vulkan: instance is not a VkInstance (expected kind Ptr, got %s)", val.Kind())
	}
	var vulkanSurface C.VkSurfaceKHR
	ret := C.glfwCreateWindowSurface(
		(C.VkInstance)(unsafe.Pointer(reflect.ValueOf(instance).Pointer())), window.data,
		(*C.VkAllocationCallbacks)(allocCallbacks), (*C.VkSurfaceKHR)(unsafe.Pointer(&vulkanSurface)))
	if ret != C.VK_SUCCESS {
		return 0, fmt.Errorf("vulkan: error creating window surface: %d", ret)
	}
	return uintptr(unsafe.Pointer(&vulkanSurface)), nil
}
