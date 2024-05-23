package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestToolbarSize(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer(), NewToolbarAction(theme.HomeIcon(), func() {}))
	assert.Equal(t, 2, len(toolbar.Items))
	size := toolbar.MinSize()

	toolbar.Items = append(toolbar.Items, &toolbarLabel{NewLabel("Hi")})
	toolbar.Refresh()
	assert.Equal(t, size.Height, toolbar.MinSize().Height)
	assert.Greater(t, toolbar.MinSize().Width, size.Width)
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

func TestToolbar_Replace(t *testing.T) {
	icon := theme.ContentCutIcon()
	toolbar := NewToolbar(NewToolbarAction(icon, func() {}))
	assert.Equal(t, 1, len(toolbar.Items))
	render := test.TempWidgetRenderer(t, toolbar)
	assert.Equal(t, icon.Name(), render.Objects()[0].(*Button).Icon.Name())

	toolbar.Items[0] = NewToolbarAction(theme.HelpIcon(), func() {})
	toolbar.Refresh()
	assert.NotEqual(t, icon, render.Objects()[0].(*Button).Icon)
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
		assert.Equal(t, fyne.NewPos(40, 0), items[1].Position())
	}
}

func TestToolbar_SetIcon(t *testing.T) {
	oldIcon := theme.CancelIcon()
	toolbarItem := NewToolbarAction(oldIcon, nil)
	newIcon := theme.QuestionIcon()
	toolbarItem.SetIcon(newIcon)
	assert.NotEqual(t, oldIcon, toolbarItem.Icon)
	assert.Equal(t, newIcon, toolbarItem.Icon)
}

type toolbarLabel struct {
	*Label
}

func (t *toolbarLabel) ToolbarObject() fyne.CanvasObject {
	return t.Label
}

func TestToolbarAction_Disable(t *testing.T) {
	testIcon := theme.InfoIcon()
	toolbarAction := NewToolbarAction(testIcon, nil)
	toolbarAction.Disable()
	assert.NotEqual(t, false, toolbarAction.Disabled())
	assert.Equal(t, true, toolbarAction.Disabled())
}

func TestToolbarAction_Enable(t *testing.T) {
	testIcon := theme.InfoIcon()
	toolbarAction := NewToolbarAction(testIcon, nil)
	toolbarAction.Disable()
	toolbarAction.Enable()
	assert.NotEqual(t, true, toolbarAction.Disabled())
	assert.Equal(t, false, toolbarAction.Disabled())
}

func TestToolbarAction_UpdateOnActivated(t *testing.T) {
	activated := false

	testIcon := theme.InfoIcon()
	toolbarAction := NewToolbarAction(testIcon, func() { activated = true })

	test.Tap(toolbarAction.ToolbarObject().(*Button))

	assert.True(t, activated)

	activated = false

	// verify that changes are synchronized as well
	toolbarAction.OnActivated = func() {}

	test.Tap(toolbarAction.ToolbarObject().(*Button))

	assert.False(t, activated)
}

func TestToolbarAction_DefaultCreation(t *testing.T) {
	testIcon := theme.InfoIcon()
	toolbarAction := ToolbarAction{Icon: testIcon}
	obj := toolbarAction.ToolbarObject()
	assert.Equal(t, testIcon, obj.(*Button).Icon)
}
