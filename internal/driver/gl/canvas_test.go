// +build !ci

package gl

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_Content(t *testing.T) {
	content := &canvas.Circle{}
	w := d.CreateWindow("Test")
	w.SetContent(content)

	assert.Equal(t, content, w.Content())
}

func TestGlCanvas_NilContent(t *testing.T) {
	w := d.CreateWindow("Test")

	assert.NotNil(t, w.Content()) // never a nil canvas so we have a sensible fallback
}

func Test_glCanvas_SetContent(t *testing.T) {
	var menuHeight int
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = widget.NewToolbar(widget.NewToolbarAction(theme.ContentCutIcon(), func() {})).MinSize().Height
	}
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        int
		expectedMenuHeight int
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, 4, 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, 4, menuHeight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := d.CreateWindow("Test").(*window)
			w.Canvas().SetScale(1)
			w.SetPadded(tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			content := canvas.NewCircle(color.Black)
			w.SetContent(content)
			w.Resize(fyne.NewSize(100, 100))
			c := w.Canvas()
			canvasSize := 100

			// wait for canvas to get its size right
			for w.canvas.Size().Width != canvasSize {
				time.Sleep(time.Millisecond * 10)
			}

			newContent := canvas.NewCircle(color.White)
			assert.Equal(t, fyne.NewPos(0, 0), newContent.Position())
			assert.Equal(t, fyne.NewSize(0, 0), newContent.Size())
			c.SetContent(newContent)
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), newContent.Position())
			assert.Equal(t, fyne.NewSize(canvasSize-2*tt.expectedPad, canvasSize-2*tt.expectedPad-tt.expectedMenuHeight), newContent.Size())
		})
	}
}

func Test_glCanvas_GrowingMinSize(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	c := w.Canvas()
	c.SetScale(1)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftObj3 := canvas.NewRectangle(color.Black)
	leftObj3.SetMinSize(fyne.NewSize(50, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2, leftObj3)
	middleObj1 := canvas.NewRectangle(color.Black)
	middleObj1.SetMinSize(fyne.NewSize(50, 50))
	middleObj2 := canvas.NewRectangle(color.Black)
	middleObj2.SetMinSize(fyne.NewSize(50, 50))
	middleObj3 := canvas.NewRectangle(color.Black)
	middleObj3.SetMinSize(fyne.NewSize(50, 50))
	middleCol := widget.NewVBox(middleObj1, middleObj2, middleObj3)
	middleColScroll := widget.NewScrollContainer(middleCol)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightObj3 := canvas.NewRectangle(color.Black)
	rightObj3.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2, rightObj3)
	rightColScroll := widget.NewScrollContainer(rightCol)
	content := widget.NewHBox(leftCol, middleColScroll, rightColScroll)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(150+4*theme.Padding(), 150+4*theme.Padding())
	// wait for canvas to get its size right
	for c.Size() != oldCanvasSize {
		time.Sleep(time.Millisecond * 10)
	}

	// child size change affects ancestors up to root
	leftObj1.SetMinSize(fyne.NewSize(60, 60))
	canvas.Refresh(leftObj1)
	for c.Size() == oldCanvasSize {
		time.Sleep(time.Millisecond * 10)
	}

	expectedCanvasSize := oldCanvasSize.Add(fyne.NewSize(10, 10))
	assert.Equal(t, expectedCanvasSize, c.Size())

	// child size change affects ancestors up to scroll
	oldCanvasSize = c.Size()
	oldRightScrollSize := rightColScroll.Size()
	oldRightColSize := rightCol.Size()
	rightObj1.SetMinSize(fyne.NewSize(50, 100))
	canvas.Refresh(rightObj1)
	for rightCol.Size() == oldRightColSize {
		time.Sleep(time.Millisecond * 10)
	}

	assert.Equal(t, oldCanvasSize, c.Size())
	assert.Equal(t, oldRightScrollSize, rightColScroll.Size())
	expectedRightColSize := oldRightColSize.Add(fyne.NewSize(0, 40))
	assert.Equal(t, expectedRightColSize, rightCol.Size())

	// 1st child size change affects ancestors up to root
	// and is not hidden by 2nd child change which affects ancestors up to scroll
	oldMiddleColSize := middleCol.Size()
	oldMiddleScrollSize := middleColScroll.Size()
	oldRightColSize = rightCol.Size()
	middleObj2.SetMinSize(fyne.NewSize(50, 100))
	rightObj2.SetMinSize(fyne.NewSize(50, 200))
	canvas.Refresh(middleObj2)
	canvas.Refresh(rightObj2)
	for middleCol.Size() == oldMiddleColSize || rightCol.Size() == oldRightColSize {
		time.Sleep(time.Millisecond * 10)
	}

	assert.Equal(t, oldCanvasSize, c.Size())
	assert.Equal(t, oldMiddleScrollSize, middleColScroll.Size())
	assert.Equal(t, oldRightScrollSize, rightColScroll.Size())
	expectedMiddleColSize := oldMiddleColSize.Add(fyne.NewSize(0, 40))
	assert.Equal(t, expectedMiddleColSize, middleCol.Size())
	expectedRightColSize = oldRightColSize.Add(fyne.NewSize(0, 150))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
}
