package theme

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icons/fyne.png
var fyneLogo []byte
var fynelogo = &fyne.StaticResource{
	StaticName:    "fyne.png",
	StaticContent: fyneLogo,
}

//go:embed icons/cancel.svg
var cancelIcon []byte
var cancelIconRes = &fyne.StaticResource{
	StaticName:    "cancel.svg",
	StaticContent: cancelIcon,
}

//go:embed icons/check.svg
var checkIcon []byte
var checkIconRes = &fyne.StaticResource{
	StaticName:    "check.svg",
	StaticContent: checkIcon,
}

//go:embed icons/delete.svg
var deleteIcon []byte
var deleteIconRes = &fyne.StaticResource{
	StaticName:    "delete.svg",
	StaticContent: deleteIcon,
}

//go:embed icons/search.svg
var searchIcon []byte
var searchIconRes = &fyne.StaticResource{
	StaticName:    "search.svg",
	StaticContent: searchIcon,
}

//go:embed icons/search-replace.svg
var searchreplaceIcon []byte
var searchreplaceIconRes = &fyne.StaticResource{
	StaticName:    "search-replace.svg",
	StaticContent: searchreplaceIcon,
}

//go:embed icons/menu.svg
var menuIcon []byte
var menuIconRes = &fyne.StaticResource{
	StaticName:    "menu.svg",
	StaticContent: menuIcon,
}

//go:embed icons/menu-expand.svg
var menuexpandIcon []byte
var menuexpandIconRes = &fyne.StaticResource{
	StaticName:    "menu-expand.svg",
	StaticContent: menuexpandIcon,
}

//go:embed icons/check-box.svg
var checkboxIcon []byte
var checkboxIconRes = &fyne.StaticResource{
	StaticName:    "check-box.svg",
	StaticContent: checkboxIcon,
}

//go:embed icons/check-box-checked.svg
var checkboxcheckedIcon []byte
var checkboxcheckedIconRes = &fyne.StaticResource{
	StaticName:    "check-box-checked.svg",
	StaticContent: checkboxcheckedIcon,
}

//go:embed icons/check-box-fill.svg
var checkboxfillIcon []byte
var checkboxfillIconRes = &fyne.StaticResource{
	StaticName:    "check-box-fill.svg",
	StaticContent: checkboxfillIcon,
}

//go:embed icons/check-box-partial.svg
var checkboxpartialIcon []byte
var checkboxpartialIconRes = &fyne.StaticResource{
	StaticName:    "check-box-partial.svg",
	StaticContent: checkboxpartialIcon,
}

//go:embed icons/radio-button.svg
var radiobuttonIcon []byte
var radiobuttonIconRes = &fyne.StaticResource{
	StaticName:    "radio-button.svg",
	StaticContent: radiobuttonIcon,
}

//go:embed icons/radio-button-checked.svg
var radiobuttoncheckedIcon []byte
var radiobuttoncheckedIconRes = &fyne.StaticResource{
	StaticName:    "radio-button-checked.svg",
	StaticContent: radiobuttoncheckedIcon,
}

//go:embed icons/radio-button-fill.svg
var radiobuttonfillIcon []byte
var radiobuttonfillIconRes = &fyne.StaticResource{
	StaticName:    "radio-button-fill.svg",
	StaticContent: radiobuttonfillIcon,
}

//go:embed icons/content-add.svg
var contentaddIcon []byte
var contentaddIconRes = &fyne.StaticResource{
	StaticName:    "content-add.svg",
	StaticContent: contentaddIcon,
}

//go:embed icons/content-remove.svg
var contentremoveIcon []byte
var contentremoveIconRes = &fyne.StaticResource{
	StaticName:    "content-remove.svg",
	StaticContent: contentremoveIcon,
}

//go:embed icons/content-cut.svg
var contentcutIcon []byte
var contentcutIconRes = &fyne.StaticResource{
	StaticName:    "content-cut.svg",
	StaticContent: contentcutIcon,
}

//go:embed icons/content-copy.svg
var contentcopyIcon []byte
var contentcopyIconRes = &fyne.StaticResource{
	StaticName:    "content-copy.svg",
	StaticContent: contentcopyIcon,
}

