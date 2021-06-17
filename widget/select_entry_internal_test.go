package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestSelectEntry_Disableable(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

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

	test.TapCanvas(c, fyne.NewPos(0, 0))
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
	test.NewApp()
	defer test.NewApp()

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
	test.NewApp()
	defer test.NewApp()

	e := NewSelectEntry([]string{"one"})
	w := test.NewWindow(e)
	defer w.Close()
	entrySize := e.MinSize()
	w.Resize(entrySize.Add(fyne.NewSize(100, 100)))
	e.Resize(entrySize)

	// open the popup
	test.Tap(e.ActionItem.(fyne.Tappable))

	// first movement
	e.Move(fyne.NewPos(10, 10))
	assert.Equal(t, fyne.NewPos(10, 10), e.Entry.Position())
	assert.Equal(t,
		fyne.NewPos(10, 10+entrySize.Height-theme.InputBorderSize()),
		e.popUp.Position(),
	)

	// second movement
	e.Move(fyne.NewPos(30, 27))
	assert.Equal(t, fyne.NewPos(30, 27), e.Entry.Position())
	assert.Equal(t,
		fyne.NewPos(30, 27+entrySize.Height-theme.InputBorderSize()),
		e.popUp.Position(),
	)
}

func TestSelectEntry_DropDownResize(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

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
	largeOptionsMinWidth := optionsMinSize(largeOptions).Width

	labelHeight := NewLabel("W").MinSize().Height

	tests := map[string]struct {
		placeholder string
		value       string
		options     []string
		want        fyne.Size
	}{
		"empty": {
			want: fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+2*theme.Padding(), labelHeight),
		},
		"empty + small options": {
			options: smallOptions,
			want:    fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+2*theme.Padding(), labelHeight),
		},
		"empty + large options": {
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth, labelHeight),
		},
		"value": {
			value: "foo", // in a scroller
			want:  fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+2*theme.Padding(), labelHeight),
		},
		"large value + small options": {
			value:   "large", // in a scroller
			options: smallOptions,
			want:    fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+2*theme.Padding(), labelHeight),
		},
		"small value + large options": {
			value:   "small", // in a scroller
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth, labelHeight),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := NewSelectEntry(tt.options)
			e.PlaceHolder = tt.placeholder
			e.Text = tt.value
			assert.Equal(t, tt.want, e.MinSize())
		})
	}
}

func TestSelectEntry_SetOptions(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

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
	test.NewApp()
	defer test.NewApp()

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

func optionsMinSize(options []string) fyne.Size {
	var labels []*Label
	for _, option := range options {
		labels = append(labels, NewLabel(option))
	}
	minWidth := float32(0)
	minHeight := float32(0)
	for _, label := range labels {
		if minWidth < label.MinSize().Width {
			minWidth = label.MinSize().Width
		}
		minHeight += label.MinSize().Height
	}
	// padding between all options
	minHeight += float32(len(labels)-1) * theme.Padding()
	return fyne.NewSize(minWidth, minHeight)
}
