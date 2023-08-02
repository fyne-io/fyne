package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/svg"
)

const (
	// IconNameCancel is the name of theme lookup for cancel icon.
	//
	// Since: 2.0
	IconNameCancel fyne.ThemeIconName = "cancel"

	// IconNameConfirm is the name of theme lookup for confirm icon.
	//
	// Since: 2.0
	IconNameConfirm fyne.ThemeIconName = "confirm"

	// IconNameDelete is the name of theme lookup for delete icon.
	//
	// Since: 2.0
	IconNameDelete fyne.ThemeIconName = "delete"

	// IconNameSearch is the name of theme lookup for search icon.
	//
	// Since: 2.0
	IconNameSearch fyne.ThemeIconName = "search"

	// IconNameSearchReplace is the name of theme lookup for search and replace icon.
	//
	// Since: 2.0
	IconNameSearchReplace fyne.ThemeIconName = "searchReplace"

	// IconNameMenu is the name of theme lookup for menu icon.
	//
	// Since: 2.0
	IconNameMenu fyne.ThemeIconName = "menu"

	// IconNameMenuExpand is the name of theme lookup for menu expansion icon.
	//
	// Since: 2.0
	IconNameMenuExpand fyne.ThemeIconName = "menuExpand"

	// IconNameCheckButtonChecked is the name of theme lookup for checked check button icon.
	//
	// Since: 2.0
	IconNameCheckButtonChecked fyne.ThemeIconName = "checked"

	// IconNameCheckButton is the name of theme lookup for  unchecked check button icon.
	//
	// Since: 2.0
	IconNameCheckButton fyne.ThemeIconName = "unchecked"

	// IconNameRadioButton is the name of theme lookup for radio button unchecked icon.
	//
	// Since: 2.0
	IconNameRadioButton fyne.ThemeIconName = "radioButton"

	// IconNameRadioButtonChecked is the name of theme lookup for radio button checked icon.
	//
	// Since: 2.0
	IconNameRadioButtonChecked fyne.ThemeIconName = "radioButtonChecked"

	// IconNameColorAchromatic is the name of theme lookup for greyscale color icon.
	//
	// Since: 2.0
	IconNameColorAchromatic fyne.ThemeIconName = "colorAchromatic"

	// IconNameColorChromatic is the name of theme lookup for full color icon.
	//
	// Since: 2.0
	IconNameColorChromatic fyne.ThemeIconName = "colorChromatic"

	// IconNameColorPalette is the name of theme lookup for color palette icon.
	//
	// Since: 2.0
	IconNameColorPalette fyne.ThemeIconName = "colorPalette"

	// IconNameContentAdd is the name of theme lookup for content add icon.
	//
	// Since: 2.0
	IconNameContentAdd fyne.ThemeIconName = "contentAdd"

	// IconNameContentRemove is the name of theme lookup for content remove icon.
	//
	// Since: 2.0
	IconNameContentRemove fyne.ThemeIconName = "contentRemove"

	// IconNameContentCut is the name of theme lookup for content cut icon.
	//
	// Since: 2.0
	IconNameContentCut fyne.ThemeIconName = "contentCut"

	// IconNameContentCopy is the name of theme lookup for content copy icon.
	//
	// Since: 2.0
	IconNameContentCopy fyne.ThemeIconName = "contentCopy"

	// IconNameContentPaste is the name of theme lookup for content paste icon.
	//
	// Since: 2.0
	IconNameContentPaste fyne.ThemeIconName = "contentPaste"

	// IconNameContentClear is the name of theme lookup for content clear icon.
	//
	// Since: 2.0
	IconNameContentClear fyne.ThemeIconName = "contentClear"

	// IconNameContentRedo is the name of theme lookup for content redo icon.
	//
	// Since: 2.0
	IconNameContentRedo fyne.ThemeIconName = "contentRedo"

	// IconNameContentUndo is the name of theme lookup for content undo icon.
	//
	// Since: 2.0
	IconNameContentUndo fyne.ThemeIconName = "contentUndo"

	// IconNameInfo is the name of theme lookup for info icon.
	//
	// Since: 2.0
	IconNameInfo fyne.ThemeIconName = "info"

	// IconNameQuestion is the name of theme lookup for question icon.
	//
	// Since: 2.0
	IconNameQuestion fyne.ThemeIconName = "question"

	// IconNameWarning is the name of theme lookup for warning icon.
	//
	// Since: 2.0
	IconNameWarning fyne.ThemeIconName = "warning"

	// IconNameError is the name of theme lookup for error icon.
	//
	// Since: 2.0
	IconNameError fyne.ThemeIconName = "error"

	// IconNameBrokenImage is the name of the theme lookup for broken-image icon.
	//
	// Since: 2.4
	IconNameBrokenImage fyne.ThemeIconName = "broken-image"

	// IconNameDocument is the name of theme lookup for document icon.
	//
	// Since: 2.0
	IconNameDocument fyne.ThemeIconName = "document"

	// IconNameDocumentCreate is the name of theme lookup for document create icon.
	//
	// Since: 2.0
	IconNameDocumentCreate fyne.ThemeIconName = "documentCreate"

	// IconNameDocumentPrint is the name of theme lookup for document print icon.
	//
	// Since: 2.0
	IconNameDocumentPrint fyne.ThemeIconName = "documentPrint"

	// IconNameDocumentSave is the name of theme lookup for document save icon.
	//
	// Since: 2.0
	IconNameDocumentSave fyne.ThemeIconName = "documentSave"

	// IconNameMoreHorizontal is the name of theme lookup for horizontal more.
	//
	// Since 2.0
	IconNameMoreHorizontal fyne.ThemeIconName = "moreHorizontal"

	// IconNameMoreVertical is the name of theme lookup for vertical more.
	//
	// Since 2.0
	IconNameMoreVertical fyne.ThemeIconName = "moreVertical"

	// IconNameMailAttachment is the name of theme lookup for mail attachment icon.
	//
	// Since: 2.0
	IconNameMailAttachment fyne.ThemeIconName = "mailAttachment"

	// IconNameMailCompose is the name of theme lookup for mail compose icon.
	//
	// Since: 2.0
	IconNameMailCompose fyne.ThemeIconName = "mailCompose"

	// IconNameMailForward is the name of theme lookup for mail forward icon.
	//
	// Since: 2.0
	IconNameMailForward fyne.ThemeIconName = "mailForward"

	// IconNameMailReply is the name of theme lookup for mail reply icon.
	//
	// Since: 2.0
	IconNameMailReply fyne.ThemeIconName = "mailReply"

	// IconNameMailReplyAll is the name of theme lookup for mail reply-all icon.
	//
	// Since: 2.0
	IconNameMailReplyAll fyne.ThemeIconName = "mailReplyAll"

	// IconNameMailSend is the name of theme lookup for mail send icon.
	//
	// Since: 2.0
	IconNameMailSend fyne.ThemeIconName = "mailSend"

	// IconNameMediaMusic is the name of theme lookup for media music icon.
	//
	// Since: 2.1
	IconNameMediaMusic fyne.ThemeIconName = "mediaMusic"

	// IconNameMediaPhoto is the name of theme lookup for media photo icon.
	//
	// Since: 2.1
	IconNameMediaPhoto fyne.ThemeIconName = "mediaPhoto"

	// IconNameMediaVideo is the name of theme lookup for media video icon.
	//
	// Since: 2.1
	IconNameMediaVideo fyne.ThemeIconName = "mediaVideo"

	// IconNameMediaFastForward is the name of theme lookup for media fast-forward icon.
	//
	// Since: 2.0
	IconNameMediaFastForward fyne.ThemeIconName = "mediaFastForward"

	// IconNameMediaFastRewind is the name of theme lookup for media fast-rewind icon.
	//
	// Since: 2.0
	IconNameMediaFastRewind fyne.ThemeIconName = "mediaFastRewind"

	// IconNameMediaPause is the name of theme lookup for media pause icon.
	//
	// Since: 2.0
	IconNameMediaPause fyne.ThemeIconName = "mediaPause"

	// IconNameMediaPlay is the name of theme lookup for media play icon.
	//
	// Since: 2.0
	IconNameMediaPlay fyne.ThemeIconName = "mediaPlay"

	// IconNameMediaRecord is the name of theme lookup for media record icon.
	//
	// Since: 2.0
	IconNameMediaRecord fyne.ThemeIconName = "mediaRecord"

	// IconNameMediaReplay is the name of theme lookup for media replay icon.
	//
	// Since: 2.0
	IconNameMediaReplay fyne.ThemeIconName = "mediaReplay"

	// IconNameMediaSkipNext is the name of theme lookup for media skip next icon.
	//
	// Since: 2.0
	IconNameMediaSkipNext fyne.ThemeIconName = "mediaSkipNext"

	// IconNameMediaSkipPrevious is the name of theme lookup for media skip previous icon.
	//
	// Since: 2.0
	IconNameMediaSkipPrevious fyne.ThemeIconName = "mediaSkipPrevious"

	// IconNameMediaStop is the name of theme lookup for media stop icon.
	//
	// Since: 2.0
	IconNameMediaStop fyne.ThemeIconName = "mediaStop"

	// IconNameMoveDown is the name of theme lookup for move down icon.
	//
	// Since: 2.0
	IconNameMoveDown fyne.ThemeIconName = "arrowDown"

	// IconNameMoveUp is the name of theme lookup for move up icon.
	//
	// Since: 2.0
	IconNameMoveUp fyne.ThemeIconName = "arrowUp"

	// IconNameNavigateBack is the name of theme lookup for navigate back icon.
	//
	// Since: 2.0
	IconNameNavigateBack fyne.ThemeIconName = "arrowBack"

	// IconNameNavigateNext is the name of theme lookup for navigate next icon.
	//
	// Since: 2.0
	IconNameNavigateNext fyne.ThemeIconName = "arrowForward"

	// IconNameArrowDropDown is the name of theme lookup for drop-down arrow icon.
	//
	// Since: 2.0
	IconNameArrowDropDown fyne.ThemeIconName = "arrowDropDown"

	// IconNameArrowDropUp is the name of theme lookup for drop-up arrow icon.
	//
	// Since: 2.0
	IconNameArrowDropUp fyne.ThemeIconName = "arrowDropUp"

	// IconNameFile is the name of theme lookup for file icon.
	//
	// Since: 2.0
	IconNameFile fyne.ThemeIconName = "file"

	// IconNameFileApplication is the name of theme lookup for file application icon.
	//
	// Since: 2.0
	IconNameFileApplication fyne.ThemeIconName = "fileApplication"

	// IconNameFileAudio is the name of theme lookup for file audio icon.
	//
	// Since: 2.0
	IconNameFileAudio fyne.ThemeIconName = "fileAudio"

	// IconNameFileImage is the name of theme lookup for file image icon.
	//
	// Since: 2.0
	IconNameFileImage fyne.ThemeIconName = "fileImage"

	// IconNameFileText is the name of theme lookup for file text icon.
	//
	// Since: 2.0
	IconNameFileText fyne.ThemeIconName = "fileText"

	// IconNameFileVideo is the name of theme lookup for file video icon.
	//
	// Since: 2.0
	IconNameFileVideo fyne.ThemeIconName = "fileVideo"

	// IconNameFolder is the name of theme lookup for folder icon.
	//
	// Since: 2.0
	IconNameFolder fyne.ThemeIconName = "folder"

	// IconNameFolderNew is the name of theme lookup for folder new icon.
	//
	// Since: 2.0
	IconNameFolderNew fyne.ThemeIconName = "folderNew"

	// IconNameFolderOpen is the name of theme lookup for folder open icon.
	//
	// Since: 2.0
	IconNameFolderOpen fyne.ThemeIconName = "folderOpen"

	// IconNameHelp is the name of theme lookup for help icon.
	//
	// Since: 2.0
	IconNameHelp fyne.ThemeIconName = "help"

	// IconNameHistory is the name of theme lookup for history icon.
	//
	// Since: 2.0
	IconNameHistory fyne.ThemeIconName = "history"

	// IconNameHome is the name of theme lookup for home icon.
	//
	// Since: 2.0
	IconNameHome fyne.ThemeIconName = "home"

	// IconNameSettings is the name of theme lookup for settings icon.
	//
	// Since: 2.0
	IconNameSettings fyne.ThemeIconName = "settings"

	// IconNameStorage is the name of theme lookup for storage icon.
	//
	// Since: 2.0
	IconNameStorage fyne.ThemeIconName = "storage"

	// IconNameUpload is the name of theme lookup for upload icon.
	//
	// Since: 2.0
	IconNameUpload fyne.ThemeIconName = "upload"

	// IconNameViewFullScreen is the name of theme lookup for view fullscreen icon.
	//
	// Since: 2.0
	IconNameViewFullScreen fyne.ThemeIconName = "viewFullScreen"

	// IconNameViewRefresh is the name of theme lookup for view refresh icon.
	//
	// Since: 2.0
	IconNameViewRefresh fyne.ThemeIconName = "viewRefresh"

	// IconNameViewZoomFit is the name of theme lookup for view zoom fit icon.
	//
	// Since: 2.0
	IconNameViewZoomFit fyne.ThemeIconName = "viewZoomFit"

	// IconNameViewZoomIn is the name of theme lookup for view zoom in icon.
	//
	// Since: 2.0
	IconNameViewZoomIn fyne.ThemeIconName = "viewZoomIn"

	// IconNameViewZoomOut is the name of theme lookup for view zoom out icon.
	//
	// Since: 2.0
	IconNameViewZoomOut fyne.ThemeIconName = "viewZoomOut"

	// IconNameViewRestore is the name of theme lookup for view restore icon.
	//
	// Since: 2.0
	IconNameViewRestore fyne.ThemeIconName = "viewRestore"

	// IconNameVisibility is the name of theme lookup for visibility icon.
	//
	// Since: 2.0
	IconNameVisibility fyne.ThemeIconName = "visibility"

	// IconNameVisibilityOff is the name of theme lookup for invisibility icon.
	//
	// Since: 2.0
	IconNameVisibilityOff fyne.ThemeIconName = "visibilityOff"

	// IconNameVolumeDown is the name of theme lookup for volume down icon.
	//
	// Since: 2.0
	IconNameVolumeDown fyne.ThemeIconName = "volumeDown"

	// IconNameVolumeMute is the name of theme lookup for volume mute icon.
	//
	// Since: 2.0
	IconNameVolumeMute fyne.ThemeIconName = "volumeMute"

	// IconNameVolumeUp is the name of theme lookup for volume up icon.
	//
	// Since: 2.0
	IconNameVolumeUp fyne.ThemeIconName = "volumeUp"

	// IconNameDownload is the name of theme lookup for download icon.
	//
	// Since: 2.0
	IconNameDownload fyne.ThemeIconName = "download"

	// IconNameComputer is the name of theme lookup for computer icon.
	//
	// Since: 2.0
	IconNameComputer fyne.ThemeIconName = "computer"

	// IconNameAccount is the name of theme lookup for account icon.
	//
	// Since: 2.1
	IconNameAccount fyne.ThemeIconName = "account"

	// IconNameLogin is the name of theme lookup for login icon.
	//
	// Since: 2.1
	IconNameLogin fyne.ThemeIconName = "login"

	// IconNameLogout is the name of theme lookup for logout icon.
	//
	// Since: 2.1
	IconNameLogout fyne.ThemeIconName = "logout"

	// IconNameList is the name of theme lookup for list icon.
	//
	// Since: 2.1
	IconNameList fyne.ThemeIconName = "list"

	// IconNameGrid is the name of theme lookup for grid icon.
	//
	// Since: 2.1
	IconNameGrid fyne.ThemeIconName = "grid"
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

		IconNameCheckButton:        NewThemedResource(checkboxIconRes),
		IconNameCheckButtonChecked: NewThemedResource(checkboxcheckedIconRes),
		"iconNameCheckButtonFill":  NewThemedResource(checkboxfillIconRes),
		IconNameRadioButton:        NewThemedResource(radiobuttonIconRes),
		IconNameRadioButtonChecked: NewThemedResource(radiobuttoncheckedIconRes),
		"iconNameRadioButtonFill":  NewThemedResource(radiobuttonfillIconRes),

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

		IconNameMoreHorizontal: NewThemedResource(morehorizontalIconRes),
		IconNameMoreVertical:   NewThemedResource(moreverticalIconRes),

		IconNameInfo:        NewThemedResource(infoIconRes),
		IconNameQuestion:    NewThemedResource(questionIconRes),
		IconNameWarning:     NewThemedResource(warningIconRes),
		IconNameError:       NewThemedResource(errorIconRes),
		IconNameBrokenImage: NewThemedResource(brokenimageIconRes),

		IconNameMailAttachment: NewThemedResource(mailattachmentIconRes),
		IconNameMailCompose:    NewThemedResource(mailcomposeIconRes),
		IconNameMailForward:    NewThemedResource(mailforwardIconRes),
		IconNameMailReply:      NewThemedResource(mailreplyIconRes),
		IconNameMailReplyAll:   NewThemedResource(mailreplyallIconRes),
		IconNameMailSend:       NewThemedResource(mailsendIconRes),

		IconNameMediaMusic:        NewThemedResource(mediamusicIconRes),
		IconNameMediaPhoto:        NewThemedResource(mediaphotoIconRes),
		IconNameMediaVideo:        NewThemedResource(mediavideoIconRes),
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

		IconNameAccount: NewThemedResource(accountIconRes),
		IconNameLogin:   NewThemedResource(loginIconRes),
		IconNameLogout:  NewThemedResource(logoutIconRes),

		IconNameList: NewThemedResource(listIconRes),
		IconNameGrid: NewThemedResource(gridIconRes),
	}
)

