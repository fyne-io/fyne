package widget

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var (
	minSizeA = NewLabel("A").MinSize()
	minSizeB = NewLabel("B").MinSize()
	minSizeC = NewLabel("C").MinSize()
	minSize1 = NewLabel("11111").MinSize()
	minSize2 = NewLabel("2222222222").MinSize()
	minSize3 = NewLabel("333333333333333").MinSize()
)

func indentation() int {
	return theme.IconInlineSize() + theme.Padding()
}

func assertTreeContentMinSize(t *testing.T, tree *Tree, expected fyne.Size) {
	t.Helper()

	// Tree.MinSize() will always be 32 x 32 as this is the Scroller's min size, instead check treeContent.MinSize()

	tr := test.WidgetRenderer(tree).(*treeRenderer)
	assert.Equal(t, expected, tr.content.MinSize())
}

func getBranch(t *testing.T, tree *Tree, uid string) (branch *branch) {
	t.Helper()
	tr := test.WidgetRenderer(tree).(*treeRenderer)
	cr := test.WidgetRenderer(tr.content).(*treeContentRenderer)
	branch = cr.branches[uid]
	assert.NotNil(t, branch)
	return branch
}

func getLeaf(t *testing.T, tree *Tree, uid string) (leaf *leaf) {
	t.Helper()
	tr := test.WidgetRenderer(tree).(*treeRenderer)
	cr := test.WidgetRenderer(tr.content).(*treeContentRenderer)
	leaf = cr.leaves[uid]
	assert.NotNil(t, leaf)
	return leaf
}

func TestTree_Indentation(t *testing.T) {
	s := 100
	data := make(map[string][]string)
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(s, s))
	assert.Equal(t, s, tree.Size().Width)
	assert.Equal(t, s, tree.Size().Height)

	widget.AddTreePath(data, "A", "B", "C")

	tree.OpenAllBranches()

	a := getBranch(t, tree, "A")
	b := getBranch(t, tree, "B")
	c := getLeaf(t, tree, "C")

	assert.Equal(t, 0*indentation(), a.Indent())
	assert.Equal(t, 1*indentation(), b.Indent())
	assert.Equal(t, 2*indentation(), c.Indent())
}

func TestTree_Resize(t *testing.T) {
	s := 100
	data := make(map[string][]string)
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(s, s))
	assert.Equal(t, s, tree.Size().Width)
	assert.Equal(t, s, tree.Size().Height)

	widget.AddTreePath(data, "A")
	widget.AddTreePath(data, "B", "C")

	tree.OpenBranch("B")

	width := fyne.Max(
		fyne.Max(
			minSizeA.Width,
			minSizeB.Width,
		),
		minSizeC.Width+indentation(),
	) + theme.IconInlineSize() + 5*theme.Padding()
	height := 2 * theme.Padding()
	height = height + fyne.Max(minSizeA.Height, theme.IconInlineSize()) + 2*theme.Padding()
	height = height + fyne.Max(minSizeB.Height, theme.IconInlineSize()) + 2*theme.Padding()
	height = height + fyne.Max(minSizeC.Height, theme.IconInlineSize()) + 2*theme.Padding()
	assertTreeContentMinSize(t, tree, fyne.NewSize(width, height))

	a := getLeaf(t, tree, "A")
	b := getBranch(t, tree, "B")
	c := getLeaf(t, tree, "C")

	assert.Equal(t, theme.Padding(), a.Position().X)
	assert.Equal(t, theme.Padding(), a.Position().Y)
	assert.Equal(t, s-(2*theme.Padding()), a.Size().Width)
	assert.Equal(t, minSizeA.Height+(2*theme.Padding()), a.Size().Height)

	assert.Equal(t, theme.Padding(), b.Position().X)
	assert.Equal(t, minSizeA.Height+(3*theme.Padding()), b.Position().Y)
	assert.Equal(t, s-(2*theme.Padding()), b.Size().Width)
	assert.Equal(t, minSizeB.Height+(2*theme.Padding()), b.Size().Height)

	assert.Equal(t, theme.Padding(), c.Position().X)
	assert.Equal(t, minSizeA.Height+minSizeB.Height+(5*theme.Padding()), c.Position().Y)
	assert.Equal(t, s-(2*theme.Padding()), c.Size().Width)
	assert.Equal(t, minSizeC.Height+(2*theme.Padding()), c.Size().Height)
}

