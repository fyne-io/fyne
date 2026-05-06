//go:build accessibility && android

package mobile

/*
#include <stdlib.h>
#include <stdint.h>

void clearAccessibilityNodes(uintptr_t jni_env, uintptr_t ctx);
void addAccessibilityNode(uintptr_t jni_env, uintptr_t ctx, int id, int role,
	const char* label, int x, int y, int width, int height, int parent_id);
void commitAccessibilityNodes(uintptr_t jni_env, uintptr_t ctx);
void setupAccessibility(uintptr_t jni_env, uintptr_t ctx);
*/
import "C"

import (
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/internal/scale"
)

const (
	androidRoleButton    = 1
	androidRoleText      = 2
	androidRoleLink      = 3
	androidRoleContainer = 4
)

type androidAccessNode struct {
	id     int
	role   int
	label  string
	x, y   int
	width  int
	height int
}

func roleNameToAndroid(role string) int {
	switch role {
	case "button":
		return androidRoleButton
	case "text":
		return androidRoleText
	case "link":
		return androidRoleLink
	default:
		return androidRoleContainer
	}
}

func (w *window) updateAccessibility() {
	nodes := w.gatherAndroidAccessNodes()

	app.RunOnJVM(func(_, jniEnv, ctx uintptr) error { //nolint:errcheck
		jenv := C.uintptr_t(jniEnv)
		jctx := C.uintptr_t(ctx)

		C.clearAccessibilityNodes(jenv, jctx)
		for _, n := range nodes {
			cLabel := C.CString(n.label)
			C.addAccessibilityNode(
				jenv, jctx,
				C.int(n.id), C.int(n.role),
				cLabel,
				C.int(n.x), C.int(n.y), C.int(n.width), C.int(n.height),
				0, // parentID always 0 (flat model)
			)
			C.free(unsafe.Pointer(cLabel))
		}
		C.commitAccessibilityNodes(jenv, jctx)
		return nil
	})
}

func (w *window) gatherAndroidAccessNodes() []androidAccessNode {
	var nodes []androidAccessNode
	nextID := 1

	if w.canvas.Content() != nil {
		w.collectAndroidNodes(w.canvas.Content(), fyne.NewPos(0, 0), &nodes, &nextID)
	}
	if w.canvas.menu != nil {
		w.collectAndroidNodes(w.canvas.menu, fyne.NewPos(0, 0), &nodes, &nextID)
	}
	for _, overlay := range w.canvas.Overlays().List() {
		w.collectAndroidNodes(overlay, fyne.NewPos(0, 0), &nodes, &nextID)
	}

	return nodes
}

func (w *window) collectAndroidNodes(
	obj fyne.CanvasObject,
	pos fyne.Position,
	nodes *[]androidAccessNode,
	nextID *int,
) {
	if obj == nil || !obj.Visible() {
		return
	}

	objPos := pos.Add(obj.Position())

	// Only add leaf roles (text, button, link) as virtual accessibility nodes.
	// Android's AccessibilityNodeProvider uses a flat model: all virtual views
	// must be direct children of the host view. Container nodes are skipped so
	// their children surface at the top level where TalkBack can reach them.
	if accessible, ok := obj.(fyne.Accessible); ok {
		role := roleNameToAndroid(string(accessible.AccessibilityRole()))
		if role != androidRoleContainer {
			*nodes = append(*nodes, androidAccessNode{
				id:     *nextID,
				role:   role,
				label:  accessible.AccessibilityLabel(),
				x:      scale.ToScreenCoordinate(w.canvas, objPos.X),
				y:      scale.ToScreenCoordinate(w.canvas, objPos.Y),
				width:  scale.ToScreenCoordinate(w.canvas, obj.Size().Width),
				height: scale.ToScreenCoordinate(w.canvas, obj.Size().Height),
			})
			(*nextID)++
		}
	}

	for _, child := range common.AccessibilityChildren(obj) {
		w.collectAndroidNodes(child, objPos, nodes, nextID)
	}
}

func (w *window) initAccessibilityForWindow() {
	app.RunOnJVM(func(_, jniEnv, ctx uintptr) error { //nolint:errcheck
		C.setupAccessibility(C.uintptr_t(jniEnv), C.uintptr_t(ctx))
		return nil
	})
}

func (w *window) cleanupAccessibilityForWindow() {
	app.RunOnJVM(func(_, jniEnv, ctx uintptr) error { //nolint:errcheck
		C.clearAccessibilityNodes(C.uintptr_t(jniEnv), C.uintptr_t(ctx))
		return nil
	})
}
