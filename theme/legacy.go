package theme

import (
	"image/color"

	"fyne.io/fyne"
)

// FromLegacy returns a 2.0 Theme created from the given LegacyTheme data.
//
// Since: 2.0.0
func FromLegacy(t fyne.LegacyTheme) fyne.Theme {
	return &legacyWrapper{old: t}
}

var _ fyne.Theme = (*legacyWrapper)(nil)

type legacyWrapper struct {
	old fyne.LegacyTheme
}

func (l *legacyWrapper) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch n {
	case Colors.Background:
		return l.old.BackgroundColor()
	case Colors.Text:
		return l.old.TextColor()
	case Colors.Button:
		return l.old.ButtonColor()
	case Colors.DisabledButton:
		return l.old.DisabledButtonColor()
	case Colors.DisabledText:
		return l.old.DisabledTextColor()
	case Colors.Focus:
		return l.old.FocusColor()
	case Colors.Hover:
		return l.old.HoverColor()
	case Colors.PlaceHolder:
		return l.old.PlaceHolderColor()
	case Colors.Primary:
		return l.old.PrimaryColor()
	case Colors.ScrollBar:
		return l.old.ScrollBarColor()
	case Colors.Shadow:
		return l.old.ShadowColor()
	default:
		return color.Transparent
	}
}

func (l *legacyWrapper) Size(n fyne.ThemeSizeName) int {
	switch n {
	case Sizes.InlineIcon:
		return l.old.IconInlineSize()
	case Sizes.Padding:
		return l.old.Padding()
	case Sizes.ScrollBar:
		return l.old.ScrollBarSize()
	case Sizes.ScrollBarSmall:
		return l.old.ScrollBarSmallSize()
	case Sizes.Text:
		return l.old.TextSize()
	default:
		return 0
	}
}

func (l *legacyWrapper) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return l.old.TextMonospaceFont()
	}
	if s.Bold {
		if s.Italic {
			return l.old.TextBoldItalicFont()
		}
		return l.old.TextBoldFont()
	}
	if s.Italic {
		return l.old.TextItalicFont()
	}
	return l.old.TextFont()
}
