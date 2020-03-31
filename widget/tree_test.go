package widget

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var (
	aMin = fyne.MeasureText("a", theme.TextSize(), fyne.TextStyle{})
	bMin = fyne.MeasureText("b", theme.TextSize(), fyne.TextStyle{})
	cMin = fyne.MeasureText("c", theme.TextSize(), fyne.TextStyle{})
	dMin = fyne.MeasureText("d", theme.TextSize(), fyne.TextStyle{})
)

func getBranch(t *testing.T, root *branch, path ...string) (branch *branch) {
	t.Helper()
	branch = root
	for _, p := range path {
		branch = branch.branches[p]
	}
	return
}

func getLeaf(t *testing.T, root *branch, path ...string) (leaf *leaf) {
	t.Helper()
	return getBranch(t, root, path[:len(path)-1]...).leaves[path[len(path)-1]]
}

func tapBranch(t *testing.T, tree *Tree, path ...string) {
	t.Helper()
	tapped := make(chan bool)
	tree.OnBranchChanged = func(path []string, open bool) {
		tapped <- open
	}
	go test.Tap(getBranch(t, tree.root, path...))
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
	go test.Tap(getLeaf(t, tree.root, path...))
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
	a := getLeaf(t, tree.root, "a")
	b := getBranch(t, tree.root, "b")
	c := getLeaf(t, tree.root, "b", "c")

	// Branches first
	assert.Equal(t, theme.Padding(), b.Position().X)
	assert.Equal(t, 0, b.Position().Y)
	assert.Equal(t, s-theme.Padding()*3, b.Size().Width)
	assert.Equal(t, bMin.Height+theme.Padding()*2, b.Size().Height)

	// Leaf position relative to parent branch
	assert.Equal(t, theme.Padding(), c.Position().X)
	assert.Equal(t, bMin.Height+theme.Padding()*2, c.Position().Y)
	assert.Equal(t, s-theme.Padding()*4, c.Size().Width)
	assert.Equal(t, cMin.Height+theme.Padding()*2, c.Size().Height)

	// Leaves second
	assert.Equal(t, theme.Padding(), a.Position().X)
	assert.Equal(t, bMin.Height+cMin.Height+theme.Padding()*4, a.Position().Y)
	assert.Equal(t, s-theme.Padding()*3, a.Size().Width)
	assert.Equal(t, aMin.Height+theme.Padding()*2, a.Size().Height)
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
		assert.Equal(t, aMin.Width+theme.Padding()*5, min.Width)
		assert.Equal(t, aMin.Height+theme.Padding()*4, min.Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b")
		tree.Add("c", "d")
		min := tree.MinSize()
		assert.Equal(t, fyne.Max(aMin.Width, cMin.Width)+theme.Padding()*5, min.Width)
		assert.Equal(t, aMin.Height+cMin.Height+theme.Padding()*6, min.Height)
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
		assert.Equal(t, expectedWidth+theme.Padding()*5, min.Width)
		assert.Equal(t, expectedHeight+theme.Padding()*10, min.Height)
	})
}

