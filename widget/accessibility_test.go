package widget_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestButton_Accessibility(t *testing.T) {
	tapped := 0
	b := widget.NewButton("Click", func() { tapped++ })

	assert.Equal(t, fyne.AccessibleRoleButton, b.AccessibilityRole())
	assert.Equal(t, "Click", b.AccessibilityLabel())
	assert.Equal(t, []fyne.AccessibleAction{fyne.AccessibleActionPress}, b.AccessibilityActions())

	assert.Empty(t, b.AccessibilityStates())
	b.Disable()
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateDisabled}, b.AccessibilityStates())
	assert.False(t, b.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 0, tapped)

	b.Enable()
	assert.True(t, b.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 1, tapped)

	assert.False(t, b.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
}

func TestHyperlink_Accessibility(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	hl := widget.NewHyperlink("Link", u)

	assert.Equal(t, fyne.AccessibleRoleLink, hl.AccessibilityRole())
	assert.Equal(t, "Link", hl.AccessibilityLabel())
	assert.Equal(t, []fyne.AccessibleAction{fyne.AccessibleActionPress}, hl.AccessibilityActions())

	tapped := 0
	hl.OnTapped = func() { tapped++ }
	assert.True(t, hl.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 1, tapped)
}

func TestCheck_Accessibility(t *testing.T) {
	changed := 0
	c := widget.NewCheck("Agree", func(bool) { changed++ })

	assert.Equal(t, fyne.AccessibleRoleCheckbox, c.AccessibilityRole())
	assert.Equal(t, "Agree", c.AccessibilityLabel())
	assert.Empty(t, c.AccessibilityStates())

	assert.True(t, c.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.True(t, c.Checked)
	assert.Equal(t, []fyne.AccessibleState{fyne.AccessibleStateChecked}, c.AccessibilityStates())
	assert.Equal(t, 1, changed)

	c.Disable()
	assert.False(t, c.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.True(t, c.Checked) // unchanged after disable
	assert.Contains(t, c.AccessibilityStates(), fyne.AccessibleStateDisabled)
}

func TestEntry_Accessibility(t *testing.T) {
	e := widget.NewEntry()
	e.PlaceHolder = "Email"
	e.Text = "user@example.com"

	assert.Equal(t, fyne.AccessibleRoleTextField, e.AccessibilityRole())
	assert.Equal(t, "Email", e.AccessibilityLabel())
	assert.Equal(t, "user@example.com", e.AccessibilityValue())

	assert.True(t, e.AccessibilitySetValue("new@example.com"))
	assert.Equal(t, "new@example.com", e.Text)

	e.Disable()
	assert.Contains(t, e.AccessibilityStates(), fyne.AccessibleStateDisabled)
	assert.False(t, e.AccessibilitySetValue("ignored"))
	assert.Equal(t, "new@example.com", e.Text)

	e.Enable()
	e.Validator = func(string) error { return nil }
	assert.Contains(t, e.AccessibilityStates(), fyne.AccessibleStateRequired)
}

func TestSelect_Accessibility(t *testing.T) {
	test.NewTempApp(t)

	s := widget.NewSelect([]string{"A", "B"}, nil)
	s.PlaceHolder = "Pick one"

	assert.Equal(t, fyne.AccessibleRoleButton, s.AccessibilityRole())
	assert.Equal(t, "Pick one", s.AccessibilityLabel())
	assert.Equal(t, "", s.AccessibilityValue())
	assert.Equal(t,
		[]fyne.AccessibleAction{fyne.AccessibleActionPress, fyne.AccessibleActionShowMenu},
		s.AccessibilityActions())

	s.SetSelected("A")
	assert.Equal(t, "A", s.AccessibilityValue())

	s.Disable()
	assert.Contains(t, s.AccessibilityStates(), fyne.AccessibleStateDisabled)
	assert.False(t, s.AccessibilityPerformAction(fyne.AccessibleActionShowMenu))
}

func TestSlider_Accessibility(t *testing.T) {
	s := widget.NewSlider(0, 10)
	s.Step = 2
	s.SetValue(4)

	assert.Equal(t, fyne.AccessibleRoleSlider, s.AccessibilityRole())
	assert.Equal(t, "4", s.AccessibilityValue())
	assert.Equal(t,
		[]fyne.AccessibleAction{
			fyne.AccessibleActionIncrement,
			fyne.AccessibleActionDecrement,
			fyne.AccessibleActionSetValue,
		},
		s.AccessibilityActions())

	assert.True(t, s.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
	assert.Equal(t, 6.0, s.Value)
	assert.True(t, s.AccessibilityPerformAction(fyne.AccessibleActionDecrement))
	assert.Equal(t, 4.0, s.Value)

	assert.True(t, s.AccessibilitySetValue("8"))
	assert.Equal(t, 8.0, s.Value)
	assert.False(t, s.AccessibilitySetValue("not-a-number"))

	s.Disable()
	assert.Contains(t, s.AccessibilityStates(), fyne.AccessibleStateDisabled)
	assert.False(t, s.AccessibilityPerformAction(fyne.AccessibleActionIncrement))
}

func TestProgressBar_Accessibility(t *testing.T) {
	p := widget.NewProgressBar()
	p.Min = 0
	p.Max = 200
	p.SetValue(50)

	assert.Equal(t, fyne.AccessibleRoleProgressBar, p.AccessibilityRole())
	assert.Equal(t, "25%", p.AccessibilityValue())

	p.TextFormatter = func() string { return "halfway" }
	p.SetValue(100)
	assert.Equal(t, "halfway", p.AccessibilityValue())
}

func TestProgressBarInfinite_Accessibility(t *testing.T) {
	p := widget.NewProgressBarInfinite()
	defer p.Stop()

	assert.Equal(t, fyne.AccessibleRoleProgressBar, p.AccessibilityRole())
	assert.Equal(t, "indeterminate", p.AccessibilityValue())
}

func TestRadioGroup_Accessibility(t *testing.T) {
	rg := widget.NewRadioGroup([]string{"A", "B"}, nil)

	assert.Equal(t, fyne.AccessibleRoleList, rg.AccessibilityRole())
	assert.Empty(t, rg.AccessibilityChildren()) // not yet rendered

	test.NewTempWindow(t, rg)
	assert.Len(t, rg.AccessibilityChildren(), 2)
}

func TestCheckGroup_Accessibility(t *testing.T) {
	cg := widget.NewCheckGroup([]string{"A", "B"}, nil)

	assert.Equal(t, fyne.AccessibleRoleList, cg.AccessibilityRole())
	assert.Empty(t, cg.AccessibilityChildren())

	test.NewTempWindow(t, cg)
	assert.Len(t, cg.AccessibilityChildren(), 2)
}

func TestForm_Accessibility(t *testing.T) {
	form := widget.NewForm(
		widget.NewFormItem("Name", widget.NewEntry()),
	)
	form.OnSubmit = func() {}

	assert.Equal(t, fyne.AccessibleRoleContainer, form.AccessibilityRole())

	test.NewTempWindow(t, form)
	children := form.AccessibilityChildren()
	assert.NotEmpty(t, children)
}

func TestList_Accessibility(t *testing.T) {
	data := []string{"A", "B", "C"}
	selected := -1
	l := widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(data[i]) },
	)
	l.OnSelected = func(id widget.ListItemID) { selected = id }

	assert.Equal(t, fyne.AccessibleRoleList, l.AccessibilityRole())
	assert.Empty(t, l.AccessibilityChildren()) // not yet rendered

	test.NewTempWindow(t, l)
	l.Resize(fyne.NewSize(200, 200))

	children := l.AccessibilityChildren()
	assert.NotEmpty(t, children)

	first, ok := children[0].(fyne.Accessible)
	assert.True(t, ok)
	assert.Equal(t, fyne.AccessibleRoleListItem, first.AccessibilityRole())
	assert.Equal(t, "A", first.AccessibilityLabel())

	actions, ok := children[0].(fyne.AccessibleActions)
	assert.True(t, ok)
	assert.True(t, actions.AccessibilityPerformAction(fyne.AccessibleActionSelect))
	assert.Equal(t, 0, selected)

	states, ok := children[0].(fyne.AccessibleStates)
	assert.True(t, ok)
	assert.Contains(t, states.AccessibilityStates(), fyne.AccessibleStateSelected)
}

func TestGridWrap_Accessibility(t *testing.T) {
	data := []string{"A", "B", "C", "D"}
	selected := -1
	g := widget.NewGridWrap(
		func() int { return len(data) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.GridWrapItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(data[i]) },
	)
	g.OnSelected = func(id widget.GridWrapItemID) { selected = id }

	assert.Equal(t, fyne.AccessibleRoleList, g.AccessibilityRole())

	test.NewTempWindow(t, g)
	g.Resize(fyne.NewSize(400, 400))

	children := g.AccessibilityChildren()
	assert.NotEmpty(t, children)

	first, ok := children[0].(fyne.Accessible)
	assert.True(t, ok)
	assert.Equal(t, fyne.AccessibleRoleListItem, first.AccessibilityRole())

	actions, ok := children[0].(fyne.AccessibleActions)
	assert.True(t, ok)
	assert.True(t, actions.AccessibilityPerformAction(fyne.AccessibleActionPress))
	assert.Equal(t, 0, selected)
}

func TestTree_Accessibility(t *testing.T) {
	data := map[string][]string{
		"":    {"A"},
		"A":   {"A.1", "A.2"},
		"A.1": {},
		"A.2": {},
	}
	tree := widget.NewTree(
		func(uid widget.TreeNodeID) []widget.TreeNodeID { return data[uid] },
		func(uid widget.TreeNodeID) bool {
			_, hasChildren := data[uid]
			return hasChildren && len(data[uid]) > 0
		},
		func(branch bool) fyne.CanvasObject { return widget.NewLabel("") },
		func(uid widget.TreeNodeID, branch bool, o fyne.CanvasObject) { o.(*widget.Label).SetText(uid) },
	)

	assert.Equal(t, fyne.AccessibleRoleTree, tree.AccessibilityRole())
	assert.Empty(t, tree.AccessibilityChildren())

	test.NewTempWindow(t, tree)
	tree.Resize(fyne.NewSize(200, 200))
	tree.OpenAllBranches()
	tree.Refresh()

	children := tree.AccessibilityChildren()
	assert.NotEmpty(t, children)

	first, ok := children[0].(fyne.Accessible)
	assert.True(t, ok)
	assert.Equal(t, fyne.AccessibleRoleTreeItem, first.AccessibilityRole())

	states, ok := children[0].(fyne.AccessibleStates)
	assert.True(t, ok)
	assert.Contains(t, states.AccessibilityStates(), fyne.AccessibleStateExpanded)
}

func TestTable_Accessibility(t *testing.T) {
	data := [][]string{{"a", "b"}, {"c", "d"}}
	tbl := widget.NewTable(
		func() (int, int) { return len(data), len(data[0]) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[id.Row][id.Col])
		},
	)

	assert.Equal(t, fyne.AccessibleRoleTable, tbl.AccessibilityRole())

	test.NewTempWindow(t, tbl)
	tbl.Resize(fyne.NewSize(200, 200))

	children := tbl.AccessibilityChildren()
	assert.NotEmpty(t, children)
}
