package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// FromLegacy returns a 2.0 Theme created from the given LegacyTheme data.
// This is a transition path and will be removed in the future (probably version 3.0).
//
// Since: 2.0
func FromLegacy(t fyne.LegacyTheme) fyne.Theme {
	return &legacyWrapper{old: t}
}

var _ fyne.Theme = (*legacyWrapper)(nil)

type legacyWrapper struct {
	old fyne.LegacyTheme
}

func (l *legacyWrapper) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case ColorNameBackground:
		return l.old.BackgroundColor()
	case ColorNameForeground:
		return l.old.TextColor()
	case ColorNameButton:
		return l.old.ButtonColor()
	case ColorNameDisabledButton:
		return l.old.DisabledButtonColor()
	case ColorNameDisabled:
		return l.old.DisabledTextColor()
	case ColorNameFocus:
		return l.old.FocusColor()
	case ColorNameHover:
		return l.old.HoverColor()
	case ColorNamePlaceHolder:
		return l.old.PlaceHolderColor()
	case ColorNamePrimary:
		return l.old.PrimaryColor()
	case ColorNameScrollBar:
		return l.old.ScrollBarColor()
	case ColorNameShadow:
		return l.old.ShadowColor()
	default:
		return DefaultTheme().Color(n, v)
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

func (l *legacyWrapper) Icon(n fyne.ThemeIconName) fyne.Resource {
	return DefaultTheme().Icon(n)
}

func (l *legacyWrapper) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case SizeNameInlineIcon:
		return float32(l.old.IconInlineSize())
	case SizeNamePadding:
		return float32(l.old.Padding())
	case SizeNameScrollBar:
		return float32(l.old.ScrollBarSize())
	case SizeNameScrollBarSmall:
		return float32(l.old.ScrollBarSmallSize())
	case SizeNameText:
		return float32(l.old.TextSize())
	default:
		return DefaultTheme().Size(n)
	}
}