func TestTree_MinSize_Icon(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		tree := NewTree()
		tree.UseArrowIcons()
		min := tree.MinSize()
		assert.Equal(t, theme.Padding()*2, min.Width)
		assert.Equal(t, theme.Padding()*2, min.Height)
	})
	t.Run("Single", func(t *testing.T) {
		tree := NewTree()
		tree.UseArrowIcons()
		tree.Add("a", "b")
		min := tree.MinSize()
		assert.Equal(t, theme.IconInlineSize()+aMin.Width+theme.Padding()*6, min.Width)
		assert.Equal(t, fyne.Max(theme.IconInlineSize(), aMin.Height)+theme.Padding()*4, min.Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.UseArrowIcons()
		tree.Add("a", "b")
		tree.Add("c", "d")
		min := tree.MinSize()
		assert.Equal(t, theme.IconInlineSize()+fyne.Max(aMin.Width, cMin.Width)+theme.Padding()*6, min.Width)
		assert.Equal(t, fyne.Max(theme.IconInlineSize(), aMin.Height+cMin.Height)+theme.Padding()*6, min.Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.UseArrowIcons()
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
		assert.Equal(t, expectedWidth+theme.Padding()*6, min.Width)
		assert.Equal(t, expectedHeight+theme.Padding()*10, min.Height)
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
		go test.Tap(getBranch(t, tree.root, "a"))
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
		go test.Tap(getLeaf(t, tree.root, "a"))
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

		a := getLeaf(t, tree.root, "a")

		assert.Equal(t, theme.Padding(), a.Position().X)
		assert.Equal(t, 0, a.Position().Y)
		assert.Equal(t, aMin.Width+theme.Padding()*2, a.Size().Width)
		assert.Equal(t, aMin.Height+theme.Padding()*2, a.Size().Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		render.Layout(render.MinSize())
		a := getLeaf(t, tree.root, "a")
		b := getLeaf(t, tree.root, "b")

		assert.Equal(t, theme.Padding(), a.Position().X)
		assert.Equal(t, 0, a.Position().Y)
		assert.Equal(t, fyne.Max(aMin.Width, bMin.Width)+theme.Padding()*2, a.Size().Width)
		assert.Equal(t, aMin.Height+theme.Padding()*2, a.Size().Height)

		assert.Equal(t, theme.Padding(), b.Position().X)
		assert.Equal(t, a.Size().Height, b.Position().Y)
		assert.Equal(t, fyne.Max(aMin.Width, bMin.Width)+theme.Padding()*2, b.Size().Width)
		assert.Equal(t, bMin.Height+theme.Padding()*2, b.Size().Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b", "c")
		tapBranch(t, tree, "b")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		render.Layout(render.MinSize())
		a := getLeaf(t, tree.root, "a")
		b := getBranch(t, tree.root, "b")
		c := getLeaf(t, tree.root, "b", "c")

		width := fyne.Max(fyne.Max(aMin.Width, bMin.Width), cMin.Width+theme.Padding())

		// Branches first
		assert.Equal(t, theme.Padding(), b.Position().X)
		assert.Equal(t, 0, b.Position().Y)
		assert.Equal(t, width+theme.Padding()*2, b.Size().Width)
		assert.Equal(t, bMin.Height+theme.Padding()*2, b.Size().Height)

		// Leaf position relative to parent branch
		assert.Equal(t, theme.Padding(), c.Position().X)
		assert.Equal(t, b.Size().Height, c.Position().Y)
		assert.Equal(t, width+theme.Padding(), c.Size().Width)
		assert.Equal(t, cMin.Height+theme.Padding()*2, c.Size().Height)

		// Leaves second
		assert.Equal(t, theme.Padding(), a.Position().X)
		assert.Equal(t, b.Position().Y+b.Size().Height+c.Size().Height, a.Position().Y)
		assert.Equal(t, width+theme.Padding()*2, a.Size().Width)
		assert.Equal(t, aMin.Height+theme.Padding()*2, a.Size().Height)
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
		assert.Equal(t, aMin.Width+theme.Padding()*5, min.Width)
		assert.Equal(t, aMin.Height+theme.Padding()*4, min.Height)
	})
	t.Run("Multiple", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b", "c")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		min := render.MinSize()
		assert.Equal(t, fyne.Max(aMin.Width, bMin.Width)+theme.Padding()*5, min.Width)
		assert.Equal(t, aMin.Height+bMin.Height+theme.Padding()*6, min.Height)
	})
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a")
		tree.Add("b", "c")
		tapBranch(t, tree, "b")
		render := test.WidgetRenderer(tree).(*treeRenderer)
		min := render.MinSize()
		width := aMin.Width + theme.Padding()*2
		width = fyne.Max(width, bMin.Width+theme.Padding()*3)
		width = fyne.Max(width, cMin.Width+theme.Padding()*4)
		assert.Equal(t, width+theme.Padding()*2, min.Width)
		assert.Equal(t, aMin.Height+bMin.Height+cMin.Height+theme.Padding()*8, min.Height)
	})
}

func TestTreeNode_getPath(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b", "c")
	a := getBranch(t, tree.root, "a")
	b := getBranch(t, tree.root, "a", "b")
	c := getLeaf(t, tree.root, "a", "b", "c")
	t.Run("Branch", func(t *testing.T) {
		path := a.getPath()
		assert.Equal(t, 1, len(path))
		assert.Equal(t, "a", path[0])

		path = b.getPath()
		assert.Equal(t, 2, len(path))
		assert.Equal(t, "a", path[0])
		assert.Equal(t, "b", path[1])
	})
	t.Run("Leaf", func(t *testing.T) {
		path := c.getPath()
		assert.Equal(t, 3, len(path))
		assert.Equal(t, "a", path[0])
		assert.Equal(t, "b", path[1])
		assert.Equal(t, "c", path[2])
	})
}

