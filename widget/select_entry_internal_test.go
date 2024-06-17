package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestSelectEntry_Disableable(t *testing.T) {
	test.NewTempApp(t)

	options := []string{"A", "B", "C"}
	e := NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	assert.False(t, e.Disabled())
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled.xml", c)

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_opened.xml", c)

	areaPos, _ := c.InteractiveArea()
	test.TapCanvas(c, areaPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_tapped_selected.xml", c)

	e.Disable()
	assert.True(t, e.Disabled())
	test.AssertRendersToMarkup(t, "select_entry/disableable_disabled.xml", c)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_disabled.xml", c, "no drop-down when disabled")

	e.Enable()
	assert.False(t, e.Disabled())
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_tapped.xml", c)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/disableable_enabled_opened.xml", c)
}

func TestSelectEntry_DropDown(t *testing.T) {
	test.NewTempApp(t)

	options := []string{"A", "B", "C"}
	e := NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "select_entry/dropdown_initial.xml", c)
	assert.Nil(t, c.Overlays().Top())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)

	test.TapCanvas(c, fyne.NewPos(50, 15+2*(theme.Padding()+e.Size().Height)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_tapped_B.xml", c)
	assert.Equal(t, "B", e.Text)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_B_opened.xml", c)

	test.TapCanvas(c, fyne.NewPos(50, 15+3*(theme.Padding()+e.Size().Height)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_tapped_C.xml", c)
	assert.Equal(t, "C", e.Text)
}

func TestSelectEntry_DropDownMove(t *testing.T) {
	test.NewTempApp(t)

	e := NewSelectEntry([]string{"one"})
	w := test.NewWindow(e)
	defer w.Close()
	entrySize := e.MinSize()
	w.Resize(entrySize.Add(fyne.NewSize(100, 100)))
	e.Resize(entrySize)
	inset, _ := w.Canvas().InteractiveArea()

	// open the popup
	test.Tap(e.ActionItem.(fyne.Tappable))

	// first movement
	e.Move(fyne.NewPos(10, 10))
	assert.Equal(t, fyne.NewPos(10, 10), e.Entry.Position())
	assert.Equal(t,
		fyne.NewPos(10, 10+entrySize.Height-theme.InputBorderSize()).Subtract(inset),
		e.popUp.Position(),
	)

	// second movement
	e.Move(fyne.NewPos(30, 27))
	assert.Equal(t, fyne.NewPos(30, 27), e.Entry.Position())
	assert.Equal(t,
		fyne.NewPos(30, 27+entrySize.Height-theme.InputBorderSize()).Subtract(inset),
		e.popUp.Position(),
	)
}

func TestSelectEntry_DropDownResize(t *testing.T) {
	test.NewTempApp(t)

	options := []string{"A", "B", "C"}
	e := NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "select_entry/dropdown_initial.xml", c)
	assert.Nil(t, c.Overlays().Top())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)

	e.Resize(e.Size().Subtract(fyne.NewSize(20, 0)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened_shrunk.xml", c)

	e.Resize(e.Size().Add(fyne.NewSize(20, 0)))
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)
}

func TestSelectEntry_MinSize(t *testing.T) {
	smallOptions := []string{"A", "B", "C"}
	largeOptions := []string{"Large Option A", "Larger Option B", "Very Large Option C"}
	labelHeight := NewLabel("W").MinSize().Height

	// since we scroll content and don't prop window open with popup all combinations should be the same min
	tests := map[string]struct {
		placeholder string
		value       string
		options     []string
	}{
		"empty": {},
		"empty + small options": {
			options: smallOptions,
		},
		"empty + large options": {
			options: largeOptions,
		},
		"value": {
			value: "foo", // in a scroller
		},
		"large value + small options": {
			value:   "large", // in a scroller
			options: smallOptions,
		},
		"small value + large options": {
			value:   "small", // in a scroller
			options: largeOptions,
		},
	}

	minSize := fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+2*theme.Padding(), labelHeight)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := NewSelectEntry(tt.options)
			e.PlaceHolder = tt.placeholder
			e.Text = tt.value
			assert.Equal(t, minSize, e.MinSize())
		})
	}
}

func TestSelectEntry_SetOptions(t *testing.T) {
	test.NewTempApp(t)

	e := NewSelectEntry([]string{"A", "B", "C"})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_opened.xml", c)
	test.TapCanvas(c, switchPos)

	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_setopts.xml", c)
}

func TestSelectEntry_SetOptions_Empty(t *testing.T) {
	test.NewTempApp(t)

	e := NewSelectEntry([]string{})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, "select_entry/dropdown_empty_setopts.xml", c)
}

func dropDownIconWidth() float32 {
	return theme.IconInlineSize() + theme.Padding()
}

func emptyTextWidth() float32 {
	return NewLabel("M").MinSize().Width
}
