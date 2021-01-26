// +build !windows !ci

package gomobile

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestMobileCanvas_DismissBar(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(canvas.NewRectangle(theme.BackgroundColor()))
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Test"))
	c.showMenu(menu)
	c.resize(fyne.NewSize(100, 100))

	assert.NotNil(t, c.menu)
	// simulate tap as the test util does not know about our menu...
	c.tapDown(fyne.NewPos(80, 20), 1)
	c.tapUp(fyne.NewPos(80, 20), 1, nil, nil, nil, nil)
	assert.Nil(t, c.menu)
}

func TestMobileCanvas_DismissMenu(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(canvas.NewRectangle(theme.BackgroundColor()))
	menu := fyne.NewMainMenu(
		fyne.NewMenu("Test", fyne.NewMenuItem("TapMe", func() {})))
	c.showMenu(menu)
	c.resize(fyne.NewSize(100, 100))

	assert.NotNil(t, c.menu)
	menuObj := c.menu.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*menuLabel)
	point := &fyne.PointEvent{Position: fyne.NewPos(10, 10)}
	menuObj.Tapped(point)

	tapMeItem := c.overlays.Top().(*internalWidget.OverlayContainer).Content.(*widget.Menu).Items[0].(fyne.Tappable)
	tapMeItem.Tapped(point)
	assert.Nil(t, c.menu)
}

func TestMobileCanvas_Menu(t *testing.T) {
	c := &mobileCanvas{}
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
	assert.Equal(t, theme.CancelIcon(), closed.Icon)

	for i := 1; i < 3; i++ {
		item, ok := menuObjects.Objects[i].(*menuLabel)
		assert.True(t, ok)
		assert.Equal(t, labels[i-1], item.menu.Label)
	}
}

func dummyWin(d *mobileDriver, title string) *window {
	ret := &window{title: title}
	d.windows = append(d.windows, ret)

	return ret
}

func TestMobileDriver_FindMenu(t *testing.T) {
	m1 := fyne.NewMainMenu(fyne.NewMenu("1"))
	m2 := fyne.NewMainMenu(fyne.NewMenu("2"))

	d := NewGoMobileDriver().(*mobileDriver)
	w1 := dummyWin(d, "top")
	w1.SetMainMenu(m1)
	assert.Equal(t, m1, d.findMenu(w1))

	w2 := dummyWin(d, "child")
	assert.Equal(t, m1, d.findMenu(w2))

	w2.SetMainMenu(m2)
	assert.Equal(t, m2, d.findMenu(w2))
}