//go:embed icons/content-paste.svg
var contentpasteIcon []byte
var contentpasteIconRes = &fyne.StaticResource{
	StaticName:    "content-paste.svg",
	StaticContent: contentpasteIcon,
}

//go:embed icons/content-redo.svg
var contentredoIcon []byte
var contentredoIconRes = &fyne.StaticResource{
	StaticName:    "content-redo.svg",
	StaticContent: contentredoIcon,
}

//go:embed icons/content-undo.svg
var contentundoIcon []byte
var contentundoIconRes = &fyne.StaticResource{
	StaticName:    "content-undo.svg",
	StaticContent: contentundoIcon,
}

//go:embed icons/color-achromatic.svg
var colorachromaticIcon []byte
var colorachromaticIconRes = &fyne.StaticResource{
	StaticName:    "color-achromatic.svg",
	StaticContent: colorachromaticIcon,
}

//go:embed icons/color-chromatic.svg
var colorchromaticIcon []byte
var colorchromaticIconRes = &fyne.StaticResource{
	StaticName:    "color-chromatic.svg",
	StaticContent: colorchromaticIcon,
}

//go:embed icons/color-palette.svg
var colorpaletteIcon []byte
var colorpaletteIconRes = &fyne.StaticResource{
	StaticName:    "color-palette.svg",
	StaticContent: colorpaletteIcon,
}

//go:embed icons/document.svg
var documentIcon []byte
var documentIconRes = &fyne.StaticResource{
	StaticName:    "document.svg",
	StaticContent: documentIcon,
}

//go:embed icons/document-create.svg
var documentcreateIcon []byte
var documentcreateIconRes = &fyne.StaticResource{
	StaticName:    "document-create.svg",
	StaticContent: documentcreateIcon,
}

//go:embed icons/document-print.svg
var documentprintIcon []byte
var documentprintIconRes = &fyne.StaticResource{
	StaticName:    "document-print.svg",
	StaticContent: documentprintIcon,
}

//go:embed icons/document-save.svg
var documentsaveIcon []byte
var documentsaveIconRes = &fyne.StaticResource{
	StaticName:    "document-save.svg",
	StaticContent: documentsaveIcon,
}

//go:embed icons/drag-corner-indicator.svg
var dragcornerindicatorIcon []byte
var dragcornerindicatorIconRes = &fyne.StaticResource{
	StaticName:    "drag-corner-indicator.svg",
	StaticContent: dragcornerindicatorIcon,
}

//go:embed icons/more-horizontal.svg
var morehorizontalIcon []byte
var morehorizontalIconRes = &fyne.StaticResource{
	StaticName:    "more-horizontal.svg",
	StaticContent: morehorizontalIcon,
}

//go:embed icons/more-vertical.svg
var moreverticalIcon []byte
var moreverticalIconRes = &fyne.StaticResource{
	StaticName:    "more-vertical.svg",
	StaticContent: moreverticalIcon,
}

//go:embed icons/info.svg
var infoIcon []byte
var infoIconRes = &fyne.StaticResource{
	StaticName:    "info.svg",
	StaticContent: infoIcon,
}

//go:embed icons/question.svg
var questionIcon []byte
var questionIconRes = &fyne.StaticResource{
	StaticName:    "question.svg",
	StaticContent: questionIcon,
}

//go:embed icons/warning.svg
var warningIcon []byte
var warningIconRes = &fyne.StaticResource{
	StaticName:    "warning.svg",
	StaticContent: warningIcon,
}

//go:embed icons/error.svg
var errorIcon []byte
var errorIconRes = &fyne.StaticResource{
	StaticName:    "error.svg",
	StaticContent: errorIcon,
}

//go:embed icons/broken-image.svg
var brokenimageIcon []byte
var brokenimageIconRes = &fyne.StaticResource{
	StaticName:    "broken-image.svg",
	StaticContent: brokenimageIcon,
}

//go:embed icons/arrow-back.svg
var arrowbackIcon []byte
var arrowbackIconRes = &fyne.StaticResource{
	StaticName:    "arrow-back.svg",
	StaticContent: arrowbackIcon,
}

//go:embed icons/arrow-down.svg
var arrowdownIcon []byte
var arrowdownIconRes = &fyne.StaticResource{
	StaticName:    "arrow-down.svg",
	StaticContent: arrowdownIcon,
}

