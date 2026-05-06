//go:build accessibility && ios

package mobile

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit

#include <stdlib.h>

void clearAccessibilityNodesIOS(void);
void addAccessibilityNodeIOS(int role, const char *label,
	float x, float y, float width, float height);
void commitAccessibilityNodesIOS(void);
void setupAccessibilityIOS(void);
*/
import "C"

import (
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/scale"
)

const (
	iosRoleButton    = 1
	iosRoleText      = 2
	iosRoleLink      = 3
	iosRoleContainer = 4
)

func roleNameToIOS(role string) int {
	switch role {
	case "button":
		return iosRoleButton
	case "text":
		return iosRoleText
	case "link":
		return iosRoleLink
	default:
		return iosRoleContainer
	}
}

func (w *window) updateAccessibility() {
	nodes := w.gatherIOSAccessNodes()

	C.clearAccessibilityNodesIOS()
	for _, n := range nodes {
		cLabel := C.CString(n.label)
		C.addAccessibilityNodeIOS(
			C.int(n.role),
			cLabel,
			C.float(n.x), C.float(n.y), C.float(n.width), C.float(n.height),
		)
		C.free(unsafe.Pointer(cLabel))
	}
	C.commitAccessibilityNodesIOS()
}

type iosAccessNode struct {
	role          int
	label         string
	x, y          float32
	width, height float32
}

func (w *window) gatherIOSAccessNodes() []iosAccessNode {
	var nodes []iosAccessNode

	if w.canvas.Content() != nil {
		w.collectIOSNodes(w.canvas.Content(), fyne.NewPos(0, 0), &nodes)
	}
	if w.canvas.menu != nil {
		w.collectIOSNodes(w.canvas.menu, fyne.NewPos(0, 0), &nodes)
	}
	for _, overlay := range w.canvas.Overlays().List() {
		w.collectIOSNodes(overlay, fyne.NewPos(0, 0), &nodes)
	}

	return nodes
}

func (w *window) collectIOSNodes(
	obj fyne.CanvasObject,
	pos fyne.Position,
	nodes *[]iosAccessNode,
) {
	if obj == nil || !obj.Visible() {
		return
	}

	objPos := pos.Add(obj.Position())

	// Only add leaf roles (text, button, link) as accessibility elements.
	// Containers are skipped so VoiceOver can navigate their children directly.
	if accessible, ok := obj.(fyne.Accessible); ok {
		role := roleNameToIOS(string(accessible.AccessibilityRole()))
		if role != iosRoleContainer {
			*nodes = append(*nodes, iosAccessNode{
				role:   role,
				label:  accessible.AccessibilityLabel(),
				x:      float32(scale.ToScreenCoordinate(w.canvas, objPos.X)),
				y:      float32(scale.ToScreenCoordinate(w.canvas, objPos.Y)),
				width:  float32(scale.ToScreenCoordinate(w.canvas, obj.Size().Width)),
				height: float32(scale.ToScreenCoordinate(w.canvas, obj.Size().Height)),
			})
		}
	}

	for _, child := range common.AccessibilityChildren(obj) {
		w.collectIOSNodes(child, objPos, nodes)
	}
}

func (w *window) initAccessibilityForWindow() {
	C.setupAccessibilityIOS()
}

func (w *window) cleanupAccessibilityForWindow() {
	C.clearAccessibilityNodesIOS()
	C.commitAccessibilityNodesIOS()
}
