package widget

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

var (
	treeSize        = float32(200)
	templateMinSize = NewLabel("Template Object").MinSize()
)

func indentation() float32 {
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

func TestTree(t *testing.T) {
	t.Run("Initializer_Empty", func(t *testing.T) {
		tree := &Tree{}
		var nodes []string
		tree.walkAll(func(uid string, branch bool, depth int) {
			nodes = append(nodes, uid)
		})
		assert.Equal(t, 0, len(nodes))
	})
	t.Run("Initializer_Populated", func(t *testing.T) {
		tree := &Tree{
			ChildUIDs: func(uid string) (children []string) {
				if uid == "" {
					children = append(children, "a", "b", "c")
				} else if uid == "c" {
					children = append(children, "d", "e", "f")
				}
				return
			},
			IsBranch: func(uid string) bool {
				return uid == "" || uid == "c"
			},
			CreateNode: func(branch bool) fyne.CanvasObject {
				return &Label{}
			},
			UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
				node.(*Label).SetText(uid)
			},
		}
		tree.OpenBranch("c")
		var branches []string
		var leaves []string
		tree.walkAll(func(uid string, branch bool, depth int) {
			if branch {
				branches = append(branches, uid)
			} else {
				leaves = append(leaves, uid)
			}
		})
		assert.Equal(t, 2, len(branches))
		assert.Equal(t, "", branches[0])
		assert.Equal(t, "c", branches[1])
		assert.Equal(t, 5, len(leaves))
		assert.Equal(t, "a", leaves[0])
		assert.Equal(t, "b", leaves[1])
		assert.Equal(t, "d", leaves[2])
		assert.Equal(t, "e", leaves[3])
		assert.Equal(t, "f", leaves[4])
	})
	t.Run("NewTree", func(t *testing.T) {
		tree := NewTree(
			func(uid string) (children []string) {
				if uid == "" {
					children = append(children, "a", "b", "c")
				} else if uid == "c" {
					children = append(children, "d", "e", "f")
				}
				return
			},
			func(uid string) bool {
				return uid == "" || uid == "c"
			},
			func(branch bool) fyne.CanvasObject {
				return &Label{}
			},
			func(uid string, branch bool, node fyne.CanvasObject) {
				node.(*Label).SetText(uid)
			},
		)
		tree.OpenBranch("c")
		var branches []string
		var leaves []string
		tree.walkAll(func(uid string, branch bool, depth int) {
			if branch {
				branches = append(branches, uid)
			} else {
				leaves = append(leaves, uid)
			}
		})
		assert.Equal(t, 2, len(branches))
		assert.Equal(t, "", branches[0])
		assert.Equal(t, "c", branches[1])
		assert.Equal(t, 5, len(leaves))
		assert.Equal(t, "a", leaves[0])
		assert.Equal(t, "b", leaves[1])
		assert.Equal(t, "d", leaves[2])
		assert.Equal(t, "e", leaves[3])
		assert.Equal(t, "f", leaves[4])
	})
	t.Run("NewTreeWithStrings", func(t *testing.T) {
		data := make(map[string][]string)
		addTreePath(data, "foo", "foobar")
		tree := NewTreeWithStrings(data)
		tree.OpenBranch("foo")
		var branches []string
		var leaves []string
		tree.walkAll(func(uid string, branch bool, depth int) {
			if branch {
				branches = append(branches, uid)
			} else {
				leaves = append(leaves, uid)
			}
		})
		assert.Equal(t, 2, len(branches))
		assert.Equal(t, "", branches[0]) // Root
		assert.Equal(t, "foo", branches[1])
		assert.Equal(t, 1, len(leaves))
		assert.Equal(t, "foobar", leaves[0])
	})
}

func TestTree_Indentation(t *testing.T) {
	data := make(map[string][]string)
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	assert.Equal(t, treeSize, tree.Size().Width)
	assert.Equal(t, treeSize, tree.Size().Height)

	addTreePath(data, "A", "B", "C")

	tree.OpenAllBranches()

	a := getBranch(t, tree, "A")
	b := getBranch(t, tree, "B")
	c := getLeaf(t, tree, "C")

	assert.Equal(t, 0*indentation(), a.Indent())
	assert.Equal(t, 1*indentation(), b.Indent())
	assert.Equal(t, 2*indentation(), c.Indent())
}

