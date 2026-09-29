package common

import "fyne.io/fyne/v2"

// AccessibilityChildren returns the canvas objects that an accessibility
// walker should descend into for obj. It prefers
// [fyne.AccessibleChildren] when the widget implements it, falling back
// to the children of a [*fyne.Container]. The returned slice may be nil
// when obj has no accessible descendants.
//
// This helper is shared by the platform-specific accessibility walkers
// in [internal/driver/glfw] and [internal/driver/mobile] so that
// recursion is consistent across desktop and mobile drivers.
//
// Since: 2.8
func AccessibilityChildren(obj fyne.CanvasObject) []fyne.CanvasObject {
	if ac, ok := obj.(fyne.AccessibleChildren); ok {
		return ac.AccessibilityChildren()
	}
	if cont, ok := obj.(*fyne.Container); ok {
		return cont.Objects
	}
	return nil
}
