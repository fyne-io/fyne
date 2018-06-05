package test

import "github.com/fyne-io/fyne/api/ui"

var dummyCanvas ui.Canvas

type testCanvas struct {
	content ui.CanvasObject
	focused ui.FocusableObject
}

func (c *testCanvas) Content() ui.CanvasObject {
	return c.content
}

func (c *testCanvas) SetContent(content ui.CanvasObject) {
	c.content = content
}

func (c *testCanvas) Refresh(ui.CanvasObject) {
}

func (c *testCanvas) Contains(ui.CanvasObject) bool {
	return true
}

func (c *testCanvas) Focus(obj ui.FocusableObject) {
	c.focused = obj
	obj.OnFocusGained()
}

func (c *testCanvas) Focused() ui.FocusableObject {
	return c.focused
}

func (c *testCanvas) Size() ui.Size {
	return ui.NewSize(10, 10)
}

func (c *testCanvas) Scale() float32 {
	return 1.0
}
func (c *testCanvas) SetScale(float32) {
}

func (c *testCanvas) SetOnKeyDown(func(*ui.KeyEvent)) {
}

// GetTestCanvas returns a reusable in-memory canvas used for testing
func GetTestCanvas() ui.Canvas {
	if dummyCanvas == nil {
		dummyCanvas = &testCanvas{}
	}

	return dummyCanvas
}
