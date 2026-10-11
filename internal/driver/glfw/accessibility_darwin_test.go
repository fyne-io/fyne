//go:build accessibility && darwin

package glfw

import (
	"context"
	"image/color"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2/canvas"
)

func TestAccessibilityNativeUpdates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "accessibility-updates")
	//nolint:gosec // Fixed repository sources, compiled into a test-owned temporary directory.
	cmd := exec.CommandContext(ctx, "clang", "-I.", "-framework", "Cocoa",
		"testdata/accessibility_updates.m", "accessibility_darwin.m", "-o", binary)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
	for _, mode := range []string{
		"attributes", "ownership", "window-isolation", "mixed-classes", "super-method", "window-lifetime",
	} {
		t.Run(mode, func(t *testing.T) {
			//nolint:gosec // Executes the just-compiled fixture with fixed subtest arguments.
			output, err := exec.CommandContext(ctx, binary, mode).CombinedOutput()
			require.NoError(t, err, "%s", output)
		})
	}
}

func accessibilityTestNodes() []accessibilityNode {
	return []accessibilityNode{
		{object: canvas.NewRectangle(color.Black), role: roleToC("group"), label: "Group", parent: -1, width: 200, height: 100},
		{object: canvas.NewRectangle(color.White), role: roleToC("button"), label: "Button", parent: 0, width: 100, height: 30},
	}
}

func TestAccessibilityTreeChanged(t *testing.T) {
	tests := []struct {
		name    string
		change  func(n []accessibilityNode)
		rebuild bool
	}{
		{"unchanged", func([]accessibilityNode) {}, false},
		{"label", func(n []accessibilityNode) { n[1].label = "Updated" }, false},
		{"frame", func(n []accessibilityNode) { n[1].x, n[1].width = 10, 30 }, false},
		{"object", func(n []accessibilityNode) { n[1].object = canvas.NewRectangle(color.White) }, true},
		{"role", func(n []accessibilityNode) { n[1].role = roleToC("text") }, true},
		{"parent", func(n []accessibilityNode) { n[1].parent = -1 }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			previous := accessibilityTestNodes()
			next := slices.Clone(previous)
			test.change(next)
			require.Equal(t, test.rebuild, accessibilityTreeChanged(previous, next))
		})
	}
	t.Run("removal", func(t *testing.T) {
		previous := accessibilityTestNodes()
		require.True(t, accessibilityTreeChanged(previous, previous[:1]))
	})
}

func TestAccessibilityElementsPreserveIdentity(t *testing.T) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes(), nil)
	defer func() { destroyAccessibilityElements(nodes) }()
	root, child := nodes[0].element, nodes[1].element
	nodes = updateAccessibilityElements(nodes, slices.Clone(nodes), nil)
	require.Equal(t, root, nodes[0].element)
	require.Equal(t, child, nodes[1].element)
	next := slices.Clone(nodes)
	next[1].label = "Updated button"
	next[1].x = 20
	next[1].width = 120
	nodes = updateAccessibilityElements(nodes, next, nil)
	require.Equal(t, root, nodes[0].element)
	require.Equal(t, child, nodes[1].element)
}

func BenchmarkAccessibilityElementsUnchanged(b *testing.B) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes(), nil)
	defer func() { destroyAccessibilityElements(nodes) }()
	next := slices.Clone(nodes)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodes = updateAccessibilityElements(nodes, next, nil)
	}
}