func TestTree_Resize(t *testing.T) {
	data := make(map[string][]string)
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	assert.Equal(t, treeSize, tree.Size().Width)
	assert.Equal(t, treeSize, tree.Size().Height)

	addTreePath(data, "A")
	addTreePath(data, "B", "C")

	tree.OpenBranch("B")
	separatorThickness := theme.SeparatorThicknessSize()

	width := templateMinSize.Width + indentation() + theme.IconInlineSize() + theme.Padding()*2
	height := fyne.Max(templateMinSize.Height, theme.IconInlineSize())*3 + separatorThickness*2
	assertTreeContentMinSize(t, tree, fyne.NewSize(width, height))

	a := getLeaf(t, tree, "A")
	b := getBranch(t, tree, "B")
	c := getLeaf(t, tree, "C")

	assert.Equal(t, float32(0), a.Position().X)
	assert.Equal(t, float32(0), a.Position().Y)
	assert.Equal(t, treeSize, a.Size().Width)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize()), a.Size().Height)

	assert.Equal(t, float32(0), b.Position().X)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize())+separatorThickness, b.Position().Y)
	assert.Equal(t, treeSize, b.Size().Width)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize()), b.Size().Height)

	assert.Equal(t, float32(0), c.Position().X)
	assert.Equal(t, 2*(fyne.Max(templateMinSize.Height, theme.IconInlineSize())+separatorThickness), c.Position().Y)
	assert.Equal(t, treeSize, c.Size().Width)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize()), c.Size().Height)
}

func TestTree_MinSize(t *testing.T) {
	separatorThickness := theme.SeparatorThicknessSize()
	t.Run("Default", func(t *testing.T) {
		tree := &Tree{}
		min := tree.MinSize()
		assert.Equal(t, float32(32), min.Width)
		assert.Equal(t, float32(32), min.Height)
	})
	t.Run("Callback", func(t *testing.T) {
		for name, tt := range map[string]struct {
			leafSize        fyne.Size
			branchSize      fyne.Size
			expectedMinSize fyne.Size
		}{
			"small": {
				fyne.NewSize(1, 1),
				fyne.NewSize(1, 1),
				fyne.NewSize(fyne.Max(1+2*theme.Padding()+theme.IconInlineSize(), float32(32)), float32(32)),
			},
			"large-leaf": {
				fyne.NewSize(100, 100),
				fyne.NewSize(1, 1),
				fyne.NewSize(100+2*theme.Padding()+theme.IconInlineSize(), 100),
			},
			"large-branch": {
				fyne.NewSize(1, 1),
				fyne.NewSize(100, 100),
				fyne.NewSize(100+2*theme.Padding()+theme.IconInlineSize(), 100),
			},
		} {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, tt.expectedMinSize, (&Tree{
					CreateNode: func(isBranch bool) fyne.CanvasObject {
						r := canvas.NewRectangle(color.Black)
						if isBranch {
							r.SetMinSize(tt.branchSize)
							r.Resize(tt.branchSize)
						} else {
							r.SetMinSize(tt.leafSize)
							r.Resize(tt.leafSize)
						}
						return r
					},
				}).MinSize())
			})
		}
	})

	for name, tt := range map[string]struct {
		items  [][]string
		opened []string
		want   fyne.Size
	}{
		"single_item": {
			items: [][]string{
				{
					"A", "11111",
				},
			},
			want: fyne.NewSize(
				templateMinSize.Width+2*theme.Padding()+theme.IconInlineSize(),
				fyne.Max(templateMinSize.Height, theme.IconInlineSize()),
			),
		},
		"single_item_opened": {
			items: [][]string{
				{
					"A", "11111",
				},
			},
			opened: []string{"A"},
			want: fyne.NewSize(
				templateMinSize.Width+indentation()+2*theme.Padding()+theme.IconInlineSize(),
				(fyne.Max(templateMinSize.Height, theme.IconInlineSize()))*2+separatorThickness,
			),
		},
		"multiple_items": {
			items: [][]string{
				{
					"A", "11111",
				},
				{
					"B", "2222222222",
				},
				{
					"B", "C", "333333333333333",
				},
			},
			want: fyne.NewSize(
				templateMinSize.Width+2*theme.Padding()+theme.IconInlineSize(),
				(fyne.Max(templateMinSize.Height, theme.IconInlineSize()))*2+separatorThickness,
			),
		},
		"multiple_items_opened": {
			items: [][]string{
				{
					"A", "11111",
				},
				{
					"B", "2222222222",
				},
				{
					"B", "C", "333333333333333",
				},
			},
			opened: []string{"A", "B", "C"},
			want: fyne.NewSize(
				templateMinSize.Width+2*indentation()+theme.IconInlineSize()+2*theme.Padding(),
				(fyne.Max(templateMinSize.Height, theme.IconInlineSize()))*6+(5*separatorThickness),
			),
		},
	} {
		t.Run(name, func(t *testing.T) {
			data := make(map[string][]string)
			for _, d := range tt.items {
				addTreePath(data, d...)
			}
			tree := NewTreeWithStrings(data)
			for _, o := range tt.opened {
				tree.OpenBranch(o)
			}

			assertTreeContentMinSize(t, tree, tt.want)
		})
	}
}

