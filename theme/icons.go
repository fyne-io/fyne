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

var cancel, confirm fyne.ThemedResource

func init() {
	cancel = &darkLightResource{cancelDark, cancelLight}
	confirm = &darkLightResource{checkDark, checkLight}
}

// CancelIcon returns a resource containing the standard cancel icon
func CancelIcon() fyne.ThemedResource {
	return cancel
}

// ConfirmIcon returns a resource containing the standard confirm icon
func ConfirmIcon() fyne.ThemedResource {
	return confirm
}
