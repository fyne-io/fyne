//go:build !ci && !mobile
// +build !ci,!mobile

package glfw

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_ChildMinSizeChangeAffectsAncestorsUpToRoot(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(100, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(100, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(100, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(100, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	content := container.NewHBox(leftCol, rightCol)
	w.SetContent(content)
	repaintWindow(w)

	oldCanvasSize := fyne.NewSize(200+3*theme.Padding(), 100+3*theme.Padding())
	assert.Equal(t, oldCanvasSize, c.Size())

	leftObj1.SetMinSize(fyne.NewSize(110, 60))
	c.Refresh(leftObj1)
	repaintWindow(w)

	expectedCanvasSize := oldCanvasSize.Add(fyne.NewSize(10, 10))
	assert.Equal(t, expectedCanvasSize, c.Size())
}

func TestGlCanvas_ChildMinSizeChangeAffectsAncestorsUpToScroll(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	rightColScroll := container.NewScroll(rightCol)
	content := container.NewHBox(leftCol, rightColScroll)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(200+3*theme.Padding(), 100+3*theme.Padding())
	w.Resize(oldCanvasSize)
	repaintWindow(w)

	// child size change affects ancestors up to scroll
	oldCanvasSize = c.Size()
	oldRightScrollSize := rightColScroll.Size()
	oldRightColSize := rightCol.Size()
	rightObj1.SetMinSize(fyne.NewSize(50, 100))
	c.Refresh(rightObj1)
	repaintWindow(w)

	assert.Equal(t, oldCanvasSize, c.Size())
	assert.Equal(t, oldRightScrollSize, rightColScroll.Size())
	expectedRightColSize := oldRightColSize.Add(fyne.NewSize(0, 50))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
}

func TestGlCanvas_ChildMinSizeChangesInDifferentScrollAffectAncestorsUpToScroll(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	leftColScroll := container.NewScroll(leftCol)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	rightColScroll := container.NewScroll(rightCol)
	content := container.NewHBox(leftColScroll, rightColScroll)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(
		2*leftColScroll.MinSize().Width+3*theme.Padding(),
		leftColScroll.MinSize().Height+2*theme.Padding(),
	)
	w.Resize(oldCanvasSize)
	repaintWindow(w)

	oldLeftColSize := leftCol.Size()
	oldLeftScrollSize := leftColScroll.Size()
	oldRightColSize := rightCol.Size()
	oldRightScrollSize := rightColScroll.Size()
	leftObj2.SetMinSize(fyne.NewSize(50, 100))
	rightObj2.SetMinSize(fyne.NewSize(50, 200))
	c.Refresh(leftObj2)
	c.Refresh(rightObj2)
	repaintWindow(w)

	assert.Equal(t, oldCanvasSize, c.Size())
	assert.Equal(t, oldLeftScrollSize, leftColScroll.Size())
	assert.Equal(t, oldRightScrollSize, rightColScroll.Size())
	expectedLeftColSize := oldLeftColSize.Add(fyne.NewSize(0, 50))
	assert.Equal(t, expectedLeftColSize, leftCol.Size())
	expectedRightColSize := oldRightColSize.Add(fyne.NewSize(0, 150))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
}

func TestGlCanvas_Content(t *testing.T) {
	content := &canvas.Circle{}
	w := createWindow("Test")
	w.SetContent(content)

	assert.Equal(t, content, w.Content())
}

func TestGlCanvas_ContentChangeWithoutMinSizeChangeDoesNotLayout(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(50, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(50, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	content := container.NewWithoutLayout(leftCol, rightCol)
	layout := &recordingLayout{}
	content.Layout = layout
	w.SetContent(content)

	repaintWindow(w)
	// clear the recorded layouts
	for layout.popLayoutEvent() != nil {
	}
	assert.Nil(t, layout.popLayoutEvent())

	leftObj1.FillColor = color.White
	rightObj1.FillColor = color.White
	rightObj2.FillColor = color.White
	c.Refresh(leftObj1)
	c.Refresh(rightObj1)
	c.Refresh(rightObj2)

	assert.Nil(t, layout.popLayoutEvent())
}

func TestGlCanvas_Focus(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	c := w.Canvas().(*glCanvas)

	ce := &focusable{id: "ce1"}
	content := container.NewVBox(ce)
	me := &focusable{id: "o2e1"}
	menuOverlay := container.NewVBox(me)
	o1e := &focusable{id: "o1e1"}
	overlay1 := container.NewVBox(o1e)
	o2e := &focusable{id: "o2e1"}
	overlay2 := container.NewVBox(o2e)
	w.SetContent(content)
	c.setMenuOverlay(menuOverlay)
	c.Overlays().Add(overlay1)
	c.Overlays().Add(overlay2)

	c.Focus(ce)
	assert.True(t, ce.focused, "focuses content object even if content is not in focus")

	c.Focus(me)
	assert.True(t, me.focused, "focuses menu object even if menu is not in focus")
	assert.True(t, ce.focused, "does not affect focus on other layer")

	c.Focus(o1e)
	assert.True(t, o1e.focused, "focuses overlay object even if menu is not in focus")
	assert.True(t, me.focused, "does not affect focus on other layer")

	c.Focus(o2e)
	assert.True(t, o2e.focused)
	assert.True(t, o1e.focused, "does not affect focus on other layer")

	foreign := &focusable{id: "o2e1"}
	c.Focus(foreign)
	assert.False(t, foreign.focused, "does not focus foreign object")
	assert.True(t, o2e.focused)
}

func TestGlCanvas_Focus_BeforeVisible(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	e := widget.NewEntry()
	c := w.Canvas().(*glCanvas)
	c.Focus(e) // this crashed in the past
}

func TestGlCanvas_Focus_SetContent(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	e := widget.NewEntry()
	w.SetContent(container.NewHBox(e))
	c := w.Canvas().(*glCanvas)
	c.Focus(e)
	assert.Equal(t, e, c.Focused())

	w.SetContent(container.NewVBox(e))
	assert.Equal(t, e, c.Focused())
}

func TestGlCanvas_FocusHandlingWhenAddingAndRemovingOverlays(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	c := w.Canvas().(*glCanvas)

	ce1 := &focusable{id: "ce1"}
	ce2 := &focusable{id: "ce2"}
	content := container.NewVBox(ce1, ce2)
	o1e1 := &focusable{id: "o1e1"}
	o1e2 := &focusable{id: "o1e2"}
	overlay1 := container.NewVBox(o1e1, o1e2)
	o2e1 := &focusable{id: "o2e1"}
	o2e2 := &focusable{id: "o2e2"}
	overlay2 := container.NewVBox(o2e1, o2e2)
	w.SetContent(content)

	assert.Nil(t, c.Focused())

	c.FocusPrevious()
	assert.Equal(t, ce2, c.Focused())
	assert.True(t, ce2.focused)

	c.Overlays().Add(overlay1)
	ctxt := "adding overlay changes focus handler but does not remove focus from content"
	assert.Nil(t, c.Focused(), ctxt)
	assert.True(t, ce2.focused, ctxt)

	c.FocusNext()
	ctxt = "changing focus affects overlay instead of content"
	assert.Equal(t, o1e1, c.Focused(), ctxt)
	assert.False(t, ce1.focused, ctxt)
	assert.True(t, ce2.focused, ctxt)
	assert.True(t, o1e1.focused, ctxt)

	c.Overlays().Add(overlay2)
	ctxt = "adding overlay changes focus handler but does not remove focus from previous overlay"
	assert.Nil(t, c.Focused(), ctxt)
	assert.True(t, o1e1.focused, ctxt)

	c.FocusPrevious()
	ctxt = "changing focus affects top overlay only"
	assert.Equal(t, o2e2, c.Focused(), ctxt)
	assert.True(t, o1e1.focused, ctxt)
	assert.False(t, o1e2.focused, ctxt)
	assert.True(t, o2e2.focused, ctxt)

	c.FocusNext()
	assert.Equal(t, o2e1, c.Focused())
	assert.False(t, o2e2.focused)
	assert.True(t, o2e1.focused)

	c.Overlays().Remove(overlay2)
	ctxt = "removing overlay restores focus handler from previous overlay but does not remove focus from removed overlay"
	assert.Equal(t, o1e1, c.Focused(), ctxt)
	assert.True(t, o2e1.focused, ctxt)
	assert.False(t, o2e2.focused, ctxt)
	assert.True(t, o1e1.focused, ctxt)

	c.FocusPrevious()
	assert.Equal(t, o1e2, c.Focused())
	assert.False(t, o1e1.focused)
	assert.True(t, o1e2.focused)

	c.Overlays().Remove(overlay1)
	ctxt = "removing last overlay restores focus handler from content but does not remove focus from removed overlay"
	assert.Equal(t, ce2, c.Focused(), ctxt)
	assert.False(t, o1e1.focused, ctxt)
	assert.True(t, o1e2.focused, ctxt)
	assert.True(t, ce2.focused, ctxt)
}

func TestGlCanvas_InsufficientSizeDoesntTriggerResizeIfSizeIsAlreadyMaxedOut(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	w.Resize(fyne.NewSize(200, 100))
	popUpContent := canvas.NewRectangle(color.Black)
	popUpContent.SetMinSize(fyne.NewSize(1000, 10))
	popUp := widget.NewPopUp(popUpContent, c)

	// This is because of a bug in PopUp size handling that will be fixed later.
	// This line will vanish then.
	popUp.Resize(popUpContent.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))

	assert.Equal(t, fyne.NewSize(1000, 10), popUpContent.Size())
	assert.Equal(t, fyne.NewSize(1000, 10).Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)), popUp.MinSize())
	assert.Equal(t, fyne.NewSize(200, 100), popUp.Size())

	repaintWindow(w)

	assert.Equal(t, fyne.NewSize(1000, 10), popUpContent.Size())
	assert.Equal(t, fyne.NewSize(200, 100), popUp.Size())
}

func TestGlCanvas_MinSizeShrinkTriggersLayout(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(100, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(100, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(100, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(100, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	content := container.NewHBox(leftCol, rightCol)
	w.SetContent(content)

	oldCanvasSize := fyne.NewSize(200+3*theme.Padding(), 100+3*theme.Padding())
	assert.Equal(t, oldCanvasSize, c.Size())
	repaintWindow(w)

	oldRightColSize := rightCol.Size()
	leftObj1.SetMinSize(fyne.NewSize(90, 40))
	rightObj1.SetMinSize(fyne.NewSize(80, 30))
	rightObj2.SetMinSize(fyne.NewSize(80, 20))
	c.Refresh(leftObj1)
	c.Refresh(rightObj1)
	c.Refresh(rightObj2)
	repaintWindow(w)

	assert.Equal(t, oldCanvasSize, c.Size())
	expectedRightColSize := oldRightColSize.Subtract(fyne.NewSize(20, 0))
	assert.Equal(t, expectedRightColSize, rightCol.Size())
	assert.Equal(t, fyne.NewSize(100, 40), leftObj1.Size())
	assert.Equal(t, fyne.NewSize(80, 30), rightObj1.Size())
	assert.Equal(t, fyne.NewSize(80, 20), rightObj2.Size())
}

func TestGlCanvas_NilContent(t *testing.T) {
	w := createWindow("Test")

	assert.NotNil(t, w.Content()) // never a nil canvas so we have a sensible fallback
}

func TestGlCanvas_PixelCoordinateAtPosition(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)

	pos := fyne.NewPos(4, 4)
	c.Lock()
	c.scale = 2.5
	c.Unlock()
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, int(10*c.texScale), x)
	assert.Equal(t, int(10*c.texScale), y)

	c.Lock()
	c.texScale = 2.0
	c.Unlock()
	x, y = c.PixelCoordinateForPosition(pos)
	assert.Equal(t, 20, x)
	assert.Equal(t, 20, y)
}

func TestGlCanvas_Resize(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	content := widget.NewLabel("Content")
	w.SetContent(content)

	size := fyne.NewSize(200, 100)
	assert.NotEqual(t, size, content.Size())

	w.Resize(size)
	assert.Equal(t, size, content.Size())
}

// TODO: this can be removed when #707 is addressed
func TestGlCanvas_ResizeWithOtherOverlay(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	content := widget.NewLabel("Content")
	over := widget.NewLabel("Over")
	w.SetContent(content)
	w.Canvas().Overlays().Add(over)
	// TODO: address #707; overlays should always be canvas size
	over.Resize(w.Canvas().Size())

	size := fyne.NewSize(200, 100)
	assert.NotEqual(t, size, content.Size())
	assert.NotEqual(t, size, over.Size())

	w.Resize(size)
	assert.Equal(t, size, content.Size(), "canvas content is resized")
	assert.Equal(t, size, over.Size(), "canvas overlay is resized")
}

func TestGlCanvas_ResizeWithOverlays(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	content := widget.NewLabel("Content")
	o1 := widget.NewLabel("o1")
	o2 := widget.NewLabel("o2")
	o3 := widget.NewLabel("o3")
	w.SetContent(content)
	w.Canvas().Overlays().Add(o1)
	w.Canvas().Overlays().Add(o2)
	w.Canvas().Overlays().Add(o3)

	size := fyne.NewSize(200, 100)
	assert.NotEqual(t, size, content.Size())
	assert.NotEqual(t, size, o1.Size())
	assert.NotEqual(t, size, o2.Size())
	assert.NotEqual(t, size, o3.Size())

	w.Resize(size)
	assert.Equal(t, size, content.Size(), "canvas content is resized")
	assert.Equal(t, size, o1.Size(), "canvas overlay 1 is resized")
	assert.Equal(t, size, o2.Size(), "canvas overlay 2 is resized")
	assert.Equal(t, size, o3.Size(), "canvas overlay 3 is resized")
}

// TODO: this can be removed when #707 is addressed
func TestGlCanvas_ResizeWithPopUpOverlay(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	content := widget.NewLabel("Content")
	over := widget.NewPopUp(widget.NewLabel("Over"), w.Canvas())
	w.SetContent(content)
	over.Show()

	size := fyne.NewSize(200, 100)
	overContentSize := over.Content.Size()
	assert.NotZero(t, overContentSize)
	assert.NotEqual(t, size, content.Size())
	assert.NotEqual(t, size, over.Size())
	assert.NotEqual(t, size, overContentSize)

	w.Resize(size)
	assert.Equal(t, size, content.Size(), "canvas content is resized")
	assert.Equal(t, size, over.Size(), "canvas overlay is resized")
	assert.Equal(t, overContentSize, over.Content.Size(), "canvas overlay content is _not_ resized")
}

func TestGlCanvas_ResizeWithModalPopUpOverlay(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)

	content := widget.NewLabel("Content")
	w.SetContent(content)

	popup := widget.NewModalPopUp(widget.NewLabel("PopUp"), w.Canvas())
	popupBgSize := fyne.NewSize(975, 575)
	popup.Show()
	popup.Resize(popupBgSize)

	winSize := fyne.NewSize(1000, 600)
	w.Resize(winSize)

	// get popup content padding dynamically
	popupContentPadding := popup.MinSize().Subtract(popup.Content.MinSize())

	assert.Equal(t, popupBgSize.Subtract(popupContentPadding), popup.Content.Size())
	assert.Equal(t, winSize, popup.Size())
}

func TestGlCanvas_Scale(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)

	c.Lock()
	c.scale = 2.5
	c.Unlock()
	assert.Equal(t, 5, int(2*c.Scale()))
}

func TestGlCanvas_SetContent(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	var menuHeight float32
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = NewMenuBar(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Empty", func() {}))), nil).MinSize().Height
	}
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        float32
		expectedMenuHeight float32
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, theme.Padding(), 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, theme.Padding(), menuHeight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createWindow("Test").(*window)
			w.SetPadded(tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			content := canvas.NewCircle(color.Black)
			canvasSize := float32(200)
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

var _ fyne.Layout = (*recordingLayout)(nil)

type recordingLayout struct {
	layoutEvents []interface{}
}

func (l *recordingLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.layoutEvents = append(l.layoutEvents, size)
}

func (l *recordingLayout) MinSize([]fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(6, 9)
}

func (l *recordingLayout) popLayoutEvent() (e interface{}) {
	e, l.layoutEvents = pop(l.layoutEvents)
	return
}
