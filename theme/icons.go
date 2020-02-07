package theme

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne"
)

// ThemedResource is a resource wrapper that will return an appropriate resource
// for the currently selected theme.
type ThemedResource struct {
	source fyne.Resource
}

// Name returns the underlying resource name (used for caching)
func (res *ThemedResource) Name() string {
	return res.source.Name()
}

// Content returns the underlying content of the correct resource for the current theme
func (res *ThemedResource) Content() []byte {
	clr := fyne.CurrentApp().Settings().Theme().IconColor()
	return colorizeResource(res.source, clr)
}

// NewThemedResource creates a resource that adapts to the current theme setting.
//
// Deprecated: NewThemedResource() will be replaced with a single parameter version in a future release
// usage of this method will break, but using the first parameter only will be a trivial change.
func NewThemedResource(src, ignored fyne.Resource) *ThemedResource {
	if ignored != nil {
		log.Println("Deprecation Warning: In version 2.0 NewThemedResource() will only accept a single StaticResource.\n" +
			"While two resources are still supported to preserve backwards compatibility, only the first resource is rendered.\n" +
			"The resource color is set by the theme's IconColor().")
	}
	return &ThemedResource{
		source: src,
	}
}

// DisabledResource is a resource wrapper that will return an appropriate resource colorized by
// the current theme's `DisabledIconColor` color.
type DisabledResource struct {
	source fyne.Resource
}

// Name returns the resource source name prefixed with `disabled_` (used for caching)
func (res *DisabledResource) Name() string {
	return fmt.Sprintf("disabled_%s", res.source.Name())
}

// Content returns the disabled style content of the correct resource for the current theme
func (res *DisabledResource) Content() []byte {
	clr := fyne.CurrentApp().Settings().Theme().DisabledIconColor()
	return colorizeResource(res.source, clr)
}

// NewDisabledResource creates a resource that adapts to the current theme's DisabledIconColor setting.
func NewDisabledResource(res fyne.Resource) *DisabledResource {
	return &DisabledResource{
		source: res,
	}
}

func colorizeResource(res fyne.Resource, clr color.Color) []byte {
	rdr := bytes.NewReader(res.Content())
	s, err := svgFromXML(rdr)
	if err != nil {
		fyne.LogError("could not load SVG, falling back to static content:", err)
		return res.Content()
	}
	if err := s.replaceFillColor(rdr, clr); err != nil {
		fyne.LogError("could not replace fill color, falling back to static content:", err)
		return res.Content()
	}
	b, err := xml.Marshal(s)
	if err != nil {
		fyne.LogError("could not marshal svg, falling back to static content:", err)
		return res.Content()
	}
	return b
}

var (
	cancel, confirm, delete, search, searchReplace, menu, menuExpand            *ThemedResource
	checked, unchecked, radioButton, radioButtonChecked                         *ThemedResource
	contentAdd, contentRemove, contentCut, contentCopy, contentPaste            *ThemedResource
	contentRedo, contentUndo, info, question, warning                           *ThemedResource
	documentCreate, documentPrint, documentSave                                 *ThemedResource
	mailAttachment, mailCompose, mailForward, mailReply, mailReplyAll, mailSend *ThemedResource
	mediaFastForward, mediaFastRewind, mediaPause, mediaPlay                    *ThemedResource
	mediaRecord, mediaReplay, mediaSkipNext, mediaSkipPrevious                  *ThemedResource
	arrowBack, arrowDown, arrowForward, arrowUp, arrowDropDown, arrowDropUp     *ThemedResource
	folder, folderNew, folderOpen, help, home, settings                         *ThemedResource
	viewFullScreen, viewRefresh, viewZoomFit, viewZoomIn, viewZoomOut           *ThemedResource
	visibility, visibilityOff, volumeDown, volumeMute, volumeUp                 *ThemedResource
)

