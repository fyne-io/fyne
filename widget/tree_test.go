package widget

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var (
	aMin = textMinSize("a", theme.TextSize(), fyne.TextStyle{})
	bMin = textMinSize("b", theme.TextSize(), fyne.TextStyle{})
	cMin = textMinSize("c", theme.TextSize(), fyne.TextStyle{})
	dMin = textMinSize("d", theme.TextSize(), fyne.TextStyle{})

	getBranchIcon = func(path []string, open bool) fyne.Resource {
		if open {
			return theme.FolderOpenIcon()
		}
		return theme.FolderIcon()
	}
	getLeafIcon = func(path []string) fyne.Resource {
		return theme.FileIcon()
	}
)

func getBranch(t *testing.T, tree *Tree, path ...string) (branch *branch) {
	t.Helper()
	branch = tree.root
	for _, p := range path {
		branch = branch.branches[p]
	}
	return
}

func getLeaf(t *testing.T, tree *Tree, path ...string) (leaf *leaf) {
	t.Helper()
	return getBranch(t, tree, path[:len(path)-1]...).leaves[path[len(path)-1]]
}

func tapBranch(t *testing.T, tree *Tree, path ...string) {
	t.Helper()
	tapped := make(chan bool)
	tree.OnBranchChanged = func(path []string, open bool) {
		tapped <- open
	}
	go test.Tap(getBranch(t, tree, path...))
	select {
	case <-tapped:
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Branch should have been changed")
	}
}

func tapLeaf(t *testing.T, tree *Tree, path ...string) {
	t.Helper()
	selected := make(chan bool)
	tree.OnLeafSelected = func(path []string) {
		selected <- true
	}
	go test.Tap(getLeaf(t, tree, path...))
	select {
	case <-selected:
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Leaf should have been selected")
	}
}

func TestTree_Add(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b", "c")
}

func TestTree_Resize(t *testing.T) {
	tree := NewTree()
	s := 100
	tree.Resize(fyne.NewSize(s, s))
	assert.Equal(t, s, tree.Size().Width)
	assert.Equal(t, s, tree.Size().Height)

	tree.Add("a")
	tree.Add("b", "c")

	tapBranch(t, tree, "b")
	a := getLeaf(t, tree, "a")
	b := getBranch(t, tree, "b")
	c := getLeaf(t, tree, "b", "c")

	// Branches first
	assert.Equal(t, theme.Padding(), b.Position().X)
	assert.Equal(t, 0, b.Position().Y)
	assert.Equal(t, s-theme.Padding()*3, b.Size().Width)
	assert.Equal(t, bMin.Height, b.Size().Height)

	// Leaf position relative to parent branch
	assert.Equal(t, theme.Padding(), c.Position().X)
	assert.Equal(t, bMin.Height, c.Position().Y)
	assert.Equal(t, s-theme.Padding()*4, c.Size().Width)
	assert.Equal(t, cMin.Height, c.Size().Height)

	// Leaves second
	assert.Equal(t, theme.Padding(), a.Position().X)
	assert.Equal(t, bMin.Height+cMin.Height, a.Position().Y)
	assert.Equal(t, s-theme.Padding()*3, a.Size().Width)
	assert.Equal(t, aMin.Height, a.Size().Height)
}

func TestTree_MinSize(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		tree := NewTree()
		min := tree.MinSize()
		assert.Equal(t, theme.Padding()*2, min.Width)
		assert.Equal(t, theme.Padding()*2, min.Height)
	})
	t.Run("Single", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b")
		min := tree.MinSize()
		assert.Equal(t, aMin.Width+theme.Padding()*3, min.Width)
		assert.Equal(t, aMin.Height+theme.Padding()*2, min.Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b")
		tree.Add("c", "d")
		min := tree.MinSize()
		assert.Equal(t, fyne.Max(aMin.Width, cMin.Width)+theme.Padding()*3, min.Width)
		assert.Equal(t, aMin.Height+cMin.Height+theme.Padding()*2, min.Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b")
		tree.Add("c", "d")

		tapBranch(t, tree, "a")
		tapBranch(t, tree, "c")

		min := tree.MinSize()
		expectedWidth := fyne.Max(
			fyne.Max(aMin.Width, bMin.Width+theme.Padding()),
			fyne.Max(cMin.Width, dMin.Width+theme.Padding()),
		)
		expectedHeight := aMin.Height + bMin.Height + cMin.Height + dMin.Height
		assert.Equal(t, expectedWidth+theme.Padding()*3, min.Width)
		assert.Equal(t, expectedHeight+theme.Padding()*2, min.Height)
	})
}