func (t *builtinTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return icons[n]
}

// ThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// for the currently selected theme.
type ThemedResource struct {
	source fyne.Resource

	// ColorName specifies which theme colour should be used to theme the resource
	//
	// Since: 2.3
	ColorName fyne.ThemeColorName
}

// NewColoredResource creates a resource that adapts to the current theme setting using
// the color named in the constructor.
//
// Since: 2.4
func NewColoredResource(src fyne.Resource, name fyne.ThemeColorName) *ThemedResource {
	return &ThemedResource{
		source:    src,
		ColorName: name,
	}
}

// NewSuccessThemedResource creates a resource that adapts to the current theme success color.
//
// Since: 2.4
func NewSuccessThemedResource(src fyne.Resource) *ThemedResource {
	return &ThemedResource{
		source:    src,
		ColorName: ColorNameSuccess,
	}
}

// NewThemedResource creates a resource that adapts to the current theme setting.
// By default this will match the foreground color, but it can be changed using the `ColorName` field.
func NewThemedResource(src fyne.Resource) *ThemedResource {
	return &ThemedResource{
		source: src,
	}
}

// NewWarningThemedResource creates a resource that adapts to the current theme warning color.
//
// Since: 2.4
func NewWarningThemedResource(src fyne.Resource) *ThemedResource {
	return &ThemedResource{
		source:    src,
		ColorName: ColorNameWarning,
	}
}

