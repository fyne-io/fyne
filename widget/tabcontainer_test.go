package widget_test

import (
	"testing"

	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/widget"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestTabContainer_CurrentTabIndex(t *testing.T) {
	tabs := widget.NewTabContainer(widget.NewTabItem("Test", widget.NewLabel("Test")))

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())
}

func TestTabContainer_CurrentTab(t *testing.T) {
	tab1 := widget.NewTabItem("Test1", widget.NewLabel("Test1"))
	tab2 := widget.NewTabItem("Test2", widget.NewLabel("Test2"))
	tabs := widget.NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())
}

func TestTabContainer_SelectTab(t *testing.T) {
	tab1 := widget.NewTabItem("Test1", widget.NewLabel("Test1"))
	tab2 := widget.NewTabItem("Test2", widget.NewLabel("Test2"))
	tabs := widget.NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())

	tabs.SelectTab(tab2)
	assert.Equal(t, tab2, tabs.CurrentTab())
}

func TestTabContainer_SelectTabIndex(t *testing.T) {
	tabs := widget.NewTabContainer(widget.NewTabItem("Test1", widget.NewLabel("Test1")),
		widget.NewTabItem("Test2", widget.NewLabel("Test2")))

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())

	tabs.SelectTabIndex(1)
	assert.Equal(t, 1, tabs.CurrentTabIndex())
}

func TestTabContainerRenderer_ApplyTheme(t *testing.T) {
	tabs := widget.NewTabContainer(widget.NewTabItem("Test1", widget.NewLabel("Test1")))
	var underline *canvas.Rectangle
	driver.WalkObjectTree(tabs, fyne.NewPos(0, 0), func(o fyne.CanvasObject, offset fyne.Position) bool {
		if u, ok := o.(*canvas.Rectangle); ok {
			underline = u
			return true
		}
		return false
	}, nil)

	r := widget.Renderer(tabs)

	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	r.ApplyTheme()
	darkColor := underline.FillColor

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	r.ApplyTheme()
	lightColor := underline.FillColor

	assert.NotEqual(t, darkColor, lightColor)
}

func TestTabContainerRenderer_Layout(t *testing.T) {
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
		widget.NewTabItemWithIcon("Text2", theme.ConfirmIcon(), canvas.NewCircle(theme.BackgroundColor())),
	)

	r := widget.Renderer(tabs)
	r.Layout(fyne.NewSize(100, 100))

	type dim struct {
		fyne.Position
		fyne.Size
	}
	images := []dim{}
	texts := []dim{}
	driver.WalkObjectTree(tabs, fyne.NewPos(0, 0), func(o fyne.CanvasObject, offset fyne.Position) bool {
		switch o.(type) {
		case *canvas.Image:
			images = append(images, dim{Position: o.Position().Add(offset), Size: o.Size()})
		case *canvas.Text:
			texts = append(texts, dim{Position: o.Position().Add(offset), Size: o.Size()})
		}
		return false
	}, nil)

	if assert.Len(t, images, 2) && assert.Len(t, texts, 2) {
		img1 := images[0]
		img2 := images[1]
		if img1.X > img2.X {
			img2, img1 = img1, img2
		}
		text1 := texts[0]
		text2 := texts[1]
		if text1.X > text2.X {
			text2, text1 = text1, text2
		}

		text1Left := text1.X
		img1Right := img1.X + img1.Width
		assert.Equal(t, img1Right, text1Left)

		text1Right := text1Left + text1.Width
		img2Left := img2.X
		assert.Equal(t, 20, img2Left-text1Right)

		img2Right := img2Left + img2.Width
		text2Left := text2.X
		assert.Equal(t, img2Right, text2Left)
	}
}
