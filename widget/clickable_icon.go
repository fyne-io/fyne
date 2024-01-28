package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type ClickableIcon struct {
	BaseWidget

	OnTapped func() `json:"-"`

	tapAnim *fyne.Animation
	icon    *canvas.Image
}

func NewClickableIcon(icon fyne.Resource, tapped func()) *ClickableIcon {
	b := &ClickableIcon{
		OnTapped: tapped,
	}

	b.ExtendBaseWidget(b)

	b.icon = canvas.NewImageFromResource(icon)
	b.icon.FillMode = canvas.ImageFillContain
	b.tapAnim = fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		k := (1 - done) / 2
		b.icon.Resize(fyne.NewSize(b.Size().Width*done, b.Size().Height*done))
		b.icon.Move(fyne.NewPos(b.Size().Width*k, b.Size().Height*k))

		canvas.Refresh(b.icon)
	})
	b.tapAnim.Curve = fyne.AnimationEaseOut
	b.SetIcon(icon)
	return b
}

func (b *ClickableIcon) CreateRenderer() fyne.WidgetRenderer {
	return b
}

func (b *ClickableIcon) SetIcon(icon fyne.Resource) {
	b.icon.Resource = icon
	b.Refresh()
}

func (b *ClickableIcon) Tapped(*fyne.PointEvent) {
	b.tapAnim.Stop()
	b.OnTapped()
	b.tapAnim.Start()
}

func (b *ClickableIcon) Layout(size fyne.Size) {
	b.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
	b.icon.Resize(b.MinSize())
}

func (b *ClickableIcon) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

func (b *ClickableIcon) Refresh() {
	if b.Visible() {
		b.icon.Show()
	} else {
		b.icon.Hide()
	}
	b.icon.Refresh()
	b.Layout(b.Size())
}

func (*ClickableIcon) Destroy() {}

func (b *ClickableIcon) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{b.icon}
}