// Name returns the underlying resource name (used for caching).
func (res *ThemedResource) Name() string {
	prefix := res.ColorName
	if prefix == "" {
		prefix = "foreground_"
	} else {
		prefix += "_"
	}

	return string(prefix) + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current text color.
func (res *ThemedResource) Content() []byte {
	name := res.ColorName
	if name == "" {
		name = ColorNameForeground
	}

	return svg.Colorize(res.source.Content(), safeColorLookup(name, currentVariant()))
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
	return "inverted_" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *InvertedThemedResource) Content() []byte {
	clr := BackgroundColor()
	return svg.Colorize(res.source.Content(), clr)
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
	return "error_" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *ErrorThemedResource) Content() []byte {
	return svg.Colorize(res.source.Content(), ErrorColor())
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
	return "primary_" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *PrimaryThemedResource) Content() []byte {
	return svg.Colorize(res.source.Content(), PrimaryColor())
}

// Original returns the underlying resource that this primary themed resource was adapted from
func (res *PrimaryThemedResource) Original() fyne.Resource {
	return res.source
}

// DisabledResource is a resource wrapper that will return an appropriate resource colorized by
// the current theme's `DisabledColor` color.
type DisabledResource struct {
	source fyne.Resource
}

// Name returns the resource source name prefixed with `disabled_` (used for caching)
func (res *DisabledResource) Name() string {
	return "disabled_" + res.source.Name()
}

// Content returns the disabled style content of the correct resource for the current theme
func (res *DisabledResource) Content() []byte {
	return svg.Colorize(res.source.Content(), DisabledColor())
}

// NewDisabledResource creates a resource that adapts to the current theme's DisabledColor setting.
func NewDisabledResource(res fyne.Resource) *DisabledResource {
	return &DisabledResource{
		source: res,
	}
}

// FyneLogo returns a resource containing the Fyne logo.
//
// Deprecated: Applications should use their own icon in most cases.
func FyneLogo() fyne.Resource {
	return fynelogo
}

// CancelIcon returns a resource containing the standard cancel icon for the current theme
func CancelIcon() fyne.Resource {
	return safeIconLookup(IconNameCancel)
}

// ConfirmIcon returns a resource containing the standard confirm icon for the current theme
func ConfirmIcon() fyne.Resource {
	return safeIconLookup(IconNameConfirm)
}

// DeleteIcon returns a resource containing the standard delete icon for the current theme
func DeleteIcon() fyne.Resource {
	return safeIconLookup(IconNameDelete)
}

// SearchIcon returns a resource containing the standard search icon for the current theme
func SearchIcon() fyne.Resource {
	return safeIconLookup(IconNameSearch)
}

// SearchReplaceIcon returns a resource containing the standard search and replace icon for the current theme
func SearchReplaceIcon() fyne.Resource {
	return safeIconLookup(IconNameSearchReplace)
}

// MenuIcon returns a resource containing the standard (mobile) menu icon for the current theme
func MenuIcon() fyne.Resource {
	return safeIconLookup(IconNameMenu)
}

// MenuExpandIcon returns a resource containing the standard (mobile) expand "submenu icon for the current theme
func MenuExpandIcon() fyne.Resource {
	return safeIconLookup(IconNameMenuExpand)
}

// CheckButtonIcon returns a resource containing the standard checkbox icon for the current theme
func CheckButtonIcon() fyne.Resource {
	return safeIconLookup(IconNameCheckButton)
}

// CheckButtonCheckedIcon returns a resource containing the standard checkbox checked icon for the current theme
func CheckButtonCheckedIcon() fyne.Resource {
	return safeIconLookup(IconNameCheckButtonChecked)
}

// RadioButtonIcon returns a resource containing the standard radio button icon for the current theme
func RadioButtonIcon() fyne.Resource {
	return safeIconLookup(IconNameRadioButton)
}

// RadioButtonCheckedIcon returns a resource containing the standard radio button checked icon for the current theme
func RadioButtonCheckedIcon() fyne.Resource {
	return safeIconLookup(IconNameRadioButtonChecked)
}

// ContentAddIcon returns a resource containing the standard content add icon for the current theme
func ContentAddIcon() fyne.Resource {
	return safeIconLookup(IconNameContentAdd)
}

// ContentRemoveIcon returns a resource containing the standard content remove icon for the current theme
func ContentRemoveIcon() fyne.Resource {
	return safeIconLookup(IconNameContentRemove)
}

// ContentClearIcon returns a resource containing the standard content clear icon for the current theme
func ContentClearIcon() fyne.Resource {
	return safeIconLookup(IconNameContentClear)
}

// ContentCutIcon returns a resource containing the standard content cut icon for the current theme
func ContentCutIcon() fyne.Resource {
	return safeIconLookup(IconNameContentCut)
}

// ContentCopyIcon returns a resource containing the standard content copy icon for the current theme
func ContentCopyIcon() fyne.Resource {
	return safeIconLookup(IconNameContentCopy)
}

// ContentPasteIcon returns a resource containing the standard content paste icon for the current theme
func ContentPasteIcon() fyne.Resource {
	return safeIconLookup(IconNameContentPaste)
}

// ContentRedoIcon returns a resource containing the standard content redo icon for the current theme
func ContentRedoIcon() fyne.Resource {
	return safeIconLookup(IconNameContentRedo)
}

// ContentUndoIcon returns a resource containing the standard content undo icon for the current theme
func ContentUndoIcon() fyne.Resource {
	return safeIconLookup(IconNameContentUndo)
}

// ColorAchromaticIcon returns a resource containing the standard achromatic color icon for the current theme
func ColorAchromaticIcon() fyne.Resource {
	return safeIconLookup(IconNameColorAchromatic)
}

// ColorChromaticIcon returns a resource containing the standard chromatic color icon for the current theme
func ColorChromaticIcon() fyne.Resource {
	return safeIconLookup(IconNameColorChromatic)
}

// ColorPaletteIcon returns a resource containing the standard color palette icon for the current theme
func ColorPaletteIcon() fyne.Resource {
	return safeIconLookup(IconNameColorPalette)
}

// DocumentIcon returns a resource containing the standard document icon for the current theme
func DocumentIcon() fyne.Resource {
	return safeIconLookup(IconNameDocument)
}

// DocumentCreateIcon returns a resource containing the standard document create icon for the current theme
func DocumentCreateIcon() fyne.Resource {
	return safeIconLookup(IconNameDocumentCreate)
}

// DocumentPrintIcon returns a resource containing the standard document print icon for the current theme
func DocumentPrintIcon() fyne.Resource {
	return safeIconLookup(IconNameDocumentPrint)
}

// DocumentSaveIcon returns a resource containing the standard document save icon for the current theme
func DocumentSaveIcon() fyne.Resource {
	return safeIconLookup(IconNameDocumentSave)
}

// MoreHorizontalIcon returns a resource containing the standard horizontal more icon for the current theme
func MoreHorizontalIcon() fyne.Resource {
	return current().Icon(IconNameMoreHorizontal)
}

// MoreVerticalIcon returns a resource containing the standard vertical more icon for the current theme
func MoreVerticalIcon() fyne.Resource {
	return current().Icon(IconNameMoreVertical)
}

// InfoIcon returns a resource containing the standard dialog info icon for the current theme
func InfoIcon() fyne.Resource {
	return safeIconLookup(IconNameInfo)
}

// QuestionIcon returns a resource containing the standard dialog question icon for the current theme
func QuestionIcon() fyne.Resource {
	return safeIconLookup(IconNameQuestion)
}

// WarningIcon returns a resource containing the standard dialog warning icon for the current theme
func WarningIcon() fyne.Resource {
	return safeIconLookup(IconNameWarning)
}

// ErrorIcon returns a resource containing the standard dialog error icon for the current theme
func ErrorIcon() fyne.Resource {
	return safeIconLookup(IconNameError)
}

// BrokenImageIcon returns a resource containing an icon to specify a broken or missing image
//
// Since: 2.4
func BrokenImageIcon() fyne.Resource {
	return safeIconLookup(IconNameBrokenImage)
}

// FileIcon returns a resource containing the appropriate file icon for the current theme
func FileIcon() fyne.Resource {
	return safeIconLookup(IconNameFile)
}

// FileApplicationIcon returns a resource containing the file icon representing application files for the current theme
func FileApplicationIcon() fyne.Resource {
	return safeIconLookup(IconNameFileApplication)
}

// FileAudioIcon returns a resource containing the file icon representing audio files for the current theme
func FileAudioIcon() fyne.Resource {
	return safeIconLookup(IconNameFileAudio)
}

// FileImageIcon returns a resource containing the file icon representing image files for the current theme
func FileImageIcon() fyne.Resource {
	return safeIconLookup(IconNameFileImage)
}

// FileTextIcon returns a resource containing the file icon representing text files for the current theme
func FileTextIcon() fyne.Resource {
	return safeIconLookup(IconNameFileText)
}

// FileVideoIcon returns a resource containing the file icon representing video files for the current theme
func FileVideoIcon() fyne.Resource {
	return safeIconLookup(IconNameFileVideo)
}

// FolderIcon returns a resource containing the standard folder icon for the current theme
func FolderIcon() fyne.Resource {
	return safeIconLookup(IconNameFolder)
}

// FolderNewIcon returns a resource containing the standard folder creation icon for the current theme
func FolderNewIcon() fyne.Resource {
	return safeIconLookup(IconNameFolderNew)
}

// FolderOpenIcon returns a resource containing the standard folder open icon for the current theme
func FolderOpenIcon() fyne.Resource {
	return safeIconLookup(IconNameFolderOpen)
}

// HelpIcon returns a resource containing the standard help icon for the current theme
func HelpIcon() fyne.Resource {
	return safeIconLookup(IconNameHelp)
}

// HistoryIcon returns a resource containing the standard history icon for the current theme
func HistoryIcon() fyne.Resource {
	return safeIconLookup(IconNameHistory)
}

// HomeIcon returns a resource containing the standard home folder icon for the current theme
func HomeIcon() fyne.Resource {
	return safeIconLookup(IconNameHome)
}

// SettingsIcon returns a resource containing the standard settings icon for the current theme
func SettingsIcon() fyne.Resource {
	return safeIconLookup(IconNameSettings)
}

// MailAttachmentIcon returns a resource containing the standard mail attachment icon for the current theme
func MailAttachmentIcon() fyne.Resource {
	return safeIconLookup(IconNameMailAttachment)
}

// MailComposeIcon returns a resource containing the standard mail compose icon for the current theme
func MailComposeIcon() fyne.Resource {
	return safeIconLookup(IconNameMailCompose)
}

// MailForwardIcon returns a resource containing the standard mail forward icon for the current theme
func MailForwardIcon() fyne.Resource {
	return safeIconLookup(IconNameMailForward)
}

// MailReplyIcon returns a resource containing the standard mail reply icon for the current theme
func MailReplyIcon() fyne.Resource {
	return safeIconLookup(IconNameMailReply)
}

// MailReplyAllIcon returns a resource containing the standard mail reply all icon for the current theme
func MailReplyAllIcon() fyne.Resource {
	return safeIconLookup(IconNameMailReplyAll)
}

// MailSendIcon returns a resource containing the standard mail send icon for the current theme
func MailSendIcon() fyne.Resource {
	return safeIconLookup(IconNameMailSend)
}

// MediaMusicIcon returns a resource containing the standard media music icon for the current theme
//
// Since: 2.1
func MediaMusicIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaMusic)
}

