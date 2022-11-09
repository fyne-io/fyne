package widget

import (
	"fmt"
	"image/color"
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
	bg := test.WidgetRenderer(button).(*buttonRenderer).buttonColor()

	button.Importance = HighImportance
	assert.NotEqual(t, bg, test.WidgetRenderer(button).(*buttonRenderer).buttonColor())
}

func TestButton_DisabledColor(t *testing.T) {
	button := NewButton("Test", nil)
	bg := test.WidgetRenderer(button).(*buttonRenderer).buttonColor()
	button.Importance = MediumImportance
	assert.Equal(t, bg, theme.ButtonColor())

	button.Disable()
	bg = test.WidgetRenderer(button).(*buttonRenderer).buttonColor()
	assert.Equal(t, bg, theme.DisabledButtonColor())
}

func TestButton_Hover_Math(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.HomeIcon(), func() {})
	render := test.WidgetRenderer(button).(*buttonRenderer)
	button.hovered = true
	// unpremultiplied over operator:
	// outA = srcA + dstA*(1-srcA)
	// outC = (srcC*srcA + dstC*dstA*(1-srcA))/outA

	// premultiplied over operator:
	// outC = srcC + dstC*(1-srcA)
	// outA = srcA + dstA*(1-srcA)

	// the buttonRenderer's buttonColor() currently calls RGBA(), turns the alpha into a 0=<a<=1 float32,
	// and blends with the premultiplied over operator, storing the result in an RGBA64.

	// We're going to call ToNRGBA instead (which returns 8-bit components), use the unpremultiplied over operator,
	// store the result in an NRGBA, then convert what buttonColor() returns to an nrgba and compare.
	bg := theme.ButtonColor()
	if render.button.Importance == HighImportance {
		bg = theme.PrimaryColor()
	}
	dstR, dstG, dstB, dstA := col.ToNRGBA(bg)
	srcR, srcG, srcB, srcA := col.ToNRGBA(theme.HoverColor())

	srcAlpha := float32(srcA) / 0xFF
	dstAlpha := float32(dstA) / 0xFF

	outAlpha := (srcAlpha + dstAlpha*(1-srcAlpha))
	outR := uint32((float32(srcR)*srcAlpha + float32(dstR)*dstAlpha*(1-srcAlpha)) / outAlpha)
	outG := uint32((float32(srcG)*srcAlpha + float32(dstG)*dstAlpha*(1-srcAlpha)) / outAlpha)
	outB := uint32((float32(srcB)*srcAlpha + float32(dstB)*dstAlpha*(1-srcAlpha)) / outAlpha)

	nrgba := color.NRGBA{R: uint8(outR), G: uint8(outG), B: uint8(outB), A: uint8(outAlpha * 0xFF)}
	bcn := color.NRGBAModel.Convert(render.buttonColor())

	assert.Equal(t, nrgba, bcn)
}

func TestButton_DisabledIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	button.Enable()
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())
}

func TestButton_DisabledIconChangeUsingSetIcon(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	// re-enable, then change the icon
	button.Enable()
	button.SetIcon(theme.SearchIcon())
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.SearchIcon().Name()))

}

func TestButton_DisabledIconChangedDirectly(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, render.icon.Resource.Name(), theme.CancelIcon().Name())

	// assert we are using the disabled original icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.CancelIcon().Name()))

	// re-enable, then change the icon
	button.Enable()
	button.Icon = theme.SearchIcon()
	render.Refresh()
	assert.Equal(t, render.icon.Resource.Name(), theme.SearchIcon().Name())

	// assert we are using the disabled new icon
	button.Disable()
	assert.Equal(t, render.icon.Resource.Name(), fmt.Sprintf("disabled_%v", theme.SearchIcon().Name()))

}

func TestButton_Focus(t *testing.T) {
	tapped := false
	button := NewButton("Test", func() {
		tapped = true
	})
	render := test.WidgetRenderer(button).(*buttonRenderer)
	assert.Equal(t, theme.ButtonColor(), render.buttonColor())

	assert.Equal(t, false, tapped)
	button.FocusGained()
	render.Refresh() // force update without waiting

	assert.Equal(t, blendColor(theme.ButtonColor(), theme.FocusColor()), render.buttonColor())
	button.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	assert.Equal(t, true, tapped)

	button.FocusLost()
	assert.Equal(t, theme.ButtonColor(), render.buttonColor())
}

func TestButtonRenderer_Layout(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)
	render.Layout(render.MinSize())

	assert.True(t, render.icon.Position().X < render.label.Position().X)
	assert.Equal(t, theme.InnerPadding(), render.icon.Position().X)
	assert.Equal(t, theme.InnerPadding(), render.MinSize().Width-render.label.Position().X-render.label.Size().Width)
}

func TestButtonRenderer_Layout_Stretch(t *testing.T) {
	button := NewButtonWithIcon("Test", theme.CancelIcon(), nil)
	button.Resize(button.MinSize().Add(fyne.NewSize(100, 100)))
	render := test.WidgetRenderer(button).(*buttonRenderer)

	textHeight := render.label.MinSize().Height
	minIconHeight := fyne.Max(theme.IconInlineSize(), textHeight)
	assert.Equal(t, 50+theme.InnerPadding(), render.icon.Position().X, "icon x")
	assert.Equal(t, 50+theme.InnerPadding(), render.icon.Position().Y, "icon y")
	assert.Equal(t, theme.IconInlineSize(), render.icon.Size().Width, "icon width")
	assert.Equal(t, minIconHeight, render.icon.Size().Height, "icon height")
	assert.Equal(t, 50+theme.InnerPadding()+theme.Padding()+theme.IconInlineSize(), render.label.Position().X, "label x")
	assert.Equal(t, render.label.MinSize().Width, render.label.Size().Width, "label size")
}

func TestButtonRenderer_Layout_NoText(t *testing.T) {
	button := NewButtonWithIcon("", theme.CancelIcon(), nil)
	render := test.WidgetRenderer(button).(*buttonRenderer)

	button.Resize(fyne.NewSize(100, 100))

	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().X)
	assert.Equal(t, 50-theme.IconInlineSize()/2, render.icon.Position().Y)
}

func TestButtonRenderer_ApplyTheme(t *testing.T) {
	button := &Button{}
	render := test.WidgetRenderer(button).(*buttonRenderer)
	textRender := test.WidgetRenderer(render.label).(*textRenderer)

	textSize := textRender.Objects()[0].(*canvas.Text).TextSize
	customTextSize := textSize
	test.WithTestTheme(t, func() {
		button.Refresh()
		customTextSize = textRender.Objects()[0].(*canvas.Text).TextSize
	})

	assert.NotEqual(t, textSize, customTextSize)
}

func TestButtonRenderer_TapAnimation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.NewTheme())

	button := NewButton("Hi", func() {})
	w := test.NewWindow(button)
	defer w.Close()
	w.Resize(fyne.NewSize(50, 50).Add(fyne.NewSize(20, 20)))
	button.Resize(fyne.NewSize(50, 50))

	render1 := test.WidgetRenderer(button).(*buttonRenderer)
	test.Tap(button)
	button.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "button/tap_animation.png", w.Canvas().Capture())

	cache.DestroyRenderer(button)
	button.Refresh()

	render2 := test.WidgetRenderer(button).(*buttonRenderer)

	assert.NotEqual(t, render1, render2)

	test.Tap(button)
	button.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "button/tap_animation.png", w.Canvas().Capture())
}
