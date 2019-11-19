package widget

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

type myWidget struct {
	BaseWidget

	refreshed chan bool
}

func (m *myWidget) Enable() {
	m.enable(m)
}

func (m *myWidget) Disable() {
	m.disable(m)
}

func (m *myWidget) Disabled() bool {
	return m.disabled
}

func (m *myWidget) Refresh() {
	m.refreshed <- true
}

func (m *myWidget) CreateRenderer() fyne.WidgetRenderer {
	m.ExtendBaseWidget(m)
	return (&Box{}).CreateRenderer()
}

func TestApplyThemeCalled(t *testing.T) {
	widget := &myWidget{refreshed: make(chan bool)}

	window := test.NewWindow(widget)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())

	func() {
		select {
		case <-widget.refreshed:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for theme apply")
		}
	}()

	close(widget.refreshed)
	window.Close()
}

func TestApplyThemeCalledChild(t *testing.T) {
	child := &myWidget{refreshed: make(chan bool)}
	parent := NewVBox(child)

	window := test.NewWindow(parent)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	func() {
		select {
		case <-child.refreshed:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for child theme apply")
		}
	}()

	close(child.refreshed)
	window.Close()
}

var (
	red   = &color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green = &color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue  = &color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

const testTextSize = 18

// testTheme is a simple theme variation used for testing that widgets adapt correctly
type testTheme struct {
}

func (testTheme) BackgroundColor() color.Color {
	return red
}

func (testTheme) ButtonColor() color.Color {
	return color.Black
}

func (testTheme) DisabledButtonColor() color.Color {
	return color.White
}

func (testTheme) HyperlinkColor() color.Color {
	return green
}

func (testTheme) TextColor() color.Color {
	return color.White
}

func (testTheme) DisabledTextColor() color.Color {
	return color.Black
}

func (testTheme) IconColor() color.Color {
	return color.White
}

func (testTheme) DisabledIconColor() color.Color {
	return color.Black
}

func (testTheme) PlaceHolderColor() color.Color {
	return blue
}

func (testTheme) PrimaryColor() color.Color {
	return green
}

func (testTheme) HoverColor() color.Color {
	return green
}

func (testTheme) FocusColor() color.Color {
	return green
}

func (testTheme) ScrollBarColor() color.Color {
	return blue
}

func (testTheme) ShadowColor() color.Color {
	return blue
}

func (testTheme) TextSize() int {
	return testTextSize
}

func (testTheme) TextFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

func (testTheme) TextBoldFont() fyne.Resource {
	return theme.DefaultTextItalicFont()
}

func (testTheme) TextItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (testTheme) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

func (testTheme) TextMonospaceFont() fyne.Resource {
	return theme.DefaultTextFont()
}

func (testTheme) Padding() int {
	return 10
}

func (testTheme) IconInlineSize() int {
	return 24
}

func (testTheme) ScrollBarSize() int {
	return 10
}

func (testTheme) ScrollBarSmallSize() int {
	return 2
}

func withTestTheme(f func()) {
	settings := fyne.CurrentApp().Settings()
	current := settings.Theme()
	settings.SetTheme(&testTheme{})
	f()
	settings.SetTheme(current)
}
