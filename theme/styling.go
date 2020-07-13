package theme

import (
	"image/color"

	"fyne.io/fyne"
)

// ButtonStyle is the colors used for Button drawing
// Each instance of a button refers to a ButtonStyle parameter with precalculated colors
// based ont the current palette. All colors can be overridden.
// THe user can create custom styles and apply them easily to different buttons.
type ButtonStyle struct {
	DisabledColor     color.Color
	EnabledColor      color.Color
	TextColor         color.Color
	DisabledTextColor color.Color
	ShadowColor       color.Color
	OutlineColor      color.Color // Set to transparent if not used
	Rounded           float64     // 1.0 is fully rounded, 0.0 is square
	TextSize          int
	Height            float64 // As a fraction of text size, typicaly 1.5 to 3.0
	MinWidth          int
	// Calculated state color modifications based on overlays
	Hovered color.Color
	Focused color.Color
	Pressed color.Color
}

// These styles are automatically calculated from the current Palette

// PrimaryRaisedButton ..
var PrimaryRaisedButton ButtonStyle

// PrimaryFlatButton ..
var PrimaryFlatButton ButtonStyle

// DefaultRaisedButton ..
var DefaultRaisedButton ButtonStyle

// DefaultFlatButton ..
var DefaultFlatButton ButtonStyle

// SecondaryRaisedButton ..
var SecondaryRaisedButton ButtonStyle

// SecondaryFlatButton ..
var SecondaryFlatButton ButtonStyle

var transparent = color.NRGBA{0, 0, 0, 0}

// Lighter is true when c1 is lighter than c2
func Lighter(c1, c2 color.Color) bool {
	r, g, b, _ := c1.RGBA()
	s1 := r + g + b
	r, g, b, _ = c2.RGBA()
	s2 := r + g + b
	return s1 > s2
}

// AddStateColors will add the Hovered, Pressed and Focused colors based on Text color and enabled color
func (b *ButtonStyle) AddStateColors() {
	if Lighter(b.EnabledColor, b.TextColor) {
		b.Hovered = Blend(b.EnabledColor, b.TextColor, 0.08)
		b.Pressed = Blend(b.EnabledColor, b.TextColor, 0.16)
		b.Focused = Blend(b.EnabledColor, b.TextColor, 0.16)
	} else {
		b.Hovered = Blend(b.EnabledColor, b.TextColor, 0.16)
		b.Pressed = Blend(b.EnabledColor, b.TextColor, 0.32)
		b.Focused = Blend(b.EnabledColor, b.TextColor, 0.32)
	}
}

func (b *ButtonStyle) updateText(tc color.Color, dt color.Color, ts int) {
	b.TextColor = tc
	b.DisabledTextColor = dt
	if ts > 0 {
		b.TextSize = ts
	}
	if b.TextSize == 0 {
		b.TextSize = DarkTheme().TextSize()
	}
}

func (b *ButtonStyle) updateColors(ec color.Color, dc color.Color, sc color.Color) {
	b.EnabledColor = ec
	b.DisabledColor = dc
	b.ShadowColor = sc
	b.AddStateColors()
}

// UpdateDefaultStyling will implement the theme by updating the global default styling structs
// Must be called when the theme is changed
func UpdateDefaultStyling(s fyne.Theme) {
	DefaultRaisedButton.updateText(s.TextColor(), s.DisabledTextColor(), s.TextSize())
	DefaultRaisedButton.updateColors(s.ButtonColor(), s.DisabledButtonColor(), s.ShadowColor())
	DefaultRaisedButton.Hovered = s.HoverColor()

	DefaultFlatButton.updateText(s.TextColor(), s.DisabledTextColor(), s.TextSize())
	DefaultFlatButton.updateColors(s.BackgroundColor(), s.DisabledButtonColor(), nil)

	if t, ok := s.(*builtinTheme); ok && t.OnPrimary != transparent {
		PrimaryRaisedButton.updateText(t.OnPrimary, s.DisabledTextColor(), s.TextSize())
	} else {
		PrimaryRaisedButton.updateText(s.TextColor(), s.DisabledTextColor(), s.TextSize())
	}
	PrimaryRaisedButton.updateColors(s.PrimaryColor(), s.DisabledButtonColor(), s.ShadowColor())

	PrimaryFlatButton.updateText(s.PrimaryColor(), s.DisabledTextColor(), s.TextSize())
	PrimaryFlatButton.updateColors(s.BackgroundColor(), s.DisabledButtonColor(), nil)

	if t, ok := s.(*builtinTheme); ok && t.OnSecondary != transparent {
		SecondaryRaisedButton.updateText(t.OnSecondary, s.DisabledTextColor(), s.TextSize())
	} else {
		SecondaryRaisedButton.updateText(s.TextColor(), s.DisabledTextColor(), s.TextSize())
	}
	if t, ok := s.(*builtinTheme); ok && t.Secondary != transparent {
		SecondaryRaisedButton.updateColors(t.Secondary, s.DisabledButtonColor(), s.ShadowColor())
	} else {
		SecondaryRaisedButton.updateColors(color.RGBA{250, 64, 129, 255}, s.DisabledButtonColor(), s.ShadowColor())
	}
	SecondaryFlatButton.updateText(color.RGBA{250, 64, 129, 255}, s.DisabledTextColor(), s.TextSize())
	SecondaryFlatButton.updateColors(s.BackgroundColor(), s.DisabledButtonColor(), nil)
}