// MediaPhotoIcon returns a resource containing the standard media photo icon for the current theme
//
// Since: 2.1
func MediaPhotoIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaPhoto)
}

// MediaVideoIcon returns a resource containing the standard media video icon for the current theme
//
// Since: 2.1
func MediaVideoIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaVideo)
}

// MediaFastForwardIcon returns a resource containing the standard media fast-forward icon for the current theme
func MediaFastForwardIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaFastForward)
}

// MediaFastRewindIcon returns a resource containing the standard media fast-rewind icon for the current theme
func MediaFastRewindIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaFastRewind)
}

// MediaPauseIcon returns a resource containing the standard media pause icon for the current theme
func MediaPauseIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaPause)
}

// MediaPlayIcon returns a resource containing the standard media play icon for the current theme
func MediaPlayIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaPlay)
}

// MediaRecordIcon returns a resource containing the standard media record icon for the current theme
func MediaRecordIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaRecord)
}

// MediaReplayIcon returns a resource containing the standard media replay icon for the current theme
func MediaReplayIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaReplay)
}

// MediaSkipNextIcon returns a resource containing the standard media skip next icon for the current theme
func MediaSkipNextIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaSkipNext)
}

// MediaSkipPreviousIcon returns a resource containing the standard media skip previous icon for the current theme
func MediaSkipPreviousIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaSkipPrevious)
}

