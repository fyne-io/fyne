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

func TestLabel_Text(t *testing.T) {
	label := &Label{Text: "Test"}
	label.Refresh()

	assert.Equal(t, "Test", label.Text)
	assert.Equal(t, "Test", textRenderTexts(label)[0].Text)
}

func TestLabel_Text_Refresht(t *testing.T) {
	label := &Label{Text: ""}

	assert.Equal(t, "", label.Text)
	assert.Equal(t, "", textRenderTexts(label)[0].Text)

	label.Text = "Test"
	label.Refresh()
	assert.Equal(t, "Test", label.Text)
	assert.Equal(t, "Test", textRenderTexts(label)[0].Text)
}

func TestLabel_SetText(t *testing.T) {
	label := &Label{Text: "Test"}
	label.SetText("Crashy")
	label.Refresh()
	label.SetText("New")

	assert.Equal(t, "New", label.Text)
	assert.Equal(t, "New", textRenderTexts(label)[0].Text)
}

func TestLabel_Alignment(t *testing.T) {
	label := &Label{Text: "Test", Alignment: fyne.TextAlignTrailing}
	label.Refresh()

	assert.Equal(t, fyne.TextAlignTrailing, textRenderTexts(label)[0].Alignment)
}

func TestLabel_Alignment_Later(t *testing.T) {
	label := &Label{Text: "Test"}
	label.Refresh()
	assert.Equal(t, fyne.TextAlignLeading, textRenderTexts(label)[0].Alignment)

	label.Alignment = fyne.TextAlignTrailing
	label.Refresh()
	assert.Equal(t, fyne.TextAlignTrailing, textRenderTexts(label)[0].Alignment)
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

func TestText_MinSizeAdjustsWithContent(t *testing.T) {
	text := NewLabel("Line 1\nLine 2\n")
	initialMin := text.MinSize()

	text.SetText("Line 1\nLine 2\nLonger Line\n")
	assert.Greater(t, text.MinSize().Width, initialMin.Width)
	assert.Greater(t, text.MinSize().Height, initialMin.Height)

	text.SetText("Line 1\nLine 2\n")
	assert.Equal(t, initialMin, text.MinSize())
	assert.Equal(t, initialMin, text.textProvider.MinSize())
}

func TestLabel_ApplyTheme(t *testing.T) {
	text := NewLabel("Line 1")
	text.Hide()

	render := Renderer(text).(*textRenderer)
	assert.Equal(t, theme.TextColor(), render.texts[0].Color)
	text.Show()
	assert.Equal(t, theme.TextColor(), render.texts[0].Color)
}
