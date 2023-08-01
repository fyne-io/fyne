package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// DefaultEmojiFont returns the font resource for the built-in emoji font.
// This may return nil if the application was packaged without an emoji font.
//
// Since: 2.4
func DefaultEmojiFont() fyne.Resource {
	return emoji
}

// DefaultTextBoldFont returns the font resource for the built-in bold font style.
func DefaultTextBoldFont() fyne.Resource {
	return bold
}

// DefaultTextBoldItalicFont returns the font resource for the built-in bold and italic font style.
func DefaultTextBoldItalicFont() fyne.Resource {
	return bolditalic
}

// DefaultTextFont returns the font resource for the built-in regular font style.
func DefaultTextFont() fyne.Resource {
	return regular
}

// DefaultTextItalicFont returns the font resource for the built-in italic font style.
func DefaultTextItalicFont() fyne.Resource {
	return italic
}

// DefaultTextMonospaceFont returns the font resource for the built-in monospace font face.
func DefaultTextMonospaceFont() fyne.Resource {
	return monospace
}

// DefaultSymbolFont returns the font resource for the built-in symbol font.
//
// Since: 2.2
func DefaultSymbolFont() fyne.Resource {
	return symbol
}

// TextBoldFont returns the font resource for the bold font style.
func TextBoldFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Bold: true})
}

// TextBoldItalicFont returns the font resource for the bold and italic font style.
func TextBoldItalicFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Bold: true, Italic: true})
}

// TextColor returns the theme's standard text color - this is actually the foreground color since 1.4.
//
// Deprecated: Use theme.ForegroundColor() colour instead.
func TextColor() color.Color {
	return safeColorLookup(ColorNameForeground, currentVariant())
}

// TextFont returns the font resource for the regular font style.
func TextFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{})
}

// TextItalicFont returns the font resource for the italic font style.
func TextItalicFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Italic: true})
}

// TextMonospaceFont returns the font resource for the monospace font face.
func TextMonospaceFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Monospace: true})
}

// SymbolFont returns the font resource for the symbol font style.
//
// Since: 2.4
func SymbolFont() fyne.Resource {
	return safeFontLookup(fyne.TextStyle{Symbol: true})
}
func safeFontLookup(s fyne.TextStyle) fyne.Resource {
	font := current().Font(s)
	if font != nil {
		return font
	}
	fyne.LogError("Loaded theme returned nil font", nil)

	if s.Monospace {
		return DefaultTextMonospaceFont()
	}
	if s.Bold {
		if s.Italic {
			return DefaultTextBoldItalicFont()
		}
		return DefaultTextBoldFont()
	}
	if s.Italic {
		return DefaultTextItalicFont()
	}
	if s.Symbol {
		return DefaultSymbolFont()
	}

	return DefaultTextFont()
}
