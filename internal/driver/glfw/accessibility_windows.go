//go:build accessibility && windows

package glfw

/*
#cgo LDFLAGS: -lole32 -loleaut32

#include <stdlib.h>
#include "accessibility_windows.h"
*/
import "C"

import (
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/scale"
)

func (w *window) updateAccessibility() {
	if w.view() == nil {
		return
	}

	hwnd := w.view().GetWin32Window()
	C.WinAccessibilitySetWindow(unsafe.Pointer(hwnd))
	C.WinAccessibilityClearElements()

	if w.canvas.Content() != nil {
		w.collectAccessibleElements(w.canvas.Content(), fyne.NewPos(0, 0))
	}

	if w.canvas.menu != nil {
		w.collectAccessibleElements(w.canvas.menu, fyne.NewPos(0, 0))
	}

	for _, overlay := range w.canvas.Overlays().List() {
		w.collectAccessibleElements(overlay, fyne.NewPos(0, 0))
	}

	C.WinAccessibilityUpdate()
}

func (w *window) collectAccessibleElements(obj fyne.CanvasObject, pos fyne.Position) {
	if obj == nil || !obj.Visible() {
		return
	}

	objPos := pos.Add(obj.Position())

	if accessible, isAccessible := obj.(fyne.Accessible); isAccessible {
		role := accessible.AccessibilityRole()
		// Use flat model: skip containers, only add leaf elements
		if role != fyne.AccessibleRoleContainer {
			label := accessible.AccessibilityLabel()

			pixelX := scale.ToScreenCoordinate(w.canvas, objPos.X)
			pixelY := scale.ToScreenCoordinate(w.canvas, objPos.Y)
			pixelW := scale.ToScreenCoordinate(w.canvas, obj.Size().Width)
			pixelH := scale.ToScreenCoordinate(w.canvas, obj.Size().Height)

			cLabel := C.CString(label)
			C.WinAccessibilityAddElement(cLabel, roleToCWin(role),
				C.double(pixelX), C.double(pixelY),
				C.double(pixelW), C.double(pixelH))
			C.free(unsafe.Pointer(cLabel))
		}
	}

	for _, child := range common.AccessibilityChildren(obj) {
		w.collectAccessibleElements(child, objPos)
	}
}

func roleToCWin(role fyne.AccessibleRole) C.WinAccessibilityRole {
	switch role {
	case fyne.AccessibleRoleButton:
		return C.WinAccessibilityRoleButton
	case fyne.AccessibleRoleCheckbox:
		return C.WinAccessibilityRoleCheckbox
	case fyne.AccessibleRoleHeading:
		return C.WinAccessibilityRoleHeading
	case fyne.AccessibleRoleImage:
		return C.WinAccessibilityRoleImage
	case fyne.AccessibleRoleLink:
		return C.WinAccessibilityRoleLink
	case fyne.AccessibleRoleList:
		return C.WinAccessibilityRoleList
	case fyne.AccessibleRoleListItem:
		return C.WinAccessibilityRoleListItem
	case fyne.AccessibleRoleProgressBar:
		return C.WinAccessibilityRoleProgressBar
	case fyne.AccessibleRoleRadio:
		return C.WinAccessibilityRoleRadio
	case fyne.AccessibleRoleSeparator:
		return C.WinAccessibilityRoleSeparator
	case fyne.AccessibleRoleSlider:
		return C.WinAccessibilityRoleSlider
	case fyne.AccessibleRoleTab:
		return C.WinAccessibilityRoleTab
	case fyne.AccessibleRoleTabList:
		return C.WinAccessibilityRoleTabList
	case fyne.AccessibleRoleTable:
		return C.WinAccessibilityRoleTable
	case fyne.AccessibleRoleText:
		return C.WinAccessibilityRoleText
	case fyne.AccessibleRoleTextField:
		return C.WinAccessibilityRoleTextField
	case fyne.AccessibleRoleTree:
		return C.WinAccessibilityRoleTree
	case fyne.AccessibleRoleTreeItem:
		return C.WinAccessibilityRoleTreeItem
	default:
		return C.WinAccessibilityRoleGroup
	}
}

func (w *window) initAccessibilityForWindow() {
	// Initialization is handled lazily in updateAccessibility
}

func (w *window) cleanupAccessibilityForWindow() {
	C.WinAccessibilityCleanup()
}
