package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestAccordion_Toggle(t *testing.T) {
	ai := NewAccordionItem("foo", NewLabel("foobar"))
	ac := NewAccordion(ai)
	ar := test.WidgetRenderer(ac).(*accordionRenderer)
	aih := ar.headers[0]
	assert.False(t, ai.Open)

	test.Tap(aih)
	assert.True(t, ai.Open)
	test.Tap(aih)
	assert.False(t, ai.Open)
}

func TestAccordionRenderer_Layout(t *testing.T) {
	ai0 := NewAccordionItem("foo0", NewLabel("foobar0"))
	ai1 := NewAccordionItem("foo1", NewLabel("foobar1"))
	ai2 := NewAccordionItem("foo2", NewLabel("foobar2"))
	aid0 := ai0.Detail
	aid1 := ai1.Detail
	aid2 := ai2.Detail
	ac := NewAccordion()
	ac.Append(ai0)
	ac.Append(ai1)
	ac.Append(ai2)

	ar := test.WidgetRenderer(ac).(*accordionRenderer)
	aih0 := ar.headers[0]
	aih1 := ar.headers[1]
	aih2 := ar.headers[2]

	t.Run("All_Closed", func(t *testing.T) {
		ac.CloseAll()
		min := ac.MinSize()
		ar.Layout(min)
		assert.Equal(t, float32(0), aih0.Position().X)
		assert.Equal(t, theme.Padding(), aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, float32(0), aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+theme.Padding()*3+1, aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, float32(0), aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+5*theme.Padding()+2, aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
	})
	t.Run("One_Closed", func(t *testing.T) {
		ac.MultiOpen = true
		ac.Close(0)
		ac.Open(1)
		ac.Open(2)
		min := ac.MinSize()
		ar.Layout(min)
		assert.Equal(t, float32(0), aih0.Position().X)
		assert.Equal(t, theme.Padding(), aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, float32(0), aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+3*theme.Padding()+1, aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, float32(0), aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+6*theme.Padding()+2, aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
		assert.Equal(t, float32(0), aid1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+4*theme.Padding()+1, aid1.Position().Y)
		assert.Equal(t, min.Width, aid1.Size().Width)
		assert.Equal(t, aid1.MinSize().Height, aid1.Size().Height)
		assert.Equal(t, float32(0), aid2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+aih2.MinSize().Height+7*theme.Padding()+2, aid2.Position().Y)
		assert.Equal(t, min.Width, aid2.Size().Width)
		assert.Equal(t, aid2.MinSize().Height, aid2.Size().Height)
	})
	t.Run("One_Opened", func(t *testing.T) {
		ac.Close(0)
		ac.Close(1)
		ac.Open(2)
		min := ac.MinSize()
		ar.Layout(min)
		assert.Equal(t, float32(0), aih0.Position().X)
		assert.Equal(t, theme.Padding(), aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, float32(0), aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+3*theme.Padding()+1, aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, float32(0), aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+5*theme.Padding()+2, aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
		assert.Equal(t, float32(0), aid2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+aih2.MinSize().Height+6*theme.Padding()+2, aid2.Position().Y)
		assert.Equal(t, min.Width, aid2.Size().Width)
		assert.Equal(t, aid2.MinSize().Height, aid2.Size().Height)
	})
	t.Run("All_Opened", func(t *testing.T) {
		ac.MultiOpen = true
		ac.OpenAll()
		min := ac.MinSize()
		ar.Layout(min)
		assert.Equal(t, float32(0), aih0.Position().X)
		assert.Equal(t, theme.Padding(), aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, float32(0), aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+4*theme.Padding()+1, aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, float32(0), aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+7*theme.Padding()+2, aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
		assert.Equal(t, float32(0), aid0.Position().X)
		assert.Equal(t, aih0.MinSize().Height+theme.Padding()*2, aid0.Position().Y)
		assert.Equal(t, min.Width, aid0.Size().Width)
		assert.Equal(t, aid0.MinSize().Height, aid0.Size().Height)
		assert.Equal(t, float32(0), aid1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+aih1.MinSize().Height+5*theme.Padding()+1, aid1.Position().Y)
		assert.Equal(t, min.Width, aid1.Size().Width)
		assert.Equal(t, aid1.MinSize().Height, aid1.Size().Height)
		assert.Equal(t, float32(0), aid2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+aih2.MinSize().Height+8*theme.Padding()+2, aid2.Position().Y)
		assert.Equal(t, min.Width, aid2.Size().Width)
		assert.Equal(t, aid2.MinSize().Height, aid2.Size().Height)
	})
}

func TestAccordionRenderer_MinSize(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		ac := NewAccordion()
		ar := test.WidgetRenderer(ac).(*accordionRenderer)
		min := ar.MinSize()
		assert.Equal(t, float32(0), min.Width)
		assert.Equal(t, float32(0), min.Height)
	})
	t.Run("Single", func(t *testing.T) {
		ai := NewAccordionItem("foo", NewLabel("foobar"))
		t.Run("Open", func(t *testing.T) {
			ac := NewAccordion()
			ac.Append(ai)
			ac.Open(0)
			ar := test.WidgetRenderer(ac).(*accordionRenderer)
			min := ar.MinSize()
			aih := ar.headers[0].MinSize()
			aid := ai.Detail.MinSize()
			assert.Equal(t, fyne.Max(aih.Width, aid.Width), min.Width)
			assert.Equal(t, aih.Height+aid.Height+theme.Padding()*3, min.Height)
		})
		t.Run("Closed", func(t *testing.T) {
			ac := NewAccordion()
			ac.Append(ai)
			ac.Close(0)
			ar := test.WidgetRenderer(ac).(*accordionRenderer)
			min := ar.MinSize()
			aih := ar.headers[0].MinSize()
			assert.Equal(t, aih.Width, min.Width)
			assert.Equal(t, aih.Height+theme.Padding()*2, min.Height)
		})
	})
	t.Run("Multiple", func(t *testing.T) {
		ai0 := NewAccordionItem("foo0", NewLabel("foobar0"))
		ai1 := NewAccordionItem("foo1", NewLabel("foobar1"))
		ai2 := NewAccordionItem("foo2", NewLabel("foobar2"))
		t.Run("One_Open", func(t *testing.T) {
			ac := NewAccordion()
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.Open(0)
			ac.Close(1)
			ac.Close(2)
			ar := test.WidgetRenderer(ac).(*accordionRenderer)
			min := ar.MinSize()
			aih0 := ar.headers[0].MinSize()
			aih1 := ar.headers[1].MinSize()
			aih2 := ar.headers[2].MinSize()
			aid0 := ai0.Detail.MinSize()
			width := fyne.Max(aih0.Width, aid0.Width)
			width = fyne.Max(width, aih1.Width)
			width = fyne.Max(width, aih2.Width)
			assert.Equal(t, width, min.Width)
			height := theme.Padding()
			height += aih0.Height
			height += theme.Padding()
			height += aid0.Height
			height += theme.Padding()*2 + 1
			height += aih1.Height
			height += theme.Padding()*2 + 1
			height += aih2.Height
			height += theme.Padding()
			assert.Equal(t, height, min.Height)
		})
		t.Run("All_Open", func(t *testing.T) {
			ac := &Accordion{
				MultiOpen: true,
			}
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.OpenAll()
			ar := test.WidgetRenderer(ac).(*accordionRenderer)
			min := ar.MinSize()
			aih0 := ar.headers[0].MinSize()
			aih1 := ar.headers[1].MinSize()
			aih2 := ar.headers[2].MinSize()
			aid0 := ai0.Detail.MinSize()
			aid1 := ai1.Detail.MinSize()
			aid2 := ai2.Detail.MinSize()
			width := fyne.Max(aih0.Width, aid0.Width)
			width = fyne.Max(width, fyne.Max(aih1.Width, aid1.Width))
			width = fyne.Max(width, fyne.Max(aih2.Width, aid2.Width))
			assert.Equal(t, width, min.Width)
			height := theme.Padding()
			height += aih0.Height
			height += theme.Padding()
			height += aid0.Height
			height += theme.Padding()*2 + 1
			height += aih1.Height
			height += theme.Padding()
			height += aid1.Height
			height += theme.Padding()*2 + 1
			height += aih2.Height
			height += theme.Padding()
			height += aid2.Height
			height += theme.Padding()
			assert.Equal(t, height, min.Height)
		})
		t.Run("One_Closed", func(t *testing.T) {
			ac := &Accordion{
				MultiOpen: true,
			}
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.Open(0)
			ac.Open(1)
			ac.Close(2)
			ar := test.WidgetRenderer(ac).(*accordionRenderer)
			min := ar.MinSize()
			aih0 := ar.headers[0].MinSize()
			aih1 := ar.headers[1].MinSize()
			aih2 := ar.headers[2].MinSize()
			aid0 := ai0.Detail.MinSize()
			aid1 := ai1.Detail.MinSize()
			width := fyne.Max(aih0.Width, aid0.Width)
			width = fyne.Max(width, fyne.Max(aih1.Width, aid1.Width))
			width = fyne.Max(width, aih2.Width)
			assert.Equal(t, width, min.Width)
			height := theme.Padding()
			height += aih0.Height
			height += theme.Padding()
			height += aid0.Height
			height += theme.Padding()*2 + 1
			height += aih1.Height
			height += theme.Padding()
			height += aid1.Height
			height += theme.Padding()*2 + 1
			height += aih2.Height
			height += theme.Padding()
			assert.Equal(t, height, min.Height)
		})
		t.Run("All_Closed", func(t *testing.T) {
			ac := NewAccordion()
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.CloseAll()
			ar := test.WidgetRenderer(ac).(*accordionRenderer)
			min := ar.MinSize()
			aih0 := ar.headers[0].MinSize()
			aih1 := ar.headers[1].MinSize()
			aih2 := ar.headers[2].MinSize()
			width := aih0.Width
			width = fyne.Max(width, aih1.Width)
			width = fyne.Max(width, aih2.Width)
			assert.Equal(t, width, min.Width)
			height := theme.Padding()
			height += aih0.Height
			height += theme.Padding()*2 + 1
			height += aih1.Height
			height += theme.Padding()*2 + 1
			height += aih2.Height
			height += theme.Padding()
			assert.Equal(t, height, min.Height)
		})
	})
}

func TestAccordionRenderer_AddRemove(t *testing.T) {
	ac := NewAccordion()
	ar := test.WidgetRenderer(ac).(*accordionRenderer)
	ac.Append(NewAccordionItem("foo0", NewLabel("foobar0")))
	ac.Append(NewAccordionItem("foo1", NewLabel("foobar1")))
	ac.Append(NewAccordionItem("foo2", NewLabel("foobar2")))

	assert.Equal(t, 3, len(ac.Items))
	assert.Equal(t, 3, len(ar.headers))
	assert.Equal(t, 2, len(ar.dividers))
	assert.True(t, ar.headers[2].Visible())
	assert.True(t, ar.dividers[1].Visible())

	ac.RemoveIndex(2)
	assert.Equal(t, 2, len(ac.Items))
	assert.False(t, ar.headers[2].Visible())
	assert.False(t, ar.dividers[1].Visible())

	ac.Append(NewAccordionItem("foo3", NewLabel("foobar3")))
	assert.Equal(t, 3, len(ac.Items))
	assert.True(t, ar.headers[2].Visible())
	assert.True(t, ar.dividers[1].Visible())
}
