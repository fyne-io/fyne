package widget_test

import (
	//"io/ioutil"
	//"os"
	//"path"
	//"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne"
	//"fyne.io/fyne/storage"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestTreeContainer(t *testing.T) {
	t.Run("Initializer", func(t *testing.T) {
		tree := &widget.TreeContainer{}
		var nodes []string
		tree.Walk(func(id string, branch bool, depth int) {
			nodes = append(nodes, id)
		})
		assert.Equal(t, 0, len(nodes))
	})
	t.Run("NewTreeOfStrings", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "foo", "foobar")
		tree := widget.NewTreeOfStrings(data)
		tree.OpenBranch("foo")
		var branches []string
		var leaves []string
		tree.Walk(func(id string, branch bool, depth int) {
			if branch {
				branches = append(branches, id)
			} else {
				leaves = append(leaves, id)
			}
		})
		assert.Equal(t, 2, len(branches))
		assert.Equal(t, "", branches[0]) // Root
		assert.Equal(t, "foo", branches[1])
		assert.Equal(t, 1, len(leaves))
		assert.Equal(t, "foobar", leaves[0])
	})
	/* TODO Not currently possible as testDriver doesn't support listable URIs:
	2020/08/31 16:49:18 Fyne error:  Unable to get lister for /var/folders/8n/1dd15fbn79s2w3l4x43v5c7w0000gn/T/test011812343
	2020/08/31 16:49:18   Cause: test driver does support creating listable URIs yet
	t.Run("NewTreeOfFiles", func(t *testing.T) {
		tempDir, err := ioutil.TempDir("", "test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)
		err = os.MkdirAll(path.Join(tempDir, "A"), os.ModePerm)
		assert.NoError(t, err)
		err = os.MkdirAll(path.Join(tempDir, "B"), os.ModePerm)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(tempDir, "B", "C"), []byte("c"), os.ModePerm)
		assert.NoError(t, err)

		tree := widget.NewTreeOfFiles(storage.NewURI("file://" + tempDir))
		tree.OpenAllBranches()
		var branches []string
		var leaves []string
		tree.Walk(func(id string, branch bool, depth int) {
			if branch {
				branches = append(branches, id)
			} else {
				leaves = append(leaves, id)
			}
		})
		assert.Equal(t, 3, len(branches))
		assert.Equal(t, tempDir, branches[0]) // Root
		assert.Equal(t, filepath.Join(tempDir, "A"), branches[1])
		assert.Equal(t, filepath.Join(tempDir, "B"), branches[2])
		assert.Equal(t, 1, len(leaves))
		assert.Equal(t, filepath.Join(tempDir, "B", "C"), leaves[0])
	})
	*/
}

func TestTreeContainer_OpenClose(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "foo", "foobar")
		tree := widget.NewTreeOfStrings(data)

		assert.False(t, tree.IsBranchOpen("foo"))

		tree.OpenBranch("foo")
		assert.True(t, tree.IsBranchOpen("foo"))

		tree.CloseBranch("foo")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.ToggleBranch("foo")
		assert.True(t, tree.IsBranchOpen("foo"))

		tree.ToggleBranch("foo")
		assert.False(t, tree.IsBranchOpen("foo"))
	})
	t.Run("Missing", func(t *testing.T) {
		data := make(map[string][]string)
		widget.AddTreePath(data, "foo", "foobar")
		tree := widget.NewTreeOfStrings(data)

		assert.False(t, tree.IsBranchOpen("foo"))

		tree.OpenBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.CloseBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.ToggleBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))

		tree.ToggleBranch("bar")
		assert.False(t, tree.IsBranchOpen("foo"))
	})
}

func TestTreeContainer_OpenCloseAll(t *testing.T) {
	data := make(map[string][]string)
	widget.AddTreePath(data, "foo0", "foobar0")
	widget.AddTreePath(data, "foo1", "foobar1")
	widget.AddTreePath(data, "foo2", "foobar2")
	tree := widget.NewTreeOfStrings(data)

	tree.OpenAllBranches()
	assert.True(t, tree.IsBranchOpen("foo0"))
	assert.True(t, tree.IsBranchOpen("foo1"))
	assert.True(t, tree.IsBranchOpen("foo2"))

	tree.CloseAllBranches()
	assert.False(t, tree.IsBranchOpen("foo0"))
	assert.False(t, tree.IsBranchOpen("foo1"))
	assert.False(t, tree.IsBranchOpen("foo2"))
}