func TestTree_MinSize(t *testing.T) {
	contentPadding := 2 * theme.Padding()
	nodePadding := 2 * theme.Padding()

	t.Run("Default", func(t *testing.T) {
		tree := &Tree{}
		min := tree.MinSize()
		assert.Equal(t, 32, min.Width)
		assert.Equal(t, 32, min.Height)
	})
	t.Run("Callback", func(t *testing.T) {
		tree := &Tree{
			NewNode: func(isBranch bool) fyne.CanvasObject {
				if isBranch {
					return NewLabel("Branch")
				}
				return NewLabel("Leaf")
			},
		}
		tMin := tree.MinSize()
		bMin := NewLabel("Branch").MinSize()
		assert.Equal(t, fyne.Max(32, bMin.Width)+nodePadding, tMin.Width)
		assert.Equal(t, fyne.Max(32, bMin.Height)+nodePadding, tMin.Height)
	})

	for name, tt := range map[string]struct {
		items  [][]string
		opened []string
		want   fyne.Size
	}{
		"single_item": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
			want: fyne.NewSize(
				minSizeA.Width+theme.Padding()+theme.IconInlineSize()+nodePadding+contentPadding,
				fyne.Max(minSizeA.Height, theme.IconInlineSize())+nodePadding+contentPadding,
			),
		},
		"single_item_opened": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
			opened: []string{"A"},
			want: fyne.NewSize(
				fyne.Max(minSizeA.Width, minSize1.Width+indentation())+theme.Padding()+theme.IconInlineSize()+nodePadding+contentPadding,
				fyne.Max(minSizeA.Height, theme.IconInlineSize())+nodePadding+minSize1.Height+nodePadding+contentPadding,
			),
		},
		"multiple_items": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
				[]string{
					"B", "2222222222",
				},
				[]string{
					"B", "C", "333333333333333",
				},
			},
			want: fyne.NewSize(
				fyne.Max(minSizeA.Width, minSizeB.Width)+theme.Padding()+theme.IconInlineSize()+nodePadding+contentPadding,
				fyne.Max(minSizeA.Height, theme.IconInlineSize())+fyne.Max(minSizeB.Height, theme.IconInlineSize())+(2*nodePadding)+contentPadding,
			),
		},
		"multiple_items_opened": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
				[]string{
					"B", "2222222222",
				},
				[]string{
					"B", "C", "333333333333333",
				},
			},
			opened: []string{"A", "B", "C"},
			want: fyne.NewSize(
				fyne.Max(
					fyne.Max(
						fyne.Max(
							minSizeA.Width,
							minSizeB.Width,
						),
						minSizeC.Width+indentation(),
					),
					fyne.Max(
						fyne.Max(
							minSize1.Width+indentation(),
							minSize2.Width+indentation(),
						),
						minSize3.Width+(2*indentation()),
					),
				)+nodePadding+contentPadding+theme.IconInlineSize()+theme.Padding(),
				fyne.Max(minSizeA.Height, theme.IconInlineSize())+fyne.Max(minSizeB.Height, theme.IconInlineSize())+fyne.Max(minSizeC.Height, theme.IconInlineSize())+minSize1.Height+minSize2.Height+minSize3.Height+(6*nodePadding)+contentPadding,
			),
		},
	} {
		t.Run(name, func(t *testing.T) {
			data := make(map[string][]string)
			for _, d := range tt.items {
				widget.AddTreePath(data, d...)
			}
			tree := NewTreeWithStrings(data)
			for _, o := range tt.opened {
				tree.OpenBranch(o)
			}

			assertTreeContentMinSize(t, tree, tt.want)
		})
	}
}

func TestTree_Tap(t *testing.T) {
	/* TODO needs "fyne.io/fyne/test".DoubleTap(obj fyne.DoubleTappable)
	t.Run("Branch", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "A", "B")
		tree := NewTreeWithStrings(data)

		tree.Refresh() // Force layout

		tapped := make(chan bool)
		tree.OnBranchOpened = func(uid string) {
			tapped <- true
		}
		go test.DoubleTap(getBranch(t, tree, "A"))
		select {
		case open := <-tapped:
			assert.True(t, open, "Branch should be open")
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been changed")
		}
	})
	*/
	t.Run("BranchIcon", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "A", "B")
		tree := NewTreeWithStrings(data)

		tree.Refresh() // Force layout

		tapped := make(chan bool)
		tree.OnBranchOpened = func(uid string) {
			tapped <- true
		}
		go test.Tap(getBranch(t, tree, "A").icon.(*branchIcon))
		select {
		case open := <-tapped:
			assert.True(t, open, "Branch should be open")
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been changed")
		}
	})
	t.Run("Leaf", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "A")
		tree := NewTreeWithStrings(data)

		tree.Refresh() // Force layout

		selected := make(chan bool)
		tree.OnNodeSelected = func(uid string, node fyne.CanvasObject) {
			selected <- true
		}
		go test.Tap(getLeaf(t, tree, "A"))
		select {
		case <-selected:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Leaf should have been selected")
		}
	})
}

