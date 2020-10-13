package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/painter/software"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestLabel_Hide(t *testing.T) {
	label := NewLabel("Test")
	label.CreateRenderer()
	label.Hide()
	label.Refresh()

	assert.True(t, label.Hidden)
	assert.False(t, label.provider.Hidden) // we don't propagate hide

	label.Show()
	assert.False(t, label.Hidden)
	assert.False(t, label.provider.Hidden)
}

func TestLabel_MinSize(t *testing.T) {
	label := NewLabel("Test")
	minA := label.MinSize()

	assert.Less(t, theme.Padding()*2, minA.Width)

	label.SetText("Longer")
	minB := label.MinSize()
	assert.Less(t, minA.Width, minB.Width)

	label.Text = "."
	minC := label.MinSize()
	assert.Greater(t, minB.Width, minC.Width)
}

func TestLabel_Resize(t *testing.T) {
	label := NewLabel("Test")
	label.CreateRenderer()
	size := fyne.NewSize(100, 20)
	label.Resize(size)

	assert.Equal(t, size, label.Size())
	assert.Equal(t, size, label.provider.Size())

	label.SetText("Longer")
	assert.Equal(t, size, label.Size())
	assert.Equal(t, size, label.provider.Size())
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
	textOneLine := NewLabel("Break")
	min := textOneLine.MinSize()
	textMultiLine := NewLabel("Bre\nak")
	min2 := textMultiLine.MinSize()

	assert.True(t, min2.Width < min.Width)
	assert.True(t, min2.Height > min.Height)

	yPos := -1
	for _, text := range test.WidgetRenderer(textMultiLine).(*textRenderer).texts {
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
	assert.Equal(t, initialMin, text.provider.MinSize())
}

func TestLabel_ApplyTheme(t *testing.T) {
	text := NewLabel("Line 1")
	text.Hide()

	render := test.WidgetRenderer(text).(*textRenderer)
	assert.Equal(t, theme.TextColor(), render.texts[0].Color)
	text.Show()
	assert.Equal(t, theme.TextColor(), render.texts[0].Color)
}

func TestLabel_CreateRendererDoesNotAffectSize(t *testing.T) {
	text := NewLabel("Hello")
	text.Resize(text.MinSize())
	size := text.Size()
	assert.NotEqual(t, fyne.NewSize(0, 0), size)
	assert.Equal(t, size, text.MinSize())

	r := text.CreateRenderer()
	assert.Equal(t, size, text.Size())
	assert.Equal(t, size, text.MinSize())
	assert.Equal(t, size, r.MinSize())
	r.Layout(size)
	assert.Equal(t, size, text.Size())
	assert.Equal(t, size, text.MinSize())
}

func TestLabel_ChangeTruncate(t *testing.T) {
	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetPadded(false)
	text := NewLabel("Hello")
	c.SetContent(text)
	c.Resize(text.MinSize())
	test.AssertImageMatches(t, "label/default.png", c.Capture())

	truncSize := text.MinSize().Subtract(fyne.NewSize(10, 0))
	text.Resize(truncSize)
	text.Wrapping = fyne.TextTruncate
	text.Refresh()
	test.AssertImageMatches(t, "label/truncate.png", c.Capture())
}
