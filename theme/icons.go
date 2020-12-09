package theme

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/color"

	"fyne.io/fyne"
)

var (
	// Icons specifies each of the known icon names that a theme can contain.
	//
	// Since 2.0.0
	Icons = struct {
		Cancel, Confirm, Delete, Search, SearchReplace, Menu, MenuExpand                fyne.ThemeIconName
		CheckButtonChecked, CheckButton, RadioButton, RadioButtonChecked                fyne.ThemeIconName
		ColorAchromatic, ColorChromatic, ColorPalette                                   fyne.ThemeIconName
		ContentAdd, ContentRemove, ContentCut, ContentCopy, ContentPaste                fyne.ThemeIconName
		ContentClear, ContentRedo, ContentUndo, Info, Question, Warning, Error          fyne.ThemeIconName
		Document, DocumentCreate, DocumentPrint, DocumentSave                           fyne.ThemeIconName
		MailAttachment, MailCompose, MailForward, MailReply, MailReplyAll, MailSend     fyne.ThemeIconName
		MediaFastForward, MediaFastRewind, MediaPause, MediaPlay                        fyne.ThemeIconName
		MediaRecord, MediaReplay, MediaSkipNext, MediaSkipPrevious, MediaStop           fyne.ThemeIconName
		NavigateBack, MoveDown, NavigateNext, MoveUp, ArrowDropDown, ArrowDropUp        fyne.ThemeIconName
		File, FileApplication, FileAudio, FileImage, FileText, FileVideo                fyne.ThemeIconName
		Folder, FolderNew, FolderOpen, Help, History, Home, Settings, Storage, Upload   fyne.ThemeIconName
		ViewFullScreen, ViewRefresh, ViewZoomFit, ViewZoomIn, ViewZoomOut, ViewRestore  fyne.ThemeIconName
		Visibility, VisibilityOff, VolumeDown, VolumeMute, VolumeUp, Download, Computer fyne.ThemeIconName
	}{
		"cancel", "confirm", "delete", "search", "searchReplace", "menu", "menuExpand",
		"checked", "unchecked", "radioButton", "radioButtonChecked",
		"colorAchromatic", "colorChromatic", "colorPalette",
		"contentAdd", "contentRemove", "contentCut", "contentCopy", "contentPaste",
		"contentClear", "contentRedo", "contentUndo", "info", "question", "warning", "errori",
		"document", "documentCreate", "documentPrint", "documentSave",
		"mailAttachment", "mailCompose", "mailForward", "mailReply", "mailReplyAll", "mailSend",
		"mediaFastForward", "mediaFastRewind", "mediaPause", "mediaPlay",
		"mediaRecord", "mediaReplay", "mediaSkipNext", "mediaSkipPrevious", "mediaStop",
		"arrowBack", "arrowDown", "arrowForward", "arrowUp", "arrowDropDown", "arrowDropUp",
		"file", "fileApplication", "fileAudio", "fileImage", "fileText", "fileVideo",
		"folder", "folderNew", "folderOpen", "help", "history", "home", "settings", "storage", "upload",
		"viewFullScreen", "viewRefresh", "viewZoomFit", "viewZoomIn", "viewZoomOut", "viewRestore",
		"visibility", "visibilityOff", "volumeDown", "volumeMute", "volumeUp", "download", "computer",
	}

	icons = map[fyne.ThemeIconName]fyne.Resource{
		Icons.Cancel:        NewThemedResource(cancelIconRes),
		Icons.Confirm:       NewThemedResource(checkIconRes),
		Icons.Delete:        NewThemedResource(deleteIconRes),
		Icons.Search:        NewThemedResource(searchIconRes),
		Icons.SearchReplace: NewThemedResource(searchreplaceIconRes),
		Icons.Menu:          NewThemedResource(menuIconRes),
		Icons.MenuExpand:    NewThemedResource(menuexpandIconRes),

		Icons.CheckButton:        NewThemedResource(checkboxblankIconRes),
		Icons.CheckButtonChecked: NewThemedResource(checkboxIconRes),
		Icons.RadioButton:        NewThemedResource(radiobuttonIconRes),
		Icons.RadioButtonChecked: NewThemedResource(radiobuttoncheckedIconRes),

		Icons.ContentAdd:    NewThemedResource(contentaddIconRes),
		Icons.ContentClear:  NewThemedResource(cancelIconRes),
		Icons.ContentRemove: NewThemedResource(contentremoveIconRes),
		Icons.ContentCut:    NewThemedResource(contentcutIconRes),
		Icons.ContentCopy:   NewThemedResource(contentcopyIconRes),
		Icons.ContentPaste:  NewThemedResource(contentpasteIconRes),
		Icons.ContentRedo:   NewThemedResource(contentredoIconRes),
		Icons.ContentUndo:   NewThemedResource(contentundoIconRes),

		Icons.ColorAchromatic: NewThemedResource(colorachromaticIconRes),
		Icons.ColorChromatic:  NewThemedResource(colorchromaticIconRes),
		Icons.ColorPalette:    NewThemedResource(colorpaletteIconRes),

		Icons.Document:       NewThemedResource(documentIconRes),
		Icons.DocumentCreate: NewThemedResource(documentcreateIconRes),
		Icons.DocumentPrint:  NewThemedResource(documentprintIconRes),
		Icons.DocumentSave:   NewThemedResource(documentsaveIconRes),

		Icons.Info:     NewThemedResource(infoIconRes),
		Icons.Question: NewThemedResource(questionIconRes),
		Icons.Warning:  NewThemedResource(warningIconRes),
		Icons.Error:    NewThemedResource(errorIconRes),

		Icons.MailAttachment: NewThemedResource(mailattachmentIconRes),
		Icons.MailCompose:    NewThemedResource(mailcomposeIconRes),
		Icons.MailForward:    NewThemedResource(mailforwardIconRes),
		Icons.MailReply:      NewThemedResource(mailreplyIconRes),
		Icons.MailReplyAll:   NewThemedResource(mailreplyallIconRes),
		Icons.MailSend:       NewThemedResource(mailsendIconRes),

		Icons.MediaFastForward:  NewThemedResource(mediafastforwardIconRes),
		Icons.MediaFastRewind:   NewThemedResource(mediafastrewindIconRes),
		Icons.MediaPause:        NewThemedResource(mediapauseIconRes),
		Icons.MediaPlay:         NewThemedResource(mediaplayIconRes),
		Icons.MediaRecord:       NewThemedResource(mediarecordIconRes),
		Icons.MediaReplay:       NewThemedResource(mediareplayIconRes),
		Icons.MediaSkipNext:     NewThemedResource(mediaskipnextIconRes),
		Icons.MediaSkipPrevious: NewThemedResource(mediaskippreviousIconRes),
		Icons.MediaStop:         NewThemedResource(mediastopIconRes),

		Icons.NavigateBack:  NewThemedResource(arrowbackIconRes),
		Icons.MoveDown:      NewThemedResource(arrowdownIconRes),
		Icons.NavigateNext:  NewThemedResource(arrowforwardIconRes),
		Icons.MoveUp:        NewThemedResource(arrowupIconRes),
		Icons.ArrowDropDown: NewThemedResource(arrowdropdownIconRes),
		Icons.ArrowDropUp:   NewThemedResource(arrowdropupIconRes),

		Icons.File:            NewThemedResource(fileIconRes),
		Icons.FileApplication: NewThemedResource(fileapplicationIconRes),
		Icons.FileAudio:       NewThemedResource(fileaudioIconRes),
		Icons.FileImage:       NewThemedResource(fileimageIconRes),
		Icons.FileText:        NewThemedResource(filetextIconRes),
		Icons.FileVideo:       NewThemedResource(filevideoIconRes),
		Icons.Folder:          NewThemedResource(folderIconRes),
		Icons.FolderNew:       NewThemedResource(foldernewIconRes),
		Icons.FolderOpen:      NewThemedResource(folderopenIconRes),
		Icons.Help:            NewThemedResource(helpIconRes),
		Icons.History:         NewThemedResource(historyIconRes),
		Icons.Home:            NewThemedResource(homeIconRes),
		Icons.Settings:        NewThemedResource(settingsIconRes),

		Icons.ViewFullScreen: NewThemedResource(viewfullscreenIconRes),
		Icons.ViewRefresh:    NewThemedResource(viewrefreshIconRes),
		Icons.ViewRestore:    NewThemedResource(viewzoomfitIconRes),
		Icons.ViewZoomFit:    NewThemedResource(viewzoomfitIconRes),
		Icons.ViewZoomIn:     NewThemedResource(viewzoominIconRes),
		Icons.ViewZoomOut:    NewThemedResource(viewzoomoutIconRes),

		Icons.Visibility:    NewThemedResource(visibilityIconRes),
		Icons.VisibilityOff: NewThemedResource(visibilityoffIconRes),

		Icons.VolumeDown: NewThemedResource(volumedownIconRes),
		Icons.VolumeMute: NewThemedResource(volumemuteIconRes),
		Icons.VolumeUp:   NewThemedResource(volumeupIconRes),

		Icons.Download: NewThemedResource(downloadIconRes),
		Icons.Computer: NewThemedResource(computerIconRes),
		Icons.Storage:  NewThemedResource(storageIconRes),
		Icons.Upload:   NewThemedResource(uploadIconRes),
	}
)

