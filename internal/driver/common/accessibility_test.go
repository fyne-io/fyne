package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type accessibleChildrenStub struct {
	canvas.Rectangle
	children []fyne.CanvasObject
}
func (s *accessibleChildrenStub) AccessibilityLabel() string { return "stub" }
func (s *accessibleChildrenStub) AccessibilityRole() fyne.AccessibleRole {
	return fyne.AccessibleRoleContainer
}
func (s *accessibleChildrenStub) AccessibilityChildren() []fyne.CanvasObject {
	return s.children
}

func TestAccessibilityChildren_PrefersAccessibleChildren(t *testing.T) {
	a := canvas.NewRectangle(nil)
	b := canvas.NewRectangle(nil)
	stub := &accessibleChildrenStub{children: []fyne.CanvasObject{a, b}}

	got := AccessibilityChildren(stub)
	assert.Equal(t, []fyne.CanvasObject{a, b}, got)
}

func TestAccessibilityChildren_FallsBackToContainer(t *testing.T) {
	a := canvas.NewRectangle(nil)
	b := canvas.NewRectangle(nil)
	cont := container.NewWithoutLayout(a, b)

	got := AccessibilityChildren(cont)
	assert.Equal(t, []fyne.CanvasObject{a, b}, got)
}

func TestAccessibilityChildren_NoAccessibleDescendants(t *testing.T) {
	rect := canvas.NewRectangle(nil)
	assert.Nil(t, AccessibilityChildren(rect))
}

func TestAccessibilityChildren_AccessibleChildrenWinsOverContainer(t *testing.T) {
	// A *fyne.Container that also implements AccessibleChildren must
	// surface AccessibilityChildren first; falling back to Objects would
	// otherwise double-traverse or skip the customised view.
	a := canvas.NewRectangle(nil)
	b := canvas.NewRectangle(nil)
	stub := &accessibleChildrenStub{children: []fyne.CanvasObject{a}}
	// Bury the stub in a Container; AccessibilityChildren should descend
	// into the container first.
	cont := container.NewWithoutLayout(stub, b)

	got := AccessibilityChildren(cont)
	assert.Equal(t, []fyne.CanvasObject{stub, b}, got)

	// And direct invocation on the stub should return its custom children.
	assert.Equal(t, []fyne.CanvasObject{a}, AccessibilityChildren(stub))
}
