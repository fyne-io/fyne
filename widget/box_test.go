package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestBoxSize(t *testing.T) {
	list := NewVBox(NewLabel("Hello"), NewLabel("World"))
	assert.Equal(t, 2, len(list.Children))
}

func TestBoxPrepend(t *testing.T) {
	list := NewVBox(NewLabel("World"))
	assert.Equal(t, 1, len(list.Children))

	label := NewLabel("Hello")
	list.Prepend(label)
	assert.Equal(t, 2, len(list.Children))
	assert.Equal(t, label, list.Children[0])
}

func TestBoxAppend(t *testing.T) {
	list := NewVBox(NewLabel("Hello"))
	assert.Equal(t, 1, len(list.Children))

	label := NewLabel("World")
	list.Append(label)
	assert.True(t, len(list.Children) == 2)
	assert.Equal(t, label, list.Children[1])
}

func TestBox_ItemPositioning(t *testing.T) {
	a := NewLabel("A")
	b := NewLabel("B")
	for name, tt := range map[string]struct {
		horizontal    bool
		wantFirstPos  fyne.Position
		wantSecondPos fyne.Position
	}{
		"horizontal": {true, fyne.NewPos(0, 0), fyne.NewPos(a.MinSize().Width+theme.Padding(), 0)},
		"vertical":   {false, fyne.NewPos(0, 0), fyne.NewPos(0, a.MinSize().Height+theme.Padding())},
	} {
		t.Run(name, func(t *testing.T) {
			box := &Box{Horizontal: tt.horizontal, Children: []fyne.CanvasObject{a, b}}
			box.ExtendBaseWidget(box)
			items := make([]fyne.CanvasObject, 0)
			for _, o := range test.LaidOutObjects(box) {
				if l, ok := o.(*Label); ok {
					items = append(items, l)
				}
			}
			if assert.Equal(t, 2, len(items)) {
				assert.Equal(t, tt.wantFirstPos, items[0].Position())
				assert.Equal(t, tt.wantSecondPos, items[1].Position())
			}
		})
	}
}