func TestTree_Select(t *testing.T) {
	data := make(map[string][]string)
	addTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)

	tree.Refresh() // Force layout

	selection := make(chan string, 1)
	tree.OnSelected = func(uid TreeNodeID) {
		selection <- uid
	}

	tree.Select("A")
	assert.Equal(t, 1, len(tree.selected))
	assert.Equal(t, "A", tree.selected[0])
	select {
	case s := <-selection:
		assert.Equal(t, "A", s)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Selection should have occurred")
	}
}

func TestTree_Select_Unselects(t *testing.T) {
	data := make(map[string][]string)
	addTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)

	tree.Refresh() // Force layout
	tree.Select("A")

	unselection := make(chan string, 1)
	tree.OnUnselected = func(uid TreeNodeID) {
		unselection <- uid
	}

	tree.Select("B")
	assert.Equal(t, 1, len(tree.selected))
	select {
	case s := <-unselection:
		assert.Equal(t, "A", s)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Selection should have been unselected")
	}

	tree.Unselect("C")
	assert.Equal(t, 1, len(tree.selected))

	tree.UnselectAll()
	assert.Equal(t, 0, len(tree.selected))
	select {
	case s := <-unselection:
		assert.Equal(t, "B", s)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Selection should have been unselected")
	}
}

func TestTree_ScrollTo(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	data := make(map[string][]string)
	addTreePath(data, "A")
	addTreePath(data, "B", "C")
	addTreePath(data, "D", "E", "F")
	tree := NewTreeWithStrings(data)
	tree.OpenBranch("D")
	tree.OpenBranch("E")

	w := test.NewWindow(tree)
	defer w.Close()

	var (
		min = getLeaf(t, tree, "A").MinSize()
		sep = theme.SeparatorThicknessSize()
		pad = theme.Padding()
	)

	// Resize tall enough to display two nodes and the separater between them
	treeHeight := 2*(min.Height) + sep
	w.Resize(fyne.Size{
		Width:  400,
		Height: treeHeight + 2*pad,
	})

	tree.ScrollTo("F")

	want := 3*min.Height + 2*sep
	assert.Equal(t, want, tree.offset.Y)
	assert.Equal(t, want, tree.scroller.Offset.Y)
}

func TestTree_ScrollToBottom(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	data := make(map[string][]string)
	addTreePath(data, "A")
	addTreePath(data, "B", "C")
	addTreePath(data, "D", "E", "F")
	tree := NewTreeWithStrings(data)
	tree.OpenBranch("B")
	tree.OpenBranch("D")
	tree.OpenBranch("E")

	w := test.NewWindow(tree)
	defer w.Close()

	var (
		min = getLeaf(t, tree, "A").MinSize()
		sep = theme.SeparatorThicknessSize()
		pad = theme.Padding()
	)

	// Resize tall enough to display two nodes and the separater between them
	treeHeight := 2*(min.Height) + sep
	w.Resize(fyne.Size{
		Width:  400,
		Height: treeHeight + 2*pad,
	})

	tree.ScrollToBottom()

	want := 4 * (min.Height + sep)
	assert.Equal(t, want, tree.offset.Y)
	assert.Equal(t, want, tree.scroller.Offset.Y)
}

func TestTree_ScrollToSelection(t *testing.T) {
	data := make(map[string][]string)
	addTreePath(data, "A")
	addTreePath(data, "B")
	addTreePath(data, "C")
	addTreePath(data, "D")
	addTreePath(data, "E")
	addTreePath(data, "F")
	tree := NewTreeWithStrings(data)

	tree.Refresh() // Force layout

	a := getLeaf(t, tree, "A")
	m := a.MinSize()
	separatorThickness := theme.SeparatorThicknessSize()

	// Make tree tall enough to display two nodes
	tree.Resize(fyne.NewSize(m.Width, m.Height*2+separatorThickness))

	// Above
	tree.scroller.Offset.Y = m.Height*3 + separatorThickness*3 // Showing "D" & "E"
	tree.Refresh()                                             // Force layout
	tree.Select("A")
	// Tree should scroll to the top to show A
	assert.Equal(t, float32(0), tree.scroller.Offset.Y)

	// Below
	tree.scroller.Offset.Y = 0 // Showing "A" & "B"
	tree.Refresh()             // Force layout
	tree.Select("F")
	// Tree should scroll to the bottom to show F
	assert.Equal(t, m.Height*4+separatorThickness*3, tree.scroller.Offset.Y)
}

