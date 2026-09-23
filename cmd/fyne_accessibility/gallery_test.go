package main

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func galleryWindow(t *testing.T, scenario string) fyne.Window {
	t.Helper()
	test.NewTempApp(t)
	w := test.NewTempWindow(t, nil)
	content, err := makeGallery(scenario, w)
	require.NoError(t, err)
	w.SetContent(content)
	w.Resize(fyne.NewSize(900, 650))
	return w
}

func accessibleObjects(root fyne.CanvasObject) []fyne.CanvasObject {
	if root == nil || !root.Visible() {
		return nil
	}
	var objects []fyne.CanvasObject
	if _, ok := root.(fyne.Accessible); ok {
		objects = append(objects, root)
	}
	for _, child := range common.AccessibilityChildren(root) {
		objects = append(objects, accessibleObjects(child)...)
	}
	return objects
}

func findAccessible(t *testing.T, w fyne.Window, role fyne.AccessibleRole, label string) fyne.CanvasObject {
	t.Helper()
	for _, object := range accessibleObjects(w.Content()) {
		accessible := object.(fyne.Accessible)
		if accessible.AccessibilityRole() == role && accessible.AccessibilityLabel() == label {
			return object
		}
	}
	t.Fatalf("missing accessible %s %q", role, label)
	return nil
}

func TestGalleryScenarios(t *testing.T) {
	for _, name := range append(scenarioNames(), "all") {
		t.Run(name, func(t *testing.T) {
			w := galleryWindow(t, name)
			assert.NotEmpty(t, accessibleObjects(w.Content()))
			if name != "all" {
				findAccessible(t, w, fyne.AccessibleRoleText, strings.ToUpper(name[:1])+name[1:]+" ready")
			}
		})
	}
	assert.False(t, validScenario("unknown"))
	content, err := makeGallery("unknown", nil)
	assert.Error(t, err)
	assert.Nil(t, content)
}

func TestGalleryControls(t *testing.T) {
	w := galleryWindow(t, "controls")
	button := findAccessible(t, w, fyne.AccessibleRoleButton, "Count presses").(fyne.AccessibleActions)
	require.True(t, button.AccessibilityPerformAction(fyne.AccessibleActionPress))
	findAccessible(t, w, fyne.AccessibleRoleText, "Pressed 1 times")
	disabled := findAccessible(t, w, fyne.AccessibleRoleButton, "Disabled button").(fyne.AccessibleActions)
	assert.False(t, disabled.AccessibilityPerformAction(fyne.AccessibleActionPress))
	check := findAccessible(t, w, fyne.AccessibleRoleCheckbox, "Enable notifications").(fyne.AccessibleActions)
	require.True(t, check.AccessibilityPerformAction(fyne.AccessibleActionPress))
	findAccessible(t, w, fyne.AccessibleRoleText, "Notifications: true")
	for _, object := range accessibleObjects(w.Content()) {
		if progress, ok := object.(*widget.ProgressBarInfinite); ok {
			assert.False(t, progress.Running(), "the gallery must not animate while the native tree is read")
		}
	}
}

func TestGalleryText(t *testing.T) {
	w := galleryWindow(t, "text")
	entry := findAccessible(t, w, fyne.AccessibleRoleTextField, "Display name").(fyne.AccessibleValueSetter)
	require.True(t, entry.AccessibilitySetValue("Grace"))
	findAccessible(t, w, fyne.AccessibleRoleText, "Name: Grace")
	disabled := findAccessible(t, w, fyne.AccessibleRoleTextField, "Read only example").(fyne.AccessibleValueSetter)
	assert.False(t, disabled.AccessibilitySetValue("changed"))
}

func TestGalleryCollections(t *testing.T) {
	w := galleryWindow(t, "collections")
	findAccessible(t, w, fyne.AccessibleRoleListItem, "List row 01")
	findAccessible(t, w, fyne.AccessibleRoleTreeItem, "Branch")
	findAccessible(t, w, fyne.AccessibleRoleText, "Cell 1,1")
	findAccessible(t, w, fyne.AccessibleRoleListItem, "Grid item 01")
	scroll := findAccessible(t, w, fyne.AccessibleRoleButton, "Scroll list to end").(*widget.Button)
	test.Tap(scroll)
	findAccessible(t, w, fyne.AccessibleRoleListItem, "List row 30")
}

func TestGalleryDynamic(t *testing.T) {
	w := galleryWindow(t, "dynamic")
	toggle := findAccessible(t, w, fyne.AccessibleRoleButton, "Toggle extra content").(fyne.AccessibleActions)
	require.True(t, toggle.AccessibilityPerformAction(fyne.AccessibleActionPress))
	findAccessible(t, w, fyne.AccessibleRoleText, "Extra content")
	require.True(t, toggle.AccessibilityPerformAction(fyne.AccessibleActionPress))
	for _, object := range accessibleObjects(w.Content()) {
		assert.NotEqual(t, "Extra content", object.(fyne.Accessible).AccessibilityLabel())
	}
	open := findAccessible(t, w, fyne.AccessibleRoleButton, "Open dialog").(fyne.AccessibleActions)
	require.True(t, open.AccessibilityPerformAction(fyne.AccessibleActionPress))
	require.NotEmpty(t, w.Canvas().Overlays().List())
	var labels []string
	for _, object := range accessibleObjects(w.Canvas().Overlays().Top()) {
		labels = append(labels, object.(fyne.Accessible).AccessibilityLabel())
	}
	assert.Contains(t, labels, "Dialog content")
	assert.Contains(t, labels, "Dismiss dialog")
}
