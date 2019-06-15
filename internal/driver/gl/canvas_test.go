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
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	var menuHeight int
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = widget.NewToolbar(widget.NewToolbarAction(theme.ContentCutIcon(), func() {})).MinSize().Height
	}
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        int
		expectedMenuHeight int
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, theme.Padding(), 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, theme.Padding(), menuHeight},
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
			canvasSize := 100
			w.SetContent(content)
			w.Resize(fyne.NewSize(canvasSize, canvasSize))

			newContent := canvas.NewCircle(color.White)
			assert.Equal(t, fyne.NewPos(0, 0), newContent.Position())
			assert.Equal(t, fyne.NewSize(0, 0), newContent.Size())
			w.SetContent(newContent)
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), newContent.Position())
			assert.Equal(t, fyne.NewSize(canvasSize-2*tt.expectedPad, canvasSize-2*tt.expectedPad-tt.expectedMenuHeight), newContent.Size())
		})
	}
}

func Test_glCanvas_ChildMinSizeChangeAffectsAncestorsUpToRoot(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	content := widget.NewHBox(leftCol, rightCol)
	w.SetContent(content)
	w.ignoreResize = true
	for c.isDirty() {
		time.Sleep(time.Millisecond * 10)
	}

	oldCanvasSize := fyne.NewSize(100+3*theme.Padding(), 100+3*theme.Padding())
	assert.Equal(t, oldCanvasSize, c.Size())

	leftObj1.SetMinSize(fyne.NewSize(60, 60))
	c.Refresh(leftObj1)
	for c.Size() == oldCanvasSize {
		time.Sleep(time.Millisecond * 10)
	}

	expectedCanvasSize := oldCanvasSize.Add(fyne.NewSize(10, 10))
	assert.Equal(t, expectedCanvasSize, c.Size())
	w.ignoreResize = false
}

func Test_glCanvas_ChildMinSizeChangeAffectsAncestorsUpToScroll(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	c.SetScale(1)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	rightColScroll := widget.NewScrollContainer(rightCol)
	content := widget.NewHBox(leftCol, rightColScroll)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(100+3*theme.Padding(), 100+3*theme.Padding())
	w.Resize(oldCanvasSize)
	w.ignoreResize = true // for some reason the window manager is intercepting and setting strange values in tests
	for c.isDirty() {
		time.Sleep(time.Millisecond * 10)
	}

	// child size change affects ancestors up to scroll
	oldCanvasSize = c.Size()
	oldRightScrollSize := rightColScroll.Size()
	oldRightColSize := rightCol.Size()
	rightObj1.SetMinSize(fyne.NewSize(50, 100))
	c.Refresh(rightObj1)
	for rightCol.Size() == oldRightColSize {
		time.Sleep(time.Millisecond * 10)
	}

	assert.Equal(t, oldCanvasSize, c.Size())
	assert.Equal(t, oldRightScrollSize, rightColScroll.Size())
	expectedRightColSize := oldRightColSize.Add(fyne.NewSize(0, 50))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
	w.ignoreResize = false
}

func Test_glCanvas_ChildMinSizeChangesInDifferentScrollAffectAncestorsUpToScroll(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	c.SetScale(1)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	leftColScroll := widget.NewScrollContainer(leftCol)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	rightColScroll := widget.NewScrollContainer(rightCol)
	content := widget.NewHBox(leftColScroll, rightColScroll)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(
		2*leftColScroll.MinSize().Width+3*theme.Padding(),
		leftColScroll.MinSize().Height+2*theme.Padding(),
	)
	w.Resize(oldCanvasSize)
	w.ignoreResize = true // for some reason the window manager is intercepting and setting strange values in tests
	for c.isDirty() {
		time.Sleep(time.Millisecond * 10)
	}

	oldLeftColSize := leftCol.Size()
	oldLeftScrollSize := leftColScroll.Size()
	oldRightColSize := rightCol.Size()
	oldRightScrollSize := rightColScroll.Size()
	leftObj2.SetMinSize(fyne.NewSize(50, 100))
	rightObj2.SetMinSize(fyne.NewSize(50, 200))
	c.Refresh(leftObj2)
	c.Refresh(rightObj2)
	for leftCol.Size() == oldLeftColSize || rightCol.Size() == oldRightColSize {
		time.Sleep(time.Millisecond * 10)
	}

	assert.Equal(t, oldCanvasSize, c.Size())
	assert.Equal(t, oldLeftScrollSize, leftColScroll.Size())
	assert.Equal(t, oldRightScrollSize, rightColScroll.Size())
	expectedLeftColSize := oldLeftColSize.Add(fyne.NewSize(0, 50))
	assert.Equal(t, expectedLeftColSize, leftCol.Size())
	expectedRightColSize := oldRightColSize.Add(fyne.NewSize(0, 150))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
	w.ignoreResize = false
}

func Test_glCanvas_MinSizeShrinkTriggersLayout(t *testing.T) {
	w := d.CreateWindow("Test").(*window)
	w.ignoreResize = true // for some reason the test is causing a WM resize event
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	content := widget.NewHBox(leftCol, rightCol)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(100+3*theme.Padding(), 100+3*theme.Padding())
	assert.Equal(t, oldCanvasSize, c.Size())
	w.ignoreResize = true // for some reason the window manager is intercepting and setting strange values in tests
	for c.isDirty() {
		time.Sleep(time.Millisecond * 10)
	}

	oldLeftObj1Size := leftObj1.Size()
	oldRightObj1Size := rightObj1.Size()
	oldRightObj2Size := rightObj2.Size()
	oldRightColSize := rightCol.Size()
	leftObj1.SetMinSize(fyne.NewSize(40, 40))
	rightObj1.SetMinSize(fyne.NewSize(30, 30))
	rightObj2.SetMinSize(fyne.NewSize(30, 20))
	c.Refresh(leftObj1)
	c.Refresh(rightObj1)
	c.Refresh(rightObj2)
	ch := make(chan bool, 1)
	go func() {
		for rightCol.Size() == oldRightColSize || leftObj1.Size() == oldLeftObj1Size ||
			rightObj1.Size() == oldRightObj1Size || rightObj2.Size() == oldRightObj2Size {
			time.Sleep(time.Millisecond * 10)
		}
		ch <- true
	}()
	select {
	case _ = <-ch:
		// all elements resized
	case <-time.After(3 * time.Second):
		t.Error("waiting for obj size change timed out")
	}

	assert.Equal(t, oldCanvasSize, c.Size())
	expectedRightColSize := oldRightColSize.Subtract(fyne.NewSize(20, 0))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
	assert.Equal(t, fyne.NewSize(50, 40), leftObj1.Size())
	assert.Equal(t, fyne.NewSize(30, 30), rightObj1.Size())
	assert.Equal(t, fyne.NewSize(30, 20), rightObj2.Size())
	w.ignoreResize = false
}
