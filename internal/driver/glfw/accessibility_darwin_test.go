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
		"attributes", "ownership", "own-method", "inherited-method",
		"window-isolation", "mixed-classes", "super-method", "window-lifetime",
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
	previous := accessibilityTestNodes()
	tests := []struct {
		name    string
		change  func([]accessibilityNode) []accessibilityNode
		rebuild bool
	}{
		{"unchanged", func(nodes []accessibilityNode) []accessibilityNode { return nodes }, false},
		{"label", func(nodes []accessibilityNode) []accessibilityNode {
			nodes[1].label = "Updated button"
			return nodes
		}, false},
		{"frame", func(nodes []accessibilityNode) []accessibilityNode {
			nodes[1].x, nodes[1].y, nodes[1].width, nodes[1].height = 10, 20, 30, 40
			return nodes
		}, false},
		{"object", func(nodes []accessibilityNode) []accessibilityNode {
			nodes[1].object = canvas.NewRectangle(color.White)
			return nodes
		}, true},
		{"role", func(nodes []accessibilityNode) []accessibilityNode {
			nodes[1].role = roleToC("text")
			return nodes
		}, true},
		{"parent", func(nodes []accessibilityNode) []accessibilityNode {
			nodes[1].parent = -1
			return nodes
		}, true},
		{"removal", func(nodes []accessibilityNode) []accessibilityNode { return nodes[:1] }, true},
		{"empty", func([]accessibilityNode) []accessibilityNode { return nil }, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.rebuild, accessibilityTreeChanged(previous, test.change(slices.Clone(previous))))
		})
	}
	t.Run("reordering", func(t *testing.T) {
		roots := slices.Clone(previous)
		roots[1].parent = -1
		next := []accessibilityNode{roots[1], roots[0]}
		require.True(t, accessibilityTreeChanged(roots, next))
	})
	require.False(t, accessibilityTreeChanged(nil, nil))
	require.True(t, accessibilityTreeChanged(nil, previous))
}

func TestAccessibilityElementsPreserveIdentity(t *testing.T) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes(), nil)
	defer func() { destroyAccessibilityElements(nodes) }()
	root, child := nodes[0].element, nodes[1].element
	for i := 0; i < 100; i++ {
		nodes = updateAccessibilityElements(nodes, slices.Clone(nodes), nil)
		require.Equal(t, root, nodes[0].element)
		require.Equal(t, child, nodes[1].element)
	}
	next := slices.Clone(nodes)
	next[1].label = "Updated button"
	next[1].x = 20
	next[1].width = 120
	nodes = updateAccessibilityElements(nodes, next, nil)
	require.Equal(t, root, nodes[0].element)
	require.Equal(t, child, nodes[1].element)
}

func TestAccessibilityElementsReplaceChangedStructure(t *testing.T) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes(), nil)
	defer func() { destroyAccessibilityElements(nodes) }()
	next := append(slices.Clone(nodes), accessibilityNode{
		object: canvas.NewRectangle(color.Black), role: roleToC("text"),
		label: "New text", parent: 0, width: 100, height: 30,
	})
	nodes = updateAccessibilityElements(nodes, next, nil)
	require.Len(t, nodes, 3)
	for _, node := range nodes {
		require.NotNil(t, node.element)
	}
	nodes = updateAccessibilityElements(nodes, nil, nil)
	require.Empty(t, nodes)
	nodes = updateAccessibilityElements(nodes, accessibilityTestNodes(), nil)
	require.Len(t, nodes, 2)
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
