package dialog

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// ColorPickerDialog is a simple dialog window that displays a color picker.
type ColorPickerDialog struct {
	*dialog
	color    color.Color
	advanced *widget.AccordionContainer
	picker   *widget.ColorAdvancedPicker
}

// SetColor updates the color of the color picker.
func (p *ColorPickerDialog) SetColor(c color.Color) {
	p.picker.SetColor(c)
}

// NewColorPicker creates a color dialog and returns the handle.
// Using the returned type you should call Show() and then set its color through SetColor().
// The callback is triggered when the user selects a color.
func NewColorPicker(title, message string, callback func(c color.Color), parent fyne.Window) *ColorPickerDialog {
	d := newDialog(title, message, theme.ColorPaletteIcon(), nil /*cancel?*/, parent)

	cb := func(c color.Color) {
		d.Hide()
		// TODO FIXME writeRecentColor(colorToString(p.color))
		callback(c)
	}

	basic := widget.NewColorBasicPicker(cb)

	greyscale := widget.NewColorGreyscalePicker(cb)

	recent := widget.NewColorRecentPicker(cb)

	p := &ColorPickerDialog{
		dialog: d,
		color:  theme.PrimaryColor(),
	}

	p.picker = widget.NewColorAdvancedPicker(theme.PrimaryColor(), func(c color.Color) {
		p.color = c
	})

	p.advanced = widget.NewAccordionContainer(widget.NewAccordionItem("Advanced", p.picker))

	contents := []fyne.CanvasObject{
		basic,
		greyscale,
	}
	if !recent.MinSize().IsZero() {
		// Add divider and recents if there are any
		contents = append(contents, canvas.NewLine(theme.ShadowColor()), recent)
	}

	d.content = widget.NewVBox(append(contents, p.advanced)...)

	d.dismiss = &widget.Button{Text: "Cancel", Icon: theme.CancelIcon(),
		OnTapped: d.Hide,
	}
	confirm := &widget.Button{Text: "Confirm", Icon: theme.ConfirmIcon(), Style: widget.PrimaryButton,
		OnTapped: func() {
			cb(p.color)
		},
	}
	d.setButtons(newButtonList(d.dismiss, confirm))
	return p
}
