package theme

import "github.com/fyne-io/fyne"

// CancelIcon returns a resource containing the standard cancel icon
func CancelIcon() *fyne.Resource {
	return cancel
}

// CheckIcon returns a resource containing the standard confirm icon
func ConfirmIcon() *fyne.Resource {
	return check
}
