package test_test

import (
	"image/color"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestAssertCanvasTappableAt(t *testing.T) {
	c := test.NewCanvas()
	b := widget.NewButton("foo", nil)
	c.SetContent(b)
	c.Resize(fyne.NewSize(300, 300))
	b.Resize(fyne.NewSize(100, 100))
	b.Move(fyne.NewPos(100, 100))

	tt := &testing.T{}
	assert.True(t, test.AssertCanvasTappableAt(tt, c, fyne.NewPos(101, 101)), "tappable found")
	assert.False(t, tt.Failed(), "test did not fail")

	tt = &testing.T{}
	assert.False(t, test.AssertCanvasTappableAt(tt, c, fyne.NewPos(99, 99)), "tappable not found")
	assert.True(t, tt.Failed(), "test failed")
}

func TestAssertRendersToMarkup(t *testing.T) {
	c := test.NewCanvas()
	c.SetContent(canvas.NewCircle(color.Black))

	markup := "<canvas padded size=\"9x9\">\n" +
		"\t<content>\n" +
		"\t\t<circle fillColor=\"rgba(0,0,0,255)\" pos=\"4,4\" size=\"1x1\"/>\n" +
		"\t</content>\n" +
		"</canvas>\n"

	t.Run("non-existing master", func(t *testing.T) {
		tt := &testing.T{}
		assert.False(t, test.AssertRendersToMarkup(tt, "non_existing_master.xml", c), "non existing master is not equal to rendered markup")
		assert.True(t, tt.Failed(), "test failed")
		assert.Equal(t, markup, readMarkup(t, "testdata/failed/non_existing_master.xml"), "markup was written to disk")
	})

	t.Run("matching master", func(t *testing.T) {
		tt := &testing.T{}
		assert.True(t, test.AssertRendersToMarkup(tt, "markup_master.xml", c), "existing master is equal to rendered markup")
		assert.False(t, tt.Failed(), "test should not fail")
	})

	t.Run("diffing master", func(t *testing.T) {
		tt := &testing.T{}
		assert.False(t, test.AssertRendersToMarkup(tt, "markup_diffing_master.xml", c), "existing master is not equal to rendered markup")
		assert.True(t, tt.Failed(), "test should fail")
		assert.Equal(t, markup, readMarkup(t, "testdata/failed/markup_diffing_master.xml"), "markup was written to disk")
	})

	if !t.Failed() {
		os.RemoveAll("testdata/failed")
	}
}

func TestDrag(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	d := &draggable{}
	c.SetContent(fyne.NewContainerWithoutLayout(d))
	c.Resize(fyne.NewSize(30, 30))
	d.Resize(fyne.NewSize(20, 20))
	d.Move(fyne.NewPos(10, 10))

	test.Drag(c, fyne.NewPos(5, 5), 10, 10)
	assert.Nil(t, d.event, "nothing happens if no draggable was found at position")
	assert.False(t, d.wasDragged)

	test.Drag(c, fyne.NewPos(15, 15), 17, 42)
	assert.Equal(t, &fyne.DragEvent{
		PointEvent: fyne.PointEvent{Position: fyne.Position{X: 5, Y: 5}},
		Dragged:    fyne.NewDelta(17, 42),
	}, d.event)
	assert.True(t, d.wasDragged)
}

func TestFocusNext(t *testing.T) {
	c := test.NewCanvas()
	f1 := &focusable{}
	f2 := &focusable{}
	f3 := &focusable{}
	c.SetContent(fyne.NewContainerWithoutLayout(f1, f2, f3))

	assert.Nil(t, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f2, c.Focused())
	assert.False(t, f1.focused)
	assert.True(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.True(t, f3.focused)

	test.FocusNext(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)
}

func TestFocusPrevious(t *testing.T) {
	c := test.NewCanvas()
	f1 := &focusable{}
	f2 := &focusable{}
	f3 := &focusable{}
	c.SetContent(fyne.NewContainerWithoutLayout(f1, f2, f3))

	assert.Nil(t, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.True(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f2, c.Focused())
	assert.False(t, f1.focused)
	assert.True(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f1, c.Focused())
	assert.True(t, f1.focused)
	assert.False(t, f2.focused)
	assert.False(t, f3.focused)

	test.FocusPrevious(c)
	assert.Equal(t, f3, c.Focused())
	assert.False(t, f1.focused)
	assert.False(t, f2.focused)
	assert.True(t, f3.focused)
}

func TestScroll(t *testing.T) {
	c := test.NewCanvas()
	c.SetPadded(false)
	s := &scrollable{}
	c.SetContent(fyne.NewContainerWithoutLayout(s))
	c.Resize(fyne.NewSize(30, 30))
	s.Resize(fyne.NewSize(20, 20))
	s.Move(fyne.NewPos(10, 10))

	test.Scroll(c, fyne.NewPos(5, 5), 10, 10)
	assert.Nil(t, s.event, "nothing happens if no scrollable was found at position")

	test.Scroll(c, fyne.NewPos(15, 15), 17, 42)
	assert.Equal(t, &fyne.ScrollEvent{Scrolled: fyne.NewDelta(17, 42)}, s.event)
}

func readMarkup(t *testing.T, path string) string {
	raw, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return string(raw)
}

var _ fyne.Draggable = (*draggable)(nil)

type draggable struct {
	widget.BaseWidget
	event      *fyne.DragEvent
	wasDragged bool
}

func (d *draggable) DragEnd() {
	d.wasDragged = true
}

func (d *draggable) Dragged(event *fyne.DragEvent) {
	d.event = event
}

var _ fyne.Focusable = (*focusable)(nil)

type focusable struct {
	widget.BaseWidget
	focused bool
}

func (f *focusable) FocusGained() {
	f.focused = true
}

func (f *focusable) FocusLost() {
	f.focused = false
}

func (f *focusable) TypedKey(event *fyne.KeyEvent) {
}

func (f *focusable) TypedRune(r rune) {
}

var _ fyne.Scrollable = (*scrollable)(nil)

type scrollable struct {
	widget.BaseWidget
	event *fyne.ScrollEvent
}

func (s *scrollable) Scrolled(event *fyne.ScrollEvent) {
	s.event = event
}
