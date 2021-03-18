package gomobile

import (
	"fyne.io/fyne/v2"
)

// Override Name on android for content://
type mobileURI struct {
	systemURI string
	fyne.URI
}
