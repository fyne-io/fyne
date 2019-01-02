package widget

import (
	"testing"

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

	added := NewToolbarAction(theme.CutIcon(), func() {})
	toolbar.Append(added)
	assert.Equal(t, 2, len(toolbar.Items))
	assert.Equal(t, added, toolbar.Items[1])
}

func TestToolbar_Prepend(t *testing.T) {
	toolbar := NewToolbar(NewToolbarSpacer())
	assert.Equal(t, 1, len(toolbar.Items))

	prepend := NewToolbarAction(theme.CutIcon(), func() {})
	toolbar.Prepend(prepend)
	assert.Equal(t, 2, len(toolbar.Items))
	assert.Equal(t, prepend, toolbar.Items[0])
}
