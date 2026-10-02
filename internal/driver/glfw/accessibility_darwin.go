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
	"fyne.io/fyne/v2/internal/scale"
)

type accessibilityNode struct {
	object              fyne.CanvasObject
	role                C.AccessibilityRole
	label               string
	x, y, width, height int
	parent              int
	element             C.AccessibilityElementRef
}

var accessibilityElements = make(map[*window][]accessibilityNode)

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

	var nodes []accessibilityNode
	if w.canvas.Content() != nil {
		nodes = w.collectAccessibilityElements(w.canvas.Content(), fyne.NewPos(0, 0), -1, nodes)
	}

	if w.canvas.menu != nil {
		nodes = w.collectAccessibilityElements(w.canvas.menu, fyne.NewPos(0, 0), -1, nodes)
	}

	for _, overlay := range w.canvas.Overlays().List() {
		nodes = w.collectAccessibilityElements(overlay, fyne.NewPos(0, 0), -1, nodes)
	}

	C.AccessibilitySetTargetWindow(w.view().GetCocoaWindow())
	accessibilityElements[w] = updateAccessibilityElements(accessibilityElements[w], nodes)
}

func (w *window) collectAccessibilityElements(
	obj fyne.CanvasObject,
	pos fyne.Position,
	parent int,
	nodes []accessibilityNode,
) []accessibilityNode {
	if obj == nil || !obj.Visible() {
		return nodes
	}

	objPos := pos.Add(obj.Position())
	accessible, ok := obj.(fyne.Accessible)
	if !ok {
		return nodes
	}

	index := len(nodes)
	nodes = append(nodes, accessibilityNode{
		object: obj,
		role:   roleToC(string(accessible.AccessibilityRole())),
		label:  accessible.AccessibilityLabel(),
		x:      scale.ToScreenCoordinate(w.canvas, objPos.X),
		y:      scale.ToScreenCoordinate(w.canvas, objPos.Y),
		width:  scale.ToScreenCoordinate(w.canvas, obj.Size().Width),
		height: scale.ToScreenCoordinate(w.canvas, obj.Size().Height),
		parent: parent,
	})
	if cont, ok := obj.(*fyne.Container); ok {
		for _, child := range cont.Objects {
			nodes = w.collectAccessibilityElements(child, objPos, index, nodes)
		}
	}

	return nodes
}

func updateAccessibilityElements(previous, next []accessibilityNode) []accessibilityNode {
	rebuild := len(previous) != len(next)
	if !rebuild {
		for i := range next {
			if previous[i].object != next[i].object || previous[i].role != next[i].role ||
				previous[i].parent != next[i].parent {
				rebuild = true
				break
			}
		}
	}
	if rebuild {
		destroyAccessibilityElements(previous)
		for i := range next {
			node := &next[i]
			label := C.CString(node.label)
			node.element = C.AccessibilityElementCreate(node.role, label, label,
				C.double(node.x), C.double(node.y), C.double(node.width), C.double(node.height),
				nil, nil, nil)
			C.free(unsafe.Pointer(label))
			if node.parent >= 0 {
				C.AccessibilityElementAddChild(next[node.parent].element, node.element)
			}
		}
		for _, node := range next {
			if node.parent < 0 {
				C.AccessibilityAttachToWindow(node.element)
			}
		}
		return next
	}
	for i := range next {
		node, old := &next[i], previous[i]
		node.element = old.element
		if node.label != old.label {
			label := C.CString(node.label)
			C.AccessibilityElementSetTitle(node.element, label)
			C.AccessibilityElementSetLabel(node.element, label)
			C.free(unsafe.Pointer(label))
		}
		if node.x != old.x || node.y != old.y || node.width != old.width || node.height != old.height {
			C.AccessibilityElementSetFrame(node.element, C.double(node.x), C.double(node.y),
				C.double(node.width), C.double(node.height))
		}
	}
	return next
}

func destroyAccessibilityElements(nodes []accessibilityNode) {
	for _, node := range nodes {
		if node.parent < 0 {
			C.AccessibilityElementDestroy(node.element)
		}
	}
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
		destroyAccessibilityElements(elements)
		delete(accessibilityElements, w)
	}
}
