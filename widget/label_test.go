package widget

import (
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestLabel_MinSize(t *testing.T) {
	label := NewLabel("Test")
	min := label.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)

	label.SetText("Longer")
	assert.True(t, label.MinSize().Width > min.Width)
}
func TestLabel_Alignment(t *testing.T) {
	label := &Label{Text: "Test", Alignment: fyne.TextAlignTrailing}
	assert.Equal(t, fyne.TextAlignTrailing, Renderer(label).(*textRenderer).texts[0].Alignment)
}

func TestText_MinSize_MultiLine(t *testing.T) {
	text := NewLabel("Break")
	min := text.MinSize()
	text = NewLabel("Bre\nak")
	min2 := text.MinSize()

	assert.True(t, min2.Width < min.Width)
	assert.True(t, min2.Height > min.Height)

	yPos := -1
	for _, text := range Renderer(text).(*textRenderer).texts {
		assert.True(t, text.Size().Height < min2.Height)
		assert.True(t, text.Position().Y > yPos)
		yPos = text.Position().Y
	}
}