func TestTree_Walk(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "A", "B", "C")
		widget.AddTreePath(data, "D", "E", "F")
		tree := NewTreeWithStrings(data)
		tree.OpenBranch("A")
		tree.OpenBranch("B")
		tree.OpenBranch("D")
		tree.OpenBranch("E")
		var branches []string
		var leaves []string
		tree.Walk(func(uid string, branch bool, depth int) {
			if branch {
				branches = append(branches, uid)
			} else {
				leaves = append(leaves, uid)
			}
		})

		assert.Equal(t, 5, len(branches))
		assert.Equal(t, 2, len(leaves))

		assert.Equal(t, "", branches[0])
		assert.Equal(t, "A", branches[1])
		assert.Equal(t, "B", branches[2])
		assert.Equal(t, "D", branches[3])
		assert.Equal(t, "E", branches[4])

		assert.Equal(t, "C", leaves[0])
		assert.Equal(t, "F", leaves[1])
	})
	t.Run("Closed", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "A", "B", "C")
		widget.AddTreePath(data, "D", "E", "F")
		tree := NewTreeWithStrings(data)
		var branches []string
		var leaves []string
		tree.Walk(func(uid string, branch bool, depth int) {
			if branch {
				branches = append(branches, uid)
			} else {
				leaves = append(leaves, uid)
			}
		})

		assert.Equal(t, 3, len(branches))
		assert.Equal(t, 0, len(leaves))

		assert.Equal(t, "", branches[0])
		assert.Equal(t, "A", branches[1])
		assert.Equal(t, "D", branches[2])
	})
}

func TestTreeNode_Hovered(t *testing.T) {
	data := make(map[string][]string)
	widget.AddTreePath(data, "A", "B", "C")
	tree := NewTreeWithStrings(data)
	tree.OpenAllBranches()
	tree.Resize(fyne.NewSize(150, 200))
	a := getBranch(t, tree, "A")
	b := getBranch(t, tree, "B")
	c := getLeaf(t, tree, "C")
	t.Run("Branch", func(t *testing.T) {
		assert.False(t, a.hovered)
		assert.False(t, b.hovered)

		a.MouseIn(&desktop.MouseEvent{})
		assert.True(t, a.hovered)
		assert.False(t, b.hovered)
		a.MouseOut()

		b.MouseIn(&desktop.MouseEvent{})
		assert.False(t, a.hovered)
		assert.True(t, b.hovered)
		b.MouseOut()

		assert.False(t, a.hovered)
		assert.False(t, b.hovered)
	})
	t.Run("Leaf", func(t *testing.T) {
		assert.False(t, c.hovered)

		c.MouseIn(&desktop.MouseEvent{})
		assert.True(t, c.hovered)
		c.MouseOut()

		assert.False(t, c.hovered)
	})
}

func TestTreeNodeRenderer_BackgroundColor(t *testing.T) {
	data := make(map[string][]string)
	widget.AddTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)
	tree.OpenAllBranches()
	t.Run("Branch", func(t *testing.T) {
		a := getBranch(t, tree, "A")
		ar := test.WidgetRenderer(a).(*treeNodeRenderer)
		assert.Equal(t, theme.BackgroundColor(), ar.BackgroundColor())
	})
	t.Run("Leaf", func(t *testing.T) {
		b := getLeaf(t, tree, "B")
		br := test.WidgetRenderer(b).(*treeNodeRenderer)
		assert.Equal(t, theme.BackgroundColor(), br.BackgroundColor())
	})
}

func TestTreeNodeRenderer_BackgroundColor_Hovered(t *testing.T) {
	data := make(map[string][]string)
	widget.AddTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)
	tree.OpenAllBranches()
	t.Run("Branch", func(t *testing.T) {
		a := getBranch(t, tree, "A")
		ar := test.WidgetRenderer(a).(*treeNodeRenderer)
		a.MouseIn(&desktop.MouseEvent{})
		/* TODO FIXME removed for now
		assert.Equal(t, theme.HoverColor(), ar.BackgroundColor())
		replaced with:*/
		assert.Equal(t, theme.HoverColor(), ar.indicator.FillColor)
		/* EMXIF ODOT */
		a.MouseOut()
		assert.Equal(t, theme.BackgroundColor(), ar.BackgroundColor())
	})
	t.Run("Leaf", func(t *testing.T) {
		b := getLeaf(t, tree, "B")
		br := test.WidgetRenderer(b).(*treeNodeRenderer)
		b.MouseIn(&desktop.MouseEvent{})
		/* TODO FIXME removed for now
		assert.Equal(t, theme.HoverColor(), br.BackgroundColor())
		replaced with:*/
		assert.Equal(t, theme.HoverColor(), br.indicator.FillColor)
		/* EMXIF ODOT */
		b.MouseOut()
		assert.Equal(t, theme.BackgroundColor(), br.BackgroundColor())
	})
}
