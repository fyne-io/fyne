package theme

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/color"

	"fyne.io/fyne"
)

// ThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// for the currently selected theme.
type ThemedResource struct {
	source fyne.Resource
}

// NewThemedResource creates a resource that adapts to the current theme setting.
func NewThemedResource(src fyne.Resource) *ThemedResource {
	return &ThemedResource{
		source: src,
	}
}

// Name returns the underlying resource name (used for caching).
func (res *ThemedResource) Name() string {
	return res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current text color.
func (res *ThemedResource) Content() []byte {
	clr := TextColor()
	return colorizeResource(res.source, clr)
}

// Error resturns a different resource for indicating an error.
func (res *ThemedResource) Error() *ErrorThemedResource {
	return NewErrorThemedResource(res)
}

// InvertedThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// for use over highlighted elements.
type InvertedThemedResource struct {
	source fyne.Resource
}

// NewInvertedThemedResource creates a resource that adapts to the current theme for use over highlighted elements.
func NewInvertedThemedResource(orig fyne.Resource) *InvertedThemedResource {
	res := &InvertedThemedResource{source: orig}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *InvertedThemedResource) Name() string {
	return "inverted-" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *InvertedThemedResource) Content() []byte {
	clr := BackgroundColor()
	return colorizeResource(res.source, clr)
}

// Original returns the underlying resource that this inverted themed resource was adapted from
func (res *InvertedThemedResource) Original() fyne.Resource {
	return res.source
}

// ErrorThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// to indicate an error.
type ErrorThemedResource struct {
	source fyne.Resource
}

// NewErrorThemedResource creates a resource that adapts to the error color for the current theme.
func NewErrorThemedResource(orig fyne.Resource) *ErrorThemedResource {
	res := &ErrorThemedResource{source: orig}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *ErrorThemedResource) Name() string {
	return "error-" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *ErrorThemedResource) Content() []byte {
	clr := &color.NRGBA{0xf4, 0x43, 0x36, 0xff} // TODO: Should be current().ErrorColor() in the future
	return colorizeResource(res.source, clr)
}

// Original returns the underlying resource that this error themed resource was adapted from
func (res *ErrorThemedResource) Original() fyne.Resource {
	return res.source
}

// PrimaryThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// to the theme primary color.
type PrimaryThemedResource struct {
	source fyne.Resource
}

// NewPrimaryThemedResource creates a resource that adapts to the primary color for the current theme.
func NewPrimaryThemedResource(orig fyne.Resource) *PrimaryThemedResource {
	res := &PrimaryThemedResource{source: orig}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *PrimaryThemedResource) Name() string {
	return "primary-" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *PrimaryThemedResource) Content() []byte {
	return colorizeResource(res.source, PrimaryColor())
}

// Original returns the underlying resource that this primary themed resource was adapted from
func (res *PrimaryThemedResource) Original() fyne.Resource {
	return res.source
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
	clr := DisabledTextColor()
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
	if err := s.replaceFillColor(clr); err != nil {
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
	cancel, confirm, delete, search, searchReplace, menu, menuExpand                *ThemedResource
	checked, unchecked, radioButton, radioButtonChecked                             *ThemedResource
	colorAchromatic, colorChromatic, colorPalette                                   *ThemedResource
	contentAdd, contentRemove, contentCut, contentCopy, contentPaste                *ThemedResource
	contentClear, contentRedo, contentUndo, info, question, warning, errori         *ThemedResource
	document, documentCreate, documentPrint, documentSave                           *ThemedResource
	mailAttachment, mailCompose, mailForward, mailReply, mailReplyAll, mailSend     *ThemedResource
	mediaFastForward, mediaFastRewind, mediaPause, mediaPlay                        *ThemedResource
	mediaRecord, mediaReplay, mediaSkipNext, mediaSkipPrevious, mediaStop           *ThemedResource
	arrowBack, arrowDown, arrowForward, arrowUp, arrowDropDown, arrowDropUp         *ThemedResource
	file, fileApplication, fileAudio, fileImage, fileText, fileVideo                *ThemedResource
	folder, folderNew, folderOpen, help, history, home, settings, storage, upload   *ThemedResource
	viewFullScreen, viewRefresh, viewZoomFit, viewZoomIn, viewZoomOut, viewRestore  *ThemedResource
	visibility, visibilityOff, volumeDown, volumeMute, volumeUp, download, computer *ThemedResource
)

func init() {
	cancel = NewThemedResource(cancelIconRes)
	confirm = NewThemedResource(checkIconRes)
	delete = NewThemedResource(deleteIconRes)
	search = NewThemedResource(searchIconRes)
	searchReplace = NewThemedResource(searchreplaceIconRes)
	menu = NewThemedResource(menuIconRes)
	menuExpand = NewThemedResource(menuexpandIconRes)

	checked = NewThemedResource(checkboxIconRes)
	unchecked = NewThemedResource(checkboxblankIconRes)
	radioButton = NewThemedResource(radiobuttonIconRes)
	radioButtonChecked = NewThemedResource(radiobuttoncheckedIconRes)

	contentAdd = NewThemedResource(contentaddIconRes)
	contentClear = NewThemedResource(cancelIconRes)
	contentRemove = NewThemedResource(contentremoveIconRes)
	contentCut = NewThemedResource(contentcutIconRes)
	contentCopy = NewThemedResource(contentcopyIconRes)
	contentPaste = NewThemedResource(contentpasteIconRes)
	contentRedo = NewThemedResource(contentredoIconRes)
	contentUndo = NewThemedResource(contentundoIconRes)

	colorAchromatic = NewThemedResource(colorachromaticIconRes)
	colorChromatic = NewThemedResource(colorchromaticIconRes)
	colorPalette = NewThemedResource(colorpaletteIconRes)

	document = NewThemedResource(documentIconRes)
	documentCreate = NewThemedResource(documentcreateIconRes)
	documentPrint = NewThemedResource(documentprintIconRes)
	documentSave = NewThemedResource(documentsaveIconRes)

	info = NewThemedResource(infoIconRes)
	question = NewThemedResource(questionIconRes)
	warning = NewThemedResource(warningIconRes)
	errori = NewThemedResource(errorIconRes)

	mailAttachment = NewThemedResource(mailattachmentIconRes)
	mailCompose = NewThemedResource(mailcomposeIconRes)
	mailForward = NewThemedResource(mailforwardIconRes)
	mailReply = NewThemedResource(mailreplyIconRes)
	mailReplyAll = NewThemedResource(mailreplyallIconRes)
	mailSend = NewThemedResource(mailsendIconRes)

	mediaFastForward = NewThemedResource(mediafastforwardIconRes)
	mediaFastRewind = NewThemedResource(mediafastrewindIconRes)
	mediaPause = NewThemedResource(mediapauseIconRes)
	mediaPlay = NewThemedResource(mediaplayIconRes)
	mediaRecord = NewThemedResource(mediarecordIconRes)
	mediaReplay = NewThemedResource(mediareplayIconRes)
	mediaSkipNext = NewThemedResource(mediaskipnextIconRes)
	mediaSkipPrevious = NewThemedResource(mediaskippreviousIconRes)
	mediaStop = NewThemedResource(mediastopIconRes)

	arrowBack = NewThemedResource(arrowbackIconRes)
	arrowDown = NewThemedResource(arrowdownIconRes)
	arrowForward = NewThemedResource(arrowforwardIconRes)
	arrowUp = NewThemedResource(arrowupIconRes)
	arrowDropDown = NewThemedResource(arrowdropdownIconRes)
	arrowDropUp = NewThemedResource(arrowdropupIconRes)

	file = NewThemedResource(fileIconRes)
	fileApplication = NewThemedResource(fileapplicationIconRes)
	fileAudio = NewThemedResource(fileaudioIconRes)
	fileImage = NewThemedResource(fileimageIconRes)
	fileText = NewThemedResource(filetextIconRes)
	fileVideo = NewThemedResource(filevideoIconRes)
	folder = NewThemedResource(folderIconRes)
	folderNew = NewThemedResource(foldernewIconRes)
	folderOpen = NewThemedResource(folderopenIconRes)
	help = NewThemedResource(helpIconRes)
	history = NewThemedResource(historyIconRes)
	home = NewThemedResource(homeIconRes)
	settings = NewThemedResource(settingsIconRes)

	viewFullScreen = NewThemedResource(viewfullscreenIconRes)
	viewRefresh = NewThemedResource(viewrefreshIconRes)
	viewRestore = NewThemedResource(viewzoomfitIconRes)
	viewZoomFit = NewThemedResource(viewzoomfitIconRes)
	viewZoomIn = NewThemedResource(viewzoominIconRes)
	viewZoomOut = NewThemedResource(viewzoomoutIconRes)

	visibility = NewThemedResource(visibilityIconRes)
	visibilityOff = NewThemedResource(visibilityoffIconRes)

	volumeDown = NewThemedResource(volumedownIconRes)
	volumeMute = NewThemedResource(volumemuteIconRes)
	volumeUp = NewThemedResource(volumeupIconRes)

	download = NewThemedResource(downloadIconRes)
	computer = NewThemedResource(computerIconRes)
	storage = NewThemedResource(storageIconRes)
	upload = NewThemedResource(uploadIconRes)
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
	return contentClear
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

// ColorAchromaticIcon returns a resource containing the standard achromatic color icon for the current theme
func ColorAchromaticIcon() fyne.Resource {
	return colorAchromatic
}

// ColorChromaticIcon returns a resource containing the standard chromatic color icon for the current theme
func ColorChromaticIcon() fyne.Resource {
	return colorChromatic
}

// ColorPaletteIcon returns a resource containing the standard color palette icon for the current theme
func ColorPaletteIcon() fyne.Resource {
	return colorPalette
}

// DocumentIcon returns a resource containing the standard document icon for the current theme
func DocumentIcon() fyne.Resource {
	return document
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

// ErrorIcon returns a resource containing the standard dialog error icon for the current theme
func ErrorIcon() fyne.Resource {
	return errori
}

// FileIcon returns a resource containing the appropriate file icon for the current theme
func FileIcon() fyne.Resource {
	return file
}

// FileApplicationIcon returns a resource containing the file icon representing application files for the current theme
func FileApplicationIcon() fyne.Resource {
	return fileApplication
}

// FileAudioIcon returns a resource containing the file icon representing audio files for the current theme
func FileAudioIcon() fyne.Resource {
	return fileAudio
}

// FileImageIcon returns a resource containing the file icon representing image files for the current theme
func FileImageIcon() fyne.Resource {
	return fileImage
}

// FileTextIcon returns a resource containing the file icon representing text files for the current theme
func FileTextIcon() fyne.Resource {
	return fileText
}

// FileVideoIcon returns a resource containing the file icon representing video files for the current theme
func FileVideoIcon() fyne.Resource {
	return fileVideo
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

// HistoryIcon returns a resource containing the standard history icon for the current theme
func HistoryIcon() fyne.Resource {
	return history
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

// MediaStopIcon returns a resource containing the standard media stop icon for the current theme
func MediaStopIcon() fyne.Resource {
	return mediaStop
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
	return viewRestore
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

// ComputerIcon returns a resource containing the standard computer icon for the current theme
func ComputerIcon() fyne.Resource {
	return computer
}

// DownloadIcon returns a resource containing the standard download icon for the current theme
func DownloadIcon() fyne.Resource {
	return download
}

// StorageIcon returns a resource containing the standard storage icon for the current theme
func StorageIcon() fyne.Resource {
	return storage
}

// UploadIcon returns a resource containing the standard upload icon for the current theme
func UploadIcon() fyne.Resource {
	return upload
}
