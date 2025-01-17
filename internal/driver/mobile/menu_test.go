//go:build !windows || !ci

package mobile

import (
	"testing"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestMobileCanvas_DismissBar(t *testing.T) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(fynecanvas.NewRectangle(theme.Color(theme.ColorNameBackground)))
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Test"))
	c.showMenu(menu)
	c.Resize(fyne.NewSize(100, 100))

	assert.NotNil(t, c.menu)
	// simulate tap as the test util does not know about our menu...
	c.tapDown(fyne.NewPos(80, 20), 1)
	c.tapUp(fyne.NewPos(80, 20), 1, nil, nil, nil, nil)
	assert.Nil(t, c.menu)
}

func TestMobileCanvas_DismissMenu(t *testing.T) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.padded = false
	c.SetContent(fynecanvas.NewRectangle(theme.Color(theme.ColorNameBackground)))
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Test", fyne.NewMenuItem("TapMe", func() {})))
	c.showMenu(menu)
	c.Resize(fyne.NewSize(100, 100))

	assert.NotNil(t, c.menu)
	menuObj := c.menu.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*menuLabel)
	point := &fyne.PointEvent{Position: fyne.NewPos(10, 10)}
	menuObj.Tapped(point)

	tapMeItem := c.Overlays().Top().(*internalWidget.OverlayContainer).Content.(*widget.PopUpMenu).Items[0].(fyne.Tappable)
	tapMeItem.Tapped(point)
	assert.Nil(t, c.menu)
}

func TestMobileCanvas_Menu(t *testing.T) {
	c := &canvas{}
	labels := []string{"File", "Edit"}
	menu := fyne.NewMainMenu(
		fyne.NewMenu(labels[0]),
		fyne.NewMenu(labels[1]))

	c.showMenu(menu)
	menuObjects := c.menu.(*fyne.Container).Objects[1].(*fyne.Container)
	assert.Equal(t, 3, len(menuObjects.Objects))
	header, ok := menuObjects.Objects[0].(*fyne.Container)
	assert.True(t, ok)
	closed, ok := header.Objects[0].(*widget.Button)
	assert.True(t, ok)
	assert.Equal(t, theme.CancelIcon().Name(), closed.Icon.Name())

	for i := 1; i < 3; i++ {
		item, ok := menuObjects.Objects[i].(*menuLabel)
		assert.True(t, ok)
		assert.Equal(t, labels[i-1], item.menu.Label)
	}
}

func TestMobileCanvas_MenuChild(t *testing.T) {
	c := &canvas{}
	c.Initialize(c, nil)
	c.Resize(fyne.NewSize(100, 200))

	child := fyne.NewMenu("Child", fyne.NewMenuItem("One", func() {}))
	parent := fyne.NewMenuItem("Parent", func() {})
	parent.ChildMenu = child
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Top", parent))

	c.showMenu(menu)
	topObj := c.menu.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*menuLabel)
	assert.Equal(t, "Top", topObj.menu.Label)

	test.Tap(topObj)
	assert.Equal(t, 1, len(c.Overlays().List()))
	rootOverlay := c.Overlays().Top()
	parentItem := rootOverlay.(*internalWidget.OverlayContainer).Content.(*widget.PopUpMenu).Items[0]
	parentDetails := test.WidgetRenderer(parentItem.(fyne.Widget))
	text := parentDetails.Objects()[1].(*fynecanvas.Text)
	assert.Equal(t, "Parent", text.Text)

	test.Tap(parentItem.(fyne.Tappable))
	assert.NotNil(t, c.menu) // still visible

	childMenu := parentItem.(interface{ Child() *widget.Menu }).Child()
	assert.NotNil(t, childMenu)
	assert.True(t, childMenu.Visible())
	assert.Equal(t, 1, len(childMenu.Items))

	childDetails := test.WidgetRenderer(childMenu.Items[0].(fyne.Widget))
	text = childDetails.Objects()[1].(*fynecanvas.Text)
	assert.Equal(t, "One", text.Text)
}

func dummyWin(d *driver, title string) *window {
	ret := &window{title: title}
	d.windows = append(d.windows, ret)

	return ret
}

func TestMobileDriver_FindMenu(t *testing.T) {
	m1 := fyne.NewMainMenu(fyne.NewMenu("1"))
	m2 := fyne.NewMainMenu(fyne.NewMenu("2"))

	d := NewGoMobileDriver().(*driver)
	w1 := dummyWin(d, "top")
	w1.SetMainMenu(m1)
	assert.Equal(t, m1, d.findMenu(w1))

	w2 := dummyWin(d, "child")
	assert.Equal(t, m1, d.findMenu(w2))

	w2.SetMainMenu(m2)
	assert.Equal(t, m2, d.findMenu(w2))
}