func (t *builtinTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return icons[n]
}

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

// FyneLogo returns a resource containing the Fyne logo
func FyneLogo() fyne.Resource {
	return fynelogo
}

// CancelIcon returns a resource containing the standard cancel icon for the current theme
func CancelIcon() fyne.Resource {
	return current().Icon(Icons.Cancel)
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() fyne.Resource {
	return current().Icon(Icons.Confirm)
}

// DeleteIcon returns a resource containing the standard delete icon for the current theme
func DeleteIcon() fyne.Resource {
	return current().Icon(Icons.Delete)
}

// SearchIcon returns a resource containing the standard search icon for the current theme
func SearchIcon() fyne.Resource {
	return current().Icon(Icons.Search)
}

// SearchReplaceIcon returns a resource containing the standard search and replace icon for the current theme
func SearchReplaceIcon() fyne.Resource {
	return current().Icon(Icons.SearchReplace)
}

// MenuIcon returns a resource containing the standard (mobile) menu icon for the current theme
func MenuIcon() fyne.Resource {
	return current().Icon(Icons.Menu)
}

// MenuExpandIcon returns a resource containing the standard (mobile) expand "submenu icon for the current theme
func MenuExpandIcon() fyne.Resource {
	return current().Icon(Icons.MenuExpand)
}

// CheckButtonIcon returns a resource containing the standard checkbox icon for the current theme
func CheckButtonIcon() fyne.Resource {
	return current().Icon(Icons.CheckButton)
}

// CheckButtonCheckedIcon returns a resource containing the standard checkbox checked icon for the current theme
func CheckButtonCheckedIcon() fyne.Resource {
	return current().Icon(Icons.CheckButtonChecked)
}

// RadioButtonIcon returns a resource containing the standard radio button icon for the current theme
func RadioButtonIcon() fyne.Resource {
	return current().Icon(Icons.RadioButton)
}

// RadioButtonCheckedIcon returns a resource containing the standard radio button checked icon for the current theme
func RadioButtonCheckedIcon() fyne.Resource {
	return current().Icon(Icons.RadioButtonChecked)
}

// ContentAddIcon returns a resource containing the standard content add icon for the current theme
func ContentAddIcon() fyne.Resource {
	return current().Icon(Icons.ContentAdd)
}

// ContentRemoveIcon returns a resource containing the standard content remove icon for the current theme
func ContentRemoveIcon() fyne.Resource {
	return current().Icon(Icons.ContentRemove)
}

// ContentClearIcon returns a resource containing the standard content clear icon for the current theme
func ContentClearIcon() fyne.Resource {
	return current().Icon(Icons.ContentClear)
}

// ContentCutIcon returns a resource containing the standard content cut icon for the current theme
func ContentCutIcon() fyne.Resource {
	return current().Icon(Icons.ContentCut)
}

// ContentCopyIcon returns a resource containing the standard content copy icon for the current theme
func ContentCopyIcon() fyne.Resource {
	return current().Icon(Icons.ContentCopy)
}

// ContentPasteIcon returns a resource containing the standard content paste icon for the current theme
func ContentPasteIcon() fyne.Resource {
	return current().Icon(Icons.ContentPaste)
}

// ContentRedoIcon returns a resource containing the standard content redo icon for the current theme
func ContentRedoIcon() fyne.Resource {
	return current().Icon(Icons.ContentRedo)
}

// ContentUndoIcon returns a resource containing the standard content undo icon for the current theme
func ContentUndoIcon() fyne.Resource {
	return current().Icon(Icons.ContentUndo)
}

// ColorAchromaticIcon returns a resource containing the standard achromatic color icon for the current theme
func ColorAchromaticIcon() fyne.Resource {
	return current().Icon(Icons.ColorAchromatic)
}

// ColorChromaticIcon returns a resource containing the standard chromatic color icon for the current theme
func ColorChromaticIcon() fyne.Resource {
	return current().Icon(Icons.ColorChromatic)
}

// ColorPaletteIcon returns a resource containing the standard color palette icon for the current theme
func ColorPaletteIcon() fyne.Resource {
	return current().Icon(Icons.ColorPalette)
}

// DocumentIcon returns a resource containing the standard document icon for the current theme
func DocumentIcon() fyne.Resource {
	return current().Icon(Icons.Document)
}

// DocumentCreateIcon returns a resource containing the standard document create icon for the current theme
func DocumentCreateIcon() fyne.Resource {
	return current().Icon(Icons.DocumentCreate)
}

// DocumentPrintIcon returns a resource containing the standard document print icon for the current theme
func DocumentPrintIcon() fyne.Resource {
	return current().Icon(Icons.DocumentPrint)
}

// DocumentSaveIcon returns a resource containing the standard document save icon for the current theme
func DocumentSaveIcon() fyne.Resource {
	return current().Icon(Icons.DocumentSave)
}

// InfoIcon returns a resource containing the standard dialog info icon for the current theme
func InfoIcon() fyne.Resource {
	return current().Icon(Icons.Info)
}

// QuestionIcon returns a resource containing the standard dialog question icon for the current theme
func QuestionIcon() fyne.Resource {
	return current().Icon(Icons.Question)
}

// WarningIcon returns a resource containing the standard dialog warning icon for the current theme
func WarningIcon() fyne.Resource {
	return current().Icon(Icons.Warning)
}

// ErrorIcon returns a resource containing the standard dialog error icon for the current theme
func ErrorIcon() fyne.Resource {
	return current().Icon(Icons.Error)
}

// FileIcon returns a resource containing the appropriate file icon for the current theme
func FileIcon() fyne.Resource {
	return current().Icon(Icons.File)
}

// FileApplicationIcon returns a resource containing the file icon representing application files for the current theme
func FileApplicationIcon() fyne.Resource {
	return current().Icon(Icons.FileApplication)
}

// FileAudioIcon returns a resource containing the file icon representing audio files for the current theme
func FileAudioIcon() fyne.Resource {
	return current().Icon(Icons.FileAudio)
}

// FileImageIcon returns a resource containing the file icon representing image files for the current theme
func FileImageIcon() fyne.Resource {
	return current().Icon(Icons.FileImage)
}

// FileTextIcon returns a resource containing the file icon representing text files for the current theme
func FileTextIcon() fyne.Resource {
	return current().Icon(Icons.FileText)
}

// FileVideoIcon returns a resource containing the file icon representing video files for the current theme
func FileVideoIcon() fyne.Resource {
	return current().Icon(Icons.FileVideo)
}

// FolderIcon returns a resource containing the standard folder icon for the current theme
func FolderIcon() fyne.Resource {
	return current().Icon(Icons.Folder)
}

// FolderNewIcon returns a resource containing the standard folder creation icon for the current theme
func FolderNewIcon() fyne.Resource {
	return current().Icon(Icons.FolderNew)
}

// FolderOpenIcon returns a resource containing the standard folder open icon for the current theme
func FolderOpenIcon() fyne.Resource {
	return current().Icon(Icons.FolderOpen)
}

// HelpIcon returns a resource containing the standard help icon for the current theme
func HelpIcon() fyne.Resource {
	return current().Icon(Icons.Help)
}

// HistoryIcon returns a resource containing the standard history icon for the current theme
func HistoryIcon() fyne.Resource {
	return current().Icon(Icons.History)
}

// HomeIcon returns a resource containing the standard home folder icon for the current theme
func HomeIcon() fyne.Resource {
	return current().Icon(Icons.Home)
}

// SettingsIcon returns a resource containing the standard settings icon for the current theme
func SettingsIcon() fyne.Resource {
	return current().Icon(Icons.Settings)
}

// MailAttachmentIcon returns a resource containing the standard mail attachment icon for the current theme
func MailAttachmentIcon() fyne.Resource {
	return current().Icon(Icons.MailAttachment)
}

// MailComposeIcon returns a resource containing the standard mail compose icon for the current theme
func MailComposeIcon() fyne.Resource {
	return current().Icon(Icons.MailCompose)
}

// MailForwardIcon returns a resource containing the standard mail forward icon for the current theme
func MailForwardIcon() fyne.Resource {
	return current().Icon(Icons.MailForward)
}

// MailReplyIcon returns a resource containing the standard mail reply icon for the current theme
func MailReplyIcon() fyne.Resource {
	return current().Icon(Icons.MailReply)
}

// MailReplyAllIcon returns a resource containing the standard mail reply all icon for the current theme
func MailReplyAllIcon() fyne.Resource {
	return current().Icon(Icons.MailReplyAll)
}

// MailSendIcon returns a resource containing the standard mail send icon for the current theme
func MailSendIcon() fyne.Resource {
	return current().Icon(Icons.MailSend)
}

// MediaFastForwardIcon returns a resource containing the standard media fast-forward icon for the current theme
func MediaFastForwardIcon() fyne.Resource {
	return current().Icon(Icons.MediaFastForward)
}

// MediaFastRewindIcon returns a resource containing the standard media fast-rewind icon for the current theme
func MediaFastRewindIcon() fyne.Resource {
	return current().Icon(Icons.MediaFastRewind)
}

// MediaPauseIcon returns a resource containing the standard media pause icon for the current theme
func MediaPauseIcon() fyne.Resource {
	return current().Icon(Icons.MediaPause)
}

// MediaPlayIcon returns a resource containing the standard media play icon for the current theme
func MediaPlayIcon() fyne.Resource {
	return current().Icon(Icons.MediaPlay)
}

// MediaRecordIcon returns a resource containing the standard media record icon for the current theme
func MediaRecordIcon() fyne.Resource {
	return current().Icon(Icons.MediaRecord)
}

// MediaReplayIcon returns a resource containing the standard media replay icon for the current theme
func MediaReplayIcon() fyne.Resource {
	return current().Icon(Icons.MediaReplay)
}

// MediaSkipNextIcon returns a resource containing the standard media skip next icon for the current theme
func MediaSkipNextIcon() fyne.Resource {
	return current().Icon(Icons.MediaSkipNext)
}

// MediaSkipPreviousIcon returns a resource containing the standard media skip previous icon for the current theme
func MediaSkipPreviousIcon() fyne.Resource {
	return current().Icon(Icons.MediaSkipPrevious)
}

// MediaStopIcon returns a resource containing the standard media stop icon for the current theme
func MediaStopIcon() fyne.Resource {
	return current().Icon(Icons.MediaStop)
}

// MoveDownIcon returns a resource containing the standard down arrow icon for the current theme
func MoveDownIcon() fyne.Resource {
	return current().Icon(Icons.MoveDown)
}

// MoveUpIcon returns a resource containing the standard up arrow icon for the current theme
func MoveUpIcon() fyne.Resource {
	return current().Icon(Icons.MoveUp)
}

// NavigateBackIcon returns a resource containing the standard backward navigation icon for the current theme
func NavigateBackIcon() fyne.Resource {
	return current().Icon(Icons.NavigateBack)
}

// NavigateNextIcon returns a resource containing the standard forward navigation icon for the current theme
func NavigateNextIcon() fyne.Resource {
	return current().Icon(Icons.NavigateNext)
}

// MenuDropDownIcon returns a resource containing the standard menu drop down icon for the current theme
func MenuDropDownIcon() fyne.Resource {
	return current().Icon(Icons.ArrowDropDown)
}

// MenuDropUpIcon returns a resource containing the standard menu drop up icon for the current theme
func MenuDropUpIcon() fyne.Resource {
	return current().Icon(Icons.ArrowDropUp)
}

// ViewFullScreenIcon returns a resource containing the standard fullscreen icon for the current theme
func ViewFullScreenIcon() fyne.Resource {
	return current().Icon(Icons.ViewFullScreen)
}

// ViewRestoreIcon returns a resource containing the standard exit fullscreen icon for the current theme
func ViewRestoreIcon() fyne.Resource {
	return current().Icon(Icons.ViewRestore)
}

// ViewRefreshIcon returns a resource containing the standard refresh icon for the current theme
func ViewRefreshIcon() fyne.Resource {
	return current().Icon(Icons.ViewRefresh)
}

// ZoomFitIcon returns a resource containing the standard zoom fit icon for the current theme
func ZoomFitIcon() fyne.Resource {
	return current().Icon(Icons.ViewZoomFit)
}

// ZoomInIcon returns a resource containing the standard zoom in icon for the current theme
func ZoomInIcon() fyne.Resource {
	return current().Icon(Icons.ViewZoomIn)
}

// ZoomOutIcon returns a resource containing the standard zoom out icon for the current theme
func ZoomOutIcon() fyne.Resource {
	return current().Icon(Icons.ViewZoomOut)
}

// VisibilityIcon returns a resource containing the standard visibity icon for the current theme
func VisibilityIcon() fyne.Resource {
	return current().Icon(Icons.Visibility)
}

// VisibilityOffIcon returns a resource containing the standard visibity off icon for the current theme
func VisibilityOffIcon() fyne.Resource {
	return current().Icon(Icons.VisibilityOff)
}

// VolumeDownIcon returns a resource containing the standard volume down icon for the current theme
func VolumeDownIcon() fyne.Resource {
	return current().Icon(Icons.VolumeDown)
}

// VolumeMuteIcon returns a resource containing the standard volume mute icon for the current theme
func VolumeMuteIcon() fyne.Resource {
	return current().Icon(Icons.VolumeMute)
}

// VolumeUpIcon returns a resource containing the standard volume up icon for the current theme
func VolumeUpIcon() fyne.Resource {
	return current().Icon(Icons.VolumeUp)
}

// ComputerIcon returns a resource containing the standard computer icon for the current theme
func ComputerIcon() fyne.Resource {
	return current().Icon(Icons.Computer)
}

// DownloadIcon returns a resource containing the standard download icon for the current theme
func DownloadIcon() fyne.Resource {
	return current().Icon(Icons.Download)
}

// StorageIcon returns a resource containing the standard storage icon for the current theme
func StorageIcon() fyne.Resource {
	return current().Icon(Icons.Storage)
}

// UploadIcon returns a resource containing the standard upload icon for the current theme
func UploadIcon() fyne.Resource {
	return current().Icon(Icons.Upload)
}
