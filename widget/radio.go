package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

type radioRenderItem struct {
	icon  *canvas.Image
	label *canvas.Text
}

type radioRenderer struct {
	items []*radioRenderItem

	objects []fyne.CanvasObject
	radio   *Radio
}

func removeDuplicates(options []string) []string {
	var result []string
	found := make(map[string]bool)

	for _, option := range options {
		if _, ok := found[option]; !ok {
			found[option] = true
			result = append(result, option)
		}
	}

	return result
}

// MinSize calculates the minimum size of a radio item.
// This is based on the contained text, the radio icon and a standard amount of padding
// between each item.
func (r *radioRenderer) MinSize() fyne.Size {
	minWidth := 0
	height := 0
	for _, item := range r.items {
		itemMin := item.label.MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
		itemMin = itemMin.Add(fyne.NewSize(theme.IconInlineSize()+theme.Padding(), 0))

		minWidth = fyne.Max(minWidth, itemMin.Width)
		height += itemMin.Height
	}

	return fyne.NewSize(minWidth, height)
}

// Layout the components of the radio widget
func (r *radioRenderer) Layout(size fyne.Size) {
	itemHeight := r.radio.itemHeight()
	y := 0
	labelSize := fyne.NewSize(size.Width-theme.IconInlineSize()-theme.Padding(), itemHeight)

	for _, item := range r.items {
		item.label.Resize(labelSize)
		item.label.Move(fyne.NewPos(theme.IconInlineSize()+theme.Padding(), y))

		item.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		item.icon.Move(fyne.NewPos(0,
			y+(labelSize.Height-theme.IconInlineSize())/2))

		y += itemHeight
	}
}

// ApplyTheme is called when the Radio may need to update it's look
func (r *radioRenderer) ApplyTheme() {
	for _, item := range r.items {
		item.label.Color = theme.TextColor()
	}

	r.Refresh()
}

func (r *radioRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *radioRenderer) Refresh() {
	r.radio.removeDuplicateOptions()

	if len(r.items) < len(r.radio.Options) {
		for i := len(r.items); i < len(r.radio.Options); i++ {
			option := r.radio.Options[i]
			icon := canvas.NewImageFromResource(theme.RadioButtonIcon())

			text := canvas.NewText(option, theme.TextColor())
			text.Alignment = fyne.TextAlignLeading

			r.objects = append(r.objects, icon, text)
			r.items = append(r.items, &radioRenderItem{icon, text})
		}
		r.Layout(r.radio.Size())
	} else if len(r.items) > len(r.radio.Options) {
		total := len(r.radio.Options)
		r.items = r.items[:total]
		r.objects = r.objects[:total*2]
	}

	for i, item := range r.items {
		option := r.radio.Options[i]
		item.label.Text = option

		if r.radio.Selected == option {
			item.icon.Resource = theme.RadioButtonCheckedIcon()
		} else {
			item.icon.Resource = theme.RadioButtonIcon()
		}
	}

	canvas.Refresh(r.radio)
}

func (r *radioRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *radioRenderer) Destroy() {
}

// Radio widget has a list of text labels and radio check icons next to each.
// Changing the selection (only one can be selected) will trigger the changed func.
type Radio struct {
	baseWidget
	Options  []string
	Selected string

	OnChanged func(string) `json:"-"`
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (r *Radio) Resize(size fyne.Size) {
	r.resize(size, r)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (r *Radio) Move(pos fyne.Position) {
	r.move(pos, r)
}

// MinSize returns the smallest size this widget can shrink to
func (r *Radio) MinSize() fyne.Size {
	return r.minSize(r)
}

// Show this widget, if it was previously hidden
func (r *Radio) Show() {
	r.show(r)
}

// Hide this widget, if it was previously visible
func (r *Radio) Hide() {
	r.hide(r)
}

// Append adds a new option to the end of a Radio widget.
func (r *Radio) Append(option string) {
	r.Options = append(r.Options, option)

	Refresh(r)
}

// Tapped is called when a pointer tapped event is captured and triggers any change handler
func (r *Radio) Tapped(event *fyne.PointEvent) {
	index := (event.Position.Y - theme.Padding()) / r.itemHeight()
	if event.Position.Y < theme.Padding() || index >= len(r.Options) { // in the padding
		return
	}
	clicked := r.Options[index]

	if r.Selected == clicked {
		r.Selected = ""
	} else {
		r.Selected = clicked
	}

	if r.OnChanged != nil {
		r.OnChanged(r.Selected)
	}
	Renderer(r).Refresh()
}

// TappedSecondary is called when a secondary pointer tapped event is captured
func (r *Radio) TappedSecondary(*fyne.PointEvent) {
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (r *Radio) CreateRenderer() fyne.WidgetRenderer {
	var items []*radioRenderItem
	var objects []fyne.CanvasObject

	for _, option := range r.Options {
		icon := canvas.NewImageFromResource(theme.RadioButtonIcon())

		text := canvas.NewText(option, theme.TextColor())
		text.Alignment = fyne.TextAlignLeading

		objects = append(objects, icon, text)
		items = append(items, &radioRenderItem{icon, text})
	}

	return &radioRenderer{items, objects, r}
}

// SetSelected sets the radio option, it can be used to set a default option.
func (r *Radio) SetSelected(option string) {
	if r.Selected == option {
		return
	}

	r.Selected = option

	Renderer(r).Refresh()
}

func (r *Radio) itemHeight() int {
	return r.MinSize().Height / len(r.Options)
}

func (r *Radio) removeDuplicateOptions() {
	r.Options = removeDuplicates(r.Options)
}

// NewRadio creates a new radio widget with the set options and change handler
func NewRadio(options []string, changed func(string)) *Radio {
	r := &Radio{
		baseWidget{},
		options,
		"",
		changed,
	}

	r.removeDuplicateOptions()

	Renderer(r).Layout(r.MinSize())
	return r
}
