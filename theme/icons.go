package theme

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/color"

	"fyne.io/fyne"
)

const (
	// IconNameCancel is the name of theme lookup for cancel icon.
	//
	// Since 2.0.0
	IconNameCancel fyne.ThemeIconName = "cancel"

	// IconNameConfirm is the name of theme lookup for confirm icon.
	//
	// Since 2.0.0
	IconNameConfirm fyne.ThemeIconName = "confirm"

	// IconNameDelete is the name of theme lookup for delete icon.
	//
	// Since 2.0.0
	IconNameDelete fyne.ThemeIconName = "delete"

	// IconNameSearch is the name of theme lookup for search icon.
	//
	// Since 2.0.0
	IconNameSearch fyne.ThemeIconName = "search"

	// IconNameSearchReplace is the name of theme lookup for search and replace icon.
	//
	// Since 2.0.0
	IconNameSearchReplace fyne.ThemeIconName = "searchReplace"

	// IconNameMenu is the name of theme lookup for menu icon.
	//
	// Since 2.0.0
	IconNameMenu fyne.ThemeIconName = "menu"

	// IconNameMenuExpand is the name of theme lookup for menu expansion icon.
	//
	// Since 2.0.0
	IconNameMenuExpand fyne.ThemeIconName = "menuExpand"

	// IconNameCheckButtonChecked is the name of theme lookup for checked check button icon.
	//
	// Since 2.0.0
	IconNameCheckButtonChecked fyne.ThemeIconName = "checked"

	// IconNameCheckButton is the name of theme lookup for  unchecked check button icon.
	//
	// Since 2.0.0
	IconNameCheckButton fyne.ThemeIconName = "unchecked"

	// IconNameRadioButton is the name of theme lookup for radio button unchecked icon.
	//
	// Since 2.0.0
	IconNameRadioButton fyne.ThemeIconName = "radioButton"

	// IconNameRadioButtonChecked is the name of theme lookup for radio button checked icon.
	//
	// Since 2.0.0
	IconNameRadioButtonChecked fyne.ThemeIconName = "radioButtonChecked"

	// IconNameColorAchromatic is the name of theme lookup for greyscale color icon.
	//
	// Since 2.0.0
	IconNameColorAchromatic fyne.ThemeIconName = "colorAchromatic"

	// IconNameColorChromatic is the name of theme lookup for full color icon.
	//
	// Since 2.0.0
	IconNameColorChromatic fyne.ThemeIconName = "colorChromatic"

	// IconNameColorPalette is the name of theme lookup for color palette icon.
	//
	// Since 2.0.0
	IconNameColorPalette fyne.ThemeIconName = "colorPalette"

	// IconNameContentAdd is the name of theme lookup for content add icon.
	//
	// Since 2.0.0
	IconNameContentAdd fyne.ThemeIconName = "contentAdd"

	// IconNameContentRemove is the name of theme lookup for content remove icon.
	//
	// Since 2.0.0
	IconNameContentRemove fyne.ThemeIconName = "contentRemove"

	// IconNameContentCut is the name of theme lookup for content cut icon.
	//
	// Since 2.0.0
	IconNameContentCut fyne.ThemeIconName = "contentCut"

	// IconNameContentCopy is the name of theme lookup for content copy icon.
	//
	// Since 2.0.0
	IconNameContentCopy fyne.ThemeIconName = "contentCopy"

	// IconNameContentPaste is the name of theme lookup for content paste icon.
	//
	// Since 2.0.0
	IconNameContentPaste fyne.ThemeIconName = "contentPaste"

	// IconNameContentClear is the name of theme lookup for content clear icon.
	//
	// Since 2.0.0
	IconNameContentClear fyne.ThemeIconName = "contentClear"

	// IconNameContentRedo is the name of theme lookup for content redo icon.
	//
	// Since 2.0.0
	IconNameContentRedo fyne.ThemeIconName = "contentRedo"

	// IconNameContentUndo is the name of theme lookup for content undo icon.
	//
	// Since 2.0.0
	IconNameContentUndo fyne.ThemeIconName = "contentUndo"

	// IconNameInfo is the name of theme lookup for info icon.
	//
	// Since 2.0.0
	IconNameInfo fyne.ThemeIconName = "info"

	// IconNameQuestion is the name of theme lookup for question icon.
	//
	// Since 2.0.0
	IconNameQuestion fyne.ThemeIconName = "question"

	// IconNameWarning is the name of theme lookup for warning icon.
	//
	// Since 2.0.0
	IconNameWarning fyne.ThemeIconName = "warning"

	// IconNameError is the name of theme lookup for error icon.
	//
	// Since 2.0.0
	IconNameError fyne.ThemeIconName = "error"

	// IconNameDocument is the name of theme lookup for document icon.
	//
	// Since 2.0.0
	IconNameDocument fyne.ThemeIconName = "document"

	// IconNameDocumentCreate is the name of theme lookup for document create icon.
	//
	// Since 2.0.0
	IconNameDocumentCreate fyne.ThemeIconName = "documentCreate"

	// IconNameDocumentPrint is the name of theme lookup for document print icon.
	//
	// Since 2.0.0
	IconNameDocumentPrint fyne.ThemeIconName = "documentPrint"

	// IconNameDocumentSave is the name of theme lookup for document save icon.
	//
	// Since 2.0.0
	IconNameDocumentSave fyne.ThemeIconName = "documentSave"

	// IconNameMailAttachment is the name of theme lookup for mail attachment icon.
	//
	// Since 2.0.0
	IconNameMailAttachment fyne.ThemeIconName = "mailAttachment"

	// IconNameMailCompose is the name of theme lookup for mail compose icon.
	//
	// Since 2.0.0
	IconNameMailCompose fyne.ThemeIconName = "mailCompose"

	// IconNameMailForward is the name of theme lookup for mail forward icon.
	//
	// Since 2.0.0
	IconNameMailForward fyne.ThemeIconName = "mailForward"

	// IconNameMailReply is the name of theme lookup for mail reply icon.
	//
	// Since 2.0.0
	IconNameMailReply fyne.ThemeIconName = "mailReply"

	// IconNameMailReplyAll is the name of theme lookup for mail reply-all icon.
	//
	// Since 2.0.0
	IconNameMailReplyAll fyne.ThemeIconName = "mailReplyAll"

	// IconNameMailSend is the name of theme lookup for mail send icon.
	//
	// Since 2.0.0
	IconNameMailSend fyne.ThemeIconName = "mailSend"

	// IconNameMediaFastForward is the name of theme lookup for media fast-forward icon.
	//
	// Since 2.0.0
	IconNameMediaFastForward fyne.ThemeIconName = "mediaFastForward"

	// IconNameMediaFastRewind is the name of theme lookup for media fast-rewind icon.
	//
	// Since 2.0.0
	IconNameMediaFastRewind fyne.ThemeIconName = "mediaFastRewind"

	// IconNameMediaPause is the name of theme lookup for media pause icon.
	//
	// Since 2.0.0
	IconNameMediaPause fyne.ThemeIconName = "mediaPause"

	// IconNameMediaPlay is the name of theme lookup for media play icon.
	//
	// Since 2.0.0
	IconNameMediaPlay fyne.ThemeIconName = "mediaPlay"

	// IconNameMediaRecord is the name of theme lookup for media record icon.
	//
	// Since 2.0.0
	IconNameMediaRecord fyne.ThemeIconName = "mediaRecord"

	// IconNameMediaReplay is the name of theme lookup for media replay icon.
	//
	// Since 2.0.0
	IconNameMediaReplay fyne.ThemeIconName = "mediaReplay"

	// IconNameMediaSkipNext is the name of theme lookup for media skip next icon.
	//
	// Since 2.0.0
	IconNameMediaSkipNext fyne.ThemeIconName = "mediaSkipNext"

	// IconNameMediaSkipPrevious is the name of theme lookup for media skip previous icon.
	//
	// Since 2.0.0
	IconNameMediaSkipPrevious fyne.ThemeIconName = "mediaSkipPrevious"

	// IconNameMediaStop is the name of theme lookup for media stop icon.
	//
	// Since 2.0.0
	IconNameMediaStop fyne.ThemeIconName = "mediaStop"

	// IconNameMoveDown is the name of theme lookup for move down icon.
	//
	// Since 2.0.0
	IconNameMoveDown fyne.ThemeIconName = "arrowDown"

	// IconNameMoveUp is the name of theme lookup for move up icon.
	//
	// Since 2.0.0
	IconNameMoveUp fyne.ThemeIconName = "arrowUp"

	// IconNameNavigateBack is the name of theme lookup for navigate back icon.
	//
	// Since 2.0.0
	IconNameNavigateBack fyne.ThemeIconName = "arrowBack"

	// IconNameNavigateNext is the name of theme lookup for navigate next icon.
	//
	// Since 2.0.0
	IconNameNavigateNext fyne.ThemeIconName = "arrowForward"

	// IconNameArrowDropDown is the name of theme lookup for drop-down arrow icon.
	//
	// Since 2.0.0
	IconNameArrowDropDown fyne.ThemeIconName = "arrowDropDown"

	// IconNameArrowDropUp is the name of theme lookup for drop-up arrow icon.
	//
	// Since 2.0.0
	IconNameArrowDropUp fyne.ThemeIconName = "arrowDropUp"

	// IconNameFile is the name of theme lookup for file icon.
	//
	// Since 2.0.0
	IconNameFile fyne.ThemeIconName = "file"

	// IconNameFileApplication is the name of theme lookup for file application icon.
	//
	// Since 2.0.0
	IconNameFileApplication fyne.ThemeIconName = "fileApplication"

	// IconNameFileAudio is the name of theme lookup for file audio icon.
	//
	// Since 2.0.0
	IconNameFileAudio fyne.ThemeIconName = "fileAudio"

	// IconNameFileImage is the name of theme lookup for file image icon.
	//
	// Since 2.0.0
	IconNameFileImage fyne.ThemeIconName = "fileImage"

	// IconNameFileText is the name of theme lookup for file text icon.
	//
	// Since 2.0.0
	IconNameFileText fyne.ThemeIconName = "fileText"

	// IconNameFileVideo is the name of theme lookup for file video icon.
	//
	// Since 2.0.0
	IconNameFileVideo fyne.ThemeIconName = "fileVideo"

	// IconNameFolder is the name of theme lookup for folder icon.
	//
	// Since 2.0.0
	IconNameFolder fyne.ThemeIconName = "folder"

	// IconNameFolderNew is the name of theme lookup for folder new icon.
	//
	// Since 2.0.0
	IconNameFolderNew fyne.ThemeIconName = "folderNew"

	// IconNameFolderOpen is the name of theme lookup for folder open icon.
	//
	// Since 2.0.0
	IconNameFolderOpen fyne.ThemeIconName = "folderOpen"

	// IconNameHelp is the name of theme lookup for help icon.
	//
	// Since 2.0.0
	IconNameHelp fyne.ThemeIconName = "help"

	// IconNameHistory is the name of theme lookup for history icon.
	//
	// Since 2.0.0
	IconNameHistory fyne.ThemeIconName = "history"

	// IconNameHome is the name of theme lookup for home icon.
	//
	// Since 2.0.0
	IconNameHome fyne.ThemeIconName = "home"

	// IconNameSettings is the name of theme lookup for settings icon.
	//
	// Since 2.0.0
	IconNameSettings fyne.ThemeIconName = "settings"

	// IconNameStorage is the name of theme lookup for storage icon.
	//
	// Since 2.0.0
	IconNameStorage fyne.ThemeIconName = "storage"

	// IconNameUpload is the name of theme lookup for upload icon.
	//
	// Since 2.0.0
	IconNameUpload fyne.ThemeIconName = "upload"

	// IconNameViewFullScreen is the name of theme lookup for view fullscreen icon.
	//
	// Since 2.0.0
	IconNameViewFullScreen fyne.ThemeIconName = "viewFullScreen"

	// IconNameViewRefresh is the name of theme lookup for view refresh icon.
	//
	// Since 2.0.0
	IconNameViewRefresh fyne.ThemeIconName = "viewRefresh"

	// IconNameViewZoomFit is the name of theme lookup for view zoom fit icon.
	//
	// Since 2.0.0
	IconNameViewZoomFit fyne.ThemeIconName = "viewZoomFit"

	// IconNameViewZoomIn is the name of theme lookup for view zoom in icon.
	//
	// Since 2.0.0
	IconNameViewZoomIn fyne.ThemeIconName = "viewZoomIn"

	// IconNameViewZoomOut is the name of theme lookup for view zoom out icon.
	//
	// Since 2.0.0
	IconNameViewZoomOut fyne.ThemeIconName = "viewZoomOut"

	// IconNameViewRestore is the name of theme lookup for view restore icon.
	//
	// Since 2.0.0
	IconNameViewRestore fyne.ThemeIconName = "viewRestore"

	// IconNameVisibility is the name of theme lookup for visibility icon.
	//
	// Since 2.0.0
	IconNameVisibility fyne.ThemeIconName = "visibility"

	// IconNameVisibilityOff is the name of theme lookup for invisibility icon.
	//
	// Since 2.0.0
	IconNameVisibilityOff fyne.ThemeIconName = "visibilityOff"

	// IconNameVolumeDown is the name of theme lookup for volume down icon.
	//
	// Since 2.0.0
	IconNameVolumeDown fyne.ThemeIconName = "volumeDown"

	// IconNameVolumeMute is the name of theme lookup for volume mute icon.
	//
	// Since 2.0.0
	IconNameVolumeMute fyne.ThemeIconName = "volumeMute"

	// IconNameVolumeUp is the name of theme lookup for volume up icon.
	//
	// Since 2.0.0
	IconNameVolumeUp fyne.ThemeIconName = "volumeUp"

	// IconNameDownload is the name of theme lookup for download icon.
	//
	// Since 2.0.0
	IconNameDownload fyne.ThemeIconName = "download"

	// IconNameComputer is the name of theme lookup for computer icon.
	//
	// Since 2.0.0
	IconNameComputer fyne.ThemeIconName = "computer"
)

