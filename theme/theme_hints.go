//go:build hints
// +build hints

package theme

var (
	fallbackColor = errorColor
	fallbackIcon  = NewErrorThemedResource(errorIconRes)
)
