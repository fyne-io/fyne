package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var (
	style      = fyne.TextStyle{}
	foo1Min    = fyne.MeasureText("foo1", theme.TextSize(), style)
	foo2Min    = fyne.MeasureText("foo2", theme.TextSize(), style)
	foo3Min    = fyne.MeasureText("foo3", theme.TextSize(), style)
	foobar1Min = fyne.MeasureText("foobar1", theme.TextSize(), style)
	foobar2Min = fyne.MeasureText("foobar2", theme.TextSize(), style)
	foobar3Min = fyne.MeasureText("foobar3", theme.TextSize(), style)
)

func addItem(t *testing.T, ac *AccordionContainer, header, detail string) {
	t.Helper()
	ac.Append(header, &canvas.Text{Text: detail, TextSize: theme.TextSize()})
}

func TestAccordionContainer_Empty(t *testing.T) {
	ac := NewAccordionContainer()
	assert.Equal(t, 0, len(ac.Items))
	min := ac.MinSize()
	assert.Equal(t, 0, min.Height)
	assert.Equal(t, 0, min.Width)
}

func TestAccordionContainer_Resize_Empty(t *testing.T) {
	ac := NewAccordionContainer()
	ac.Resize(fyne.NewSize(10, 10))
	size := ac.Size()
	assert.Equal(t, 10, size.Height)
	assert.Equal(t, 10, size.Width)
}

func TestAccordionContainer_Open(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo", "foobar")
	ac.Open(0)
	assert.True(t, ac.Items[0].open)
}

func TestAccordionContainer_OpenAll(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo1", "foobar1")
	addItem(t, ac, "foo2", "foobar2")
	addItem(t, ac, "foo3", "foobar3")
	ac.OpenAll()
	assert.True(t, ac.Items[0].open)
	assert.True(t, ac.Items[1].open)
	assert.True(t, ac.Items[2].open)
}

func TestAccordionContainer_Close(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo", "foobar")
	ac.Close(0)
	assert.False(t, ac.Items[0].open)
}

func TestAccordionContainer_CloseAll(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo1", "foobar1")
	addItem(t, ac, "foo2", "foobar2")
	addItem(t, ac, "foo3", "foobar3")
	ac.CloseAll()
	assert.False(t, ac.Items[0].open)
	assert.False(t, ac.Items[1].open)
	assert.False(t, ac.Items[2].open)
}

func TestAccordionContainerRenderer_Layout(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo1", "foobar1")
	addItem(t, ac, "foo2", "foobar2")
	addItem(t, ac, "foo3", "foobar3")

	rac := test.WidgetRenderer(ac).(*accordionContainerRenderer)

	checkSizeAndPosition := func(t *testing.T, min fyne.Size) {
		t.Helper()
		item0Pos := ac.Items[0].Position()
		item1Pos := ac.Items[1].Position()
		item2Pos := ac.Items[2].Position()
		item0Size := ac.Items[0].Size()
		item1Size := ac.Items[1].Size()
		item2Size := ac.Items[2].Size()
		assert.Equal(t, theme.Padding(), item0Pos.X)
		assert.Equal(t, theme.Padding(), item0Pos.Y)
		assert.Equal(t, min.Width-theme.Padding()*2, item0Size.Width)
		assert.Equal(t, ac.Items[0].MinSize().Height, item0Size.Height)
		assert.Equal(t, theme.Padding(), item1Pos.X)
		assert.Equal(t, theme.Padding()+ac.Items[0].MinSize().Height, item1Pos.Y)
		assert.Equal(t, min.Width-theme.Padding()*2, item1Size.Width)
		assert.Equal(t, ac.Items[1].MinSize().Height, item1Size.Height)
		assert.Equal(t, theme.Padding(), item2Pos.X)
		assert.Equal(t, theme.Padding()+ac.Items[0].MinSize().Height+ac.Items[1].MinSize().Height, item2Pos.Y)
		assert.Equal(t, min.Width-theme.Padding()*2, item2Size.Width)
		assert.Equal(t, ac.Items[2].MinSize().Height, item2Size.Height)
	}

	t.Run("AllClosed", func(t *testing.T) {
		ac.CloseAll()
		min := ac.MinSize()
		rac.Layout(min)
		checkSizeAndPosition(t, min)
	})
	t.Run("OneClosed", func(t *testing.T) {
		ac.Open(0)
		ac.Open(1)
		ac.Close(2)
		min := ac.MinSize()
		rac.Layout(min)
		checkSizeAndPosition(t, min)
	})
	t.Run("OneOpened", func(t *testing.T) {
		ac.Open(0)
		ac.Close(1)
		ac.Close(2)
		min := ac.MinSize()
		rac.Layout(min)
		checkSizeAndPosition(t, min)
	})
	t.Run("AllOpened", func(t *testing.T) {
		ac.OpenAll()
		min := ac.MinSize()
		rac.Layout(min)
		checkSizeAndPosition(t, min)
	})
}

