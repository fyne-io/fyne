package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestLabel_Binding(t *testing.T) {
	label := NewLabel("Init")
	assert.Equal(t, "Init", label.Text)

	str := binding.NewString()
	label.Bind(str)
	waitForBinding()
	assert.Equal(t, "", label.Text)

	str.Set("Updated")
	waitForBinding()
	assert.Equal(t, "Updated", label.Text)

	label.Unbind()
	waitForBinding()
	assert.Equal(t, "Updated", label.Text)
}

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

	assert.Less(t, theme.InnerPadding(), minA.Width)

	label.SetText("Longer")
	minB := label.MinSize()
	assert.Less(t, minA.Width, minB.Width)

	label.Text = "."
	label.Refresh()
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
	assert.Equal(t, "Test", labelTextRenderTexts(label)[0].Text)
}

func TestLabel_Text_Refresh(t *testing.T) {
	label := &Label{Text: ""}

	assert.Equal(t, "", label.Text)
	assert.Equal(t, "", labelTextRenderTexts(label)[0].Text)

	label.Text = "Test"
	label.Refresh()
	assert.Equal(t, "Test", label.Text)
	assert.Equal(t, "Test", labelTextRenderTexts(label)[0].Text)
}

func TestLabel_SetText(t *testing.T) {
	label := &Label{Text: "Test"}
	label.SetText("Crashy")
	label.Refresh()
	label.SetText("New")

	assert.Equal(t, "New", label.Text)
	assert.Equal(t, "New", labelTextRenderTexts(label)[0].Text)
}

func TestLabel_Alignment(t *testing.T) {
	label := &Label{Text: "Test", Alignment: fyne.TextAlignTrailing}
	label.Refresh()

	assert.Equal(t, fyne.TextAlignTrailing, labelTextRenderTexts(label)[0].Alignment)
}

func TestLabel_Alignment_Later(t *testing.T) {
	label := &Label{Text: "Test"}
	label.Refresh()
	assert.Equal(t, fyne.TextAlignLeading, labelTextRenderTexts(label)[0].Alignment)

	label.Alignment = fyne.TextAlignTrailing
	label.Refresh()
	assert.Equal(t, fyne.TextAlignTrailing, labelTextRenderTexts(label)[0].Alignment)
}

func TestText_MinSize_MultiLine(t *testing.T) {
	textOneLine := NewLabel("Break")
	min := textOneLine.MinSize()
	textMultiLine := NewLabel("Bre\nak")
	rich := test.TempWidgetRenderer(t, textMultiLine).Objects()[0].(*RichText)
	min2 := textMultiLine.MinSize()

	assert.Less(t, min2.Width, min.Width)
	assert.Greater(t, min2.Height, min.Height)

	yPos := float32(-1)
	for _, text := range test.TempWidgetRenderer(t, rich).(*textRenderer).Objects() {
		assert.Less(t, text.Size().Height, min2.Height)
		assert.Greater(t, text.Position().Y, yPos)
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
	rich := test.TempWidgetRenderer(t, text).Objects()[0].(*RichText)

	render := test.TempWidgetRenderer(t, rich).(*textRenderer)
	assert.Equal(t, theme.Color(theme.ColorNameForeground), render.Objects()[0].(*canvas.Text).Color)
	text.Show()
	assert.Equal(t, theme.Color(theme.ColorNameForeground), render.Objects()[0].(*canvas.Text).Color)
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
	test.NewTempApp(t)

	c := test.NewCanvasWithPainter(software.NewPainter())
	c.SetPadded(false)
	text := NewLabel("Hello")
	c.SetContent(text)
	c.Resize(text.MinSize())
	test.AssertRendersToMarkup(t, "label/default.xml", c)

	truncSize := text.MinSize().Subtract(fyne.NewSize(10, 0))
	text.Resize(truncSize)
	text.Truncation = fyne.TextTruncateClip
	text.Refresh()
	test.AssertRendersToMarkup(t, "label/truncate.xml", c)
}

func TestLabel_Select(t *testing.T) {
	l := NewLabel("Hello")
	l.Selectable = true

	assert.Empty(t, l.SelectedText())
	assert.Equal(t, 2, len(test.WidgetRenderer(l).Objects()))

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	sel.MouseDown(&desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 10)}})
	sel.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 15, DY: 0},
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)}})
	sel.DragEnd()
	sel.MouseUp(&desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)}})
	assert.Equal(t, "el", l.SelectedText())

	sel.TypedShortcut(&fyne.ShortcutCopy{})
	assert.Equal(t, "el", fyne.CurrentApp().Clipboard().Content())

	l.Selectable = false
	l.Refresh()
	assert.Equal(t, 1, len(test.WidgetRenderer(l).Objects()))
	assert.Empty(t, l.SelectedText())
}

