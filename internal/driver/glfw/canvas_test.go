// +build !ci
// +build !mobile

package glfw

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_ChildMinSizeChangeAffectsAncestorsUpToRoot(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(100, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(100, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(100, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(100, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	content := widget.NewHBox(leftCol, rightCol)
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
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	rightColScroll := widget.NewScrollContainer(rightCol)
	content := widget.NewHBox(leftCol, rightColScroll)
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
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(50, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(50, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	content := fyne.NewContainer(leftCol, rightCol)
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
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(100, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(100, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	content := widget.NewHBox(leftCol, rightCol)
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
	c.scale = 2.5
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, int(10*c.texScale), x)
	assert.Equal(t, int(10*c.texScale), y)

	c.texScale = 2.0
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
	w.Canvas().SetOverlay(over)
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
	// TODO: address #707; overlays should always be canvas size
	o1.Resize(w.Canvas().Size())
	w.Canvas().Overlays().Add(o2)
	// TODO: address #707; overlays should always be canvas size
	o2.Resize(w.Canvas().Size())
	w.Canvas().Overlays().Add(o3)
	// TODO: address #707; overlays should always be canvas size
	o3.Resize(w.Canvas().Size())

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
	w.Canvas().Overlays().Add(over)

	size := fyne.NewSize(200, 100)
	overContentSize := over.Content.Size()
	assert.NotEqual(t, size, content.Size())
	assert.NotEqual(t, size, over.Size())
	assert.NotEqual(t, size, overContentSize)

	w.Resize(size)
	assert.Equal(t, size, content.Size(), "canvas content is resized")
	assert.Equal(t, size, over.Size(), "canvas overlay is resized")
	assert.Equal(t, overContentSize, over.Content.Size(), "canvas overlay content is _not_ resized")
}

func TestGlCanvas_Scale(t *testing.T) {
	w := createWindow("Test").(*window)
	c := w.Canvas().(*glCanvas)

	c.scale = 2.5
	assert.Equal(t, 5, int(2*c.Scale()))
}

func TestGlCanvas_SetContent(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	var menuHeight int
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = NewMenuBar(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Empty", func() {}))), nil).MinSize().Height
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
			w := createWindow("Test").(*window)
			w.SetPadded(tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			content := canvas.NewCircle(color.Black)
			canvasSize := 200
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

func TestGlCanvas_walkTree(t *testing.T) {
	leftObj1 := canvas.NewRectangle(color.Gray16{Y: 1})
	leftObj2 := canvas.NewRectangle(color.Gray16{Y: 2})
	leftCol := &modifiableBox{Box: widget.NewVBox(leftObj1, leftObj2)}
	rightObj1 := canvas.NewRectangle(color.Gray16{Y: 10})
	rightObj2 := canvas.NewRectangle(color.Gray16{Y: 20})
	rightCol := &modifiableBox{Box: widget.NewVBox(rightObj1, rightObj2)}
	content := fyne.NewContainer(leftCol, rightCol)
	content.Move(fyne.NewPos(17, 42))
	leftCol.Move(fyne.NewPos(300, 400))
	leftObj1.Move(fyne.NewPos(1, 2))
	leftObj2.Move(fyne.NewPos(20, 30))
	rightObj1.Move(fyne.NewPos(500, 600))
	rightObj2.Move(fyne.NewPos(60, 70))
	rightCol.Move(fyne.NewPos(7, 8))

	tree := &renderCacheTree{root: &renderCacheNode{obj: content}}
	c := newCanvas()

	type nodeInfo struct {
		obj                                     fyne.CanvasObject
		lastBeforeCallIndex, lastAfterCallIndex int
	}
	updateInfoBefore := func(node *renderCacheNode, index int) {
		pd, _ := node.painterData.(nodeInfo)
		if (pd != nodeInfo{}) && pd.obj != node.obj {
			panic("node cache does not match node obj - nodes should not be reused for different objects")
		}
		pd.obj = node.obj
		pd.lastBeforeCallIndex = index
		node.painterData = pd
	}
	updateInfoAfter := func(node *renderCacheNode, index int) {
		pd := node.painterData.(nodeInfo)
		if pd.obj != node.obj {
			panic("node cache does not match node obj - nodes should not be reused for different objects")
		}
		pd.lastAfterCallIndex = index
		node.painterData = pd
	}

	//
	// test that first walk calls the hooks correctly
	//
	type beforeCall struct {
		obj    fyne.CanvasObject
		parent fyne.CanvasObject
		pos    fyne.Position
	}
	beforeCalls := []beforeCall{}
	type afterCall struct {
		obj    fyne.CanvasObject
		parent fyne.CanvasObject
	}
	afterCalls := []afterCall{}

	var i int
	c.walkTree(tree, func(node *renderCacheNode, pos fyne.Position) {
		var parent fyne.CanvasObject
		if node.parent != nil {
			parent = node.parent.obj
		}
		i++
		updateInfoBefore(node, i)
		beforeCalls = append(beforeCalls, beforeCall{obj: node.obj, parent: parent, pos: pos})
	}, func(node *renderCacheNode) {
		var parent fyne.CanvasObject
		if node.parent != nil {
			parent = node.parent.obj
		}
		i++
		updateInfoAfter(node, i)
		node.minSize.Height = node.obj.Position().Y
		afterCalls = append(afterCalls, afterCall{obj: node.obj, parent: parent})
	})

	assert.Equal(t, []beforeCall{
		{obj: content, pos: fyne.NewPos(17, 42)},
		{obj: leftCol, parent: content, pos: fyne.NewPos(317, 442)},
		{obj: leftObj1, parent: leftCol, pos: fyne.NewPos(318, 444)},
		{obj: leftObj2, parent: leftCol, pos: fyne.NewPos(337, 472)},
		{obj: rightCol, parent: content, pos: fyne.NewPos(24, 50)},
		{obj: rightObj1, parent: rightCol, pos: fyne.NewPos(524, 650)},
		{obj: rightObj2, parent: rightCol, pos: fyne.NewPos(84, 120)},
	}, beforeCalls, "calls before children hook with the correct node and position")
	assert.Equal(t, []afterCall{
		{obj: leftObj1, parent: leftCol},
		{obj: leftObj2, parent: leftCol},
		{obj: leftCol, parent: content},
		{obj: rightObj1, parent: rightCol},
		{obj: rightObj2, parent: rightCol},
		{obj: rightCol, parent: content},
		{obj: content},
	}, afterCalls, "calls after children hook with the correct node")

	//
	// test that second walk gives access to the cache
	//
	secondRunBeforePainterData := []nodeInfo{}
	secondRunAfterPainterData := []nodeInfo{}
	nodes := []*renderCacheNode{}

	c.walkTree(tree, func(node *renderCacheNode, pos fyne.Position) {
		secondRunBeforePainterData = append(secondRunBeforePainterData, node.painterData.(nodeInfo))
		nodes = append(nodes, node)
	}, func(node *renderCacheNode) {
		secondRunAfterPainterData = append(secondRunAfterPainterData, node.painterData.(nodeInfo))
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 13},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightObj2, lastBeforeCallIndex: 11, lastAfterCallIndex: 12},
	}, secondRunBeforePainterData, "second run uses cached nodes")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightObj2, lastBeforeCallIndex: 11, lastAfterCallIndex: 12},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 13},
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
	}, secondRunAfterPainterData, "second run uses cached nodes")
	leftObj1Node := nodes[2]
	leftObj2Node := nodes[3]
	assert.Equal(t, leftObj2Node, leftObj1Node.nextSibling, "correct sibling relation")
	assert.Nil(t, leftObj2Node.nextSibling, "no surplus nodes")
	rightObj1Node := nodes[5]
	rightObj2Node := nodes[6]
	assert.Equal(t, rightObj2Node, rightObj1Node.nextSibling, "correct sibling relation")
	rightColNode := nodes[4]
	assert.Nil(t, rightColNode.nextSibling, "no surplus nodes")

	//
	// test that removal, replacement and adding at the end of a children list works
	//
	leftCol.deleteAt(1)
	leftNewObj2 := canvas.NewRectangle(color.Gray16{Y: 3})
	leftCol.Append(leftNewObj2)
	rightCol.deleteAt(1)
	thirdCol := widget.NewVBox()
	content.AddObject(thirdCol)
	thirdRunBeforePainterData := []nodeInfo{}
	thirdRunAfterPainterData := []nodeInfo{}

	i = 0
	c.walkTree(tree, func(node *renderCacheNode, pos fyne.Position) {
		i++
		updateInfoBefore(node, i)
		thirdRunBeforePainterData = append(thirdRunBeforePainterData, node.painterData.(nodeInfo))
	}, func(node *renderCacheNode) {
		i++
		updateInfoAfter(node, i)
		thirdRunAfterPainterData = append(thirdRunAfterPainterData, node.painterData.(nodeInfo))
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 0}, // new node for replaced obj
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 13},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 0}, // new node for third column
	}, thirdRunBeforePainterData, "third run uses cached nodes if possible")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6}, // new node for replaced obj
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 11},
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 13}, // new node for third column
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
	}, thirdRunAfterPainterData, "third run uses cached nodes if possible")
	assert.NotEqual(t, leftObj2Node, leftObj1Node.nextSibling, "new node for replaced object")
	assert.Nil(t, rightObj1Node.nextSibling, "node for removed object has been removed, too")
	assert.NotNil(t, rightColNode.nextSibling, "new node for new object")

	//
	// test that insertion at the beginnning or in the middle of a children list
	// removes all following siblings and their subtrees
	//
	leftNewObj2a := canvas.NewRectangle(color.Gray16{Y: 4})
	leftCol.insert(leftNewObj2a, 1)
	rightNewObj0 := canvas.NewRectangle(color.Gray16{Y: 30})
	rightCol.Prepend(rightNewObj0)
	fourthRunBeforePainterData := []nodeInfo{}
	fourthRunAfterPainterData := []nodeInfo{}
	nodes = []*renderCacheNode{}

	i = 0
	c.walkTree(tree, func(node *renderCacheNode, pos fyne.Position) {
		i++
		updateInfoBefore(node, i)
		fourthRunBeforePainterData = append(fourthRunBeforePainterData, node.painterData.(nodeInfo))
		nodes = append(nodes, node)
	}, func(node *renderCacheNode) {
		i++
		updateInfoAfter(node, i)
		fourthRunAfterPainterData = append(fourthRunAfterPainterData, node.painterData.(nodeInfo))
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2a, lastBeforeCallIndex: 5, lastAfterCallIndex: 0}, // new node for inserted obj
		{obj: leftNewObj2, lastBeforeCallIndex: 7, lastAfterCallIndex: 0},  // new node because of tail cut
		{obj: rightCol, lastBeforeCallIndex: 10, lastAfterCallIndex: 11},
		{obj: rightNewObj0, lastBeforeCallIndex: 11, lastAfterCallIndex: 0}, // new node for inserted obj
		{obj: rightObj1, lastBeforeCallIndex: 13, lastAfterCallIndex: 0},    // new node because of tail cut
		{obj: thirdCol, lastBeforeCallIndex: 16, lastAfterCallIndex: 13},
	}, fourthRunBeforePainterData, "fourth run uses cached nodes if possible")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2a, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: leftNewObj2, lastBeforeCallIndex: 7, lastAfterCallIndex: 8},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 9},
		{obj: rightNewObj0, lastBeforeCallIndex: 11, lastAfterCallIndex: 12},
		{obj: rightObj1, lastBeforeCallIndex: 13, lastAfterCallIndex: 14},
		{obj: rightCol, lastBeforeCallIndex: 10, lastAfterCallIndex: 15},
		{obj: thirdCol, lastBeforeCallIndex: 16, lastAfterCallIndex: 17},
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 18},
	}, fourthRunAfterPainterData, "fourth run uses cached nodes if possible")
	// check cache tree integrity
	// content node
	assert.Equal(t, content, nodes[0].obj)
	assert.Equal(t, leftCol, nodes[0].firstChild.obj)
	assert.Nil(t, nodes[0].nextSibling)
	// leftCol node
	assert.Equal(t, leftCol, nodes[1].obj)
	assert.Equal(t, leftObj1, nodes[1].firstChild.obj)
	assert.Equal(t, rightCol, nodes[1].nextSibling.obj)
	// leftObj1 node
	assert.Equal(t, leftObj1, nodes[2].obj)
	assert.Nil(t, nodes[2].firstChild)
	assert.Equal(t, leftNewObj2a, nodes[2].nextSibling.obj)
	// leftNewObj2a node
	assert.Equal(t, leftNewObj2a, nodes[3].obj)
	assert.Nil(t, nodes[3].firstChild)
	assert.Equal(t, leftNewObj2, nodes[3].nextSibling.obj)
	// leftNewObj2 node
	assert.Equal(t, leftNewObj2, nodes[4].obj)
	assert.Nil(t, nodes[4].firstChild)
	assert.Nil(t, nodes[4].nextSibling)
	// rightCol node
	assert.Equal(t, rightCol, nodes[5].obj)
	assert.Equal(t, rightNewObj0, nodes[5].firstChild.obj)
	assert.Equal(t, thirdCol, nodes[5].nextSibling.obj)
	// rightNewObj0 node
	assert.Equal(t, rightNewObj0, nodes[6].obj)
	assert.Nil(t, nodes[6].firstChild)
	assert.Equal(t, rightObj1, nodes[6].nextSibling.obj)
	// rightObj1 node
	assert.Equal(t, rightObj1, nodes[7].obj)
	assert.Nil(t, nodes[7].firstChild)
	assert.Nil(t, nodes[7].nextSibling)
	// thirdCol node
	assert.Equal(t, thirdCol, nodes[8].obj)
	assert.Nil(t, nodes[8].firstChild)
	assert.Nil(t, nodes[8].nextSibling)

	//
	// test that removal at the beginning or in the middle of a children list
	// removes all following siblings and their subtrees
	//
	leftCol.deleteAt(1)
	rightCol.deleteAt(0)
	fifthRunBeforePainterData := []nodeInfo{}
	fifthRunAfterPainterData := []nodeInfo{}
	nodes = []*renderCacheNode{}

	i = 0
	c.walkTree(tree, func(node *renderCacheNode, pos fyne.Position) {
		i++
		updateInfoBefore(node, i)
		fifthRunBeforePainterData = append(fifthRunBeforePainterData, node.painterData.(nodeInfo))
		nodes = append(nodes, node)
	}, func(node *renderCacheNode) {
		i++
		updateInfoAfter(node, i)
		fifthRunAfterPainterData = append(fifthRunAfterPainterData, node.painterData.(nodeInfo))
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 18},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 9},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 0}, // new node because of tail cut
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 15},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 0}, // new node because of tail cut
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 17},
	}, fifthRunBeforePainterData, "fifth run uses cached nodes if possible")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 11},
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 13},
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
	}, fifthRunAfterPainterData, "fifth run uses cached nodes if possible")
	// check cache tree integrity
	// content node
	assert.Equal(t, content, nodes[0].obj)
	assert.Equal(t, leftCol, nodes[0].firstChild.obj)
	assert.Nil(t, nodes[0].nextSibling)
	// leftCol node
	assert.Equal(t, leftCol, nodes[1].obj)
	assert.Equal(t, leftObj1, nodes[1].firstChild.obj)
	assert.Equal(t, rightCol, nodes[1].nextSibling.obj)
	// leftObj1 node
	assert.Equal(t, leftObj1, nodes[2].obj)
	assert.Nil(t, nodes[2].firstChild)
	assert.Equal(t, leftNewObj2, nodes[2].nextSibling.obj)
	// leftNewObj2 node
	assert.Equal(t, leftNewObj2, nodes[3].obj)
	assert.Nil(t, nodes[3].firstChild)
	assert.Nil(t, nodes[3].nextSibling)
	// rightCol node
	assert.Equal(t, rightCol, nodes[4].obj)
	assert.Equal(t, rightObj1, nodes[4].firstChild.obj)
	assert.Equal(t, thirdCol, nodes[4].nextSibling.obj)
	// rightObj1 node
	assert.Equal(t, rightObj1, nodes[5].obj)
	assert.Nil(t, nodes[5].firstChild)
	assert.Nil(t, nodes[5].nextSibling)
	// thirdCol node
	assert.Equal(t, thirdCol, nodes[6].obj)
	assert.Nil(t, nodes[6].firstChild)
	assert.Nil(t, nodes[6].nextSibling)
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

type modifiableBox struct {
	*widget.Box
}

func (b *modifiableBox) Append(object fyne.CanvasObject) {
	b.Children = append(b.Children, object)
	widget.Renderer(b).Refresh()
}

func (b *modifiableBox) deleteAt(index int) {
	if index < len(b.Children)-1 {
		b.Children = append(b.Children[:index], b.Children[index+1:]...)
	} else {
		b.Children = b.Children[:index]
	}
	widget.Renderer(b).Refresh()
}

func (b *modifiableBox) insert(object fyne.CanvasObject, index int) {
	tail := append([]fyne.CanvasObject{object}, b.Children[index:]...)
	b.Children = append(b.Children[:index], tail...)
	widget.Renderer(b).Refresh()
}

func (b *modifiableBox) Prepend(object fyne.CanvasObject) {
	b.Children = append([]fyne.CanvasObject{object}, b.Children...)
	widget.Renderer(b).Refresh()
}
