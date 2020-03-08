package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			want: fyne.NewSize(emptyTextWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
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
			want:  widget.NewLabel("foo").MinSize().Add(fyne.NewSize(4*theme.Padding(), 2*theme.Padding())),
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
			want:        widget.NewLabel("example").MinSize().Add(fyne.NewSize(4*theme.Padding(), 2*theme.Padding())),
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

func TestSelectEntry_DropDown(t *testing.T) {
	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	c := w.Canvas()

	assert.Nil(t, c.Overlays().Top())

	canvasItems := test.InspectCanvasItems(c.Content(), c.Size())
	var dropDownSwitch fyne.Tappable
	for _, item := range canvasItems {
		if item.Type == "*widget.dropDownSwitch" {
			dropDownSwitch = item.Object.(fyne.Tappable)
			break
		}
	}
	require.NotNil(t, dropDownSwitch, "drop down switch not found")

	test.Tap(dropDownSwitch)
	require.NotNil(t, c.Overlays().Top(), "drop down didn't open")
	require.IsType(t, &widget.PopUp{}, c.Overlays().Top(), "drop down is not a *widget.PopUp")

	popUp := c.Overlays().Top().(*widget.PopUp)
	entryMinWidth := dropDownIconWidth() + emptyTextWidth() + 4*theme.Padding()
	assert.Equal(t, optionsMinSize(options).Union(fyne.NewSize(entryMinWidth-2*theme.Padding(), 0)).Add(fyne.NewSize(0, 2*theme.Padding())), popUp.Content.Size())
	assert.Equal(t, options, popUpOptions(popUp), "drop down menu texts don't match SelectEntry options")

	tapPopUpItem(popUp, 1)
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "B", e.Text)

	test.Tap(dropDownSwitch)
	popUp = c.Overlays().Top().(*widget.PopUp)
	tapPopUpItem(popUp, 2)
	assert.Nil(t, c.Overlays().Top())
	assert.Equal(t, "C", e.Text)
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
	return fyne.NewSize(minWidth, minHeight)
}

func popUpOptions(popUp *widget.PopUp) []string {
	var texts []string
	for _, item := range test.InspectCanvasItems(popUp.Content, popUp.Content.Size()) {
		if item.Type == "*canvas.Text" {
			texts = append(texts, item.Object.(*canvas.Text).Text)
		}
	}
	return texts
}

func tapPopUpItem(p *widget.PopUp, i int) {
	var items []fyne.Tappable
	for _, item := range test.InspectCanvasItems(p.Content, p.Content.Size()) {
		if t, ok := item.Object.(fyne.Tappable); ok {
			items = append(items, t)
		}
	}
	test.Tap(items[i])
}
