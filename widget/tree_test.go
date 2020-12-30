package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

var treeData = map[string][]string{
	"":    {"foo"},
	"foo": {"foobar"},
}

func TestTree_OpenClose(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		tree := widget.NewTreeWithStrings(treeData)

		closed := make(chan string, 1)
		opened := make(chan string, 1)
		tree.OnBranchClosed = func(uid widget.TreeNodeID) {
			closed <- uid
		}
		tree.OnBranchOpened = func(uid widget.TreeNodeID) {
			opened <- uid
		}

		assert.False(t, tree.IsBranchOpen("foo"))

		tree.OpenBranch("foo")
		assert.True(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-opened:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been opened")
		}

		tree.CloseBranch("foo")
		assert.False(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-closed:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been closed")
		}

		tree.ToggleBranch("foo")
		assert.True(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-opened:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been opened")
		}

		tree.ToggleBranch("foo")
		assert.False(t, tree.IsBranchOpen("foo"))

		select {
		case s := <-closed:
			assert.Equal(t, "foo", s)
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Branch should have been closed")
		}
	})
	t.Run("Missing", func(t *testing.T) {
		tree := widget.NewTreeWithStrings(treeData)

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

func TestTree_OpenCloseAll(t *testing.T) {
	data := map[string][]string{
		"":     {"foo0", "foo1", "foo2"},
		"foo0": {"foobar0"},
		"foo1": {"foobar1"},
		"foo2": {"foobar2"},
	}
	tree := widget.NewTreeWithStrings(data)

	tree.OpenAllBranches()
	assert.True(t, tree.IsBranchOpen("foo0"))
	assert.True(t, tree.IsBranchOpen("foo1"))
	assert.True(t, tree.IsBranchOpen("foo2"))

	tree.CloseAllBranches()
	assert.False(t, tree.IsBranchOpen("foo0"))
	assert.False(t, tree.IsBranchOpen("foo1"))
	assert.False(t, tree.IsBranchOpen("foo2"))
}

func TestTree_Layout(t *testing.T) {
	test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		items    map[string][]string
		selected widget.TreeNodeID
		opened   []widget.TreeNodeID
	}{
		"single_leaf": {
			items: map[string][]string{
				"": {"11111"},
			},
		},
		"single_leaf_selected": {
			items: map[string][]string{
				"": {"11111"},
			},
			selected: "11111",
		},
		"single_branch": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
		},
		"single_branch_selected": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			selected: "A",
		},
		"single_branch_opened": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			opened: []string{"A"},
		},
		"single_branch_opened_selected": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			opened:   []string{"A"},
			selected: "A",
		},
		"single_branch_opened_leaf_selected": {
			items: map[string][]string{
				"":  {"A"},
				"A": {"11111"},
			},
			opened:   []string{"A"},
			selected: "11111",
		},
		"multiple": {
			items: map[string][]string{
				"":  {"A", "B", "44444444444444444444"},
				"A": {"11111"},
				"B": {"2222222222"},
			},
		},
		"multiple_selected": {
			items: map[string][]string{
				"":  {"A", "B", "44444444444444444444"},
				"A": {"11111"},
				"B": {"2222222222"},
			},
			selected: "44444444444444444444",
		},
		"multiple_leaf": {
			items: map[string][]string{
				"": {"11111", "2222222222", "333333333333333", "44444444444444444444"},
			},
		},
		"multiple_leaf_selected": {
			items: map[string][]string{
				"": {"11111", "2222222222", "333333333333333", "44444444444444444444"},
			},
			selected: "2222222222",
		},
		"multiple_branch": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
		},
		"multiple_branch_selected": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			selected: "B",
		},
		"multiple_branch_opened": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			opened: []string{"A", "B", "C"},
		},
		"multiple_branch_opened_selected": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			opened:   []string{"A", "B", "C"},
			selected: "B",
		},
		"multiple_branch_opened_leaf_selected": {
			items: map[string][]string{
				"":  {"A", "B"},
				"A": {"11111"},
				"B": {"2222222222", "C"},
				"C": {"333333333333333"},
			},
			opened:   []string{"A", "B", "C"},
			selected: "2222222222",
		},
	} {
		t.Run(name, func(t *testing.T) {
			tree := widget.NewTreeWithStrings(tt.items)
			for _, o := range tt.opened {
				tree.OpenBranch(o)
			}
			tree.Select(tt.selected)

			window := test.NewWindow(tree)
			defer window.Close()
			window.Resize(fyne.NewSize(200, 300))

			tree.Refresh() // Force layout

			test.AssertImageMatches(t, "tree/layout_"+name+".png", window.Canvas().Capture())
		})
	}
}

func TestTree_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	tree := widget.NewTreeWithStrings(treeData)
	tree.OpenBranch("foo")

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(220, 220))

	tree.Refresh() // Force layout

	test.AssertImageMatches(t, "tree/theme_initial.png", window.Canvas().Capture())

	test.WithTestTheme(t, func() {
		tree.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "tree/theme_changed.png", window.Canvas().Capture())
	})
}

func TestTree_Move(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	tree := widget.NewTreeWithStrings(treeData)
	tree.OpenBranch("foo")

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(220, 220))

	test.AssertImageMatches(t, "tree/move_initial.png", window.Canvas().Capture())

	tree.Move(fyne.NewPos(20, 20))
	test.AssertImageMatches(t, "tree/move_moved.png", window.Canvas().Capture())
}

func TestTree_Refresh(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	value := "Foo Leaf"
	tree := widget.NewTreeWithStrings(treeData)
	tree.UpdateNode = func(uid widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		if uid == "foobar" {
			node.(*widget.Label).SetText(value)
			assert.False(t, branch)
		} else {
			node.(*widget.Label).SetText(uid)
			assert.True(t, branch)
		}
	}
	tree.OpenBranch("foo")

	window := test.NewWindow(tree)
	defer window.Close()
	window.Resize(fyne.NewSize(220, 220))

	test.AssertImageMatches(t, "tree/refresh_initial.png", window.Canvas().Capture())

	value = "Replaced"
	tree.Refresh()

	test.AssertImageMatches(t, "tree/refresh_replaced.png", window.Canvas().Capture())
}