//go:embed icons/arrow-forward.svg
var arrowforwardIcon []byte
var arrowforwardIconRes = &fyne.StaticResource{
	StaticName:    "arrow-forward.svg",
	StaticContent: arrowforwardIcon,
}

//go:embed icons/arrow-up.svg
var arrowupIcon []byte
var arrowupIconRes = &fyne.StaticResource{
	StaticName:    "arrow-up.svg",
	StaticContent: arrowupIcon,
}

//go:embed icons/arrow-drop-down.svg
var arrowdropdownIcon []byte
var arrowdropdownIconRes = &fyne.StaticResource{
	StaticName:    "arrow-drop-down.svg",
	StaticContent: arrowdropdownIcon,
}

//go:embed icons/arrow-drop-up.svg
var arrowdropupIcon []byte
var arrowdropupIconRes = &fyne.StaticResource{
	StaticName:    "arrow-drop-up.svg",
	StaticContent: arrowdropupIcon,
}

//go:embed icons/file.svg
var fileIcon []byte
var fileIconRes = &fyne.StaticResource{
	StaticName:    "file.svg",
	StaticContent: fileIcon,
}

//go:embed icons/file-application.svg
var fileapplicationIcon []byte
var fileapplicationIconRes = &fyne.StaticResource{
	StaticName:    "file-application.svg",
	StaticContent: fileapplicationIcon,
}

//go:embed icons/file-audio.svg
var fileaudioIcon []byte
var fileaudioIconRes = &fyne.StaticResource{
	StaticName:    "file-audio.svg",
	StaticContent: fileaudioIcon,
}

//go:embed icons/file-image.svg
var fileimageIcon []byte
var fileimageIconRes = &fyne.StaticResource{
	StaticName:    "file-image.svg",
	StaticContent: fileimageIcon,
}

//go:embed icons/file-text.svg
var filetextIcon []byte
var filetextIconRes = &fyne.StaticResource{
	StaticName:    "file-text.svg",
	StaticContent: filetextIcon,
}

//go:embed icons/file-video.svg
var filevideoIcon []byte
var filevideoIconRes = &fyne.StaticResource{
	StaticName:    "file-video.svg",
	StaticContent: filevideoIcon,
}

//go:embed icons/folder.svg
var folderIcon []byte
var folderIconRes = &fyne.StaticResource{
	StaticName:    "folder.svg",
	StaticContent: folderIcon,
}

//go:embed icons/folder-new.svg
var foldernewIcon []byte
var foldernewIconRes = &fyne.StaticResource{
	StaticName:    "folder-new.svg",
	StaticContent: foldernewIcon,
}

//go:embed icons/folder-open.svg
var folderopenIcon []byte
var folderopenIconRes = &fyne.StaticResource{
	StaticName:    "folder-open.svg",
	StaticContent: folderopenIcon,
}

//go:embed icons/help.svg
var helpIcon []byte
var helpIconRes = &fyne.StaticResource{
	StaticName:    "help.svg",
	StaticContent: helpIcon,
}

//go:embed icons/history.svg
var historyIcon []byte
var historyIconRes = &fyne.StaticResource{
	StaticName:    "history.svg",
	StaticContent: historyIcon,
}

//go:embed icons/home.svg
var homeIcon []byte
var homeIconRes = &fyne.StaticResource{
	StaticName:    "home.svg",
	StaticContent: homeIcon,
}

//go:embed icons/settings.svg
var settingsIcon []byte
var settingsIconRes = &fyne.StaticResource{
	StaticName:    "settings.svg",
	StaticContent: settingsIcon,
}

//go:embed icons/mail-attachment.svg
var mailattachmentIcon []byte
var mailattachmentIconRes = &fyne.StaticResource{
	StaticName:    "mail-attachment.svg",
	StaticContent: mailattachmentIcon,
}

//go:embed icons/mail-compose.svg
var mailcomposeIcon []byte
var mailcomposeIconRes = &fyne.StaticResource{
	StaticName:    "mail-compose.svg",
	StaticContent: mailcomposeIcon,
}

