package theme

import "github.com/fyne-io/fyne"

type darkLightResource struct {
	dark  *fyne.Resource
	light *fyne.Resource
}

func (res *darkLightResource) CurrentResource() *fyne.Resource {
	if fyne.GetSettings().Theme() == "light" {
		return res.light
	}

	return res.dark
}

var cancel, confirm, checked, unchecked fyne.ThemedResource

func init() {
	cancel = &darkLightResource{cancelDark, cancelLight}
	confirm = &darkLightResource{checkDark, checkLight}
	checked = &darkLightResource{checkboxDark, checkboxLight}
	unchecked = &darkLightResource{checkboxblankDark, checkboxblankLight}
}

// CancelIcon returns a resource containing the standard cancel icon for the current theme
func CancelIcon() fyne.ThemedResource {
	return cancel
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() fyne.ThemedResource {
	return confirm
}

// CheckedIcon returns a resource containing the standard checkbox icon for the current theme
func CheckedIcon() fyne.ThemedResource {
	return checked
}

// UncheckedIcon returns a resource containing the standard checkbox unchecked icon for the current theme
func UncheckedIcon() fyne.ThemedResource {
	return unchecked
}
