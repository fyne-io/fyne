package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestAccordionContainer_Toggle(t *testing.T) {
	ai := NewAccordionItem("foo", NewLabel("foobar"))
	ac := NewAccordionContainer(ai)
	acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
	aih := acr.headers[0]
	assert.False(t, ai.Open)

	test.Tap(aih)
	assert.True(t, ai.Open)
	test.Tap(aih)
	assert.False(t, ai.Open)
}

func TestAccordionContainerRenderer_Layout(t *testing.T) {
	ai0 := NewAccordionItem("foo0", NewLabel("foobar0"))
	ai1 := NewAccordionItem("foo1", NewLabel("foobar1"))
	ai2 := NewAccordionItem("foo2", NewLabel("foobar2"))
	aid0 := ai0.Detail
	aid1 := ai1.Detail
	aid2 := ai2.Detail
	ac := NewAccordionContainer()
	ac.Append(ai0)
	ac.Append(ai1)
	ac.Append(ai2)

	acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
	aih0 := acr.headers[0]
	aih1 := acr.headers[1]
	aih2 := acr.headers[2]

	t.Run("All_Closed", func(t *testing.T) {
		ac.CloseAll()
		min := ac.MinSize()
		acr.Layout(min)
		assert.Equal(t, 0, aih0.Position().X)
		assert.Equal(t, 0, aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, 0, aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+theme.Padding(), aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, 0, aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+2*theme.Padding(), aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
	})
	t.Run("One_Closed", func(t *testing.T) {
		ac.MultiOpen = true
		ac.Close(0)
		ac.Open(1)
		ac.Open(2)
		min := ac.MinSize()
		acr.Layout(min)
		assert.Equal(t, 0, aih0.Position().X)
		assert.Equal(t, 0, aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, 0, aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+theme.Padding(), aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, 0, aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+3*theme.Padding(), aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
		assert.Equal(t, 0, aid1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+2*theme.Padding(), aid1.Position().Y)
		assert.Equal(t, min.Width, aid1.Size().Width)
		assert.Equal(t, aid1.MinSize().Height, aid1.Size().Height)
		assert.Equal(t, 0, aid2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+aih2.MinSize().Height+4*theme.Padding(), aid2.Position().Y)
		assert.Equal(t, min.Width, aid2.Size().Width)
		assert.Equal(t, aid2.MinSize().Height, aid2.Size().Height)
	})
	t.Run("One_Opened", func(t *testing.T) {
		ac.Close(0)
		ac.Close(1)
		ac.Open(2)
		min := ac.MinSize()
		acr.Layout(min)
		assert.Equal(t, 0, aih0.Position().X)
		assert.Equal(t, 0, aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, 0, aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+1*theme.Padding(), aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, 0, aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+2*theme.Padding(), aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
		assert.Equal(t, 0, aid2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aih1.MinSize().Height+aih2.MinSize().Height+3*theme.Padding(), aid2.Position().Y)
		assert.Equal(t, min.Width, aid2.Size().Width)
		assert.Equal(t, aid2.MinSize().Height, aid2.Size().Height)
	})
	t.Run("All_Opened", func(t *testing.T) {
		ac.MultiOpen = true
		ac.OpenAll()
		min := ac.MinSize()
		acr.Layout(min)
		assert.Equal(t, 0, aih0.Position().X)
		assert.Equal(t, 0, aih0.Position().Y)
		assert.Equal(t, min.Width, aih0.Size().Width)
		assert.Equal(t, aih0.MinSize().Height, aih0.Size().Height)
		assert.Equal(t, 0, aih1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+2*theme.Padding(), aih1.Position().Y)
		assert.Equal(t, min.Width, aih1.Size().Width)
		assert.Equal(t, aih1.MinSize().Height, aih1.Size().Height)
		assert.Equal(t, 0, aih2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+4*theme.Padding(), aih2.Position().Y)
		assert.Equal(t, min.Width, aih2.Size().Width)
		assert.Equal(t, aih2.MinSize().Height, aih2.Size().Height)
		assert.Equal(t, 0, aid0.Position().X)
		assert.Equal(t, aih0.MinSize().Height+theme.Padding(), aid0.Position().Y)
		assert.Equal(t, min.Width, aid0.Size().Width)
		assert.Equal(t, aid0.MinSize().Height, aid0.Size().Height)
		assert.Equal(t, 0, aid1.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+aih1.MinSize().Height+3*theme.Padding(), aid1.Position().Y)
		assert.Equal(t, min.Width, aid1.Size().Width)
		assert.Equal(t, aid1.MinSize().Height, aid1.Size().Height)
		assert.Equal(t, 0, aid2.Position().X)
		assert.Equal(t, aih0.MinSize().Height+aid0.MinSize().Height+aih1.MinSize().Height+aid1.MinSize().Height+aih2.MinSize().Height+5*theme.Padding(), aid2.Position().Y)
		assert.Equal(t, min.Width, aid2.Size().Width)
		assert.Equal(t, aid2.MinSize().Height, aid2.Size().Height)
	})
}

func TestAccordionContainerRenderer_MinSize(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		ac := NewAccordionContainer()
		acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
		min := acr.MinSize()
		assert.Equal(t, 0, min.Width)
		assert.Equal(t, 0, min.Height)
	})
	t.Run("Single", func(t *testing.T) {
		ai := NewAccordionItem("foo", NewLabel("foobar"))
		t.Run("Open", func(t *testing.T) {
			ac := NewAccordionContainer()
			ac.Append(ai)
			ac.Open(0)
			acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
			min := acr.MinSize()
			aih := acr.headers[0].MinSize()
			aid := ai.Detail.MinSize()
			assert.Equal(t, fyne.Max(aih.Width, aid.Width), min.Width)
			assert.Equal(t, aih.Height+aid.Height+theme.Padding(), min.Height)
		})
		t.Run("Closed", func(t *testing.T) {
			ac := NewAccordionContainer()
			ac.Append(ai)
			ac.Close(0)
			acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
			min := acr.MinSize()
			aih := acr.headers[0].MinSize()
			assert.Equal(t, aih.Width, min.Width)
			assert.Equal(t, aih.Height, min.Height)
		})
	})
	t.Run("Multiple", func(t *testing.T) {
		ai0 := NewAccordionItem("foo0", NewLabel("foobar0"))
		ai1 := NewAccordionItem("foo1", NewLabel("foobar1"))
		ai2 := NewAccordionItem("foo2", NewLabel("foobar2"))
		t.Run("One_Open", func(t *testing.T) {
			ac := NewAccordionContainer()
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.Open(0)
			ac.Close(1)
			ac.Close(2)
			acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
			min := acr.MinSize()
			aih0 := acr.headers[0].MinSize()
			aih1 := acr.headers[1].MinSize()
			aih2 := acr.headers[2].MinSize()
			aid0 := ai0.Detail.MinSize()
			width := fyne.Max(aih0.Width, aid0.Width)
			width = fyne.Max(width, aih1.Width)
			width = fyne.Max(width, aih2.Width)
			assert.Equal(t, width, min.Width)
			height := aih0.Height
			height += theme.Padding()
			height += aid0.Height
			height += theme.Padding()
			height += aih1.Height
			height += theme.Padding()
			height += aih2.Height
			assert.Equal(t, height, min.Height)
		})
		t.Run("All_Open", func(t *testing.T) {
			ac := &AccordionContainer{
				MultiOpen: true,
			}
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.OpenAll()
			acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
			min := acr.MinSize()
			aih0 := acr.headers[0].MinSize()
			aih1 := acr.headers[1].MinSize()
			aih2 := acr.headers[2].MinSize()
			aid0 := ai0.Detail.MinSize()
			aid1 := ai1.Detail.MinSize()
			aid2 := ai2.Detail.MinSize()
			width := fyne.Max(aih0.Width, aid0.Width)
			width = fyne.Max(width, fyne.Max(aih1.Width, aid1.Width))
			width = fyne.Max(width, fyne.Max(aih2.Width, aid2.Width))
			assert.Equal(t, width, min.Width)
			height := aih0.Height
			height += theme.Padding()
			height += aid0.Height
			height += theme.Padding()
			height += aih1.Height
			height += theme.Padding()
			height += aid1.Height
			height += theme.Padding()
			height += aih2.Height
			height += theme.Padding()
			height += aid2.Height
			assert.Equal(t, height, min.Height)
		})
		t.Run("One_Closed", func(t *testing.T) {
			ac := &AccordionContainer{
				MultiOpen: true,
			}
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.Open(0)
			ac.Open(1)
			ac.Close(2)
			acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
			min := acr.MinSize()
			aih0 := acr.headers[0].MinSize()
			aih1 := acr.headers[1].MinSize()
			aih2 := acr.headers[2].MinSize()
			aid0 := ai0.Detail.MinSize()
			aid1 := ai1.Detail.MinSize()
			width := fyne.Max(aih0.Width, aid0.Width)
			width = fyne.Max(width, fyne.Max(aih1.Width, aid1.Width))
			width = fyne.Max(width, aih2.Width)
			assert.Equal(t, width, min.Width)
			height := aih0.Height
			height += theme.Padding()
			height += aid0.Height
			height += theme.Padding()
			height += aih1.Height
			height += theme.Padding()
			height += aid1.Height
			height += theme.Padding()
			height += aih2.Height
			assert.Equal(t, height, min.Height)
		})
		t.Run("All_Closed", func(t *testing.T) {
			ac := NewAccordionContainer()
			ac.Append(ai0)
			ac.Append(ai1)
			ac.Append(ai2)
			ac.CloseAll()
			acr := test.WidgetRenderer(ac).(*accordionContainerRenderer)
			min := acr.MinSize()
			aih0 := acr.headers[0].MinSize()
			aih1 := acr.headers[1].MinSize()
			aih2 := acr.headers[2].MinSize()
			width := aih0.Width
			width = fyne.Max(width, aih1.Width)
			width = fyne.Max(width, aih2.Width)
			assert.Equal(t, width, min.Width)
			height := aih0.Height
			height += theme.Padding()
			height += aih1.Height
			height += theme.Padding()
			height += aih2.Height
			assert.Equal(t, height, min.Height)
		})
	})
}
