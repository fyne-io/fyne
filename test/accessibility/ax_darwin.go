//go:build darwin && cgo && accessibility && axintegration

package accessibility

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void* axApplication(int pid) { return (void*)AXUIElementCreateApplication(pid); }
static void axRetain(void* ref) { CFRetain(ref); }
static void axRelease(void* ref) { CFRelease(ref); }
static int axEqual(void* a, void* b) { return CFEqual(a, b); }
static int axTimeout(void* ref) { return AXUIElementSetMessagingTimeout(ref, 1.0); }
static int axAttribute(void* ref, const char* name, void** result) {
	CFStringRef key = CFStringCreateWithCString(NULL, name, kCFStringEncodingUTF8);
	CFTypeRef value = NULL;
	AXError err = AXUIElementCopyAttributeValue(ref, key, &value);
	CFRelease(key);
	*result = (void*)value;
	return err;
}
static int axAttributes(void* ref, void** result) {
	const void* keys[] = {
		kAXRoleAttribute, kAXSubroleAttribute, kAXTitleAttribute, kAXDescriptionAttribute,
		kAXValueAttribute, kAXEnabledAttribute, kAXSelectedAttribute, kAXExpandedAttribute,
		kAXSizeAttribute, kAXPositionAttribute, kAXRoleDescriptionAttribute
	};
	CFArrayRef names = CFArrayCreate(NULL, keys, 11, &kCFTypeArrayCallBacks);
	CFArrayRef values = NULL;
	AXError err = AXUIElementCopyMultipleAttributeValues(ref, names, 0, &values);
	CFRelease(names);
	*result = (void*)values;
	return err;
}
static int axAttributeError(void* ref) {
	if (CFGetTypeID(ref) == AXValueGetTypeID() && AXValueGetType(ref) == kAXValueAXErrorType) {
		AXError err;
		if (AXValueGetValue(ref, kAXValueAXErrorType, &err)) return err;
		return kAXErrorFailure;
	}
	return kAXErrorSuccess;
}
static char* axScalar(void* ref) {
	CFTypeID type = CFGetTypeID(ref);
	if (type == CFBooleanGetTypeID()) {
		return strdup(CFBooleanGetValue(ref) ? "true" : "false");
	}
	if (type == CFNumberGetTypeID()) {
		double value;
		if (!CFNumberGetValue(ref, kCFNumberDoubleType, &value)) return NULL;
		char text[64];
		snprintf(text, sizeof(text), "%.17g", value);
		return strdup(text);
	}
	if (type != CFStringGetTypeID()) return NULL;
	CFIndex size = CFStringGetMaximumSizeForEncoding(CFStringGetLength(ref), kCFStringEncodingUTF8) + 1;
	char* text = malloc(size);
	if (!text) return NULL;
	if (!CFStringGetCString(ref, text, size, kCFStringEncodingUTF8)) {
		free(text);
		return NULL;
	}
	return text;
}
static long axArrayCount(void* ref) {
	if (CFGetTypeID(ref) != CFArrayGetTypeID()) return -1;
	return CFArrayGetCount(ref);
}
static void* axArrayItem(void* ref, long index) {
	return (void*)CFArrayGetValueAtIndex(ref, index);
}
static int axIsElement(void* ref) { return CFGetTypeID(ref) == AXUIElementGetTypeID(); }
static int axGeometry(void* ref, int size, double* x, double* y) {
	if (CFGetTypeID(ref) != AXValueGetTypeID()) return 0;
	if (size) {
		CGSize value;
		if (!AXValueGetValue(ref, kAXValueCGSizeType, &value)) return 0;
		*x = value.width;
		*y = value.height;
	} else {
		CGPoint value;
		if (!AXValueGetValue(ref, kAXValueCGPointType, &value)) return 0;
		*x = value.x;
		*y = value.y;
	}
	return 1;
}
static int axActions(void* ref, void** result) {
	CFArrayRef value = NULL;
	AXError err = AXUIElementCopyActionNames(ref, &value);
	*result = (void*)value;
	return err;
}
static int axPerform(void* ref, const char* action) {
	CFStringRef name = CFStringCreateWithCString(NULL, action, kCFStringEncodingUTF8);
	AXError err = AXUIElementPerformAction(ref, name);
	CFRelease(name);
	return err;
}
static int axSetValue(void* ref, const char* text) {
	CFStringRef value = CFStringCreateWithCString(NULL, text, kCFStringEncodingUTF8);
	AXError err = AXUIElementSetAttributeValue(ref, kAXValueAttribute, value);
	CFRelease(value);
	return err;
}
static int axValueSettable(void* ref, int* result) {
	Boolean settable = false;
	AXError err = AXUIElementIsAttributeSettable(ref, kAXValueAttribute, &settable);
	*result = settable;
	return err;
}
*/
import "C"