func TestTreeNode_isHovered(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b", "c")
	a := getBranch(t, tree.root, "a")
	b := getBranch(t, tree.root, "a", "b")
	c := getLeaf(t, tree.root, "a", "b", "c")
	t.Run("Branch", func(t *testing.T) {
		assert.False(t, a.isHovered())
		assert.False(t, b.isHovered())

		a.MouseIn(&desktop.MouseEvent{})
		assert.True(t, a.isHovered())
		assert.False(t, b.isHovered())
		a.MouseOut()

		b.MouseIn(&desktop.MouseEvent{})
		assert.False(t, a.isHovered())
		assert.True(t, b.isHovered())
		b.MouseOut()

		assert.False(t, a.isHovered())
		assert.False(t, b.isHovered())
	})
	t.Run("Leaf", func(t *testing.T) {
		assert.False(t, c.isHovered())

		c.MouseIn(&desktop.MouseEvent{})
		assert.True(t, c.isHovered())
		c.MouseOut()

		assert.False(t, c.isHovered())
	})
}

func TestTreeNode_walk(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b", "c")
		tree.Add("d", "e", "f")
		getBranch(t, tree.root, "a").setOpen(true)
		getBranch(t, tree.root, "a", "b").setOpen(true)
		getBranch(t, tree.root, "d").setOpen(true)
		getBranch(t, tree.root, "d", "e").setOpen(true)
		var branches []*branch
		var leaves []*leaf
		tree.root.walk(0, func(d int, b *branch) {
			branches = append(branches, b)
		}, func(d int, l *leaf) {
			leaves = append(leaves, l)
		})

		assert.Equal(t, 5, len(branches))
		assert.Equal(t, 2, len(leaves))

		assert.Equal(t, "", branches[0].getText())
		assert.Equal(t, "a", branches[1].getText())
		assert.Equal(t, "b", branches[2].getText())
		assert.Equal(t, "d", branches[3].getText())
		assert.Equal(t, "e", branches[4].getText())

		assert.Equal(t, "c", leaves[0].getText())
		assert.Equal(t, "f", leaves[1].getText())
	})
	t.Run("Closed", func(t *testing.T) {
		tree := NewTree()
		tree.Add("a", "b", "c")
		tree.Add("d", "e", "f")
		var branches []*branch
		var leaves []*leaf
		tree.root.walk(0, func(d int, b *branch) {
			branches = append(branches, b)
		}, func(d int, l *leaf) {
			leaves = append(leaves, l)
		})

		assert.Equal(t, 3, len(branches))
		assert.Equal(t, 0, len(leaves))

		assert.Equal(t, "", branches[0].getText())
		assert.Equal(t, "a", branches[1].getText())
		assert.Equal(t, "d", branches[2].getText())
	})
}

func TestTreeBranch_add(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		root := newBranch(nil, nil, "")
		root.add()
		assert.Equal(t, 0, len(root.branches))
		assert.Equal(t, 0, len(root.leaves))
	})
	t.Run("Single", func(t *testing.T) {
		root := newBranch(nil, nil, "")
		root.add("a")
		assert.Equal(t, 0, len(root.branches))
		assert.Equal(t, 1, len(root.leaves))
		assert.Equal(t, "a", root.leafNames[0])
	})
	t.Run("Multiple", func(t *testing.T) {
		root := newBranch(nil, nil, "")
		root.add("a", "b")
		root.add("c")
		assert.Equal(t, 1, len(root.branches))
		assert.Equal(t, 1, len(root.leaves))
		assert.Equal(t, "a", root.branchNames[0])
		assert.Equal(t, "c", root.leafNames[0])
	})
	t.Run("Multiple_Sort", func(t *testing.T) {
		root := newBranch(nil, nil, "")
		root.add("c")
		root.add("b")
		root.add("a")
		assert.Equal(t, 0, len(root.branches))
		assert.Equal(t, 3, len(root.leaves))
		assert.Equal(t, "a", root.leafNames[0])
		assert.Equal(t, "b", root.leafNames[1])
		assert.Equal(t, "c", root.leafNames[2])
	})
}

func TestTreeBranch_setOpen(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		root := newBranch(&Tree{}, nil, "")
		root.add("a", "b")
		a := getBranch(t, root, "a")
		b := getLeaf(t, root, "a", "b")
		root.open = true
		a.open = true
		root.update()
		a.update()
		assert.True(t, a.Visible())
		assert.True(t, b.Visible())
	})
	t.Run("Closed", func(t *testing.T) {
		root := newBranch(nil, nil, "")
		root.add("a", "b")
		a := getBranch(t, root, "a")
		b := getLeaf(t, root, "a", "b")
		assert.False(t, a.Visible())
		assert.False(t, b.Visible())
	})
}

