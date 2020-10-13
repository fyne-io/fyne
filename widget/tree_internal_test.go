package widget

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var (
	treeSize        = 200
	templateMinSize = NewLabel("Template Object").MinSize()
	doublePadding   = 2 * theme.Padding()
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
	t.Run("NewTreeWithFiles", func(t *testing.T) {
		tempDir, err := ioutil.TempDir("", "test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)
		err = os.MkdirAll(path.Join(tempDir, "A"), os.ModePerm)
		assert.NoError(t, err)
		err = os.MkdirAll(path.Join(tempDir, "B"), os.ModePerm)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(tempDir, "B", "C"), []byte("c"), os.ModePerm)
		assert.NoError(t, err)

		root := storage.NewURI("file://" + tempDir)
		tree := NewTreeWithFiles(root)
		tree.OpenAllBranches()
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
		assert.Equal(t, root.String(), branches[0]) // Root
		b1, err := storage.Child(root, "A")
		assert.NoError(t, err)
		assert.Equal(t, b1.String(), branches[1])
		b2, err := storage.Child(root, "B")
		assert.NoError(t, err)
		assert.Equal(t, b2.String(), branches[2])
		assert.Equal(t, 1, len(leaves))
		l1, err := storage.Child(root, "B")
		assert.NoError(t, err)
		l1, err = storage.Child(l1, "C")
		assert.NoError(t, err)
		assert.Equal(t, l1.String(), leaves[0])
	})
	t.Run("NewTreeWithStrings", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "foo", "foobar")
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
	data := make(map[string][]string)
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	assert.Equal(t, treeSize, tree.Size().Width)
	assert.Equal(t, treeSize, tree.Size().Height)

	widget.AddTreePath(data, "A")
	widget.AddTreePath(data, "B", "C")

	tree.OpenBranch("B")

	width := templateMinSize.Width + indentation() + theme.IconInlineSize() + doublePadding + theme.Padding()
	height := (fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding)*3 + treeDividerHeight*2
	assertTreeContentMinSize(t, tree, fyne.NewSize(width, height))

	a := getLeaf(t, tree, "A")
	b := getBranch(t, tree, "B")
	c := getLeaf(t, tree, "C")

	assert.Equal(t, 0, a.Position().X)
	assert.Equal(t, 0, a.Position().Y)
	assert.Equal(t, treeSize, a.Size().Width)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding, a.Size().Height)

	assert.Equal(t, 0, b.Position().X)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding+treeDividerHeight, b.Position().Y)
	assert.Equal(t, treeSize, b.Size().Width)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding, b.Size().Height)

	assert.Equal(t, 0, c.Position().X)
	assert.Equal(t, 2*(fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding+treeDividerHeight), c.Position().Y)
	assert.Equal(t, treeSize, c.Size().Width)
	assert.Equal(t, fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding, c.Size().Height)
}

func TestTree_MinSize(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		tree := &Tree{}
		min := tree.MinSize()
		assert.Equal(t, 32, min.Width)
		assert.Equal(t, 32, min.Height)
	})
	t.Run("Callback", func(t *testing.T) {
		tree := &Tree{
			CreateNode: func(isBranch bool) fyne.CanvasObject {
				if isBranch {
					return NewLabel("Branch")
				}
				return NewLabel("Leaf")
			},
		}
		tMin := tree.MinSize()
		bMin := newBranch(tree, NewLabel("Branch")).MinSize()
		assert.Equal(t, fyne.Max(32, bMin.Width), tMin.Width)
		assert.Equal(t, fyne.Max(32, bMin.Height), tMin.Height)
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
				templateMinSize.Width+theme.Padding()+theme.IconInlineSize()+doublePadding,
				fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding,
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
				templateMinSize.Width+indentation()+theme.Padding()+theme.IconInlineSize()+doublePadding,
				(fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding)*2+treeDividerHeight,
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
				templateMinSize.Width+theme.Padding()+theme.IconInlineSize()+doublePadding,
				(fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding)*2+treeDividerHeight,
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
				templateMinSize.Width+2*indentation()+doublePadding+theme.IconInlineSize()+theme.Padding(),
				(fyne.Max(templateMinSize.Height, theme.IconInlineSize())+doublePadding)*6+(5*treeDividerHeight),
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
		tree.OnNodeSelected = func(uid string) {
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
		widget.AddTreePath(data, "A", "B", "C")
		widget.AddTreePath(data, "D", "E", "F")
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
	widget.AddTreePath(data, "A", "B", "C")
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
	widget.AddTreePath(data, "A", "B")
	tree := NewTreeWithStrings(data)
	tree.Resize(fyne.NewSize(treeSize, treeSize))
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
	tree.Resize(fyne.NewSize(treeSize, treeSize))
	tree.OpenAllBranches()
	t.Run("Branch", func(t *testing.T) {
		a := getBranch(t, tree, "A")
		ar := test.WidgetRenderer(a).(*treeNodeRenderer)
		a.MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, theme.HoverColor(), ar.indicator.FillColor)
		a.MouseOut()
		assert.Equal(t, theme.BackgroundColor(), ar.BackgroundColor())
	})
	t.Run("Leaf", func(t *testing.T) {
		b := getLeaf(t, tree, "B")
		br := test.WidgetRenderer(b).(*treeNodeRenderer)
		b.MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, theme.HoverColor(), br.indicator.FillColor)
		b.MouseOut()
		assert.Equal(t, theme.BackgroundColor(), br.BackgroundColor())
	})
}
