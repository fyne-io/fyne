//go:build accessibility && darwin

package glfw

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit -framework ApplicationServices

#include <stdbool.h>
#include <stdlib.h>
#include "accessibility_darwin.h"
*/
import "C"

import (
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/scale"
)

var accessibilityElements = make(map[*window][]C.AccessibilityElementRef)

func roleToC(role string) C.AccessibilityRole {
	switch role {
	case "button":
		return C.AccessibilityRoleButton
	case "text":
		return C.AccessibilityRoleStaticText
	case "link":
		return C.AccessibilityRoleButton
	default:
		return C.AccessibilityRoleGroup
	}
}

func (w *window) updateAccessibility() {
	if w.view() == nil {
		return
	}

	// Clean up old elements (recursive cleanup happens in C bridge)
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

	// Only attach root elements to window
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
		role := string(accessible.AccessibilityRole())

		pixelX := scale.ToScreenCoordinate(w.canvas, objPos.X)
		pixelY := scale.ToScreenCoordinate(w.canvas, objPos.Y)
		pixelWidth := scale.ToScreenCoordinate(w.canvas, obj.Size().Width)
		pixelHeight := scale.ToScreenCoordinate(w.canvas, obj.Size().Height)

		cLabel := C.CString(label)
		cTitle := C.CString(label)
		defer C.free(unsafe.Pointer(cLabel))
		defer C.free(unsafe.Pointer(cTitle))

		currentElement = C.AccessibilityElementCreate(
			roleToC(role),
			cTitle, cLabel,
			C.double(pixelX), C.double(pixelY),
			C.double(pixelWidth), C.double(pixelHeight),
			nil, nil, nil,
		)

		if parent != nil {
			C.AccessibilityElementAddChild(parent, currentElement)
		} else {
			result = append(result, currentElement)
		}
	}

	// Recurse via AccessibleChildren if available, otherwise the standard
	// Container.Objects path. Non-Accessible objects still recurse so we
	// don't lose subtrees rooted under decorative wrappers.
	for _, child := range common.AccessibilityChildren(obj) {
		childResults := w.collectAccessibilityElements(child, objPos, currentElement, depth+1)
		if parent == nil && currentElement == parent {
			// We didn't create an element here; child roots bubble up.
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
