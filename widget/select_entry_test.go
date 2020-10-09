package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestSelectEntry_Disableable(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	assert.False(t, e.Disabled())
	test.AssertImageMatches(t, "select_entry/disableable_enabled.png", c.Capture())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/disableable_enabled_opened.png", c.Capture())

	test.TapCanvas(c, fyne.NewPos(0, 0))
	test.AssertImageMatches(t, "select_entry/disableable_enabled.png", c.Capture())

	e.Disable()
	assert.True(t, e.Disabled())
	test.AssertImageMatches(t, "select_entry/disableable_disabled.png", c.Capture())

	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/disableable_disabled.png", c.Capture(), "no drop-down when disabled")

	e.Enable()
	assert.False(t, e.Disabled())
	test.AssertImageMatches(t, "select_entry/disableable_enabled.png", c.Capture())

	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/disableable_enabled_opened.png", c.Capture())
}

func TestSelectEntry_DropDown(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	test.AssertImageMatches(t, "select_entry/dropdown_initial.png", c.Capture())
	assert.Nil(t, c.Overlays().Top())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/dropdown_empty_opened.png", c.Capture())

	test.TapCanvas(c, fyne.NewPos(50, 15+2*(theme.Padding()+e.Size().Height)))
	test.AssertImageMatches(t, "select_entry/dropdown_tapped_B.png", c.Capture())
	assert.Equal(t, "B", e.Text)

	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/dropdown_B_opened.png", c.Capture())

	test.TapCanvas(c, fyne.NewPos(50, 15+3*(theme.Padding()+e.Size().Height)))
	test.AssertImageMatches(t, "select_entry/dropdown_tapped_C.png", c.Capture())
	assert.Equal(t, "C", e.Text)
}

func TestSelectEntry_DropDownResize(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	test.AssertImageMatches(t, "select_entry/dropdown_initial.png", c.Capture())
	assert.Nil(t, c.Overlays().Top())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/dropdown_empty_opened.png", c.Capture())

	e.Resize(e.Size().Subtract(fyne.NewSize(20, 0)))
	test.AssertImageMatches(t, "select_entry/dropdown_empty_opened_shrunk.png", c.Capture())

	e.Resize(e.Size().Add(fyne.NewSize(20, 0)))
	test.AssertImageMatches(t, "select_entry/dropdown_empty_opened.png", c.Capture())
}

func TestSelectEntry_MinSize(t *testing.T) {
	smallOptions := []string{"A", "B", "C"}

	largeOptions := []string{"Large Option A", "Larger Option B", "Very Large Option C"}
	largeOptionsMinWidth := optionsMinSize(largeOptions).Width

	minTextHeight := widget.NewLabel("W").MinSize().Height

	tests := map[string]struct {
		placeholder string
		value       string
		options     []string
		want        fyne.Size
	}{
		"empty": {
			want: fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"empty + small options": {
			options: smallOptions,
			want:    fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"empty + large options": {
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"value": {
			value: "foo",
			want:  widget.NewLabel("foo").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"large value + small options": {
			value:   "large",
			options: smallOptions,
			want:    widget.NewLabel("large").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"small value + large options": {
			value:   "small",
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"placeholder": {
			placeholder: "example",
			want:        widget.NewLabel("example").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"large placeholder + small options": {
			placeholder: "large",
			options:     smallOptions,
			want:        widget.NewLabel("large").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"small placeholder + large options": {
			placeholder: "small",
			options:     largeOptions,
			want:        fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := widget.NewSelectEntry(tt.options)
			e.PlaceHolder = tt.placeholder
			e.Text = tt.value
			assert.Equal(t, tt.want, e.MinSize())
		})
	}
}

func TestSelectEntry_SetOptions(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	e := widget.NewSelectEntry([]string{"A", "B", "C"})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/dropdown_empty_opened.png", c.Capture())
	test.TapCanvas(c, switchPos)

	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/dropdown_empty_setopts.png", c.Capture())
}

func TestSelectEntry_SetOptions_Empty(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	e := widget.NewSelectEntry([]string{})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertImageMatches(t, "select_entry/dropdown_empty_setopts.png", c.Capture())
}

func dropDownIconWidth() int {
	dropDownIconWidth := theme.IconInlineSize() + theme.Padding()
	return dropDownIconWidth
}

func emptyTextWidth() int {
	return widget.NewLabel("M").MinSize().Width
}

func optionsMinSize(options []string) fyne.Size {
	var labels []*widget.Label
	for _, option := range options {
		labels = append(labels, widget.NewLabel(option))
	}
	minWidth := 0
	minHeight := 0
	for _, label := range labels {
		if minWidth < label.MinSize().Width {
			minWidth = label.MinSize().Width
		}
		minHeight += label.MinSize().Height
	}
	// padding between all options
	minHeight += (len(labels) - 1) * theme.Padding()
	return fyne.NewSize(minWidth, minHeight)
}