func TestAccordionContainerRenderer_MinSize(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo1", "foobar1")
	addItem(t, ac, "foo2", "foobar2")
	addItem(t, ac, "foo3", "foobar3")

	rac := test.WidgetRenderer(ac).(*accordionContainerRenderer)

	headerWidth := fyne.Max(foo1Min.Width, foo2Min.Width)
	headerWidth = fyne.Max(headerWidth, foo3Min.Width)
	headerWidth += theme.IconInlineSize() + theme.Padding()

	t.Run("AllClosed", func(t *testing.T) {
		ac.CloseAll()
		min := rac.MinSize()
		height := foo1Min.Height + foo2Min.Height + foo3Min.Height
		assert.Equal(t, headerWidth+theme.Padding()*2, min.Width)
		assert.Equal(t, height+theme.Padding()*2, min.Height)
	})
	t.Run("OneClosed", func(t *testing.T) {
		ac.Open(0)
		ac.Open(1)
		ac.Close(2)
		min := rac.MinSize()
		width := fyne.Max(headerWidth, foobar1Min.Width)
		width = fyne.Max(width, foobar2Min.Width)
		height := foo1Min.Height + foo2Min.Height + foo3Min.Height
		height += foobar1Min.Height + foobar2Min.Height
		assert.Equal(t, width+theme.Padding()*2, min.Width)
		assert.Equal(t, height+theme.Padding()*2, min.Height)
	})
	t.Run("OneOpened", func(t *testing.T) {
		ac.Open(0)
		ac.Close(1)
		ac.Close(2)
		min := rac.MinSize()
		width := fyne.Max(headerWidth, foobar1Min.Width)
		height := foo1Min.Height + foo2Min.Height + foo3Min.Height
		height += foobar1Min.Height
		assert.Equal(t, width+theme.Padding()*2, min.Width)
		assert.Equal(t, height+theme.Padding()*2, min.Height)
	})
	t.Run("AllOpened", func(t *testing.T) {
		ac.OpenAll()
		min := rac.MinSize()
		width := fyne.Max(headerWidth, foobar1Min.Width)
		width = fyne.Max(width, foobar2Min.Width)
		width = fyne.Max(width, foobar3Min.Width)
		height := foo1Min.Height + foo2Min.Height + foo3Min.Height
		height += foobar1Min.Height + foobar2Min.Height + foobar3Min.Height
		assert.Equal(t, width+theme.Padding()*2, min.Width)
		assert.Equal(t, height+theme.Padding()*2, min.Height)
	})
}

func Test_accordionItemHeaderRenderer_BackgroundColor(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo", "foobar")
	acih := ac.Items[0].Header
	acihr := test.WidgetRenderer(acih).(*accordionItemHeaderRenderer)
	acihr.createCanvasObjects()
	assert.Equal(t, theme.BackgroundColor(), acihr.BackgroundColor())
}

func Test_accordionItemHeaderRenderer_BackgroundColor_Hovered(t *testing.T) {
	ac := NewAccordionContainer()
	addItem(t, ac, "foo", "foobar")
	acih := ac.Items[0].Header
	acihr := test.WidgetRenderer(acih).(*accordionItemHeaderRenderer)
	acihr.createCanvasObjects()
	acih.MouseIn(&desktop.MouseEvent{})
	assert.Equal(t, theme.HoverColor(), acihr.BackgroundColor())
	acih.MouseOut()
	assert.Equal(t, theme.BackgroundColor(), acihr.BackgroundColor())
}
