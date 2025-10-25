package widget

import (
	"fmt"
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestButton_Style(t *testing.T) {
	button := NewButton("Test", nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	render.applyTheme()
	bg := render.background.FillColor

	button.Importance = HighImportance
	render.applyTheme()
	assert.NotEqual(t, bg, render.background.FillColor)
}

func TestButton_DisabledColor(t *testing.T) {
	button := NewButton("Test", nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	render.applyTheme()
	bg := render.background.FillColor
	button.Importance = MediumImportance
	render.applyTheme()
	assert.Equal(t, bg, render.background.FillColor)

	button.Disable()
	render.applyTheme()
	assert.Equal(t, theme.Color(theme.ColorNameDisabledButton), render.background.FillColor)
}

func TestButton_Hover_Math(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.HomeIcon(), func() {})
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	button.hovered = true
	// unpremultiplied over operator:
	// outA = srcA + dstA*(1-srcA)
	// outC = (srcC*srcA + dstC*dstA*(1-srcA))/outA

	// premultiplied over operator:
	// outC = srcC + dstC*(1-srcA)
	// outA = srcA + dstA*(1-srcA)

	// the buttonRenderer currently calls RGBA(), turns the alpha into a 0=<a<=1 float32,
	// and blends with the premultiplied over operator, storing the result in an RGBA64.

	// We're going to call ToNRGBA instead (which returns 8-bit components), use the unpremultiplied over operator,
	// store the result in an NRGBA, then convert what buttonColor() returns to an nrgba and compare.
	bg := theme.Color(theme.ColorNameButton)
	if render.button.Importance == HighImportance {
		bg = theme.Color(theme.ColorNamePrimary)
	}
	dstR, dstG, dstB, dstA := col.ToNRGBA(bg)
	srcR, srcG, srcB, srcA := col.ToNRGBA(theme.Color(theme.ColorNameHover))

	srcAlpha := float32(srcA) / 0xFF
	dstAlpha := float32(dstA) / 0xFF

	outAlpha := (srcAlpha + dstAlpha*(1-srcAlpha))
	outR := uint32((float32(srcR)*srcAlpha + float32(dstR)*dstAlpha*(1-srcAlpha)) / outAlpha)
	outG := uint32((float32(srcG)*srcAlpha + float32(dstG)*dstAlpha*(1-srcAlpha)) / outAlpha)
	outB := uint32((float32(srcB)*srcAlpha + float32(dstB)*dstAlpha*(1-srcAlpha)) / outAlpha)

	nrgba := color.NRGBA{R: uint8(outR), G: uint8(outG), B: uint8(outB), A: uint8(outAlpha * 0xFF)}
	render.applyTheme()
	bcn := color.NRGBAModel.Convert(render.background.FillColor)

	assert.Equal(t, nrgba, bcn)
}

func TestButton_DisabledIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	button.Disable()
	assert.True(t, strings.HasPrefix(render.icon.Resource.Name(), "disabled_"))

	button.Enable()
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())
}

func TestButton_DisabledIconChangeUsingSetIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	cancelBaseName := strings.TrimPrefix(theme.CancelIcon().Name(), "foreground_")
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", cancelBaseName))

	// re-enable, then change the icon
	button.Enable()
	button.SetIcon(theme.SearchIcon())
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	searchBaseName := strings.TrimPrefix(theme.SearchIcon().Name(), "foreground_")
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", searchBaseName))
}

func TestButton_DisabledIconChangedDirectly(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	cancelBaseName := strings.TrimPrefix(theme.CancelIcon().Name(), "foreground_")
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", cancelBaseName))

	// re-enable, then change the icon
	button.Enable()
	button.Icon = theme.SearchIcon()
	render.Refresh()
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	searchBaseName := strings.TrimPrefix(theme.SearchIcon().Name(), "foreground_")
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", searchBaseName))
}

func TestButton_Focus(t *testing.T) {
	tapped := false
	button := NewButton("Test", func() {
		tapped = true
	})
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	render.applyTheme()
	assert.Equal(t, theme.Color(theme.ColorNameButton), render.background.FillColor)

	assert.False(t, tapped)
	button.FocusGained()
	render.Refresh() // force update without waiting

	assert.Equal(t, blendColor(theme.Color(theme.ColorNameButton), theme.Color(theme.ColorNameFocus)), render.background.FillColor)
	button.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.True(t, tapped)

	button.FocusLost()
	render.applyTheme()
	assert.Equal(t, theme.Color(theme.ColorNameButton), render.background.FillColor)
}

func TestButtonRenderer_Layout(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	render.Layout(render.MinSize())

	assert.Less(t, render.icon.Position().X, render.label.Position().X)
	assert.Equal(t, theme.InnerPadding(), render.icon.Position().X)
	assert.Equal(t, theme.InnerPadding(), render.MinSize().Width-render.label.Position().X-render.label.Size().Width)
}

func TestButtonRenderer_Layout_Stretch(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	button.Resize(button.MinSize().Add(fyne.NewSize(100, 100)))
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)

	textHeight := render.label.MinSize().Height
	minIconHeight := max(theme.IconInlineSize(), textHeight)
	assert.Equal(t, 50+theme.InnerPadding(), render.icon.Position().X, "icon x")
	assert.Equal(t, 50+theme.InnerPadding(), render.icon.Position().Y, "icon y")
	assert.Equal(t, theme.IconInlineSize(), render.icon.Size().Width, "icon width")
	assert.Equal(t, minIconHeight, render.icon.Size().Height, "icon height")
	assert.Equal(t, 50+theme.InnerPadding()+theme.Padding()+theme.IconInlineSize(), render.label.Position().X, "label x")
	assert.Equal(t, render.label.MinSize().Width, render.label.Size().Width, "label size")
}

func TestButtonRenderer_Layout_NoText(t *testing.T) {
	button := NewButtonWithIcon("", theme.CancelIcon(), nil)
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)

	button.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().X)
	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().Y)
}

func TestButtonRenderer_ApplyTheme(t *testing.T) {
	button := &Button{}
	render := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	textRender := test.TempWidgetRenderer(t, render.label).(*textRenderer)

	textSize := textRender.Objects()[0].(*canvas.Text).TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		button.Refresh()
		customTextSize = textRender.Objects()[0].(*canvas.Text).TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func TestButtonRenderer_TapAnimation(t *testing.T) {
	test.NewTempApp(t)
	test.ApplyTheme(t, test.NewTheme())

	button := NewButton("Hi", func() {})
	w := test.NewWindow(button)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50).Add(fyne.NewSize(20, 20)))
	button.Resize(fyne.NewSize(50, 50))

	render1 := test.TempWidgetRenderer(t, button).(*buttonRenderer)
	test.Tap(button)
	button.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "button/tap_animation.png", w.Canvas().Capture())

	cache.DestroyRenderer(button)
	button.Refresh()

	render2 := test.TempWidgetRenderer(t, button).(*buttonRenderer)

	assert.NotEqual(t, render1, render2)

	test.Tap(button)
	button.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "button/tap_animation.png", w.Canvas().Capture())
}
