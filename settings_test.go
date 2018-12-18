package fyne

import (
	"image/color"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetTheme(t *testing.T) {
	light := Theme(&dummyTheme{})
	GlobalSettings().SetTheme(light)

	assert.Equal(t, light, GlobalSettings().Theme())
}

func TestListenerCallback(t *testing.T) {
	listener := make(chan Settings)
	GlobalSettings().AddChangeListener(listener)
	GlobalSettings().SetTheme(&dummyTheme{})

	func() {
		select {
		case <-listener:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for callback")
		}
	}()
}

type dummyTheme struct {
}

func (dummyTheme) BackgroundColor() color.Color {
	return color.White
}

func (dummyTheme) ButtonColor() color.Color {
	return color.Black
}

func (dummyTheme) TextColor() color.Color {
	return color.Black
}

func (dummyTheme) PrimaryColor() color.Color {
	return color.Black
}

func (dummyTheme) FocusColor() color.Color {
	return color.Black
}

func (dummyTheme) TextSize() int {
	return 1
}

func (dummyTheme) TextFont() Resource {
	return nil
}

func (dummyTheme) TextBoldFont() Resource {
	return nil
}

func (dummyTheme) TextItalicFont() Resource {
	return nil
}

func (dummyTheme) TextBoldItalicFont() Resource {
	return nil
}

func (dummyTheme) TextMonospaceFont() Resource {
	return nil
}

func (dummyTheme) Padding() int {
	return 1
}

func (dummyTheme) IconInlineSize() int {
	return 1
}