//go:embed icons/mail-forward.svg
var mailforwardIcon []byte
var mailforwardIconRes = &fyne.StaticResource{
	StaticName:    "mail-forward.svg",
	StaticContent: mailforwardIcon,
}

//go:embed icons/mail-reply.svg
var mailreplyIcon []byte
var mailreplyIconRes = &fyne.StaticResource{
	StaticName:    "mail-reply.svg",
	StaticContent: mailreplyIcon,
}

//go:embed icons/mail-reply_all.svg
var mailreplyallIcon []byte
var mailreplyallIconRes = &fyne.StaticResource{
	StaticName:    "mail-reply_all.svg",
	StaticContent: mailreplyallIcon,
}

//go:embed icons/mail-send.svg
var mailsendIcon []byte
var mailsendIconRes = &fyne.StaticResource{
	StaticName:    "mail-send.svg",
	StaticContent: mailsendIcon,
}

//go:embed icons/media-music.svg
var mediamusicIcon []byte
var mediamusicIconRes = &fyne.StaticResource{
	StaticName:    "media-music.svg",
	StaticContent: mediamusicIcon,
}

//go:embed icons/media-photo.svg
var mediaphotoIcon []byte
var mediaphotoIconRes = &fyne.StaticResource{
	StaticName:    "media-photo.svg",
	StaticContent: mediaphotoIcon,
}

//go:embed icons/media-video.svg
var mediavideoIcon []byte
var mediavideoIconRes = &fyne.StaticResource{
	StaticName:    "media-video.svg",
	StaticContent: mediavideoIcon,
}

//go:embed icons/media-fast-forward.svg
var mediafastforwardIcon []byte
var mediafastforwardIconRes = &fyne.StaticResource{
	StaticName:    "media-fast-forward.svg",
	StaticContent: mediafastforwardIcon,
}

//go:embed icons/media-fast-rewind.svg
var mediafastrewindIcon []byte
var mediafastrewindIconRes = &fyne.StaticResource{
	StaticName:    "media-fast-rewind.svg",
	StaticContent: mediafastrewindIcon,
}

//go:embed icons/media-pause.svg
var mediapauseIcon []byte
var mediapauseIconRes = &fyne.StaticResource{
	StaticName:    "media-pause.svg",
	StaticContent: mediapauseIcon,
}

//go:embed icons/media-play.svg
var mediaplayIcon []byte
var mediaplayIconRes = &fyne.StaticResource{
	StaticName:    "media-play.svg",
	StaticContent: mediaplayIcon,
}

//go:embed icons/media-record.svg
var mediarecordIcon []byte
var mediarecordIconRes = &fyne.StaticResource{
	StaticName:    "media-record.svg",
	StaticContent: mediarecordIcon,
}

//go:embed icons/media-replay.svg
var mediareplayIcon []byte
var mediareplayIconRes = &fyne.StaticResource{
	StaticName:    "media-replay.svg",
	StaticContent: mediareplayIcon,
}

//go:embed icons/media-skip-next.svg
var mediaskipnextIcon []byte
var mediaskipnextIconRes = &fyne.StaticResource{
	StaticName:    "media-skip-next.svg",
	StaticContent: mediaskipnextIcon,
}

//go:embed icons/media-skip-previous.svg
var mediaskippreviousIcon []byte
var mediaskippreviousIconRes = &fyne.StaticResource{
	StaticName:    "media-skip-previous.svg",
	StaticContent: mediaskippreviousIcon,
}

//go:embed icons/media-stop.svg
var mediastopIcon []byte
var mediastopIconRes = &fyne.StaticResource{
	StaticName:    "media-stop.svg",
	StaticContent: mediastopIcon,
}

//go:embed icons/view-fullscreen.svg
var viewfullscreenIcon []byte
var viewfullscreenIconRes = &fyne.StaticResource{
	StaticName:    "view-fullscreen.svg",
	StaticContent: viewfullscreenIcon,
}

//go:embed icons/view-refresh.svg
var viewrefreshIcon []byte
var viewrefreshIconRes = &fyne.StaticResource{
	StaticName:    "view-refresh.svg",
	StaticContent: viewrefreshIcon,
}

//go:embed icons/view-zoom-fit.svg
var viewzoomfitIcon []byte
var viewzoomfitIconRes = &fyne.StaticResource{
	StaticName:    "view-zoom-fit.svg",
	StaticContent: viewzoomfitIcon,
}

