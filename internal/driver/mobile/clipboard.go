package mobile

import (
	"fyne.io/fyne/v2"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = (*mobileClipboard)(nil)

var MobileClipboard mobileClipboard

// mobileClipboard represents the system mobileClipboard
type mobileClipboard struct {
}
