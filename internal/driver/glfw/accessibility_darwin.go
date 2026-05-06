//go:build accessibility && darwin

package glfw

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit -framework ApplicationServices

#include <stdbool.h>
#include <stdlib.h>
#include "accessibility_darwin.h"

extern int fyneAccessibilityPerformAction(void* ctx, int actionCode, char* setValueArg);
extern void fyneAccessibilityDestroyContext(void* ctx);

// Static trampolines: cgo cannot take the address of an exported Go function
// directly through `C.f` syntax, so wrap each export in a C function and
// expose the function pointer via accessor helpers below.
static int actionTrampoline(void* ctx, int actionCode, const char* setValueArg) {
    return fyneAccessibilityPerformAction(ctx, actionCode, (char*)setValueArg);
}

static void destroyTrampoline(void* ctx) {
    fyneAccessibilityDestroyContext(ctx);
}

static AccessibilityActionCallback getActionCallback(void) {
    return actionTrampoline;
}

static AccessibilityContextDestroy getDestroyCallback(void) {
    return destroyTrampoline;
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/scale"
)

var accessibilityElements = make(map[*window][]C.AccessibilityElementRef)

func roleToC(role fyne.AccessibleRole) C.AccessibilityRole {
	switch role {
	case fyne.AccessibleRoleButton:
		return C.AccessibilityRoleButton
	case fyne.AccessibleRoleCheckbox:
		return C.AccessibilityRoleCheckbox
	case fyne.AccessibleRoleContainer:
		return C.AccessibilityRoleContainer
	case fyne.AccessibleRoleHeading:
		return C.AccessibilityRoleHeading
	case fyne.AccessibleRoleImage:
		return C.AccessibilityRoleImage
	case fyne.AccessibleRoleLink:
		return C.AccessibilityRoleLink
	case fyne.AccessibleRoleList:
		return C.AccessibilityRoleList
	case fyne.AccessibleRoleListItem:
		return C.AccessibilityRoleListItem
	case fyne.AccessibleRoleProgressBar:
		return C.AccessibilityRoleProgressBar
	case fyne.AccessibleRoleRadio:
		return C.AccessibilityRoleRadio
	case fyne.AccessibleRoleSeparator:
		return C.AccessibilityRoleSeparator
	case fyne.AccessibleRoleSlider:
		return C.AccessibilityRoleSlider
	case fyne.AccessibleRoleTab:
		return C.AccessibilityRoleTab
	case fyne.AccessibleRoleTabList:
		return C.AccessibilityRoleTabList
	case fyne.AccessibleRoleTable:
		return C.AccessibilityRoleTable
	case fyne.AccessibleRoleText:
		return C.AccessibilityRoleStaticText
	case fyne.AccessibleRoleTextField:
		return C.AccessibilityRoleTextField
	case fyne.AccessibleRoleTree:
		return C.AccessibilityRoleTree
	case fyne.AccessibleRoleTreeItem:
		return C.AccessibilityRoleTreeItem
	default:
		return C.AccessibilityRoleGroup
	}
}

// stateMaskFor builds the C state bitmask from a widget's reported states.
func stateMaskFor(obj fyne.CanvasObject) int {
	provider, ok := obj.(fyne.AccessibleStates)
	if !ok {
		return 0
	}
	var mask int
	for _, s := range provider.AccessibilityStates() {
		switch s {
		case fyne.AccessibleStateChecked:
			mask |= int(C.AccessibilityStateMaskChecked)
		case fyne.AccessibleStateDisabled:
			mask |= int(C.AccessibilityStateMaskDisabled)
		case fyne.AccessibleStateExpanded:
			mask |= int(C.AccessibilityStateMaskExpanded)
		case fyne.AccessibleStateFocused:
			mask |= int(C.AccessibilityStateMaskFocused)
		case fyne.AccessibleStateInvalid:
			mask |= int(C.AccessibilityStateMaskInvalid)
		case fyne.AccessibleStateRequired:
			mask |= int(C.AccessibilityStateMaskRequired)
		case fyne.AccessibleStateSelected:
			mask |= int(C.AccessibilityStateMaskSelected)
		}
	}
	return mask
}

