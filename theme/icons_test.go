package theme_test

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func init() {
	fyne.SetCurrentApp(&themedApp{})
}

func helperNewStaticResource() *fyne.StaticResource {
	return &fyne.StaticResource{
		StaticName: "content-remove.svg",
		StaticContent: []byte{
			60, 115, 118, 103, 32, 120, 109, 108, 110, 115, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 119, 51, 46, 111, 114, 103, 47, 50, 48, 48, 48, 47, 115, 118, 103, 34, 32, 119, 105, 100, 116, 104, 61, 34, 50, 52, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 50, 52, 34, 32, 118, 105, 101, 119, 66, 111, 120, 61, 34, 48, 32, 48, 32, 50, 52, 32, 50, 52, 34, 62, 60, 112, 97, 116, 104, 32, 102, 105, 108, 108, 61, 34, 35, 102, 102, 102, 102, 102, 102, 34, 32, 100, 61, 34, 77, 49, 57, 32, 49, 51, 72, 53, 118, 45, 50, 104, 49, 52, 118, 50, 122, 34, 47, 62, 60, 112, 97, 116, 104, 32, 100, 61, 34, 77, 48, 32, 48, 104, 50, 52, 118, 50, 52, 72, 48, 122, 34, 32, 102, 105, 108, 108, 61, 34, 110, 111, 110, 101, 34, 47, 62, 60, 47, 115, 118, 103, 62,
		},
	}
}

