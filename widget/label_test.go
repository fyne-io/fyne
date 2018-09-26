package widget

import "testing"

import "github.com/stretchr/testify/assert"

import "github.com/fyne-io/fyne"
import _ "github.com/fyne-io/fyne/test"
import "github.com/fyne-io/fyne/theme"

func TestLabel_MinSize(t *testing.T) {
	label := NewLabel("Test")
	min := label.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)

	label.SetText("Longer")
	assert.True(t, label.MinSize().Width > min.Width)
}

func TestLabel_Alignment(t *testing.T) {
	label := &Label{Text: "Test", Alignment: fyne.TextAlignTrailing}

	assert.Equal(t, fyne.TextAlignTrailing, label.Renderer().(*labelRenderer).texts[0].Alignment)
}

func TestText_MinSize_Multiline(t *testing.T) {
	text := NewLabel("Break")
	min := text.MinSize()

	text = NewLabel("Bre\nak")
	min2 := text.MinSize()
	assert.True(t, min2.Width < min.Width)
	assert.True(t, min2.Height > min.Height)

	yPos := -1
	for _, text := range text.Renderer().(*labelRenderer).texts {
		assert.True(t, text.CurrentSize().Height < min2.Height)
		assert.True(t, text.CurrentPosition().Y > yPos)
		yPos = text.CurrentPosition().Y
	}
}