// MediaStopIcon returns a resource containing the standard media stop icon for the current theme
func MediaStopIcon() fyne.Resource {
	return safeIconLookup(IconNameMediaStop)
}

// MoveDownIcon returns a resource containing the standard down arrow icon for the current theme
func MoveDownIcon() fyne.Resource {
	return safeIconLookup(IconNameMoveDown)
}

// MoveUpIcon returns a resource containing the standard up arrow icon for the current theme
func MoveUpIcon() fyne.Resource {
	return safeIconLookup(IconNameMoveUp)
}

// NavigateBackIcon returns a resource containing the standard backward navigation icon for the current theme
func NavigateBackIcon() fyne.Resource {
	return safeIconLookup(IconNameNavigateBack)
}

// NavigateNextIcon returns a resource containing the standard forward navigation icon for the current theme
func NavigateNextIcon() fyne.Resource {
	return safeIconLookup(IconNameNavigateNext)
}

// MenuDropDownIcon returns a resource containing the standard menu drop down icon for the current theme
func MenuDropDownIcon() fyne.Resource {
	return safeIconLookup(IconNameArrowDropDown)
}

// MenuDropUpIcon returns a resource containing the standard menu drop up icon for the current theme
func MenuDropUpIcon() fyne.Resource {
	return safeIconLookup(IconNameArrowDropUp)
}

