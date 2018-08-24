package theme

import "github.com/fyne-io/fyne"

type darkLightResource struct {
	dark  *fyne.StaticResource
	light *fyne.StaticResource
}

func (res *darkLightResource) Name() string {
	if fyne.GetSettings().Theme() == "light" {
		return res.light.StaticName
	}

	return res.dark.StaticName
}

func (res *darkLightResource) Content() []byte {
	if fyne.GetSettings().Theme() == "light" {
		return res.light.StaticContent
	}

	return res.dark.StaticContent
}

func (res *darkLightResource) CachePath() string {
	if fyne.GetSettings().Theme() == "light" {
		return res.light.CachePath()
	}

	return res.dark.CachePath()
}

var cancel, confirm, checked, unchecked *darkLightResource
var contentCut, contentCopy, contentPaste *darkLightResource

func init() {
	cancel = &darkLightResource{cancelDark, cancelLight}
	confirm = &darkLightResource{checkDark, checkLight}
	checked = &darkLightResource{checkboxDark, checkboxLight}
	unchecked = &darkLightResource{checkboxblankDark, checkboxblankLight}

	contentCut = &darkLightResource{contentcutDark, contentcutLight}
	contentCopy = &darkLightResource{contentcopyDark, contentcopyLight}
	contentPaste = &darkLightResource{contentpasteDark, contentpasteLight}
}

// CancelIcon returns a resource containing the standard cancel icon for the current theme
func CancelIcon() fyne.Resource {
	return cancel
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() fyne.Resource {
	return confirm
}

// CheckedIcon returns a resource containing the standard checkbox icon for the current theme
func CheckedIcon() fyne.Resource {
	return checked
}

// UncheckedIcon returns a resource containing the standard checkbox unchecked icon for the current theme
func UncheckedIcon() fyne.Resource {
	return unchecked
}

// CutIcon returns a resource containing the standard content cut icon for the current theme
func CutIcon() fyne.Resource {
	return contentCut
}

// CopyIcon returns a resource containing the standard content copy icon for the current theme
func CopyIcon() fyne.Resource {
	return contentCopy
}

// PasteIcon returns a resource containing the standard content paste icon for the current theme
func PasteIcon() fyne.Resource {
	return contentPaste
}
