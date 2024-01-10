package software

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	internalTest "fyne.io/fyne/v2/internal/test"
	"fyne.io/fyne/v2/test"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestRender(t *testing.T) {
	obj := widget.NewLabel("Hi")
	test.AssertImageMatches(t, "label_dark.png", Render(obj, internalTest.DarkTheme(theme.DefaultTheme())))
	test.AssertImageMatches(t, "label_light.png", Render(obj, internalTest.LightTheme(theme.DefaultTheme())))
}

func TestRender_State(t *testing.T) {
	obj := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {})
	test.AssertImageMatches(t, "button.png", Render(obj, internalTest.DarkTheme(theme.DefaultTheme())))

	obj.Importance = widget.HighImportance
	obj.Refresh()
	test.AssertImageMatches(t, "button_important.png", Render(obj, internalTest.DarkTheme(theme.DefaultTheme())))
}

func TestRender_Focus(t *testing.T) {
	obj := widget.NewEntry()
	test.AssertImageMatches(t, "entry.png", Render(obj, internalTest.DarkTheme(theme.DefaultTheme())))

	obj.FocusGained()
	test.AssertImageMatches(t, "entry_focus.png", Render(obj, internalTest.DarkTheme(theme.DefaultTheme())))
}

func TestRenderCanvas(t *testing.T) {
	obj := container.NewAppTabs(
		container.NewTabItem("Tab 1", container.NewVBox(
			widget.NewLabel("Label"),
			widget.NewButton("Button", func() {}),
		)))

	c := NewCanvas()
	c.SetContent(obj)

	if fyne.CurrentDevice().IsMobile() {
		test.AssertImageMatches(t, "canvas_mobile.png", RenderCanvas(c, internalTest.LightTheme(theme.DefaultTheme())))
	} else {
		test.AssertImageMatches(t, "canvas.png", RenderCanvas(c, internalTest.LightTheme(theme.DefaultTheme())))
	}
}

func TestRender_ImageSize(t *testing.T) {
	image := canvas.NewImageFromFile("../../theme/icons/fyne.png")
	image.FillMode = canvas.ImageFillOriginal
	bg := canvas.NewCircle(color.NRGBA{255, 0, 0, 128})
	bg.StrokeColor = color.White
	bg.StrokeWidth = 5

	c := container.NewStack(image, bg)

	test.AssertImageMatches(t, "image_size.png", Render(c, internalTest.LightTheme(theme.DefaultTheme())))
}