func helperLoadRes(t *testing.T, name string) fyne.Resource {
	path := filepath.Join("testdata", name) // pathObj relative to this file
	res, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func TestIconThemeChangeDoesNotChangeName(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	cancel := theme.CancelIcon()
	name := cancel.Name()

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	assert.Equal(t, name, cancel.Name())
}

func TestIconThemeChangeContent(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	cancel := theme.CancelIcon()
	content := cancel.Content()

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	assert.NotEqual(t, content, cancel.Content())
}

func TestIconThemeResourceWithViewboxOnly(t *testing.T) {
	r1 := &fyne.StaticResource{
		StaticName:    "unsplash.svg",
		StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M7.5 6.75V0h9v6.75h-9zm9 3.75H24V24H0V10.5h7.5v6.75h9V10.5z"></path></svg>`),
	}
	r2 := theme.NewThemedResource(r1)
	re := regexp.MustCompile(`\s+fill\S+`)
	content := re.ReplaceAll(r2.Content(), []byte(""))

	assert.Equal(t, r1.Content(), content)
}

func TestNewThemedResource_StaticResourceSupport(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	custom := theme.NewThemedResource(helperNewStaticResource())
	content := custom.Content()
	name := custom.Name()

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	assert.NotEqual(t, content, custom.Content())
	assert.Equal(t, name, custom.Name())
}

func TestNewColoredResource(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	source := helperNewStaticResource()
	custom := theme.NewColoredResource(source, theme.ColorNameSuccess)

	assert.Equal(t, theme.ColorNameSuccess, custom.ColorName)
	assert.Equal(t, custom.Name(), fmt.Sprintf("success_%v", source.Name()))

	custom = theme.NewColoredResource(source, theme.ColorNamePrimary)
	assert.Equal(t, custom.Name(), fmt.Sprintf("primary_%v", source.Name()))
}

func TestNewDisabledResource(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	source := helperNewStaticResource()
	custom := theme.NewDisabledResource(source)
	name := custom.Name()

	assert.Equal(t, name, fmt.Sprintf("disabled_%v", source.Name()))
}

func TestThemedResource_Invert(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	inverted := theme.NewInvertedThemedResource(staticResource)
	assert.Equal(t, "inverted_"+staticResource.Name(), inverted.Name())
}

func TestThemedResource_Name(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	themedResource := theme.NewThemedResource(staticResource)
	assert.Equal(t, "foreground_"+staticResource.Name(), themedResource.Name())
}

func TestThemedResource_Content_NoGroupsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	themedResource := theme.NewThemedResource(staticResource)
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupPathFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	staticResource := helperLoadRes(t, "check_GroupPaths.svg")
	themedResource := theme.NewThemedResource(staticResource)
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupRectFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	staticResource := helperLoadRes(t, "info_GroupRects.svg")
	themedResource := theme.NewThemedResource(staticResource)
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupPolygonsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	staticResource := helperLoadRes(t, "warning_GroupPolygons.svg")
	themedResource := theme.NewThemedResource(staticResource)
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

// a black svg object omits the fill tag, this checks if it still properly updated
func TestThemedResource_Content_BlackFillIsUpdated(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	staticResource := helperLoadRes(t, "cancel_PathsBlackFill.svg")
	themedResource := theme.NewThemedResource(staticResource)
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Error(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	source := helperNewStaticResource()
	custom := theme.NewThemedResource(source)
	custom.ColorName = theme.ColorNameError

	assert.Equal(t, custom.Name(), fmt.Sprintf("error_%v", source.Name()))
	custom2 := theme.NewErrorThemedResource(source)
	assert.Equal(t, custom2.Name(), fmt.Sprintf("error_%v", source.Name()))
}

func TestThemedResource_Success(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	source := helperNewStaticResource()
	custom := theme.NewThemedResource(source)
	custom.ColorName = theme.ColorNameSuccess

	assert.Equal(t, custom.Name(), fmt.Sprintf("success_%v", source.Name()))

	custom = theme.NewSuccessThemedResource(source)
	assert.Equal(t, custom.Name(), fmt.Sprintf("success_%v", source.Name()))
}

func TestThemedResource_Warning(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	source := helperNewStaticResource()
	custom := theme.NewThemedResource(source)
	custom.ColorName = theme.ColorNameWarning

	assert.Equal(t, custom.Name(), fmt.Sprintf("warning_%v", source.Name()))
	custom = theme.NewWarningThemedResource(source)
	assert.Equal(t, custom.Name(), fmt.Sprintf("warning_%v", source.Name()))
}

func TestDisabledResource_Name(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	disabledResource := theme.NewDisabledResource(staticResource)
	assert.Equal(t, fmt.Sprintf("disabled_%v", staticResource.Name()), disabledResource.Name())
}

func TestDisabledResource_Content_NoGroupsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	disabledResource := theme.NewDisabledResource(staticResource)
	assert.NotEqual(t, staticResource.Content(), disabledResource.Content())
}

// Test Asset Sources
func Test_FyneLogo_FileSource(t *testing.T) {
	result := theme.FyneLogo().Name()
	assert.Equal(t, "fyne.png", result)
}

func Test_CancelIcon_FileSource(t *testing.T) {
	result := theme.CancelIcon().Name()
	assert.Equal(t, "foreground_cancel.svg", result)
}

func Test_ConfirmIcon_FileSource(t *testing.T) {
	result := theme.ConfirmIcon().Name()
	assert.Equal(t, "foreground_check.svg", result)
}

func Test_DeleteIcon_FileSource(t *testing.T) {
	result := theme.DeleteIcon().Name()
	assert.Equal(t, "foreground_delete.svg", result)
}

func Test_SearchIcon_FileSource(t *testing.T) {
	result := theme.SearchIcon().Name()
	assert.Equal(t, "foreground_search.svg", result)
}

func Test_SearchReplaceIcon_FileSource(t *testing.T) {
	result := theme.SearchReplaceIcon().Name()
	assert.Equal(t, "foreground_search-replace.svg", result)
}

func Test_CheckButtonIcon_FileSource(t *testing.T) {
	result := theme.CheckButtonIcon().Name()
	assert.Equal(t, "foreground_check-box.svg", result)
}

func Test_CheckButtonCheckedIcon_FileSource(t *testing.T) {
	result := theme.CheckButtonCheckedIcon().Name()
	assert.Equal(t, "foreground_check-box-checked.svg", result)
}

func Test_RadioButtonIcon_FileSource(t *testing.T) {
	result := theme.RadioButtonIcon().Name()
	assert.Equal(t, "foreground_radio-button.svg", result)
}

func Test_RadioButtonCheckedIcon_FileSource(t *testing.T) {
	result := theme.RadioButtonCheckedIcon().Name()
	assert.Equal(t, "foreground_radio-button-checked.svg", result)
}

func Test_ContentAddIcon_FileSource(t *testing.T) {
	result := theme.ContentAddIcon().Name()
	assert.Equal(t, "foreground_content-add.svg", result)
}

func Test_ContentRemoveIcon_FileSource(t *testing.T) {
	result := theme.ContentRemoveIcon().Name()
	assert.Equal(t, "foreground_content-remove.svg", result)
}

func Test_ContentClearIcon_FileSource(t *testing.T) {
	result := theme.ContentClearIcon().Name()
	assert.Equal(t, "foreground_cancel.svg", result)
}

func Test_ContentCutIcon_FileSource(t *testing.T) {
	result := theme.ContentCutIcon().Name()
	assert.Equal(t, "foreground_content-cut.svg", result)
}

func Test_ContentCopyIcon_FileSource(t *testing.T) {
	result := theme.ContentCopyIcon().Name()
	assert.Equal(t, "foreground_content-copy.svg", result)
}

func Test_ContentPasteIcon_FileSource(t *testing.T) {
	result := theme.ContentPasteIcon().Name()
	assert.Equal(t, "foreground_content-paste.svg", result)
}

func Test_ContentRedoIcon_FileSource(t *testing.T) {
	result := theme.ContentRedoIcon().Name()
	assert.Equal(t, "foreground_content-redo.svg", result)
}

func Test_ContentUndoIcon_FileSource(t *testing.T) {
	result := theme.ContentUndoIcon().Name()
	assert.Equal(t, "foreground_content-undo.svg", result)
}

func Test_DocumentCreateIcon_FileSource(t *testing.T) {
	result := theme.DocumentCreateIcon().Name()
	assert.Equal(t, "foreground_document-create.svg", result)
}

func Test_DocumentPrintIcon_FileSource(t *testing.T) {
	result := theme.DocumentPrintIcon().Name()
	assert.Equal(t, "foreground_document-print.svg", result)
}

func Test_DocumentSaveIcon_FileSource(t *testing.T) {
	result := theme.DocumentSaveIcon().Name()
	assert.Equal(t, "foreground_document-save.svg", result)
}

func Test_InfoIcon_FileSource(t *testing.T) {
	result := theme.InfoIcon().Name()
	assert.Equal(t, "foreground_info.svg", result)
}

func Test_QuestionIcon_FileSource(t *testing.T) {
	result := theme.QuestionIcon().Name()
	assert.Equal(t, "foreground_question.svg", result)
}

func Test_WarningIcon_FileSource(t *testing.T) {
	result := theme.WarningIcon().Name()
	assert.Equal(t, "foreground_warning.svg", result)
}

func Test_BrokenImageIcon_FileSource(t *testing.T) {
	result := theme.BrokenImageIcon().Name()
	assert.Equal(t, "foreground_broken-image.svg", result)
}

func Test_FolderIcon_FileSource(t *testing.T) {
	result := theme.FolderIcon().Name()
	assert.Equal(t, "foreground_folder.svg", result)
}

func Test_FolderNewIcon_FileSource(t *testing.T) {
	result := theme.FolderNewIcon().Name()
	assert.Equal(t, "foreground_folder-new.svg", result)
}

func Test_FolderOpenIcon_FileSource(t *testing.T) {
	result := theme.FolderOpenIcon().Name()
	assert.Equal(t, "foreground_folder-open.svg", result)
}

func Test_HelpIcon_FileSource(t *testing.T) {
	result := theme.HelpIcon().Name()
	assert.Equal(t, "foreground_help.svg", result)
}

func Test_HomeIcon_FileSource(t *testing.T) {
	result := theme.HomeIcon().Name()
	assert.Equal(t, "foreground_home.svg", result)
}

func Test_SettingsIcon_FileSource(t *testing.T) {
	result := theme.SettingsIcon().Name()
	assert.Equal(t, "foreground_settings.svg", result)
}

func Test_MailAttachmentIcon_FileSource(t *testing.T) {
	result := theme.MailAttachmentIcon().Name()
	assert.Equal(t, "foreground_mail-attachment.svg", result)
}

func Test_MailComposeIcon_FileSource(t *testing.T) {
	result := theme.MailComposeIcon().Name()
	assert.Equal(t, "foreground_mail-compose.svg", result)
}

func Test_MailForwardIcon_FileSource(t *testing.T) {
	result := theme.MailForwardIcon().Name()
	assert.Equal(t, "foreground_mail-forward.svg", result)
}

func Test_MailReplyIcon_FileSource(t *testing.T) {
	result := theme.MailReplyIcon().Name()
	assert.Equal(t, "foreground_mail-reply.svg", result)
}

func Test_MailReplyAllIcon_FileSource(t *testing.T) {
	result := theme.MailReplyAllIcon().Name()
	assert.Equal(t, "foreground_mail-reply_all.svg", result)
}

func Test_MailSendIcon_FileSource(t *testing.T) {
	result := theme.MailSendIcon().Name()
	assert.Equal(t, "foreground_mail-send.svg", result)
}

func Test_MoveDownIcon_FileSource(t *testing.T) {
	result := theme.MoveDownIcon().Name()
	assert.Equal(t, "foreground_arrow-down.svg", result)
}

func Test_MoveUpIcon_FileSource(t *testing.T) {
	result := theme.MoveUpIcon().Name()
	assert.Equal(t, "foreground_arrow-up.svg", result)
}

func Test_NavigateBackIcon_FileSource(t *testing.T) {
	result := theme.NavigateBackIcon().Name()
	assert.Equal(t, "foreground_arrow-back.svg", result)
}

func Test_NavigateNextIcon_FileSource(t *testing.T) {
	result := theme.NavigateNextIcon().Name()
	assert.Equal(t, "foreground_arrow-forward.svg", result)
}

func Test_ViewFullScreenIcon_FileSource(t *testing.T) {
	result := theme.ViewFullScreenIcon().Name()
	assert.Equal(t, "foreground_view-fullscreen.svg", result)
}

func Test_ViewRestoreIcon_FileSource(t *testing.T) {
	result := theme.ViewRestoreIcon().Name()
	assert.Equal(t, "foreground_view-zoom-fit.svg", result)
}

func Test_ViewRefreshIcon_FileSource(t *testing.T) {
	result := theme.ViewRefreshIcon().Name()
	assert.Equal(t, "foreground_view-refresh.svg", result)
}

func Test_ZoomFitIcon_FileSource(t *testing.T) {
	result := theme.ZoomFitIcon().Name()
	assert.Equal(t, "foreground_view-zoom-fit.svg", result)
}

func Test_ZoomInIcon_FileSource(t *testing.T) {
	result := theme.ZoomInIcon().Name()
	assert.Equal(t, "foreground_view-zoom-in.svg", result)
}

func Test_ZoomOutIcon_FileSource(t *testing.T) {
	result := theme.ZoomOutIcon().Name()
	assert.Equal(t, "foreground_view-zoom-out.svg", result)
}

func Test_VisibilityIcon_FileSource(t *testing.T) {
	result := theme.VisibilityIcon().Name()
	assert.Equal(t, "foreground_visibility.svg", result)
}

func Test_VisibilityOffIcon_FileSource(t *testing.T) {
	result := theme.VisibilityOffIcon().Name()
	assert.Equal(t, "foreground_visibility-off.svg", result)
}

func Test_AccountIcon_FileSource(t *testing.T) {
	result := theme.AccountIcon().Name()
	assert.Equal(t, "foreground_account.svg", result)
}

func Test_LoginIcon_FileSource(t *testing.T) {
	result := theme.LoginIcon().Name()
	assert.Equal(t, "foreground_login.svg", result)
}

func Test_LogoutIcon_FileSource(t *testing.T) {
	result := theme.LogoutIcon().Name()
	assert.Equal(t, "foreground_logout.svg", result)
}

func Test_ListIcon_FileSource(t *testing.T) {
	result := theme.ListIcon().Name()
	assert.Equal(t, "foreground_list.svg", result)
}

func Test_GridIcon_FileSource(t *testing.T) {
	result := theme.GridIcon().Name()
	assert.Equal(t, "foreground_grid.svg", result)
}