// actionMaskFor builds the C action bitmask from a widget's reported actions.
func actionMaskFor(obj fyne.CanvasObject) (int, bool) {
	provider, ok := obj.(fyne.AccessibleActions)
	if !ok {
		return 0, false
	}
	var mask int
	for _, a := range provider.AccessibilityActions() {
		switch a {
		case fyne.AccessibleActionPress:
			mask |= int(C.AccessibilityActionMaskPress)
		case fyne.AccessibleActionIncrement:
			mask |= int(C.AccessibilityActionMaskIncrement)
		case fyne.AccessibleActionDecrement:
			mask |= int(C.AccessibilityActionMaskDecrement)
		case fyne.AccessibleActionShowMenu:
			mask |= int(C.AccessibilityActionMaskShowMenu)
		case fyne.AccessibleActionSelect:
			mask |= int(C.AccessibilityActionMaskSelect)
		case fyne.AccessibleActionSetValue:
			mask |= int(C.AccessibilityActionMaskSetValue)
		}
	}
	if _, hasSetter := obj.(fyne.AccessibleValueSetter); hasSetter {
		mask |= int(C.AccessibilityActionMaskSetValue)
	}
	return mask, mask != 0
}

func actionFromC(code int) (fyne.AccessibleAction, bool) {
	switch C.int(code) {
	case C.AccessibilityActionPress:
		return fyne.AccessibleActionPress, true
	case C.AccessibilityActionIncrement:
		return fyne.AccessibleActionIncrement, true
	case C.AccessibilityActionDecrement:
		return fyne.AccessibleActionDecrement, true
	case C.AccessibilityActionShowMenu:
		return fyne.AccessibleActionShowMenu, true
	case C.AccessibilityActionSelect:
		return fyne.AccessibleActionSelect, true
	case C.AccessibilityActionSetValue:
		return fyne.AccessibleActionSetValue, true
	}
	return "", false
}

//export fyneAccessibilityPerformAction
func fyneAccessibilityPerformAction(ctx unsafe.Pointer, actionCode C.int, setValueArg *C.char) C.int {
	if ctx == nil {
		return 0
	}
	defer func() {
		// Swallow panics from misbehaving widgets so we never crash AppKit.
		_ = recover()
	}()

	handle := cgo.Handle(uintptr(ctx))
	value, ok := handle.Value().(fyne.CanvasObject)
	if !ok || value == nil {
		return 0
	}

	act, known := actionFromC(int(actionCode))
	if !known {
		return 0
	}

	if act == fyne.AccessibleActionSetValue {
		setter, ok := value.(fyne.AccessibleValueSetter)
		if !ok {
			return 0
		}
		var arg string
		if setValueArg != nil {
			arg = C.GoString(setValueArg)
		}
		if setter.AccessibilitySetValue(arg) {
			return 1
		}
		return 0
	}

	performer, ok := value.(fyne.AccessibleActions)
	if !ok {
		return 0
	}
	if performer.AccessibilityPerformAction(act) {
		return 1
	}
	return 0
}

//export fyneAccessibilityDestroyContext
func fyneAccessibilityDestroyContext(ctx unsafe.Pointer) {
	if ctx == nil {
		return
	}
	cgo.Handle(uintptr(ctx)).Delete()
}