func init() {
	cancel = NewThemedResource(cancelIconRes, nil)
	confirm = NewThemedResource(checkIconRes, nil)
	delete = NewThemedResource(deleteIconRes, nil)
	search = NewThemedResource(searchIconRes, nil)
	searchReplace = NewThemedResource(searchreplaceIconRes, nil)
	menu = NewThemedResource(menuIconRes, nil)
	menuExpand = NewThemedResource(menuexpandIconRes, nil)

	checked = NewThemedResource(checkboxIconRes, nil)
	unchecked = NewThemedResource(checkboxblankIconRes, nil)
	radioButton = NewThemedResource(radiobuttonIconRes, nil)
	radioButtonChecked = NewThemedResource(radiobuttoncheckedIconRes, nil)

	contentAdd = NewThemedResource(contentaddIconRes, nil)
	contentRemove = NewThemedResource(contentremoveIconRes, nil)
	contentCut = NewThemedResource(contentcutIconRes, nil)
	contentCopy = NewThemedResource(contentcopyIconRes, nil)
	contentPaste = NewThemedResource(contentpasteIconRes, nil)
	contentRedo = NewThemedResource(contentredoIconRes, nil)
	contentUndo = NewThemedResource(contentundoIconRes, nil)

	documentCreate = NewThemedResource(documentcreateIconRes, nil)
	documentPrint = NewThemedResource(documentprintIconRes, nil)
	documentSave = NewThemedResource(documentsaveIconRes, nil)

	info = NewThemedResource(infoIconRes, nil)
	question = NewThemedResource(questionIconRes, nil)
	warning = NewThemedResource(warningIconRes, nil)

	mailAttachment = NewThemedResource(mailattachmentIconRes, nil)
	mailCompose = NewThemedResource(mailcomposeIconRes, nil)
	mailForward = NewThemedResource(mailforwardIconRes, nil)
	mailReply = NewThemedResource(mailreplyIconRes, nil)
	mailReplyAll = NewThemedResource(mailreplyallIconRes, nil)
	mailSend = NewThemedResource(mailsendIconRes, nil)

	mediaFastForward = NewThemedResource(mediafastforwardIconRes, nil)
	mediaFastRewind = NewThemedResource(mediafastrewindIconRes, nil)
	mediaPause = NewThemedResource(mediapauseIconRes, nil)
	mediaPlay = NewThemedResource(mediaplayIconRes, nil)
	mediaRecord = NewThemedResource(mediarecordIconRes, nil)
	mediaReplay = NewThemedResource(mediareplayIconRes, nil)
	mediaSkipNext = NewThemedResource(mediaskipnextIconRes, nil)
	mediaSkipPrevious = NewThemedResource(mediaskippreviousIconRes, nil)

	arrowBack = NewThemedResource(arrowbackIconRes, nil)
	arrowDown = NewThemedResource(arrowdownIconRes, nil)
	arrowForward = NewThemedResource(arrowforwardIconRes, nil)
	arrowUp = NewThemedResource(arrowupIconRes, nil)
	arrowDropDown = NewThemedResource(arrowdropdownIconRes, nil)
	arrowDropUp = NewThemedResource(arrowdropupIconRes, nil)

	folder = NewThemedResource(folderIconRes, nil)
	folderNew = NewThemedResource(foldernewIconRes, nil)
	folderOpen = NewThemedResource(folderopenIconRes, nil)
	help = NewThemedResource(helpIconRes, nil)
	home = NewThemedResource(homeIconRes, nil)
	settings = NewThemedResource(settingsIconRes, nil)

	viewFullScreen = NewThemedResource(viewfullscreenIconRes, nil)
	viewRefresh = NewThemedResource(viewrefreshIconRes, nil)
	viewZoomFit = NewThemedResource(viewzoomfitIconRes, nil)
	viewZoomIn = NewThemedResource(viewzoominIconRes, nil)
	viewZoomOut = NewThemedResource(viewzoomoutIconRes, nil)

	visibility = NewThemedResource(visibilityIconRes, nil)
	visibilityOff = NewThemedResource(visibilityoffIconRes, nil)

	volumeDown = NewThemedResource(volumedownIconRes, nil)
	volumeMute = NewThemedResource(volumemuteIconRes, nil)
	volumeUp = NewThemedResource(volumeupIconRes, nil)
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

// SearchReplaceIcon returns a resource containing the standard search and replace icon for the current theme
func SearchReplaceIcon() fyne.Resource {
	return searchReplace
}

// MenuIcon returns a resource containing the standard (mobile) menu icon for the current theme
func MenuIcon() fyne.Resource {
	return menu
}

// MenuExpandIcon returns a resource containing the standard (mobile) expand "submenu icon for the current theme
func MenuExpandIcon() fyne.Resource {
	return menuExpand
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

// ContentAddIcon returns a resource containing the standard content add icon for the current theme
func ContentAddIcon() fyne.Resource {
	return contentAdd
}

// ContentRemoveIcon returns a resource containing the standard content remove icon for the current theme
func ContentRemoveIcon() fyne.Resource {
	return contentRemove
}

// ContentClearIcon returns a resource containing the standard content clear icon for the current theme
func ContentClearIcon() fyne.Resource {
	return cancel
}

// ContentCutIcon returns a resource containing the standard content cut icon for the current theme
func ContentCutIcon() fyne.Resource {
	return contentCut
}

// ContentCopyIcon returns a resource containing the standard content copy icon for the current theme
func ContentCopyIcon() fyne.Resource {
	return contentCopy
}

// ContentPasteIcon returns a resource containing the standard content paste icon for the current theme
func ContentPasteIcon() fyne.Resource {
	return contentPaste
}

// ContentRedoIcon returns a resource containing the standard content redo icon for the current theme
func ContentRedoIcon() fyne.Resource {
	return contentRedo
}

// ContentUndoIcon returns a resource containing the standard content undo icon for the current theme
func ContentUndoIcon() fyne.Resource {
	return contentUndo
}

// DocumentCreateIcon returns a resource containing the standard document create icon for the current theme
func DocumentCreateIcon() fyne.Resource {
	return documentCreate
}

// DocumentPrintIcon returns a resource containing the standard document print icon for the current theme
func DocumentPrintIcon() fyne.Resource {
	return documentPrint
}

// DocumentSaveIcon returns a resource containing the standard document save icon for the current theme
func DocumentSaveIcon() fyne.Resource {
	return documentSave
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

// FolderOpenIcon returns a resource containing the standard folder open icon for the current theme
func FolderOpenIcon() fyne.Resource {
	return folderOpen
}

// HelpIcon returns a resource containing the standard help icon for the current theme
func HelpIcon() fyne.Resource {
	return help
}

// HomeIcon returns a resource containing the standard home folder icon for the current theme
func HomeIcon() fyne.Resource {
	return home
}

// SettingsIcon returns a resource containing the standard settings icon for the current theme
func SettingsIcon() fyne.Resource {
	return settings
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

// MediaFastForwardIcon returns a resource containing the standard media fast-forward icon for the current theme
func MediaFastForwardIcon() fyne.Resource {
	return mediaFastForward
}

// MediaFastRewindIcon returns a resource containing the standard media fast-rewind icon for the current theme
func MediaFastRewindIcon() fyne.Resource {
	return mediaFastRewind
}

// MediaPauseIcon returns a resource containing the standard media pause icon for the current theme
func MediaPauseIcon() fyne.Resource {
	return mediaPause
}

// MediaPlayIcon returns a resource containing the standard media play icon for the current theme
func MediaPlayIcon() fyne.Resource {
	return mediaPlay
}

// MediaRecordIcon returns a resource containing the standard media record icon for the current theme
func MediaRecordIcon() fyne.Resource {
	return mediaRecord
}

// MediaReplayIcon returns a resource containing the standard media replay icon for the current theme
func MediaReplayIcon() fyne.Resource {
	return mediaReplay
}

// MediaSkipNextIcon returns a resource containing the standard media skip next icon for the current theme
func MediaSkipNextIcon() fyne.Resource {
	return mediaSkipNext
}

// MediaSkipPreviousIcon returns a resource containing the standard media skip previous icon for the current theme
func MediaSkipPreviousIcon() fyne.Resource {
	return mediaSkipPrevious
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

// MenuDropDownIcon returns a resource containing the standard menu drop down icon for the current theme
func MenuDropDownIcon() fyne.Resource {
	return arrowDropDown
}

// MenuDropUpIcon returns a resource containing the standard menu drop up icon for the current theme
func MenuDropUpIcon() fyne.Resource {
	return arrowDropUp
}

// ViewFullScreenIcon returns a resource containing the standard fullscreen icon for the current theme
func ViewFullScreenIcon() fyne.Resource {
	return viewFullScreen
}

// ViewRestoreIcon returns a resource containing the standard exit fullscreen icon for the current theme
func ViewRestoreIcon() fyne.Resource {
	return viewZoomFit
}

// ViewRefreshIcon returns a resource containing the standard refresh icon for the current theme
func ViewRefreshIcon() fyne.Resource {
	return viewRefresh
}

// ZoomFitIcon returns a resource containing the standard zoom fit icon for the current theme
func ZoomFitIcon() fyne.Resource {
	return viewZoomFit
}

// ZoomInIcon returns a resource containing the standard zoom in icon for the current theme
func ZoomInIcon() fyne.Resource {
	return viewZoomIn
}

// ZoomOutIcon returns a resource containing the standard zoom out icon for the current theme
func ZoomOutIcon() fyne.Resource {
	return viewZoomOut
}

// VisibilityIcon returns a resource containing the standard visibity icon for the current theme
func VisibilityIcon() fyne.Resource {
	return visibility
}

// VisibilityOffIcon returns a resource containing the standard visibity off icon for the current theme
func VisibilityOffIcon() fyne.Resource {
	return visibilityOff
}

// VolumeDownIcon returns a resource containing the standard volume down icon for the current theme
func VolumeDownIcon() fyne.Resource {
	return volumeDown
}

// VolumeMuteIcon returns a resource containing the standard volume mute icon for the current theme
func VolumeMuteIcon() fyne.Resource {
	return volumeMute
}

// VolumeUpIcon returns a resource containing the standard volume up icon for the current theme
func VolumeUpIcon() fyne.Resource {
	return volumeUp
}
