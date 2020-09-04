package theme

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/test"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	fyne.SetCurrentApp(&themedApp{})
}

func helperNewStaticResource() *fyne.StaticResource {
	return &fyne.StaticResource{
		StaticName: "content-remove.svg",
		StaticContent: []byte{
			60, 115, 118, 103, 32, 120, 109, 108, 110, 115, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 119, 51, 46, 111, 114, 103, 47, 50, 48, 48, 48, 47, 115, 118, 103, 34, 32, 119, 105, 100, 116, 104, 61, 34, 50, 52, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 50, 52, 34, 32, 118, 105, 101, 119, 66, 111, 120, 61, 34, 48, 32, 48, 32, 50, 52, 32, 50, 52, 34, 62, 60, 112, 97, 116, 104, 32, 102, 105, 108, 108, 61, 34, 35, 102, 102, 102, 102, 102, 102, 34, 32, 100, 61, 34, 77, 49, 57, 32, 49, 51, 72, 53, 118, 45, 50, 104, 49, 52, 118, 50, 122, 34, 47, 62, 60, 112, 97, 116, 104, 32, 100, 61, 34, 77, 48, 32, 48, 104, 50, 52, 118, 50, 52, 72, 48, 122, 34, 32, 102, 105, 108, 108, 61, 34, 110, 111, 110, 101, 34, 47, 62, 60, 47, 115, 118, 103, 62},
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

func helperDrawSVG(t *testing.T, data []byte) image.Image {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(data))
	require.NoError(t, err, "failed to read SVG data")

	width := int(icon.ViewBox.W) * 2
	height := int(icon.ViewBox.H) * 2
	icon.SetTarget(0, 0, float64(width), float64(height))
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1)
	return img
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
	name := custom.Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, custom.Content())
	assert.Equal(t, name, custom.Name())
}

func TestNewThemedResource_OneStaticResourceSupport(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	custom := NewThemedResource(helperNewStaticResource(), nil)
	content := custom.Content()
	name := custom.Name()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, content, custom.Content())
	assert.Equal(t, name, custom.Name())
}

func TestNewDisabledResource(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	source := helperNewStaticResource()
	custom := NewDisabledResource(source)
	name := custom.Name()

	assert.Equal(t, name, fmt.Sprintf("disabled_%v", source.Name()))
}

func TestThemedResource_Invert(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	inverted := NewInvertedThemedResource(staticResource)
	assert.Equal(t, "inverted-"+staticResource.Name(), inverted.Name())
}

func TestThemedResource_Name(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.Equal(t, staticResource.Name(), themedResource.Name())
}

func TestThemedResource_Content_NoGroupsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupPathFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "check_GroupPaths.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupRectFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "info_GroupRects.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestThemedResource_Content_GroupPolygonsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "warning_GroupPolygons.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

// a black svg object omits the fill tag, this checks it it still properly updated
func TestThemedResource_Content_BlackFillIsUpdated(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "cancel_PathsBlackFill.svg")
	themedResource := &ThemedResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), themedResource.Content())
}

func TestDisabledResource_Name(t *testing.T) {
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	disabledResource := &DisabledResource{
		source: staticResource,
	}
	assert.Equal(t, fmt.Sprintf("disabled_%v", staticResource.Name()), disabledResource.Name())
}

func TestDisabledResource_Content_NoGroupsFile(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	staticResource := helperLoadRes(t, "cancel_Paths.svg")
	disabledResource := &DisabledResource{
		source: staticResource,
	}
	assert.NotEqual(t, staticResource.Content(), disabledResource.Content())
}

func TestColorizeResource(t *testing.T) {
	tests := map[string]struct {
		svgFile   string
		color     color.Color
		wantImage string
	}{
		"paths": {
			svgFile:   "cancel_Paths.svg",
			color:     color.RGBA{R: 100, G: 100, A: 200},
			wantImage: "colorized/paths.png",
		},
		"circles": {
			svgFile:   "circles.svg",
			color:     color.RGBA{R: 100, B: 100, A: 200},
			wantImage: "colorized/circles.png",
		},
		"polygons": {
			svgFile:   "polygons.svg",
			color:     color.RGBA{G: 100, B: 100, A: 200},
			wantImage: "colorized/polygons.png",
		},
		"rects": {
			svgFile:   "rects.svg",
			color:     color.RGBA{R: 100, G: 100, B: 100, A: 200},
			wantImage: "colorized/rects.png",
		},
		"group of paths": {
			svgFile:   "check_GroupPaths.svg",
			color:     color.RGBA{R: 100, G: 100, A: 100},
			wantImage: "colorized/group_paths.png",
		},
		"group of circles": {
			svgFile:   "group_circles.svg",
			color:     color.RGBA{R: 100, B: 100, A: 100},
			wantImage: "colorized/group_circles.png",
		},
		"group of polygons": {
			svgFile:   "warning_GroupPolygons.svg",
			color:     color.RGBA{G: 100, B: 100, A: 100},
			wantImage: "colorized/group_polygons.png",
		},
		"group of rects": {
			svgFile:   "info_GroupRects.svg",
			color:     color.RGBA{R: 100, G: 100, B: 100, A: 100},
			wantImage: "colorized/group_rects.png",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			src := helperLoadRes(t, tt.svgFile)
			got := helperDrawSVG(t, colorizeResource(src, tt.color))
			test.AssertImageMatches(t, tt.wantImage, got)
		})
	}
}

