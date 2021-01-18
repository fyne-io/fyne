package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
)

type extendedLabel struct {
	Label
}

func newExtendedLabel(text string) *extendedLabel {
	label := &extendedLabel{}
	label.Text = text
	label.ExtendBaseWidget(label)
	return label
}

func TestLabel_Extended_SetText(t *testing.T) {
	label := newExtendedLabel("Start")
	objs := cache.Renderer(label).Objects()
	assert.Equal(t, 1, len(objs))
	assert.Equal(t, "Start", objs[0].(*canvas.Text).Text)

	label.SetText("Replace")
	assert.Equal(t, "Replace", objs[0].(*canvas.Text).Text)
}