import (
	"fmt"
	"strings"
	"time"
	"unsafe"
)

type axError int

func (e axError) Error() string {
	return fmt.Sprintf("AX error %d", int(e))
}

func axResult(code C.int) error {
	if code == C.kAXErrorSuccess {
		return nil
	}
	return axError(code)
}

func trusted() bool {
	return C.AXIsProcessTrusted() != 0
}

type axElement struct {
	ref unsafe.Pointer
}

func application(pid int) (axElement, error) {
	element := axElement{ref: C.axApplication(C.int(pid))}
	if err := axResult(C.axTimeout(element.ref)); err != nil {
		element.close()
		return axElement{}, err
	}
	return element, nil
}

func (e axElement) close() {
	C.axRelease(e.ref)
}

func (e axElement) attribute(name string) (unsafe.Pointer, error) {
	key := C.CString(name)
	defer C.free(unsafe.Pointer(key))
	var value unsafe.Pointer
	code := C.axAttribute(e.ref, key, &value)
	if code == C.kAXErrorAttributeUnsupported || code == C.kAXErrorNoValue {
		return nil, nil // Optional attributes are absent on many standard AppKit elements.
	}
	return value, axResult(code)
}

func scalar(value unsafe.Pointer) (string, error) {
	text := C.axScalar(value)
	if text == nil {
		return "", fmt.Errorf("expected a string, number, or boolean accessibility attribute")
	}
	defer C.free(unsafe.Pointer(text))
	return C.GoString(text), nil
}

func (e axElement) text(name string) (string, error) {
	value, err := e.attribute(name)
	if err != nil || value == nil {
		return "", err
	}
	defer C.axRelease(value)
	return scalar(value)
}

func (e axElement) elements(name string) ([]axElement, error) {
	value, err := e.attribute(name)
	if err != nil || value == nil {
		return nil, err
	}
	defer C.axRelease(value)
	count := C.axArrayCount(value)
	if count < 0 {
		return nil, fmt.Errorf("%s is not an array", name)
	}
	var elements []axElement
	for i := C.long(0); i < count; i++ {
		ref := C.axArrayItem(value, i)
		if C.axIsElement(ref) == 0 {
			for _, element := range elements {
				element.close()
			}
			return nil, fmt.Errorf("%s contains a non-element", name)
		}
		C.axRetain(ref)
		elements = append(elements, axElement{ref: ref})
		if err := axResult(C.axTimeout(ref)); err != nil {
			for _, element := range elements {
				element.close()
			}
			return nil, err
		}
	}
	return elements, nil
}

