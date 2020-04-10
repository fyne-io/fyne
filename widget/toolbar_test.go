package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestToolbarSize(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer(), NewToolbarSpacer())
	assert.Equal(t, 2, len(toolbar.Items))
}

func TestToolbar_Apppend(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer())
	assert.Equal(t, 1, len(toolbar.Items))

	added := NewToolbarAction(theme.ContentCutIcon(), func() {})
	toolbar.Append(added)
	assert.Equal(t, 2, len(toolbar.Items))
	assert.Equal(t, added, toolbar.Items[1])
}

func TestToolbar_Prepend(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer())
	assert.Equal(t, 1, len(toolbar.Items))

	prepend := NewToolbarAction(theme.ContentCutIcon(), func() {})
	toolbar.Prepend(prepend)
	assert.Equal(t, 2, len(toolbar.Items))
	assert.Equal(t, prepend, toolbar.Items[0])
}

func TestToolbar_ItemPositioning(t *testing.T) {
	toolbar := &Toolbar{
		Items: []ToolbarItem{
			NewToolbarAction(theme.ContentCopyIcon(), func() {}),
			NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		},
	}
	toolbar.ExtendBaseWidget(toolbar)
	toolbar.Refresh()
	var items []fyne.CanvasObject
	for _, o := range test.LaidOutObjects(toolbar) {
		if b, ok := o.(*Button); ok {
			items = append(items, b)
		}
	}
	if assert.Equal(t, 2, len(items)) {
		assert.Equal(t, fyne.NewPos(0, 0), items[0].Position())
		assert.Equal(t, fyne.NewPos(32, 0), items[1].Position())
	}
}
