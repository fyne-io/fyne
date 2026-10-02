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
	cmd := exec.CommandContext(ctx, "clang", "-I.", "-framework", "Cocoa",
		"testdata/accessibility_updates.m", "accessibility_darwin.m", "-o", binary)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
	for _, mode := range []string{"attributes", "own-method", "inherited-method"} {
		t.Run(mode, func(t *testing.T) {
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

func TestAccessibilityElementsPreserveIdentity(t *testing.T) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes())
	defer func() { destroyAccessibilityElements(nodes) }()
	root, child := nodes[0].element, nodes[1].element
	for i := 0; i < 100; i++ {
		nodes = updateAccessibilityElements(nodes, slices.Clone(nodes))
		require.Equal(t, root, nodes[0].element)
		require.Equal(t, child, nodes[1].element)
	}
	next := slices.Clone(nodes)
	next[1].label = "Updated button"
	next[1].x = 20
	next[1].width = 120
	nodes = updateAccessibilityElements(nodes, next)
	require.Equal(t, root, nodes[0].element)
	require.Equal(t, child, nodes[1].element)
}

func TestAccessibilityElementsReplaceChangedStructure(t *testing.T) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes())
	defer func() { destroyAccessibilityElements(nodes) }()
	next := append(slices.Clone(nodes), accessibilityNode{
		object: canvas.NewRectangle(color.Black), role: roleToC("text"),
		label: "New text", parent: 0, width: 100, height: 30,
	})
	nodes = updateAccessibilityElements(nodes, next)
	require.Len(t, nodes, 3)
	for _, node := range nodes {
		require.NotNil(t, node.element)
	}
	nodes = updateAccessibilityElements(nodes, nil)
	require.Empty(t, nodes)
	nodes = updateAccessibilityElements(nodes, accessibilityTestNodes())
	require.Len(t, nodes, 2)
}

func BenchmarkAccessibilityElementsUnchanged(b *testing.B) {
	nodes := updateAccessibilityElements(nil, accessibilityTestNodes())
	defer func() { destroyAccessibilityElements(nodes) }()
	next := slices.Clone(nodes)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodes = updateAccessibilityElements(nodes, next)
	}
}
