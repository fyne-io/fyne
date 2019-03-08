package theme

import "fyne.io/fyne"

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

// NewThemedResource creates a resource that adapts to the current theme setting.
// It is currently a simple pairing of a dark and light variant of the same resource.
func NewThemedResource(dark, light fyne.Resource) *ThemedResource {
	return &ThemedResource{dark, light}
}

var (
	cancel, confirm, delete, search                                             *ThemedResource
	checked, unchecked, radioButton, radioButtonChecked                         *ThemedResource
	contentCut, contentCopy, contentPaste                                       *ThemedResource
	info, question, warning                                                     *ThemedResource
	mailAttachment, mailCompose, mailForward, mailReply, mailReplyAll, mailSend *ThemedResource
	arrowBack, arrowDown, arrowForward, arrowUp                                 *ThemedResource
	folder, folderNew, help, home                                               *ThemedResource
)

func init() {
	cancel = &ThemedResource{cancelDark, cancelLight}
	confirm = &ThemedResource{checkDark, checkLight}
	delete = &ThemedResource{deleteDark, deleteLight}
	search = &ThemedResource{searchDark, searchLight}

	checked = &ThemedResource{checkboxDark, checkboxLight}
	unchecked = &ThemedResource{checkboxblankDark, checkboxblankLight}
	radioButton = &ThemedResource{radiobuttonDark, radiobuttonLight}
	radioButtonChecked = &ThemedResource{radiobuttoncheckedDark, radiobuttoncheckedLight}

	contentCut = &ThemedResource{contentcutDark, contentcutLight}
	contentCopy = &ThemedResource{contentcopyDark, contentcopyLight}
	contentPaste = &ThemedResource{contentpasteDark, contentpasteLight}

	info = &ThemedResource{infoDark, infoLight}
	question = &ThemedResource{questionDark, questionLight}
	warning = &ThemedResource{warningDark, warningLight}

	mailAttachment = &ThemedResource{mailattachmentDark, mailattachmentLight}
	mailCompose = &ThemedResource{mailcomposeDark, mailcomposeLight}
	mailForward = &ThemedResource{mailforwardDark, mailforwardLight}
	mailReply = &ThemedResource{mailreplyDark, mailreplyLight}
	mailReplyAll = &ThemedResource{mailreplyallDark, mailreplyallLight}
	mailSend = &ThemedResource{mailsendDark, mailsendLight}

	arrowBack = &ThemedResource{arrowbackDark, arrowbackLight}
	arrowDown = &ThemedResource{arrowdownDark, arrowdownLight}
	arrowForward = &ThemedResource{arrowforwardDark, arrowforwardLight}
	arrowUp = &ThemedResource{arrowupDark, arrowupLight}

	folder = &ThemedResource{folderDark, folderLight}
	folderNew = &ThemedResource{foldernewDark, foldernewLight}
	help = &ThemedResource{helpDark, helpLight}
	home = &ThemedResource{homeDark, homeLight}
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

// SearchIcon returns a resource containing the standard search icon for the current theme
func SearchIcon() fyne.Resource {
	return search
}

// CheckButtonIcon returns a resource containing the standard checkbox icon for the current theme
func CheckButtonIcon() fyne.Resource {
	return unchecked
}

// CheckButtonCheckedIcon returns a resource containing the standard checkbox checked icon for the current theme
func CheckButtonCheckedIcon() fyne.Resource {
	return checked
}

// RadioButtonIcon returns a resource containing the standard radio button icon for the current theme
func RadioButtonIcon() fyne.Resource {
	return radioButton
}

// RadioButtonCheckedIcon returns a resource containing the standard radio button checked icon for the current theme
func RadioButtonCheckedIcon() fyne.Resource {
	return radioButtonChecked
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

// FolderIcon returns a resource containing the standard folder icon for the current theme
func FolderIcon() fyne.Resource {
	return folder
}

// FolderNewIcon returns a resource containing the standard folder creation icon for the current theme
func FolderNewIcon() fyne.Resource {
	return folderNew
}

// HelpIcon returns a resource containing the standard help icon for the current theme
func HelpIcon() fyne.Resource {
	return help
}

// HomeIcon returns a resource containing the standard home folder icon for the current theme
func HomeIcon() fyne.Resource {
	return home
}

// MailAttachmentIcon returns a resource containing the standard mail attachment icon for the current theme
func MailAttachmentIcon() fyne.Resource {
	return mailAttachment
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

// MailSendIcon returns a resource containing the standard mail send icon for the current theme
func MailSendIcon() fyne.Resource {
	return mailSend
}

// MoveDownIcon returns a resource containing the standard down arrow icon for the current theme
func MoveDownIcon() fyne.Resource {
	return arrowDown
}

// MoveUpIcon returns a resource containing the standard up arrow icon for the current theme
func MoveUpIcon() fyne.Resource {
	return arrowUp
}

// NavigateBackIcon returns a resource containing the standard backward navigation icon for the current theme
func NavigateBackIcon() fyne.Resource {
	return arrowBack
}

// NavigateNextIcon returns a resource containing the standard forward navigation icon for the current theme
func NavigateNextIcon() fyne.Resource {
	return arrowForward
}