func TestLabel_SelectWord(t *testing.T) {
	l := NewLabel("Hello")
	l.Selectable = true

	assert.Empty(t, l.SelectedText())

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	sel.DoubleTapped(&fyne.PointEvent{Position: fyne.NewPos(15, 10)})
	assert.Equal(t, "Hello", l.SelectedText())
}

func TestLabel_SelectLine(t *testing.T) {
	l := NewLabel("Longer line")
	l.Selectable = true

	assert.Empty(t, l.SelectedText())

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	pointEvent := fyne.PointEvent{Position: fyne.NewPos(15, 10)}
	tapEvent := &desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
		PointEvent: pointEvent}
	sel.DoubleTapped(&pointEvent)
	sel.MouseDown(tapEvent)
	sel.MouseUp(tapEvent)

	assert.Equal(t, "Longer line", l.SelectedText())
}

func TestNewLabelWithData(t *testing.T) {
	str := binding.NewString()
	str.Set("Init")

	label := NewLabelWithData(str)
	waitForBinding()
	assert.Equal(t, "Init", label.Text)
}

func TestLabelImportance(t *testing.T) {
	test.NewTempApp(t)
	test.ApplyTheme(t, test.Theme())

	lbl := NewLabel("hello, fyne")
	w := test.NewWindow(lbl)
	defer w.Close()

	test.AssertImageMatches(t, "label/label_importance_medium.png", w.Canvas().Capture())

	lbl.Importance = LowImportance
	lbl.Refresh()
	test.AssertImageMatches(t, "label/label_importance_low.png", w.Canvas().Capture())

	lbl.Importance = MediumImportance
	lbl.Refresh()
	test.AssertImageMatches(t, "label/label_importance_medium.png", w.Canvas().Capture())

	lbl.Importance = HighImportance
	lbl.Refresh()
	test.AssertImageMatches(t, "label/label_importance_high.png", w.Canvas().Capture())

	lbl.Importance = DangerImportance
	lbl.Refresh()
	test.AssertImageMatches(t, "label/label_importance_danger.png", w.Canvas().Capture())

	lbl.Importance = WarningImportance
	lbl.Refresh()
	test.AssertImageMatches(t, "label/label_importance_warning.png", w.Canvas().Capture())

	lbl.Importance = SuccessImportance
	lbl.Refresh()
	test.AssertImageMatches(t, "label/label_importance_success.png", w.Canvas().Capture())
}

func TestLabelSizeNameWithSelection(t *testing.T) {
	l := NewLabel("Hello")
	l.Selectable = true

	w := test.NewWindow(l)
	defer w.Close()

	assert.Empty(t, l.SelectedText())
	assert.Equal(t, 2, len(test.WidgetRenderer(l).Objects()))

	sel := test.WidgetRenderer(l).Objects()[0].(*focusSelectable)
	sel.MouseDown(&desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(15, 10)}})
	sel.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 15, DY: 0},
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)}})
	sel.DragEnd()
	sel.MouseUp(&desktop.MouseEvent{Button: desktop.MouseButtonPrimary,
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(30, 10)}})
	assert.Equal(t, "el", l.SelectedText())

	test.AssertRendersToImage(t, "label/label_selection_defaultsize.png", w.Canvas())

	l.SizeName = theme.SizeNameHeadingText
	l.Refresh()
	w.SetContent(l) // updates window size

	test.AssertRendersToImage(t, "label/label_selection_headersize.png", w.Canvas())
}

func labelTextRenderTexts(p fyne.Widget) []*canvas.Text {
	rich := cache.Renderer(p).Objects()[0].(*RichText)
	return richTextRenderTexts(rich)
}
