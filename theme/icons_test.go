package theme

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func helperNewStaticResource() *fyne.StaticResource {
	return &fyne.StaticResource{
		StaticName: "content-remove.svg",
		StaticContent: []byte{
			60, 115, 118, 103, 32, 120, 109, 108, 110, 115, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 119, 51, 46, 111, 114, 103, 47, 50, 48, 48, 48, 47, 115, 118, 103, 34, 32, 119, 105, 100, 116, 104, 61, 34, 50, 52, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 50, 52, 34, 32, 118, 105, 101, 119, 66, 111, 120, 61, 34, 48, 32, 48, 32, 50, 52, 32, 50, 52, 34, 62, 60, 112, 97, 116, 104, 32, 102, 105, 108, 108, 61, 34, 35, 102, 102, 102, 102, 102, 102, 34, 32, 100, 61, 34, 77, 49, 57, 32, 49, 51, 72, 53, 118, 45, 50, 104, 49, 52, 118, 50, 122, 34, 47, 62, 60, 112, 97, 116, 104, 32, 100, 61, 34, 77, 48, 32, 48, 104, 50, 52, 118, 50, 52, 72, 48, 122, 34, 32, 102, 105, 108, 108, 61, 34, 110, 111, 110, 101, 34, 47, 62, 60, 47, 115, 118, 103, 62},
	}
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // pathObj relative to this file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func TestIconThemeChangeDoesNotChangeName(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	cancel := CancelIcon()
	name := cancel.Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, name, cancel.Name())
}

func TestIconThemeChangeContent(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	cancel := CancelIcon()
	content := cancel.Content()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, cancel.Content())
}

func TestNewThemedResource_TwoStaticResourceSupport(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	custom := NewThemedResource(helperNewStaticResource(), helperNewStaticResource())
	content := custom.Content()
	nm := custom.Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, custom.Content())
	assert.Equal(t, nm, custom.Name())
}

func TestNewThemedResource_OneStaticResourceSupport(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	custom := NewThemedResource(helperNewStaticResource(), nil)
	content := custom.Content()
	nm := custom.Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, custom.Content())
	assert.Equal(t, nm, custom.Name())
}

func TestDynamicResource_Name(t *testing.T) {
	sr := fyne.NewStaticResource("cancel_Paths.svg",
		helperLoadBytes(t, "cancel_Paths.svg"))
	dr := &DynamicResource{
		BaseResource: sr,
	}
	assert.Equal(t, sr.Name(), dr.Name())
}

func TestDynamicResource_Content_NoGroupsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	sr := fyne.NewStaticResource("cancel_Paths.svg",
		helperLoadBytes(t, "cancel_Paths.svg"))
	dr := &DynamicResource{
		BaseResource: sr,
	}
	assert.NotEqual(t, sr.Content(), dr.Content())
}

func TestDynamicResource_Content_GroupPathFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	sr := fyne.NewStaticResource("check_GroupPaths.svg",
		helperLoadBytes(t, "check_GroupPaths.svg"))
	dr := &DynamicResource{
		BaseResource: sr,
	}
	assert.NotEqual(t, sr.Content(), dr.Content())
}

func TestDynamicResource_Content_GroupRectFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	sr := fyne.NewStaticResource("info_GroupRects.svg",
		helperLoadBytes(t, "info_GroupRects.svg"))
	dr := &DynamicResource{
		BaseResource: sr,
	}
	assert.NotEqual(t, sr.Content(), dr.Content())
}

func TestDynamicResource_Content_GroupPolygonsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	sr := fyne.NewStaticResource("warning_GroupPolygons.svg",
		helperLoadBytes(t, "warning_GroupPolygons.svg"))
	dr := &DynamicResource{
		BaseResource: sr,
	}
	assert.NotEqual(t, sr.Content(), dr.Content())
}

// Test Asset Sources
func Test_FyneLogo(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := FyneLogo().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "fyne.png", result)
}

func Test_CancelIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := CancelIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "cancel.svg", result)
}

func Test_ConfirmIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ConfirmIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "check.svg", result)
}

func Test_DeleteIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := DeleteIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "delete.svg", result)
}

func Test_SearchIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := SearchIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "search.svg", result)
}

func Test_SearchReplaceIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := SearchReplaceIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "search-replace.svg", result)
}

func Test_CheckButtonIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := CheckButtonIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "check-box-blank.svg", result)
}

func Test_CheckButtonCheckedIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := CheckButtonCheckedIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "check-box.svg", result)
}

func Test_RadioButtonIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := RadioButtonIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "radio-button.svg", result)
}

func Test_RadioButtonCheckedIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := RadioButtonCheckedIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "radio-button-checked.svg", result)
}

func Test_ContentAddIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentAddIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-add.svg", result)
}

func Test_ContentRemoveIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentRemoveIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-remove.svg", result)
}

func Test_ContentClearIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentClearIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "cancel.svg", result)
}

func Test_ContentCutIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentCutIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-cut.svg", result)
}

func Test_ContentCopyIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentCopyIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-copy.svg", result)
}

func Test_ContentPasteIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentPasteIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-paste.svg", result)
}

func Test_ContentRedoIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentRedoIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-redo.svg", result)
}

func Test_ContentUndoIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ContentUndoIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "content-undo.svg", result)
}

func Test_DocumentCreateIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := DocumentCreateIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "document-create.svg", result)
}

func Test_DocumentPrintIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := DocumentPrintIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "document-print.svg", result)
}

func Test_DocumentSaveIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := DocumentSaveIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "document-save.svg", result)
}

func Test_InfoIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := InfoIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "info.svg", result)
}

func Test_QuestionIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := QuestionIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "question.svg", result)
}

func Test_WarningIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := WarningIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "warning.svg", result)
}

func Test_FolderIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := FolderIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "folder.svg", result)
}

func Test_FolderNewIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := FolderNewIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "folder-new.svg", result)
}

func Test_FolderOpenIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := FolderOpenIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "folder-open.svg", result)
}

func Test_HelpIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := HelpIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "help.svg", result)
}

func Test_HomeIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := HomeIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "home.svg", result)
}

func Test_MailAttachmentIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MailAttachmentIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "mail-attachment.svg", result)
}

func Test_MailComposeIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MailComposeIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "mail-compose.svg", result)
}

func Test_MailForwardIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MailForwardIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "mail-forward.svg", result)
}

func Test_MailReplyIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MailReplyIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "mail-reply.svg", result)
}

func Test_MailReplyAllIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MailReplyAllIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "mail-reply_all.svg", result)
}

func Test_MailSendIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MailSendIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "mail-send.svg", result)
}

func Test_MoveDownIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MoveDownIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "arrow-down.svg", result)
}

func Test_MoveUpIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := MoveUpIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "arrow-up.svg", result)
}

func Test_NavigateBackIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := NavigateBackIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "arrow-back.svg", result)
}

func Test_NavigateNextIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := NavigateNextIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "arrow-forward.svg", result)
}

func Test_ViewFullScreenIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ViewFullScreenIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "view-fullscreen.svg", result)
}

func Test_ViewRestoreIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ViewRestoreIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "view-zoom-fit.svg", result)
}

func Test_ViewRefreshIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ViewRefreshIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "view-refresh.svg", result)
}

func Test_ZoomFitIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ZoomFitIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "view-zoom-fit.svg", result)
}

func Test_ZoomInIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ZoomInIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "view-zoom-in.svg", result)
}

func Test_ZoomOutIcon(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	result := ZoomOutIcon().Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.Equal(t, "view-zoom-out.svg", result)
}