func TestTreeContainer_Layout(t *testing.T) {
	test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		items    [][]string
		selected string
		// TODO hovered  string
		opened []string
	}{
		"single_leaf": {
			items: [][]string{
				[]string{
					"11111",
				},
			},
		},
		/*
			"single_leaf_hovered": {
				items: [][]string{
					[]string{
						"11111",
					},
				},
				hovered: "11111",
			},
		*/
		"single_leaf_selected": {
			items: [][]string{
				[]string{
					"11111",
				},
			},
			selected: "11111",
		},
		"single_branch": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
		},
		/*
			"single_branch_hovered": {
				items: [][]string{
					[]string{
						"A", "11111",
					},
				},
				hovered: "A",
			},
		*/
		"single_branch_selected": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
			selected: "A",
		},
		"single_branch_opened": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
			opened: []string{"A"},
		},
		"single_branch_opened_selected": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
			opened:   []string{"A"},
			selected: "A",
		},
		/*
			"single_branch_opened_hovered": {
				items: [][]string{
					[]string{
						"A", "11111",
					},
				},
				opened:  []string{"A"},
				hovered: "A",
			},
		*/
		"single_branch_opened_leaf_selected": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
			},
			opened:   []string{"A"},
			selected: "11111",
		},
		/*
			"single_branch_opened_leaf_hovered": {
				items: [][]string{
					[]string{
						"A", "11111",
					},
				},
				opened:  []string{"A"},
				hovered: "11111",
			},
		*/
		"multiple": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
				[]string{
					"B", "2222222222",
				},
				[]string{
					"44444444444444444444",
				},
			},
		},
		/*
			"multiple_hovered": {
				items: [][]string{
					[]string{
						"A", "11111",
					},
					[]string{
						"B", "2222222222",
					},
					[]string{
						"44444444444444444444",
					},
				},
				hovered: "44444444444444444444",
			},
		*/
		"multiple_selected": {
			items: [][]string{
				[]string{
					"A", "11111",
				},
				[]string{
					"B", "2222222222",
				},
				[]string{
					"44444444444444444444",
				},
			},
			selected: "44444444444444444444",
		},
		"multiple_leaf": {
			items: [][]string{
				[]string{
					"11111",
				},
				[]string{
					"2222222222",
				},
				[]string{
					"333333333333333",
				},
				[]string{
					"44444444444444444444",
				},
			},
		},
		/*
			"multiple_leaf_hovered": {
				items: [][]string{
					[]string{
						"11111",
					},
					[]string{
						"2222222222",
					},
					[]string{
						"333333333333333",
					},
					[]string{
						"44444444444444444444",
					},
				},
				hovered: "2222222222",
			},
		*/
		"multiple_leaf_selected": {
			items: [][]string{
				[]string{
					"11111",
				},
				[]string{
					"2222222222",
				},
				[]string{
					"333333333333333",
				},
				[]string{
					"44444444444444444444",
				},
			},
			selected: "2222222222",
		},
		"multiple_branch": {
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
		},
		/*
			"multiple_branch_hovered": {
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
				hovered: "B",
			},
		*/
		"multiple_branch_selected": {
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
			selected: "B",
		},
		"multiple_branch_opened": {
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
		},
		"multiple_branch_opened_selected": {
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
			opened:   []string{"A", "B", "C"},
			selected: "B",
		},
		/*
			"multiple_branch_opened_hovered": {
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
				opened:  []string{"A", "B", "C"},
				hovered: "B",
			},
		*/
		"multiple_branch_opened_leaf_selected": {
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
			opened:   []string{"A", "B", "C"},
			selected: "2222222222",
		},
		/*
			"multiple_branch_opened_leaf_hovered": {
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
				opened:  []string{"A", "B", "C"},
				hovered: "2222222222",
			},
		*/
	} {
		t.Run(name, func(t *testing.T) {
			data := make(map[string][]string)
			for _, d := range tt.items {
				widget.AddTreePath(data, d...)
			}
			tree := widget.NewTreeOfStrings(data)
			for _, o := range tt.opened {
				tree.OpenBranch(o)
			}
			tree.SetSelectedNode(tt.selected)

			// TODO tree.Hovered = tt.hovered

			window := test.NewWindow(tree)
			window.Resize(tree.MinSize().Max(fyne.NewSize(200, 300)))

			test.AssertImageMatches(t, "tree/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}

func TestTree_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	data := make(map[string][]string)
	widget.AddTreePath(data, "foo", "foobar")
	tree := widget.NewTreeOfStrings(data)
	tree.OpenBranch("foo")

	w := test.NewWindow(tree)
	defer w.Close()
	w.Resize(fyne.NewSize(220, 220))

	test.AssertImageMatches(t, "tree/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		tree.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "tree/theme_changed.png", w.Canvas().Capture())
	})
}

func TestTree_Move(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	data := make(map[string][]string)
	widget.AddTreePath(data, "foo", "foobar")
	tree := widget.NewTreeOfStrings(data)
	tree.OpenBranch("foo")

	w := test.NewWindow(tree)
	defer w.Close()
	w.Resize(fyne.NewSize(220, 220))
	test.AssertImageMatches(t, "tree/move_initial.png", w.Canvas().Capture())

	tree.Move(fyne.NewPos(20, 20))
	test.AssertImageMatches(t, "tree/move_moved.png", w.Canvas().Capture())
}
