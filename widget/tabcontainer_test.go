package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestTabContainer_CurrentTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test", Content: NewLabel("Test")})

	assert.Equal(t, 1, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())
}

func TestTabContainer_CurrentTab(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())
}

func TestTabContainer_SelectTab(t *testing.T) {
	tab1 := &TabItem{Text: "Test1", Content: NewLabel("Test1")}
	tab2 := &TabItem{Text: "Test2", Content: NewLabel("Test2")}
	tabs := NewTabContainer(tab1, tab2)

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, tab1, tabs.CurrentTab())

	tabs.SelectTab(tab2)
	assert.Equal(t, tab2, tabs.CurrentTab())
}

func TestTabContainer_SelectTabIndex(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")},
		&TabItem{Text: "Test2", Content: NewLabel("Test2")})

	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, 0, tabs.CurrentTabIndex())

	tabs.SelectTabIndex(1)
	assert.Equal(t, 1, tabs.CurrentTabIndex())
}

func TestTabContainer_ApplyTheme(t *testing.T) {
	tabs := NewTabContainer(&TabItem{Text: "Test1", Content: NewLabel("Test1")})
	underline := Renderer(tabs).(*tabContainerRenderer).line
	barColor := underline.FillColor

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	Renderer(tabs).ApplyTheme()
	assert.NotEqual(t, barColor, underline.FillColor)
}

func TestTabContainerRenderer_Layout(t *testing.T) {
	tabs := NewTabContainer(
		NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
		NewTabItemWithIcon("Text2", theme.ConfirmIcon(), canvas.NewCircle(theme.BackgroundColor())),
	)

	r := Renderer(tabs)
	r.Layout(fyne.NewSize(100, 100))

	images, texts := collectImagesAndTexts(r.Objects(), fyne.NewPos(0, 0))

	if assert.Len(t, images, 2) && assert.Len(t, texts, 2) {
		img1 := images[0]
		img2 := images[1]
		if img1.Position().X > img2.Position().X {
			img2, img1 = img1, img2
		}
		text1 := texts[0]
		text2 := texts[1]
		if text1.Position().X > text2.Position().X {
			text2, text1 = text1, text2
		}

		text1Left := text1.Position().X
		img1Right := img1.Position().X + img1.Size().Width
		assert.Equal(t, img1Right, text1Left)

		text1Right := text1Left + text1.Size().Width
		img2Left := img2.Position().X
		assert.Equal(t, 20, img2Left-text1Right)

		img2Right := img2Left + img2.Size().Width
		text2Left := text2.Position().X
		assert.Equal(t, img2Right, text2Left)
	}
}

func collectImagesAndTexts(objects []fyne.CanvasObject, offset fyne.Position) ([]*canvas.Image, []*canvas.Text) {
	images := []*canvas.Image{}
	texts := []*canvas.Text{}
	for _, o := range objects {
		switch o.(type) {
		case fyne.Widget:
			i, l := collectImagesAndTexts(Renderer(o.(fyne.Widget)).Objects(), o.Position())
			images = append(images, i...)
			texts = append(texts, l...)
		case *canvas.Image:
			o.Move(o.Position().Add(offset))
			images = append(images, o.(*canvas.Image))
		case *canvas.Text:
			o.Move(o.Position().Add(offset))
			texts = append(texts, o.(*canvas.Text))
		}
	}
	return images, texts
}