func (w *window) updateAccessibility() {
	if w.view() == nil {
		return
	}

	if oldElements, ok := accessibilityElements[w]; ok {
		for _, elem := range oldElements {
			C.AccessibilityElementDestroy(elem)
		}
		delete(accessibilityElements, w)
	}

	var rootElements []C.AccessibilityElementRef
	if w.canvas.Content() != nil {
		contentRoots := w.collectAccessibilityElements(w.canvas.Content(), fyne.NewPos(0, 0), nil, 0)
		rootElements = append(rootElements, contentRoots...)
	}

	if w.canvas.menu != nil {
		menuRoots := w.collectAccessibilityElements(w.canvas.menu, fyne.NewPos(0, 0), nil, 0)
		rootElements = append(rootElements, menuRoots...)
	}

	for _, overlay := range w.canvas.Overlays().List() {
		overlayRoots := w.collectAccessibilityElements(overlay, fyne.NewPos(0, 0), nil, 0)
		rootElements = append(rootElements, overlayRoots...)
	}

	if w.view() != nil {
		nsWindow := w.view().GetCocoaWindow()
		C.AccessibilitySetTargetWindow(nsWindow)
	}

	for _, rootElem := range rootElements {
		C.AccessibilityAttachToWindow(rootElem)
	}
	accessibilityElements[w] = rootElements
}

func (w *window) collectAccessibilityElements(
	obj fyne.CanvasObject,
	pos fyne.Position,
	parent C.AccessibilityElementRef,
	depth int,
) []C.AccessibilityElementRef {
	if obj == nil || !obj.Visible() {
		return nil
	}

	objPos := pos.Add(obj.Position())
	var result []C.AccessibilityElementRef
	currentElement := parent

	if accessible, ok := obj.(fyne.Accessible); ok {
		label := accessible.AccessibilityLabel()
		role := accessible.AccessibilityRole()

		pixelX := scale.ToScreenCoordinate(w.canvas, objPos.X)
		pixelY := scale.ToScreenCoordinate(w.canvas, objPos.Y)
		pixelWidth := scale.ToScreenCoordinate(w.canvas, obj.Size().Width)
		pixelHeight := scale.ToScreenCoordinate(w.canvas, obj.Size().Height)

		cLabel := C.CString(label)
		cTitle := C.CString(label)
		defer C.free(unsafe.Pointer(cLabel))
		defer C.free(unsafe.Pointer(cTitle))

		// cgo.Handle the live widget so the Obj-C side can route actions
		// back. The handle is freed by destroyTrampoline when -dealloc fires.
		handle := cgo.NewHandle(obj)
		ctx := unsafe.Pointer(uintptr(handle))

		currentElement = C.AccessibilityElementCreate(
			roleToC(role),
			cTitle, cLabel,
			C.double(pixelX), C.double(pixelY),
			C.double(pixelWidth), C.double(pixelHeight),
			nil,
			C.getActionCallback(),
			ctx,
			C.getDestroyCallback(),
		)

		if valued, ok := obj.(fyne.AccessibleValue); ok {
			cValue := C.CString(valued.AccessibilityValue())
			C.AccessibilityElementSetValue(currentElement, cValue)
			C.free(unsafe.Pointer(cValue))
		}

		C.AccessibilityElementSetStates(currentElement, C.int(stateMaskFor(obj)))
		if mask, ok := actionMaskFor(obj); ok {
			C.AccessibilityElementSetSupportedActions(currentElement, C.int(mask))
		}

		if parent != nil {
			C.AccessibilityElementAddChild(parent, currentElement)
		} else {
			result = append(result, currentElement)
		}
	}

	for _, child := range common.AccessibilityChildren(obj) {
		childResults := w.collectAccessibilityElements(child, objPos, currentElement, depth+1)
		if parent == nil && currentElement == parent {
			result = append(result, childResults...)
		}
	}

	return result
}

func (w *window) initAccessibilityForWindow() {
	if w.view() == nil {
		return
	}
}

func (w *window) cleanupAccessibilityForWindow() {
	if w.view() == nil {
		return
	}

	if elements, ok := accessibilityElements[w]; ok {
		for _, elem := range elements {
			C.AccessibilityElementDestroy(elem)
		}
		delete(accessibilityElements, w)
	}
}