func (e axElement) actions() ([]string, error) {
	var value unsafe.Pointer
	code := C.axActions(e.ref, &value)
	if code == C.kAXErrorNotImplemented || code == C.kAXErrorActionUnsupported {
		return nil, nil
	}
	if err := axResult(code); err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("AXCopyActionNames returned no array")
	}
	defer C.axRelease(value)
	count := C.axArrayCount(value)
	if count < 0 {
		return nil, fmt.Errorf("AX actions are not an array")
	}
	var actions []string
	for i := C.long(0); i < count; i++ {
		action, err := scalar(C.axArrayItem(value, i))
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func (e axElement) perform(action string) error {
	name := C.CString(action)
	defer C.free(unsafe.Pointer(name))
	return axResult(C.axPerform(e.ref, name))
}

func (e axElement) setValue(text string) error {
	var settable C.int
	if err := axResult(C.axValueSettable(e.ref, &settable)); err != nil {
		return err
	}
	if settable == 0 {
		return fmt.Errorf("AXValue is not settable")
	}
	value := C.CString(text)
	defer C.free(unsafe.Pointer(value))
	return axResult(C.axSetValue(e.ref, value))
}

func (e axElement) related(name string, other axElement) (bool, error) {
	value, err := e.attribute(name)
	if err != nil || value == nil {
		return false, err
	}
	defer C.axRelease(value)
	return C.axEqual(value, other.ref) != 0, nil
}

type axNode struct {
	element    axElement
	attributes map[string]string
	depth      int
	width      float64
	height     float64
	x          float64
	y          float64
	positioned bool
}

func (n *axNode) matches(role, label string) bool {
	return n.attributes["AXRole"] == role &&
		(n.attributes["AXTitle"] == label || n.attributes["AXDescription"] == label)
}

type axSnapshot struct {
	nodes []*axNode
}

func (s *axSnapshot) close() {
	for _, node := range s.nodes {
		node.element.close()
	}
}

func (s *axSnapshot) find(role, label string) *axNode {
	return s.findWithAttributes(role, label, nil)
}

func (s *axSnapshot) findWithAttributes(role, label string, attributes map[string]string) *axNode {
	for _, node := range s.nodes {
		if !node.matches(role, label) {
			continue
		}
		matches := true
		for name, value := range attributes {
			if node.attributes[name] != value {
				matches = false
				break
			}
		}
		if matches {
			return node
		}
	}
	return nil
}

func (s *axSnapshot) String() string {
	var out strings.Builder
	for _, node := range s.nodes {
		fmt.Fprintf(&out, "%s%v bounds=(%.0f,%.0f %.0fx%.0f)\n", strings.Repeat("  ", node.depth),
			node.attributes, node.x, node.y, node.width, node.height)
	}
	return out.String()
}

func (s *axSnapshot) walk(element axElement, depth int, deadline time.Time) error {
	if time.Now().After(deadline) {
		return fmt.Errorf("accessibility traversal deadline exceeded")
	}
	if depth > 32 || len(s.nodes) >= 1000 {
		return fmt.Errorf("accessibility tree exceeds traversal limit (possible cycle)")
	}
	C.axRetain(element.ref)
	node := &axNode{element: element, depth: depth, attributes: make(map[string]string)}
	s.nodes = append(s.nodes, node)
	if err := node.readAttributes(); err != nil {
		return err
	}
	children, err := element.elements("AXChildren")
	if err != nil {
		return err
	}
	defer func() {
		for _, child := range children {
			child.close()
		}
	}()
	for _, child := range children {
		if err := s.walk(child, depth+1, deadline); err != nil {
			return err
		}
	}
	return nil
}

func (n *axNode) readAttributes() error {
	// One IPC request per element avoids both slow walks and mixed attributes
	// from different repaints.
	var values unsafe.Pointer
	if err := axResult(C.axAttributes(n.element.ref, &values)); err != nil {
		return err
	}
	if values == nil {
		return fmt.Errorf("AXCopyMultipleAttributeValues returned no array")
	}
	defer C.axRelease(values)
	names := []string{"AXRole", "AXSubrole", "AXTitle", "AXDescription", "AXValue", "AXEnabled", "AXSelected", "AXExpanded", "AXSize", "AXPosition", "AXRoleDescription"}
	if int(C.axArrayCount(values)) != len(names) {
		return fmt.Errorf("unexpected AX attribute array length")
	}
	for i, name := range names {
		value := C.axArrayItem(values, C.long(i))
		code := C.axAttributeError(value)
		if code == C.kAXErrorAttributeUnsupported || code == C.kAXErrorNoValue {
			continue
		}
		if err := axResult(code); err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if name == "AXSize" || name == "AXPosition" {
			var x, y C.double
			var size C.int
			if name == "AXSize" {
				size = 1
			}
			if C.axGeometry(value, size, &x, &y) == 0 {
				return fmt.Errorf("%s has incorrect geometry type", name)
			}
			if size == 1 {
				n.width, n.height = float64(x), float64(y)
			} else {
				n.x, n.y, n.positioned = float64(x), float64(y), true
			}
		} else {
			text, err := scalar(value)
			if err != nil {
				return fmt.Errorf("read %s: %w", name, err)
			}
			n.attributes[name] = text
		}
	}
	if n.attributes["AXRole"] == "" {
		return fmt.Errorf("element has no AXRole")
	}
	return nil
}

func snapshot(app axElement, title string, deadline time.Time) (*axSnapshot, error) {
	windows, err := app.elements("AXWindows")
	if err != nil {
		return nil, fmt.Errorf("read AXWindows: %w", err)
	}
	defer func() {
		for _, window := range windows {
			window.close()
		}
	}()
	for _, window := range windows {
		name, err := window.text("AXTitle")
		if err != nil {
			return nil, err
		}
		if name != title {
			continue
		}
		tree := &axSnapshot{}
		if err := tree.walk(window, 0, deadline); err != nil {
			tree.close()
			return nil, err
		}
		return tree, nil
	}
	return nil, fmt.Errorf("window %q not found in AXWindows", title)
}