var (
	icons = map[fyne.ThemeIconName]fyne.Resource{
		IconNameCancel:        NewThemedResource(cancelIconRes),
		IconNameConfirm:       NewThemedResource(checkIconRes),
		IconNameDelete:        NewThemedResource(deleteIconRes),
		IconNameSearch:        NewThemedResource(searchIconRes),
		IconNameSearchReplace: NewThemedResource(searchreplaceIconRes),
		IconNameMenu:          NewThemedResource(menuIconRes),
		IconNameMenuExpand:    NewThemedResource(menuexpandIconRes),

		IconNameCheckButton:        NewThemedResource(checkboxblankIconRes),
		IconNameCheckButtonChecked: NewThemedResource(checkboxIconRes),
		IconNameRadioButton:        NewThemedResource(radiobuttonIconRes),
		IconNameRadioButtonChecked: NewThemedResource(radiobuttoncheckedIconRes),

		IconNameContentAdd:    NewThemedResource(contentaddIconRes),
		IconNameContentClear:  NewThemedResource(cancelIconRes),
		IconNameContentRemove: NewThemedResource(contentremoveIconRes),
		IconNameContentCut:    NewThemedResource(contentcutIconRes),
		IconNameContentCopy:   NewThemedResource(contentcopyIconRes),
		IconNameContentPaste:  NewThemedResource(contentpasteIconRes),
		IconNameContentRedo:   NewThemedResource(contentredoIconRes),
		IconNameContentUndo:   NewThemedResource(contentundoIconRes),

		IconNameColorAchromatic: NewThemedResource(colorachromaticIconRes),
		IconNameColorChromatic:  NewThemedResource(colorchromaticIconRes),
		IconNameColorPalette:    NewThemedResource(colorpaletteIconRes),

		IconNameDocument:       NewThemedResource(documentIconRes),
		IconNameDocumentCreate: NewThemedResource(documentcreateIconRes),
		IconNameDocumentPrint:  NewThemedResource(documentprintIconRes),
		IconNameDocumentSave:   NewThemedResource(documentsaveIconRes),

		IconNameInfo:     NewThemedResource(infoIconRes),
		IconNameQuestion: NewThemedResource(questionIconRes),
		IconNameWarning:  NewThemedResource(warningIconRes),
		IconNameError:    NewThemedResource(errorIconRes),

		IconNameMailAttachment: NewThemedResource(mailattachmentIconRes),
		IconNameMailCompose:    NewThemedResource(mailcomposeIconRes),
		IconNameMailForward:    NewThemedResource(mailforwardIconRes),
		IconNameMailReply:      NewThemedResource(mailreplyIconRes),
		IconNameMailReplyAll:   NewThemedResource(mailreplyallIconRes),
		IconNameMailSend:       NewThemedResource(mailsendIconRes),

		IconNameMediaFastForward:  NewThemedResource(mediafastforwardIconRes),
		IconNameMediaFastRewind:   NewThemedResource(mediafastrewindIconRes),
		IconNameMediaPause:        NewThemedResource(mediapauseIconRes),
		IconNameMediaPlay:         NewThemedResource(mediaplayIconRes),
		IconNameMediaRecord:       NewThemedResource(mediarecordIconRes),
		IconNameMediaReplay:       NewThemedResource(mediareplayIconRes),
		IconNameMediaSkipNext:     NewThemedResource(mediaskipnextIconRes),
		IconNameMediaSkipPrevious: NewThemedResource(mediaskippreviousIconRes),
		IconNameMediaStop:         NewThemedResource(mediastopIconRes),

		IconNameNavigateBack:  NewThemedResource(arrowbackIconRes),
		IconNameMoveDown:      NewThemedResource(arrowdownIconRes),
		IconNameNavigateNext:  NewThemedResource(arrowforwardIconRes),
		IconNameMoveUp:        NewThemedResource(arrowupIconRes),
		IconNameArrowDropDown: NewThemedResource(arrowdropdownIconRes),
		IconNameArrowDropUp:   NewThemedResource(arrowdropupIconRes),

		IconNameFile:            NewThemedResource(fileIconRes),
		IconNameFileApplication: NewThemedResource(fileapplicationIconRes),
		IconNameFileAudio:       NewThemedResource(fileaudioIconRes),
		IconNameFileImage:       NewThemedResource(fileimageIconRes),
		IconNameFileText:        NewThemedResource(filetextIconRes),
		IconNameFileVideo:       NewThemedResource(filevideoIconRes),
		IconNameFolder:          NewThemedResource(folderIconRes),
		IconNameFolderNew:       NewThemedResource(foldernewIconRes),
		IconNameFolderOpen:      NewThemedResource(folderopenIconRes),
		IconNameHelp:            NewThemedResource(helpIconRes),
		IconNameHistory:         NewThemedResource(historyIconRes),
		IconNameHome:            NewThemedResource(homeIconRes),
		IconNameSettings:        NewThemedResource(settingsIconRes),

		IconNameViewFullScreen: NewThemedResource(viewfullscreenIconRes),
		IconNameViewRefresh:    NewThemedResource(viewrefreshIconRes),
		IconNameViewRestore:    NewThemedResource(viewzoomfitIconRes),
		IconNameViewZoomFit:    NewThemedResource(viewzoomfitIconRes),
		IconNameViewZoomIn:     NewThemedResource(viewzoominIconRes),
		IconNameViewZoomOut:    NewThemedResource(viewzoomoutIconRes),

		IconNameVisibility:    NewThemedResource(visibilityIconRes),
		IconNameVisibilityOff: NewThemedResource(visibilityoffIconRes),

		IconNameVolumeDown: NewThemedResource(volumedownIconRes),
		IconNameVolumeMute: NewThemedResource(volumemuteIconRes),
		IconNameVolumeUp:   NewThemedResource(volumeupIconRes),

		IconNameDownload: NewThemedResource(downloadIconRes),
		IconNameComputer: NewThemedResource(computerIconRes),
		IconNameStorage:  NewThemedResource(storageIconRes),
		IconNameUpload:   NewThemedResource(uploadIconRes),
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

// Error returns a different resource for indicating an error.
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
	return colorizeResource(res.source, ErrorColor())
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
	return current().Icon(IconNameCancel)
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() fyne.Resource {
	return current().Icon(IconNameConfirm)
}

// DeleteIcon returns a resource containing the standard delete icon for the current theme
func DeleteIcon() fyne.Resource {
	return current().Icon(IconNameDelete)
}

// SearchIcon returns a resource containing the standard search icon for the current theme
func SearchIcon() fyne.Resource {
	return current().Icon(IconNameSearch)
}

// SearchReplaceIcon returns a resource containing the standard search and replace icon for the current theme
func SearchReplaceIcon() fyne.Resource {
	return current().Icon(IconNameSearchReplace)
}

// MenuIcon returns a resource containing the standard (mobile) menu icon for the current theme
func MenuIcon() fyne.Resource {
	return current().Icon(IconNameMenu)
}

// MenuExpandIcon returns a resource containing the standard (mobile) expand "submenu icon for the current theme
func MenuExpandIcon() fyne.Resource {
	return current().Icon(IconNameMenuExpand)
}

// CheckButtonIcon returns a resource containing the standard checkbox icon for the current theme
func CheckButtonIcon() fyne.Resource {
	return current().Icon(IconNameCheckButton)
}

// CheckButtonCheckedIcon returns a resource containing the standard checkbox checked icon for the current theme
func CheckButtonCheckedIcon() fyne.Resource {
	return current().Icon(IconNameCheckButtonChecked)
}

// RadioButtonIcon returns a resource containing the standard radio button icon for the current theme
func RadioButtonIcon() fyne.Resource {
	return current().Icon(IconNameRadioButton)
}

// RadioButtonCheckedIcon returns a resource containing the standard radio button checked icon for the current theme
func RadioButtonCheckedIcon() fyne.Resource {
	return current().Icon(IconNameRadioButtonChecked)
}

// ContentAddIcon returns a resource containing the standard content add icon for the current theme
func ContentAddIcon() fyne.Resource {
	return current().Icon(IconNameContentAdd)
}

// ContentRemoveIcon returns a resource containing the standard content remove icon for the current theme
func ContentRemoveIcon() fyne.Resource {
	return current().Icon(IconNameContentRemove)
}

// ContentClearIcon returns a resource containing the standard content clear icon for the current theme
func ContentClearIcon() fyne.Resource {
	return current().Icon(IconNameContentClear)
}

// ContentCutIcon returns a resource containing the standard content cut icon for the current theme
func ContentCutIcon() fyne.Resource {
	return current().Icon(IconNameContentCut)
}

// ContentCopyIcon returns a resource containing the standard content copy icon for the current theme
func ContentCopyIcon() fyne.Resource {
	return current().Icon(IconNameContentCopy)
}

// ContentPasteIcon returns a resource containing the standard content paste icon for the current theme
func ContentPasteIcon() fyne.Resource {
	return current().Icon(IconNameContentPaste)
}

// ContentRedoIcon returns a resource containing the standard content redo icon for the current theme
func ContentRedoIcon() fyne.Resource {
	return current().Icon(IconNameContentRedo)
}

// ContentUndoIcon returns a resource containing the standard content undo icon for the current theme
func ContentUndoIcon() fyne.Resource {
	return current().Icon(IconNameContentUndo)
}

// ColorAchromaticIcon returns a resource containing the standard achromatic color icon for the current theme
func ColorAchromaticIcon() fyne.Resource {
	return current().Icon(IconNameColorAchromatic)
}

// ColorChromaticIcon returns a resource containing the standard chromatic color icon for the current theme
func ColorChromaticIcon() fyne.Resource {
	return current().Icon(IconNameColorChromatic)
}

// ColorPaletteIcon returns a resource containing the standard color palette icon for the current theme
func ColorPaletteIcon() fyne.Resource {
	return current().Icon(IconNameColorPalette)
}

// DocumentIcon returns a resource containing the standard document icon for the current theme
func DocumentIcon() fyne.Resource {
	return current().Icon(IconNameDocument)
}

// DocumentCreateIcon returns a resource containing the standard document create icon for the current theme
func DocumentCreateIcon() fyne.Resource {
	return current().Icon(IconNameDocumentCreate)
}

// DocumentPrintIcon returns a resource containing the standard document print icon for the current theme
func DocumentPrintIcon() fyne.Resource {
	return current().Icon(IconNameDocumentPrint)
}

// DocumentSaveIcon returns a resource containing the standard document save icon for the current theme
func DocumentSaveIcon() fyne.Resource {
	return current().Icon(IconNameDocumentSave)
}

// InfoIcon returns a resource containing the standard dialog info icon for the current theme
func InfoIcon() fyne.Resource {
	return current().Icon(IconNameInfo)
}

// QuestionIcon returns a resource containing the standard dialog question icon for the current theme
func QuestionIcon() fyne.Resource {
	return current().Icon(IconNameQuestion)
}

// WarningIcon returns a resource containing the standard dialog warning icon for the current theme
func WarningIcon() fyne.Resource {
	return current().Icon(IconNameWarning)
}

// ErrorIcon returns a resource containing the standard dialog error icon for the current theme
func ErrorIcon() fyne.Resource {
	return current().Icon(IconNameError)
}

// FileIcon returns a resource containing the appropriate file icon for the current theme
func FileIcon() fyne.Resource {
	return current().Icon(IconNameFile)
}

// FileApplicationIcon returns a resource containing the file icon representing application files for the current theme
func FileApplicationIcon() fyne.Resource {
	return current().Icon(IconNameFileApplication)
}

// FileAudioIcon returns a resource containing the file icon representing audio files for the current theme
func FileAudioIcon() fyne.Resource {
	return current().Icon(IconNameFileAudio)
}

// FileImageIcon returns a resource containing the file icon representing image files for the current theme
func FileImageIcon() fyne.Resource {
	return current().Icon(IconNameFileImage)
}

// FileTextIcon returns a resource containing the file icon representing text files for the current theme
func FileTextIcon() fyne.Resource {
	return current().Icon(IconNameFileText)
}

// FileVideoIcon returns a resource containing the file icon representing video files for the current theme
func FileVideoIcon() fyne.Resource {
	return current().Icon(IconNameFileVideo)
}

// FolderIcon returns a resource containing the standard folder icon for the current theme
func FolderIcon() fyne.Resource {
	return current().Icon(IconNameFolder)
}

// FolderNewIcon returns a resource containing the standard folder creation icon for the current theme
func FolderNewIcon() fyne.Resource {
	return current().Icon(IconNameFolderNew)
}

// FolderOpenIcon returns a resource containing the standard folder open icon for the current theme
func FolderOpenIcon() fyne.Resource {
	return current().Icon(IconNameFolderOpen)
}

// HelpIcon returns a resource containing the standard help icon for the current theme
func HelpIcon() fyne.Resource {
	return current().Icon(IconNameHelp)
}

// HistoryIcon returns a resource containing the standard history icon for the current theme
func HistoryIcon() fyne.Resource {
	return current().Icon(IconNameHistory)
}

// HomeIcon returns a resource containing the standard home folder icon for the current theme
func HomeIcon() fyne.Resource {
	return current().Icon(IconNameHome)
}

// SettingsIcon returns a resource containing the standard settings icon for the current theme
func SettingsIcon() fyne.Resource {
	return current().Icon(IconNameSettings)
}

// MailAttachmentIcon returns a resource containing the standard mail attachment icon for the current theme
func MailAttachmentIcon() fyne.Resource {
	return current().Icon(IconNameMailAttachment)
}

// MailComposeIcon returns a resource containing the standard mail compose icon for the current theme
func MailComposeIcon() fyne.Resource {
	return current().Icon(IconNameMailCompose)
}

// MailForwardIcon returns a resource containing the standard mail forward icon for the current theme
func MailForwardIcon() fyne.Resource {
	return current().Icon(IconNameMailForward)
}

// MailReplyIcon returns a resource containing the standard mail reply icon for the current theme
func MailReplyIcon() fyne.Resource {
	return current().Icon(IconNameMailReply)
}

// MailReplyAllIcon returns a resource containing the standard mail reply all icon for the current theme
func MailReplyAllIcon() fyne.Resource {
	return current().Icon(IconNameMailReplyAll)
}

// MailSendIcon returns a resource containing the standard mail send icon for the current theme
func MailSendIcon() fyne.Resource {
	return current().Icon(IconNameMailSend)
}

// MediaFastForwardIcon returns a resource containing the standard media fast-forward icon for the current theme
func MediaFastForwardIcon() fyne.Resource {
	return current().Icon(IconNameMediaFastForward)
}

// MediaFastRewindIcon returns a resource containing the standard media fast-rewind icon for the current theme
func MediaFastRewindIcon() fyne.Resource {
	return current().Icon(IconNameMediaFastRewind)
}

// MediaPauseIcon returns a resource containing the standard media pause icon for the current theme
func MediaPauseIcon() fyne.Resource {
	return current().Icon(IconNameMediaPause)
}

// MediaPlayIcon returns a resource containing the standard media play icon for the current theme
func MediaPlayIcon() fyne.Resource {
	return current().Icon(IconNameMediaPlay)
}

// MediaRecordIcon returns a resource containing the standard media record icon for the current theme
func MediaRecordIcon() fyne.Resource {
	return current().Icon(IconNameMediaRecord)
}

// MediaReplayIcon returns a resource containing the standard media replay icon for the current theme
func MediaReplayIcon() fyne.Resource {
	return current().Icon(IconNameMediaReplay)
}

// MediaSkipNextIcon returns a resource containing the standard media skip next icon for the current theme
func MediaSkipNextIcon() fyne.Resource {
	return current().Icon(IconNameMediaSkipNext)
}

// MediaSkipPreviousIcon returns a resource containing the standard media skip previous icon for the current theme
func MediaSkipPreviousIcon() fyne.Resource {
	return current().Icon(IconNameMediaSkipPrevious)
}

// MediaStopIcon returns a resource containing the standard media stop icon for the current theme
func MediaStopIcon() fyne.Resource {
	return current().Icon(IconNameMediaStop)
}

// MoveDownIcon returns a resource containing the standard down arrow icon for the current theme
func MoveDownIcon() fyne.Resource {
	return current().Icon(IconNameMoveDown)
}

// MoveUpIcon returns a resource containing the standard up arrow icon for the current theme
func MoveUpIcon() fyne.Resource {
	return current().Icon(IconNameMoveUp)
}

// NavigateBackIcon returns a resource containing the standard backward navigation icon for the current theme
func NavigateBackIcon() fyne.Resource {
	return current().Icon(IconNameNavigateBack)
}

// NavigateNextIcon returns a resource containing the standard forward navigation icon for the current theme
func NavigateNextIcon() fyne.Resource {
	return current().Icon(IconNameNavigateNext)
}

// MenuDropDownIcon returns a resource containing the standard menu drop down icon for the current theme
func MenuDropDownIcon() fyne.Resource {
	return current().Icon(IconNameArrowDropDown)
}

// MenuDropUpIcon returns a resource containing the standard menu drop up icon for the current theme
func MenuDropUpIcon() fyne.Resource {
	return current().Icon(IconNameArrowDropUp)
}

// ViewFullScreenIcon returns a resource containing the standard fullscreen icon for the current theme
func ViewFullScreenIcon() fyne.Resource {
	return current().Icon(IconNameViewFullScreen)
}

// ViewRestoreIcon returns a resource containing the standard exit fullscreen icon for the current theme
func ViewRestoreIcon() fyne.Resource {
	return current().Icon(IconNameViewRestore)
}

// ViewRefreshIcon returns a resource containing the standard refresh icon for the current theme
func ViewRefreshIcon() fyne.Resource {
	return current().Icon(IconNameViewRefresh)
}

// ZoomFitIcon returns a resource containing the standard zoom fit icon for the current theme
func ZoomFitIcon() fyne.Resource {
	return current().Icon(IconNameViewZoomFit)
}

// ZoomInIcon returns a resource containing the standard zoom in icon for the current theme
func ZoomInIcon() fyne.Resource {
	return current().Icon(IconNameViewZoomIn)
}

// ZoomOutIcon returns a resource containing the standard zoom out icon for the current theme
func ZoomOutIcon() fyne.Resource {
	return current().Icon(IconNameViewZoomOut)
}

// VisibilityIcon returns a resource containing the standard visibity icon for the current theme
func VisibilityIcon() fyne.Resource {
	return current().Icon(IconNameVisibility)
}

// VisibilityOffIcon returns a resource containing the standard visibity off icon for the current theme
func VisibilityOffIcon() fyne.Resource {
	return current().Icon(IconNameVisibilityOff)
}

// VolumeDownIcon returns a resource containing the standard volume down icon for the current theme
func VolumeDownIcon() fyne.Resource {
	return current().Icon(IconNameVolumeDown)
}

// VolumeMuteIcon returns a resource containing the standard volume mute icon for the current theme
func VolumeMuteIcon() fyne.Resource {
	return current().Icon(IconNameVolumeMute)
}

// VolumeUpIcon returns a resource containing the standard volume up icon for the current theme
func VolumeUpIcon() fyne.Resource {
	return current().Icon(IconNameVolumeUp)
}

// ComputerIcon returns a resource containing the standard computer icon for the current theme
func ComputerIcon() fyne.Resource {
	return current().Icon(IconNameComputer)
}

// DownloadIcon returns a resource containing the standard download icon for the current theme
func DownloadIcon() fyne.Resource {
	return current().Icon(IconNameDownload)
}

// StorageIcon returns a resource containing the standard storage icon for the current theme
func StorageIcon() fyne.Resource {
	return current().Icon(IconNameStorage)
}

// UploadIcon returns a resource containing the standard upload icon for the current theme
func UploadIcon() fyne.Resource {
	return current().Icon(IconNameUpload)
}