func TestTree_MinSize_Icon(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		tree := NewTree()
		tree.GetBranchIcon = getBranchIcon
		tree.GetLeafIcon = getLeafIcon
		min := tree.MinSize()
		assert.Equal(t, theme.Padding()*2, min.Width)
		assert.Equal(t, theme.Padding()*2, min.Height)
	})
	t.Run("Single", func(t *testing.T) {
		tree := NewTree()
		tree.GetBranchIcon = getBranchIcon
		tree.GetLeafIcon = getLeafIcon
		tree.Add("a", "b")
		min := tree.MinSize()
		assert.Equal(t, theme.IconInlineSize()+aMin.Width+theme.Padding()*3, min.Width)
		assert.Equal(t, fyne.Max(theme.IconInlineSize(), aMin.Height)+theme.Padding()*2, min.Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.GetBranchIcon = getBranchIcon
		tree.GetLeafIcon = getLeafIcon
		tree.Add("a", "b")
		tree.Add("c", "d")
		min := tree.MinSize()
		assert.Equal(t, theme.IconInlineSize()+fyne.Max(aMin.Width, cMin.Width)+theme.Padding()*3, min.Width)
		assert.Equal(t, fyne.Max(theme.IconInlineSize(), aMin.Height+cMin.Height)+theme.Padding()*2, min.Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.GetBranchIcon = getBranchIcon
		tree.GetLeafIcon = getLeafIcon
		tree.Add("a", "b")
		tree.Add("c", "d")

		tapBranch(t, tree, "a")
		tapBranch(t, tree, "c")

		min := tree.MinSize()
		expectedWidth := theme.IconInlineSize() + fyne.Max(
			fyne.Max(aMin.Width, bMin.Width+theme.Padding()),
			fyne.Max(cMin.Width, dMin.Width+theme.Padding()),
		)
		expectedHeight := fyne.Max(theme.IconInlineSize(), aMin.Height) +
			fyne.Max(theme.IconInlineSize(), bMin.Height) +
			fyne.Max(theme.IconInlineSize(), cMin.Height) +
			fyne.Max(theme.IconInlineSize(), dMin.Height)
		assert.Equal(t, expectedWidth+theme.Padding()*3, min.Width)
		assert.Equal(t, expectedHeight+theme.Padding()*2, min.Height)
	})
}

func TestTree_Tap(t *testing.T) {
	t.Run("Branch", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b")

		tapped := make(chan bool)
		tree.OnBranchChanged = func(path []string, open bool) {
			tapped <- open
		}
		go test.Tap(getBranch(t, tree, "a"))
		select {
		case open := <-tapped:
			assert.True(t, open, "Branch should be open")
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been changed")
		}
	})
	t.Run("Leaf", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")

		selected := make(chan bool)
		tree.OnLeafSelected = func(path []string) {
			selected <- true
		}
		go test.Tap(getLeaf(t, tree, "a"))
		select {
		case <-selected:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Leaf should have been selected")
		}
	})
}

func TestTreeRenderer_Layout(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		tree := NewTree()
		render := test.WidgetRenderer(tree).(*treeRenderer)
		render.Layout(render.MinSize())
	})
	t.Run("Single", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		render.Layout(render.MinSize())

		a := getLeaf(t, tree, "a")

		assert.Equal(t, theme.Padding(), a.Position().X)
		assert.Equal(t, 0, a.Position().Y)
		assert.Equal(t, aMin.Width, a.Size().Width)
		assert.Equal(t, aMin.Height, a.Size().Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		render.Layout(render.MinSize())
		a := getLeaf(t, tree, "a")
		b := getLeaf(t, tree, "b")

		assert.Equal(t, theme.Padding(), a.Position().X)
		assert.Equal(t, 0, a.Position().Y)
		assert.Equal(t, fyne.Max(aMin.Width, bMin.Width), a.Size().Width)
		assert.Equal(t, aMin.Height, a.Size().Height)

		assert.Equal(t, theme.Padding(), b.Position().X)
		assert.Equal(t, a.Size().Height, b.Position().Y)
		assert.Equal(t, fyne.Max(aMin.Width, bMin.Width), b.Size().Width)
		assert.Equal(t, bMin.Height, b.Size().Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b", "c")
		tapBranch(t, tree, "b")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		render.Layout(render.MinSize())
		a := getLeaf(t, tree, "a")
		b := getBranch(t, tree, "b")
		c := getLeaf(t, tree, "b", "c")

		width := fyne.Max(fyne.Max(aMin.Width, bMin.Width), cMin.Width+theme.Padding())

		// Branches first
		assert.Equal(t, theme.Padding(), b.Position().X)
		assert.Equal(t, 0, b.Position().Y)
		assert.Equal(t, width, b.Size().Width)
		assert.Equal(t, bMin.Height, b.Size().Height)

		// Leaf position relative to parent branch
		assert.Equal(t, theme.Padding(), c.Position().X)
		assert.Equal(t, b.Size().Height, c.Position().Y)
		assert.Equal(t, width-theme.Padding(), c.Size().Width)
		assert.Equal(t, cMin.Height, c.Size().Height)

		// Leaves second
		assert.Equal(t, theme.Padding(), a.Position().X)
		assert.Equal(t, b.Position().Y+b.Size().Height+c.Size().Height, a.Position().Y)
		assert.Equal(t, width, a.Size().Width)
		assert.Equal(t, aMin.Height, a.Size().Height)
	})
}

func TestTreeRenderer_MinSize(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		tree := NewTree()
		render := test.WidgetRenderer(tree).(*treeRenderer)
		min := render.MinSize()
		assert.Equal(t, theme.Padding()*2, min.Width)
		assert.Equal(t, theme.Padding()*2, min.Height)
	})
	t.Run("Single", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		min := render.MinSize()
		assert.Equal(t, aMin.Width+theme.Padding()*3, min.Width)
		assert.Equal(t, aMin.Height+theme.Padding()*2, min.Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b", "c")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		min := render.MinSize()
		assert.Equal(t, fyne.Max(aMin.Width, bMin.Width)+theme.Padding()*3, min.Width)
		assert.Equal(t, aMin.Height+bMin.Height+theme.Padding()*2, min.Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b", "c")
		tapBranch(t, tree, "b")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		min := render.MinSize()
		width := aMin.Width + theme.Padding()
		width = fyne.Max(width, bMin.Width+theme.Padding())
		width = fyne.Max(width, cMin.Width+theme.Padding()*2)
		assert.Equal(t, width+theme.Padding()*2, min.Width)
		assert.Equal(t, aMin.Height+bMin.Height+cMin.Height+theme.Padding()*2, min.Height)
	})
}

func TestTreeNodeRenderer_Layout(t *testing.T) {
	//
}

func TestTreeNodeRenderer_MinSize(t *testing.T) {
	//
}
