package software_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/test"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestRender(t *testing.T) {
	obj := widget.NewLabel("Hi")
	test.AssertImageMatches(t, "label.png", software.Render(obj, test.Theme()))
	painter.ClearFontCache() // avoid side effects of the cause of #4937
	test.AssertImageMatches(t, "label_ugly_theme.png", software.Render(obj, test.NewTheme()))
	painter.ClearFontCache() // avoid side effects of the cause of #4937
}

func TestRender_State(t *testing.T) {
	obj := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {})
	test.AssertImageMatches(t, "button.png", software.Render(obj, test.Theme()))

	obj.Importance = widget.HighImportance
	obj.Refresh()
	test.AssertImageMatches(t, "button_important.png", software.Render(obj, test.Theme()))
}

func TestRender_Focus(t *testing.T) {
	obj := widget.NewEntry()
	test.AssertImageMatches(t, "entry.png", software.Render(obj, test.Theme()))

	obj.FocusGained()
	test.AssertImageMatches(t, "entry_focus.png", software.Render(obj, test.Theme()))
}

func TestRenderCanvas(t *testing.T) {
	obj := container.NewAppTabs(
		container.NewTabItem("Tab 1", container.NewVBox(
			widget.NewLabel("Label"),
			widget.NewButton("Button", func() {}),
		)),
	)

	c := software.NewCanvas()
	c.SetContent(obj)

	if fyne.CurrentDevice().IsMobile() {
		test.AssertImageMatches(t, "canvas_mobile.png", software.RenderCanvas(c, test.Theme()))
	} else {
		test.AssertImageMatches(t, "canvas.png", software.RenderCanvas(c, test.Theme()))
	}
}

func TestRender_ImageSize(t *testing.T) {
	image := canvas.NewImageFromFile("../../theme/icons/fyne.png")
	image.FillMode = canvas.ImageFillOriginal
	bg := canvas.NewCircle(color.NRGBA{255, 0, 0, 128})
	bg.StrokeColor = color.White
	bg.StrokeWidth = 5

	c := container.NewStack(image, bg)

	test.AssertImageMatches(t, "image_size.png", software.Render(c, test.Theme()))
}