func TestTreeNodeRenderer_Layout(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b")
	a := getBranch(t, tree.root, "a")
	b := getLeaf(t, tree.root, "a", "b")
	ar := test.WidgetRenderer(a).(*treeNodeRenderer)
	br := test.WidgetRenderer(b).(*treeNodeRenderer)
	s := 100
	ar.Layout(fyne.NewSize(s, s))
	br.Layout(fyne.NewSize(s, s))
	assert.Nil(t, ar.image)
	assert.Equal(t, s-theme.Padding()*2, ar.text.Size().Width)
	assert.Equal(t, aMin.Height, ar.text.Size().Height)
	assert.Nil(t, br.image)
	assert.Equal(t, s-theme.Padding()*2, br.text.Size().Width)
	assert.Equal(t, bMin.Height, br.text.Size().Height)
}

func TestTreeNodeRenderer_Layout_Icon(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b")
	tree.UseArrowIcons()
	a := getBranch(t, tree.root, "a")
	b := getLeaf(t, tree.root, "a", "b")
	ar := test.WidgetRenderer(a).(*treeNodeRenderer)
	br := test.WidgetRenderer(b).(*treeNodeRenderer)
	s := 100
	ar.createCanvasObjects()
	ar.Layout(fyne.NewSize(s, s))
	br.createCanvasObjects()
	br.Layout(fyne.NewSize(s, s))
	assert.Equal(t, theme.IconInlineSize(), ar.image.Size().Width)
	assert.Equal(t, theme.IconInlineSize(), ar.image.Size().Height)
	assert.Equal(t, s-theme.IconInlineSize()-theme.Padding()*3, ar.text.Size().Width)
	assert.Equal(t, aMin.Height, ar.text.Size().Height)
	assert.Equal(t, theme.IconInlineSize(), br.image.Size().Width)
	assert.Equal(t, theme.IconInlineSize(), br.image.Size().Height)
	assert.Equal(t, s-theme.IconInlineSize()-theme.Padding()*3, br.text.Size().Width)
	assert.Equal(t, bMin.Height, br.text.Size().Height)
}

func TestTreeNodeRenderer_MinSize(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b")
	a := getBranch(t, tree.root, "a")
	b := getLeaf(t, tree.root, "a", "b")
	ar := test.WidgetRenderer(a).(*treeNodeRenderer)
	br := test.WidgetRenderer(b).(*treeNodeRenderer)
	ar.createCanvasObjects()
	br.createCanvasObjects()
	assert.Equal(t, aMin.Width+theme.Padding()*2, ar.MinSize().Width)
	assert.Equal(t, aMin.Height+theme.Padding()*2, ar.MinSize().Height)
	assert.Equal(t, bMin.Width+theme.Padding()*2, br.MinSize().Width)
	assert.Equal(t, bMin.Height+theme.Padding()*2, br.MinSize().Height)
}

func TestTreeNodeRenderer_MinSize_Icon(t *testing.T) {
	tree := NewTree()
	tree.Add("a", "b")
	tree.UseArrowIcons()
	a := getBranch(t, tree.root, "a")
	b := getLeaf(t, tree.root, "a", "b")
	ar := test.WidgetRenderer(a).(*treeNodeRenderer)
	br := test.WidgetRenderer(b).(*treeNodeRenderer)
	ar.createCanvasObjects()
	br.createCanvasObjects()
	assert.Equal(t, aMin.Width+theme.IconInlineSize()+theme.Padding()*3, ar.MinSize().Width)
	assert.Equal(t, fyne.Max(aMin.Height, theme.IconInlineSize())+theme.Padding()*2, ar.MinSize().Height)
	assert.Equal(t, bMin.Width+theme.IconInlineSize()+theme.Padding()*3, br.MinSize().Width)
	assert.Equal(t, fyne.Max(bMin.Height, theme.IconInlineSize())+theme.Padding()*2, br.MinSize().Height)
}

func TestTreeNodeRenderer_BackgroundColor(t *testing.T) {
	tree := NewTree()
	tree.Add("a")
	tree.UseArrowIcons()
	a := getLeaf(t, tree.root, "a")
	ar := test.WidgetRenderer(a).(*treeNodeRenderer)
	ar.createCanvasObjects()
	assert.Equal(t, theme.BackgroundColor(), ar.BackgroundColor())
}

func TestTreeNodeRenderer_BackgroundColor_Hovered(t *testing.T) {
	tree := NewTree()
	tree.Add("a")
	tree.UseArrowIcons()
	a := getLeaf(t, tree.root, "a")
	a.MouseIn(&desktop.MouseEvent{})
	ar := test.WidgetRenderer(a).(*treeNodeRenderer)
	ar.createCanvasObjects()
	assert.Equal(t, theme.HoverColor(), ar.BackgroundColor())
	a.MouseOut()
	assert.Equal(t, theme.BackgroundColor(), ar.BackgroundColor())
}