// Test Asset Sources
func Test_FyneLogo_FileSource(t *testing.T) {
	result := FyneLogo().Name()
	assert.Equal(t, "fyne.png", result)
}

func Test_CancelIcon_FileSource(t *testing.T) {
	result := CancelIcon().Name()
	assert.Equal(t, "cancel.svg", result)
}

func Test_ConfirmIcon_FileSource(t *testing.T) {
	result := ConfirmIcon().Name()
	assert.Equal(t, "check.svg", result)
}

func Test_DeleteIcon_FileSource(t *testing.T) {
	result := DeleteIcon().Name()
	assert.Equal(t, "delete.svg", result)
}

func Test_SearchIcon_FileSource(t *testing.T) {
	result := SearchIcon().Name()
	assert.Equal(t, "search.svg", result)
}

func Test_SearchReplaceIcon_FileSource(t *testing.T) {
	result := SearchReplaceIcon().Name()
	assert.Equal(t, "search-replace.svg", result)
}

func Test_CheckButtonIcon_FileSource(t *testing.T) {
	result := CheckButtonIcon().Name()
	assert.Equal(t, "check-box-blank.svg", result)
}

func Test_CheckButtonCheckedIcon_FileSource(t *testing.T) {
	result := CheckButtonCheckedIcon().Name()
	assert.Equal(t, "check-box.svg", result)
}

func Test_RadioButtonIcon_FileSource(t *testing.T) {
	result := RadioButtonIcon().Name()
	assert.Equal(t, "radio-button.svg", result)
}

func Test_RadioButtonCheckedIcon_FileSource(t *testing.T) {
	result := RadioButtonCheckedIcon().Name()
	assert.Equal(t, "radio-button-checked.svg", result)
}

func Test_ContentAddIcon_FileSource(t *testing.T) {
	result := ContentAddIcon().Name()
	assert.Equal(t, "content-add.svg", result)
}

func Test_ContentRemoveIcon_FileSource(t *testing.T) {
	result := ContentRemoveIcon().Name()
	assert.Equal(t, "content-remove.svg", result)
}

func Test_ContentClearIcon_FileSource(t *testing.T) {
	result := ContentClearIcon().Name()
	assert.Equal(t, "cancel.svg", result)
}

func Test_ContentCutIcon_FileSource(t *testing.T) {
	result := ContentCutIcon().Name()
	assert.Equal(t, "content-cut.svg", result)
}

func Test_ContentCopyIcon_FileSource(t *testing.T) {
	result := ContentCopyIcon().Name()
	assert.Equal(t, "content-copy.svg", result)
}

func Test_ContentPasteIcon_FileSource(t *testing.T) {
	result := ContentPasteIcon().Name()
	assert.Equal(t, "content-paste.svg", result)
}

func Test_ContentRedoIcon_FileSource(t *testing.T) {
	result := ContentRedoIcon().Name()
	assert.Equal(t, "content-redo.svg", result)
}

func Test_ContentUndoIcon_FileSource(t *testing.T) {
	result := ContentUndoIcon().Name()
	assert.Equal(t, "content-undo.svg", result)
}

func Test_DocumentCreateIcon_FileSource(t *testing.T) {
	result := DocumentCreateIcon().Name()
	assert.Equal(t, "document-create.svg", result)
}

func Test_DocumentPrintIcon_FileSource(t *testing.T) {
	result := DocumentPrintIcon().Name()
	assert.Equal(t, "document-print.svg", result)
}

func Test_DocumentSaveIcon_FileSource(t *testing.T) {
	result := DocumentSaveIcon().Name()
	assert.Equal(t, "document-save.svg", result)
}