func TestTree_ScrollToTop(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	data := make(map[string][]string)
	addTreePath(data, "A")
	addTreePath(data, "B", "C")
	addTreePath(data, "D", "E", "F")
	tree := NewTreeWithStrings(data)
	tree.OpenBranch("D")
	tree.OpenBranch("E")

	w := test.NewWindow(tree)
	defer w.Close()

	tree.ScrollTo("F")

	tree.ScrollToTop()
	assert.Equal(t, float32(0), tree.offset.Y)
	assert.Equal(t, float32(0), tree.scroller.Offset.Y)
}

func TestTree_Tap(t *testing.T) {
	t.Run("Branch", func(t *testing.T) {
		data := make(map[string][]string)
		addTreePath(data, "A", "B")
		tree := NewTreeWithStrings(data)

		tree.Refresh() // Force layout

		selected := make(chan bool)
		tree.OnSelected = func(uid string) {
			selected <- true
		}
		go test.Tap(getBranch(t, tree, "A"))
		select {
		case <-selected:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been selected")
		}
	})
	t.Run("BranchIcon", func(t *testing.T) {
		data := make(map[string][]string)
		addTreePath(data, "A", "B")
		tree := NewTreeWithStrings(data)

		tree.Refresh() // Force layout

		tapped := make(chan bool)
		tree.OnBranchOpened = func(uid TreeNodeID) {
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
		addTreePath(data, "A")
		tree := NewTreeWithStrings(data)

		tree.Refresh() // Force layout

		selected := make(chan bool)
		tree.OnSelected = func(uid TreeNodeID) {
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

func TestTree_Unselect(t *testing.T) {
	data := make(map[string][]string)
	addTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)

	tree.Refresh() // Force layout
	tree.Select("A")

	unselection := make(chan string, 1)
	tree.OnUnselected = func(uid TreeNodeID) {
		unselection <- uid
	}

	tree.Unselect("A")
	assert.Equal(t, 0, len(tree.selected))
	select {
	case s := <-unselection:
		assert.Equal(t, "A", s)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Selection should have been cleared")
	}
}

func TestTree_Walk(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		data := make(map[string][]string)
		addTreePath(data, "A", "B", "C")
		addTreePath(data, "D", "E", "F")
		tree := NewTreeWithStrings(data)
		tree.OpenBranch("A")
		tree.OpenBranch("B")
		tree.OpenBranch("D")
		tree.OpenBranch("E")
		var branches []string
		var leaves []string
		tree.walkAll(func(uid string, branch bool, depth int) {
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
		addTreePath(data, "A", "B", "C")
		addTreePath(data, "D", "E", "F")
		tree := NewTreeWithStrings(data)
		var branches []string
		var leaves []string
		tree.walkAll(func(uid string, branch bool, depth int) {
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
	addTreePath(data, "A", "B", "C")
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	tree.OpenAllBranches()
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
	addTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	tree.OpenAllBranches()
	t.Run("Branch", func(t *testing.T) {
		a := getBranch(t, tree, "A")
		ar := test.WidgetRenderer(a).(*treeNodeRenderer)
		assert.Equal(t, theme.HoverColor(), ar.background.FillColor)
		assert.False(t, ar.background.Visible())
	})
	t.Run("Leaf", func(t *testing.T) {
		b := getLeaf(t, tree, "B")
		br := test.WidgetRenderer(b).(*treeNodeRenderer)
		assert.Equal(t, theme.HoverColor(), br.background.FillColor)
		assert.False(t, br.background.Visible())
	})
}

func TestTreeNodeRenderer_BackgroundColor_Hovered(t *testing.T) {
	data := make(map[string][]string)
	addTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	tree.OpenAllBranches()
	t.Run("Branch", func(t *testing.T) {
		a := getBranch(t, tree, "A")
		ar := test.WidgetRenderer(a).(*treeNodeRenderer)
		a.MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, theme.HoverColor(), ar.background.FillColor)
		assert.True(t, ar.background.Visible())
		a.MouseOut()
		assert.Equal(t, theme.HoverColor(), ar.background.FillColor)
		assert.False(t, ar.background.Visible())
	})
	t.Run("Leaf", func(t *testing.T) {
		b := getLeaf(t, tree, "B")
		br := test.WidgetRenderer(b).(*treeNodeRenderer)
		b.MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, theme.HoverColor(), br.background.FillColor)
		assert.True(t, br.background.Visible())
		b.MouseOut()
		assert.Equal(t, theme.HoverColor(), br.background.FillColor)
		assert.False(t, br.background.Visible())
	})
}

// addTreePath adds the given path to the given parent->children map
func addTreePath(data map[string][]string, path ...string) {
	parent := ""
	for _, p := range path {
		children := data[parent]
		add := true
		for _, c := range children {
			if c == p {
				add = false
				break
			}
		}
		if add {
			data[parent] = append(children, p)
		}
		parent = p
	}
}