//go:embed icons/view-zoom-in.svg
var viewzoominIcon []byte
var viewzoominIconRes = &fyne.StaticResource{
	StaticName:    "view-zoom-in.svg",
	StaticContent: viewzoominIcon,
}

//go:embed icons/view-zoom-out.svg
var viewzoomoutIcon []byte
var viewzoomoutIconRes = &fyne.StaticResource{
	StaticName:    "view-zoom-out.svg",
	StaticContent: viewzoomoutIcon,
}

//go:embed icons/volume-down.svg
var volumedownIcon []byte
var volumedownIconRes = &fyne.StaticResource{
	StaticName:    "volume-down.svg",
	StaticContent: volumedownIcon,
}

//go:embed icons/volume-mute.svg
var volumemuteIcon []byte
var volumemuteIconRes = &fyne.StaticResource{
	StaticName:    "volume-mute.svg",
	StaticContent: volumemuteIcon,
}

//go:embed icons/volume-up.svg
var volumeupIcon []byte
var volumeupIconRes = &fyne.StaticResource{
	StaticName:    "volume-up.svg",
	StaticContent: volumeupIcon,
}

//go:embed icons/visibility.svg
var visibilityIcon []byte
var visibilityIconRes = &fyne.StaticResource{
	StaticName:    "visibility.svg",
	StaticContent: visibilityIcon,
}

//go:embed icons/visibility-off.svg
var visibilityoffIcon []byte
var visibilityoffIconRes = &fyne.StaticResource{
	StaticName:    "visibility-off.svg",
	StaticContent: visibilityoffIcon,
}

//go:embed icons/download.svg
var downloadIcon []byte
var downloadIconRes = &fyne.StaticResource{
	StaticName:    "download.svg",
	StaticContent: downloadIcon,
}

//go:embed icons/computer.svg
var computerIcon []byte
var computerIconRes = &fyne.StaticResource{
	StaticName:    "computer.svg",
	StaticContent: computerIcon,
}

//go:embed icons/desktop.svg
var desktopIcon []byte
var desktopIconRes = &fyne.StaticResource{
	StaticName:    "desktop.svg",
	StaticContent: desktopIcon,
}

//go:embed icons/storage.svg
var storageIcon []byte
var storageIconRes = &fyne.StaticResource{
	StaticName:    "storage.svg",
	StaticContent: storageIcon,
}

//go:embed icons/upload.svg
var uploadIcon []byte
var uploadIconRes = &fyne.StaticResource{
	StaticName:    "upload.svg",
	StaticContent: uploadIcon,
}

//go:embed icons/account.svg
var accountIcon []byte
var accountIconRes = &fyne.StaticResource{
	StaticName:    "account.svg",
	StaticContent: accountIcon,
}

//go:embed icons/calendar.svg
var calendarIcon []byte
var calendarIconRes = &fyne.StaticResource{
	StaticName:    "calendar.svg",
	StaticContent: calendarIcon,
}

//go:embed icons/login.svg
var loginIcon []byte
var loginIconRes = &fyne.StaticResource{
	StaticName:    "login.svg",
	StaticContent: loginIcon,
}

//go:embed icons/logout.svg
var logoutIcon []byte
var logoutIconRes = &fyne.StaticResource{
	StaticName:    "logout.svg",
	StaticContent: logoutIcon,
}

//go:embed icons/list.svg
var listIcon []byte
var listIconRes = &fyne.StaticResource{
	StaticName:    "list.svg",
	StaticContent: listIcon,
}

//go:embed icons/grid.svg
var gridIcon []byte
var gridIconRes = &fyne.StaticResource{
	StaticName:    "grid.svg",
	StaticContent: gridIcon,
}

//go:embed icons/maximize.svg
var maximizeIcon []byte
var maximizeIconRes = &fyne.StaticResource{
	StaticName:    "maximize.svg",
	StaticContent: maximizeIcon,
}

//go:embed icons/minimize.svg
var minimizeIcon []byte
var minimizeIconRes = &fyne.StaticResource{
	StaticName:    "minimize.svg",
	StaticContent: minimizeIcon,
}
