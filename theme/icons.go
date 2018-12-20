package theme

import "github.com/fyne-io/fyne"

// ThemedResource is a resource wrapper that will return an appropriate resource
// for the currently selected theme. In this implementation it chooses between a dark
// and light alternative to match the current setting.
type ThemedResource struct {
	dark  fyne.Resource
	light fyne.Resource
}

func isThemeLight() bool {
	return fyne.CurrentApp().Settings().Theme().BackgroundColor() == lightBackground
}

// Name returns the underlrying resource name (used for caching)
func (res *ThemedResource) Name() string {
	if isThemeLight() {
		return res.light.Name()
	}

	return res.dark.Name()
}

// Content returns the underlying content of the correct resource for the current theme
func (res *ThemedResource) Content() []byte {
	if isThemeLight() {
		return res.light.Content()
	}

	return res.dark.Content()
}

// CachePath returns the cachepath of the correct resource for the current theme setting
func (res *ThemedResource) CachePath() string {
	if isThemeLight() {
		return res.light.CachePath()
	}

	return res.dark.CachePath()
}

// NewThemedResource creates a resource that adapts to the current theme setting.
// It is currently a simple pairing of a dark and light variant of the same resource.
func NewThemedResource(dark, light fyne.Resource) *ThemedResource {
	return &ThemedResource{dark, light}
}

var cancel, confirm, delete, checked, unchecked *ThemedResource
var contentCut, contentCopy, contentPaste *ThemedResource
var info, question, warning *ThemedResource
var mailCompose, mailForward, mailReply, mailReplyAll *ThemedResource
var arrowBack, arrowForward fyne.Resource

func init() {
	cancel = &ThemedResource{cancelDark, cancelLight}
	confirm = &ThemedResource{checkDark, checkLight}
	delete = &ThemedResource{deleteDark, deleteLight}
	checked = &ThemedResource{checkboxDark, checkboxLight}
	unchecked = &ThemedResource{checkboxblankDark, checkboxblankLight}

	contentCut = &ThemedResource{contentcutDark, contentcutLight}
	contentCopy = &ThemedResource{contentcopyDark, contentcopyLight}
	contentPaste = &ThemedResource{contentpasteDark, contentpasteLight}

	info = &ThemedResource{infoDark, infoLight}
	question = &ThemedResource{questionDark, questionLight}
	warning = &ThemedResource{warningDark, warningLight}

	mailCompose = &ThemedResource{mailcomposeDark, mailcomposeLight}
	mailForward = &ThemedResource{mailforwardDark, mailforwardLight}
	mailReply = &ThemedResource{mailreplyDark, mailreplyLight}
	mailReplyAll = &ThemedResource{mailreplyallDark, mailreplyallLight}

	arrowBack = &ThemedResource{arrowbackDark, arrowbackLight}
	arrowForward = &ThemedResource{arrowforwardDark, arrowforwardLight}
}

// FyneLogo returns a resource containing the Fyne logo
func FyneLogo() fyne.Resource {
	return fynelogo
}

// CancelIcon returns a resource containing the standard cancel icon for the current theme
func CancelIcon() fyne.Resource {
	return cancel
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() fyne.Resource {
	return confirm
}

// DeleteIcon returns a resource containing the standard delete icon for the current theme
func DeleteIcon() fyne.Resource {
	return delete
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

// InfoIcon returns a resource containing the standard dialog info icon for the current theme
func InfoIcon() fyne.Resource {
	return info
}

// QuestionIcon returns a resource containing the standard dialog question icon for the current theme
func QuestionIcon() fyne.Resource {
	return question
}

// WarningIcon returns a resource containing the standard dialog warning icon for the current theme
func WarningIcon() fyne.Resource {
	return warning
}

// MailComposeIcon returns a resource containing the standard mail compose icon for the current theme
func MailComposeIcon() fyne.Resource {
	return mailCompose
}

// MailForwardIcon returns a resource containing the standard mail forward icon for the current theme
func MailForwardIcon() fyne.Resource {
	return mailForward
}

// MailReplyIcon returns a resource containing the standard mail reply icon for the current theme
func MailReplyIcon() fyne.Resource {
	return mailReply
}

// MailReplyAllIcon returns a resource containing the standard mail reply all icon for the current theme
func MailReplyAllIcon() fyne.Resource {
	return mailReplyAll
}

// NavigateBackIcon returns a resource containing the standard backward navigation icon for the current theme
func NavigateBackIcon() fyne.Resource {
	return arrowBack
}

// NavigateNextIcon returns a resource containing the standard forward navigation icon for the current theme
func NavigateNextIcon() fyne.Resource {
	return arrowForward
}