func Test_InfoIcon_FileSource(t *testing.T) {
	result := InfoIcon().Name()
	assert.Equal(t, "info.svg", result)
}

func Test_QuestionIcon_FileSource(t *testing.T) {
	result := QuestionIcon().Name()
	assert.Equal(t, "question.svg", result)
}

func Test_WarningIcon_FileSource(t *testing.T) {
	result := WarningIcon().Name()
	assert.Equal(t, "warning.svg", result)
}

func Test_FolderIcon_FileSource(t *testing.T) {
	result := FolderIcon().Name()
	assert.Equal(t, "folder.svg", result)
}

func Test_FolderNewIcon_FileSource(t *testing.T) {
	result := FolderNewIcon().Name()
	assert.Equal(t, "folder-new.svg", result)
}

func Test_FolderOpenIcon_FileSource(t *testing.T) {
	result := FolderOpenIcon().Name()
	assert.Equal(t, "folder-open.svg", result)
}

func Test_HelpIcon_FileSource(t *testing.T) {
	result := HelpIcon().Name()
	assert.Equal(t, "help.svg", result)
}

func Test_HomeIcon_FileSource(t *testing.T) {
	result := HomeIcon().Name()
	assert.Equal(t, "home.svg", result)
}

func Test_SettingsIcon_FileSource(t *testing.T) {
	result := SettingsIcon().Name()
	assert.Equal(t, "settings.svg", result)
}

func Test_MailAttachmentIcon_FileSource(t *testing.T) {
	result := MailAttachmentIcon().Name()
	assert.Equal(t, "mail-attachment.svg", result)
}

func Test_MailComposeIcon_FileSource(t *testing.T) {
	result := MailComposeIcon().Name()
	assert.Equal(t, "mail-compose.svg", result)
}

func Test_MailForwardIcon_FileSource(t *testing.T) {
	result := MailForwardIcon().Name()
	assert.Equal(t, "mail-forward.svg", result)
}

func Test_MailReplyIcon_FileSource(t *testing.T) {
	result := MailReplyIcon().Name()
	assert.Equal(t, "mail-reply.svg", result)
}

func Test_MailReplyAllIcon_FileSource(t *testing.T) {
	result := MailReplyAllIcon().Name()
	assert.Equal(t, "mail-reply_all.svg", result)
}

func Test_MailSendIcon_FileSource(t *testing.T) {
	result := MailSendIcon().Name()
	assert.Equal(t, "mail-send.svg", result)
}

func Test_MoveDownIcon_FileSource(t *testing.T) {
	result := MoveDownIcon().Name()
	assert.Equal(t, "arrow-down.svg", result)
}

func Test_MoveUpIcon_FileSource(t *testing.T) {
	result := MoveUpIcon().Name()
	assert.Equal(t, "arrow-up.svg", result)
}

func Test_NavigateBackIcon_FileSource(t *testing.T) {
	result := NavigateBackIcon().Name()
	assert.Equal(t, "arrow-back.svg", result)
}

func Test_NavigateNextIcon_FileSource(t *testing.T) {
	result := NavigateNextIcon().Name()
	assert.Equal(t, "arrow-forward.svg", result)
}

func Test_ViewFullScreenIcon_FileSource(t *testing.T) {
	result := ViewFullScreenIcon().Name()
	assert.Equal(t, "view-fullscreen.svg", result)
}

func Test_ViewRestoreIcon_FileSource(t *testing.T) {
	result := ViewRestoreIcon().Name()
	assert.Equal(t, "view-zoom-fit.svg", result)
}

func Test_ViewRefreshIcon_FileSource(t *testing.T) {
	result := ViewRefreshIcon().Name()
	assert.Equal(t, "view-refresh.svg", result)
}

func Test_ZoomFitIcon_FileSource(t *testing.T) {
	result := ZoomFitIcon().Name()
	assert.Equal(t, "view-zoom-fit.svg", result)
}

func Test_ZoomInIcon_FileSource(t *testing.T) {
	result := ZoomInIcon().Name()
	assert.Equal(t, "view-zoom-in.svg", result)
}

func Test_ZoomOutIcon_FileSource(t *testing.T) {
	result := ZoomOutIcon().Name()
	assert.Equal(t, "view-zoom-out.svg", result)
}

func Test_VisibilityIcon_FileSource(t *testing.T) {
	result := VisibilityIcon().Name()
	assert.Equal(t, "visibility.svg", result)
}

func Test_VisibilityOffIcon_FileSource(t *testing.T) {
	result := VisibilityOffIcon().Name()
	assert.Equal(t, "visibility-off.svg", result)
}