// ViewFullScreenIcon returns a resource containing the standard fullscreen icon for the current theme
func ViewFullScreenIcon() fyne.Resource {
	return safeIconLookup(IconNameViewFullScreen)
}

// ViewRestoreIcon returns a resource containing the standard exit fullscreen icon for the current theme
func ViewRestoreIcon() fyne.Resource {
	return safeIconLookup(IconNameViewRestore)
}

// ViewRefreshIcon returns a resource containing the standard refresh icon for the current theme
func ViewRefreshIcon() fyne.Resource {
	return safeIconLookup(IconNameViewRefresh)
}

// ZoomFitIcon returns a resource containing the standard zoom fit icon for the current theme
func ZoomFitIcon() fyne.Resource {
	return safeIconLookup(IconNameViewZoomFit)
}

// ZoomInIcon returns a resource containing the standard zoom in icon for the current theme
func ZoomInIcon() fyne.Resource {
	return safeIconLookup(IconNameViewZoomIn)
}

// ZoomOutIcon returns a resource containing the standard zoom out icon for the current theme
func ZoomOutIcon() fyne.Resource {
	return safeIconLookup(IconNameViewZoomOut)
}

// VisibilityIcon returns a resource containing the standard visibility icon for the current theme
func VisibilityIcon() fyne.Resource {
	return safeIconLookup(IconNameVisibility)
}

