package gomobile

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	canv "fyne.io/fyne/canvas"
	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestCanvas_Resize(t *testing.T) {
	c := &canvas{content: canv.NewRectangle(color.White), padded: true, scale: 2}
	screenSize := fyne.NewSize(480, 640)
	c.Resize(screenSize)

	theme := fyne.CurrentApp().Settings().Theme()
	assert.Equal(t, screenSize.Width-theme.Padding()*2, c.Content().Size().Width)
	assert.Greater(t, screenSize.Height-theme.Padding()*2, c.Content().Size().Height) // a status bar...

	assert.Equal(t, theme.Padding(), c.Content().Position().X)
	assert.Less(t, theme.Padding(), c.Content().Position().Y)
}