// VisibilityOffIcon returns a resource containing the standard visibility off icon for the current theme
func VisibilityOffIcon() fyne.Resource {
	return safeIconLookup(IconNameVisibilityOff)
}

// VolumeDownIcon returns a resource containing the standard volume down icon for the current theme
func VolumeDownIcon() fyne.Resource {
	return safeIconLookup(IconNameVolumeDown)
}

// VolumeMuteIcon returns a resource containing the standard volume mute icon for the current theme
func VolumeMuteIcon() fyne.Resource {
	return safeIconLookup(IconNameVolumeMute)
}

// VolumeUpIcon returns a resource containing the standard volume up icon for the current theme
func VolumeUpIcon() fyne.Resource {
	return safeIconLookup(IconNameVolumeUp)
}

// ComputerIcon returns a resource containing the standard computer icon for the current theme
func ComputerIcon() fyne.Resource {
	return safeIconLookup(IconNameComputer)
}

// DownloadIcon returns a resource containing the standard download icon for the current theme
func DownloadIcon() fyne.Resource {
	return safeIconLookup(IconNameDownload)
}

// StorageIcon returns a resource containing the standard storage icon for the current theme
func StorageIcon() fyne.Resource {
	return safeIconLookup(IconNameStorage)
}

// UploadIcon returns a resource containing the standard upload icon for the current theme
func UploadIcon() fyne.Resource {
	return safeIconLookup(IconNameUpload)
}

// AccountIcon returns a resource containing the standard account icon for the current theme
func AccountIcon() fyne.Resource {
	return safeIconLookup(IconNameAccount)
}

// LoginIcon returns a resource containing the standard login icon for the current theme
func LoginIcon() fyne.Resource {
	return safeIconLookup(IconNameLogin)
}

// LogoutIcon returns a resource containing the standard logout icon for the current theme
func LogoutIcon() fyne.Resource {
	return safeIconLookup(IconNameLogout)
}

// ListIcon returns a resource containing the standard list icon for the current theme
func ListIcon() fyne.Resource {
	return safeIconLookup(IconNameList)
}

// GridIcon returns a resource containing the standard grid icon for the current theme
func GridIcon() fyne.Resource {
	return safeIconLookup(IconNameGrid)
}

func safeIconLookup(n fyne.ThemeIconName) fyne.Resource {
	icon := current().Icon(n)
	if icon == nil {
		fyne.LogError("Loaded theme returned nil icon", nil)
		return fallbackIcon
	}
	return icon
}
