package test

import (
	"fmt"
	"image/color"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	col "fyne.io/fyne/v2/internal/color"
	intdriver "fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

type markupRenderer struct {
	indentation int
	w           strings.Builder
}

// snapshot creates a new snapshot of the current render tree.
func snapshot(c fyne.Canvas) string {
	r := markupRenderer{}
	r.writeCanvas(c)
	return r.w.String()
}

func (r *markupRenderer) setAlignmentAttr(attrs map[string]*string, name string, a fyne.TextAlign) {
	var value string
	switch a {
	case fyne.TextAlignLeading:
		// default mode, don’t add an attr
	case fyne.TextAlignCenter:
		value = "center"
	case fyne.TextAlignTrailing:
		value = "trailing"
	default:
		value = fmt.Sprintf("unknown alignment: %d", a)
	}
	r.setStringAttr(attrs, name, value)
}

func (r *markupRenderer) setBoolAttr(attrs map[string]*string, name string, b bool) {
	if !b {
		return
	}
	attrs[name] = nil
}

func (r *markupRenderer) setColorAttr(attrs map[string]*string, name string, c color.Color) {
	r.setColorAttrWithDefault(attrs, name, c, color.Transparent)
}

func (r *markupRenderer) setColorAttrWithDefault(attrs map[string]*string, name string, c color.Color, d color.Color) {
	if c == nil || c == d {
		return
	}

	if value := knownColor(c); value != "" {
		r.setStringAttr(attrs, name, value)
		return
	}

	rd, g, b, a := col.ToNRGBA(c)
	r.setStringAttr(attrs, name, fmt.Sprintf("rgba(%d,%d,%d,%d)", uint8(rd), uint8(g), uint8(b), uint8(a)))
}

func (r *markupRenderer) setFillModeAttr(attrs map[string]*string, name string, m fynecanvas.ImageFill) {
	var fillMode string
	switch m {
	case fynecanvas.ImageFillStretch:
		// default mode, don’t add an attr
	case fynecanvas.ImageFillContain:
		fillMode = "contain"
	case fynecanvas.ImageFillOriginal:
		fillMode = "original"
	default:
		fillMode = fmt.Sprintf("unknown fill mode: %d", m)
	}
	r.setStringAttr(attrs, name, fillMode)
}

func (r *markupRenderer) setFloatAttr(attrs map[string]*string, name string, f float64) {
	r.setFloatAttrWithDefault(attrs, name, f, 0)
}

func (r *markupRenderer) setFloatAttrWithDefault(attrs map[string]*string, name string, f float64, d float64) {
	if f == d {
		return
	}
	value := fmt.Sprintf("%g", f)
	attrs[name] = &value
}

func (r *markupRenderer) setFloatPosAttr(attrs map[string]*string, name string, x, y float64) {
	if x == 0 && y == 0 {
		return
	}
	value := fmt.Sprintf("%g,%g", x, y)
	attrs[name] = &value
}

func (r *markupRenderer) setSizeAttrWithDefault(attrs map[string]*string, name string, i float32, d float32) {
	if int(i) == int(d) {
		return
	}
	value := fmt.Sprintf("%d", int(i))
	attrs[name] = &value
}

func (r *markupRenderer) setPosAttr(attrs map[string]*string, name string, pos fyne.Position) {
	if int(pos.X) == 0 && int(pos.Y) == 0 {
		return
	}
	value := fmt.Sprintf("%d,%d", int(pos.X), int(pos.Y))
	attrs[name] = &value
}

func (r *markupRenderer) setResourceAttr(attrs map[string]*string, name string, rsc fyne.Resource) {
	if rsc == nil {
		return
	}

	named := false
	if value := knownResource(rsc); value != "" {
		r.setStringAttr(attrs, name, value)
		named = true
	}

	var variant string
	switch t := rsc.(type) {
	case *theme.DisabledResource:
		variant = "disabled"
	case *theme.ErrorThemedResource:
		variant = "error"
	case *theme.InvertedThemedResource:
		variant = "inverted"
	case *theme.PrimaryThemedResource:
		variant = "primary"
	case *theme.ThemedResource:
		variant = string(t.ColorName)
		if variant == "" {
			variant = "foreground"
		}
	default:
		r.setStringAttr(attrs, name, rsc.Name())
		return
	}

	if !named {
		// That’s some magic to access the private `source` field of the themed resource.
		v := reflect.ValueOf(rsc).Elem().Field(0)
		src := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface().(fyne.Resource)
		r.setResourceAttr(attrs, name, src)
	}
	r.setStringAttr(attrs, "themed", variant)
}

func (r *markupRenderer) setScaleModeAttr(attrs map[string]*string, name string, m fynecanvas.ImageScale) {
	var scaleMode string
	switch m {
	case fynecanvas.ImageScaleSmooth:
		// default mode, don’t add an attr
	case fynecanvas.ImageScalePixels:
		scaleMode = "pixels"
	default:
		scaleMode = fmt.Sprintf("unknown scale mode: %d", m)
	}
	r.setStringAttr(attrs, name, scaleMode)
}

func (r *markupRenderer) setSizeAttr(attrs map[string]*string, name string, size fyne.Size) {
	value := fmt.Sprintf("%dx%d", int(size.Width), int(size.Height))
	attrs[name] = &value
}

func (r *markupRenderer) setStringAttr(attrs map[string]*string, name string, s string) {
	if s == "" {
		return
	}
	attrs[name] = &s
}

func (r *markupRenderer) writeCanvas(c fyne.Canvas) {
	attrs := map[string]*string{}
	r.setSizeAttr(attrs, "size", c.Size())
	if tc, ok := c.(WindowlessCanvas); ok {
		r.setBoolAttr(attrs, "padded", tc.Padded())
	}
	r.writeTag("canvas", false, attrs)
	r.w.WriteRune('\n')
	r.indentation++
	r.writeTag("content", false, nil)
	r.w.WriteRune('\n')
	r.indentation++
	intdriver.WalkVisibleObjectTree(c.Content(), r.writeCanvasObject, r.writeCloseCanvasObject)
	r.indentation--
	r.writeIndent()
	r.writeCloseTag("content")
	for _, o := range c.Overlays().List() {
		r.writeTag("overlay", false, nil)
		r.w.WriteRune('\n')
		r.indentation++
		intdriver.WalkVisibleObjectTree(o, r.writeCanvasObject, r.writeCloseCanvasObject)
		r.indentation--
		r.writeIndent()
		r.writeCloseTag("overlay")
	}
	r.indentation--
	r.writeIndent()
	r.writeCloseTag("canvas")
}

func (r *markupRenderer) writeCanvasObject(obj fyne.CanvasObject, _, _ fyne.Position, _ fyne.Size) bool {
	attrs := map[string]*string{}
	r.setPosAttr(attrs, "pos", obj.Position())
	r.setSizeAttr(attrs, "size", obj.Size())
	switch o := obj.(type) {
	case *fynecanvas.Circle:
		r.writeCircle(o, attrs)
	case *fynecanvas.Image:
		r.writeImage(o, attrs)
	case *fynecanvas.Line:
		r.writeLine(o, attrs)
	case *fynecanvas.LinearGradient:
		r.writeLinearGradient(o, attrs)
	case *fynecanvas.RadialGradient:
		r.writeRadialGradient(o, attrs)
	case *fynecanvas.Raster:
		r.writeRaster(o, attrs)
	case *fynecanvas.Rectangle:
		r.writeRectangle(o, attrs)
	case *fynecanvas.Text:
		r.writeText(o, attrs)
	case *fyne.Container:
		r.writeContainer(o, attrs)
	case fyne.Widget:
		r.writeWidget(o, attrs)
	case *layout.Spacer:
		r.writeSpacer(o, attrs)
	default:
		panic(fmt.Sprint("please add support for", reflect.TypeOf(o)))
	}

	return false
}

func (r *markupRenderer) writeCircle(c *fynecanvas.Circle, attrs map[string]*string) {
	r.setColorAttr(attrs, "fillColor", c.FillColor)
	r.setColorAttr(attrs, "strokeColor", c.StrokeColor)
	r.setFloatAttr(attrs, "strokeWidth", float64(c.StrokeWidth))
	r.writeTag("circle", true, attrs)
}

func (r *markupRenderer) writeCloseCanvasObject(o fyne.CanvasObject, _ fyne.Position, _ fyne.CanvasObject) {
	switch o.(type) {
	case *fyne.Container:
		r.indentation--
		r.writeIndent()
		r.writeCloseTag("container")
	case fyne.Widget:
		r.indentation--
		r.writeIndent()
		r.writeCloseTag("widget")
	}
}

func (r *markupRenderer) writeCloseTag(name string) {
	r.w.WriteString("</")
	r.w.WriteString(name)
	r.w.WriteString(">\n")
}

func (r *markupRenderer) writeContainer(_ *fyne.Container, attrs map[string]*string) {
	r.writeTag("container", false, attrs)
	r.w.WriteRune('\n')
	r.indentation++
}

func (r *markupRenderer) writeIndent() {
	for i := 0; i < r.indentation; i++ {
		r.w.WriteRune('\t')
	}
}

func (r *markupRenderer) writeImage(i *fynecanvas.Image, attrs map[string]*string) {
	r.setStringAttr(attrs, "file", i.File)
	r.setResourceAttr(attrs, "rsc", i.Resource)
	if i.File == "" && i.Resource == nil {
		r.setBoolAttr(attrs, "img", i.Image != nil)
	}
	r.setFloatAttr(attrs, "translucency", i.Translucency)
	r.setFillModeAttr(attrs, "fillMode", i.FillMode)
	r.setScaleModeAttr(attrs, "scaleMode", i.ScaleMode)
	if i.Size().Width == theme.IconInlineSize() && i.Size().Height == i.Size().Width {
		r.setStringAttr(attrs, "size", "iconInlineSize")
	}
	r.writeTag("image", true, attrs)
}

func (r *markupRenderer) writeLine(l *fynecanvas.Line, attrs map[string]*string) {
	r.setColorAttr(attrs, "strokeColor", l.StrokeColor)
	r.setFloatAttrWithDefault(attrs, "strokeWidth", float64(l.StrokeWidth), 1)
	r.writeTag("line", true, attrs)
}

func (r *markupRenderer) writeLinearGradient(g *fynecanvas.LinearGradient, attrs map[string]*string) {
	r.setColorAttr(attrs, "startColor", g.StartColor)
	r.setColorAttr(attrs, "endColor", g.EndColor)
	r.setFloatAttr(attrs, "angle", g.Angle)
	r.writeTag("linearGradient", true, attrs)
}

func (r *markupRenderer) writeRadialGradient(g *fynecanvas.RadialGradient, attrs map[string]*string) {
	r.setColorAttr(attrs, "startColor", g.StartColor)
	r.setColorAttr(attrs, "endColor", g.EndColor)
	r.setFloatPosAttr(attrs, "centerOffset", g.CenterOffsetX, g.CenterOffsetY)
	r.writeTag("radialGradient", true, attrs)
}

func (r *markupRenderer) writeRaster(rst *fynecanvas.Raster, attrs map[string]*string) {
	r.setFloatAttr(attrs, "translucency", rst.Translucency)
	r.writeTag("raster", true, attrs)
}

func (r *markupRenderer) writeRectangle(rct *fynecanvas.Rectangle, attrs map[string]*string) {
	r.setColorAttr(attrs, "fillColor", rct.FillColor)
	r.setColorAttr(attrs, "strokeColor", rct.StrokeColor)
	r.setFloatAttr(attrs, "strokeWidth", float64(rct.StrokeWidth))
	r.setFloatAttr(attrs, "radius", float64(rct.CornerRadius))
	r.writeTag("rectangle", true, attrs)
}

func (r *markupRenderer) writeSpacer(_ *layout.Spacer, attrs map[string]*string) {
	r.writeTag("spacer", true, attrs)
}

func (r *markupRenderer) writeTag(name string, isEmpty bool, attrs map[string]*string) {
	r.writeIndent()
	r.w.WriteRune('<')
	r.w.WriteString(name)
	for _, key := range sortedKeys(attrs) {
		r.w.WriteRune(' ')
		r.w.WriteString(key)
		if attrs[key] != nil {
			r.w.WriteString("=\"")
			r.w.WriteString(*attrs[key])
			r.w.WriteRune('"')
		}

	}
	if isEmpty {
		r.w.WriteString("/>\n")
	} else {
		r.w.WriteRune('>')
	}
}

func (r *markupRenderer) writeText(t *fynecanvas.Text, attrs map[string]*string) {
	r.setColorAttrWithDefault(attrs, "color", t.Color, theme.Color(theme.ColorNameForeground))
	r.setAlignmentAttr(attrs, "alignment", t.Alignment)
	r.setSizeAttrWithDefault(attrs, "textSize", t.TextSize, theme.TextSize())
	r.setBoolAttr(attrs, "bold", t.TextStyle.Bold)
	r.setBoolAttr(attrs, "italic", t.TextStyle.Italic)
	r.setBoolAttr(attrs, "monospace", t.TextStyle.Monospace)
	r.writeTag("text", false, attrs)
	r.w.WriteString(t.Text)
	r.writeCloseTag("text")
}

func (r *markupRenderer) writeWidget(w fyne.Widget, attrs map[string]*string) {
	r.setStringAttr(attrs, "type", reflect.TypeOf(w).String())
	r.writeTag("widget", false, attrs)
	r.w.WriteRune('\n')
	r.indentation++
}

func nrgbaColor(c color.Color) color.NRGBA {
	// using ColorToNRGBA to avoid problems with colors with 16-bit components or alpha values that aren't 0 or the maximum possible alpha value
	r, g, b, a := col.ToNRGBA(c)
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

func knownColor(c color.Color) string {
	return map[color.Color]string{
		nrgbaColor(theme.Color(theme.ColorNameBackground)):          "background",
		nrgbaColor(theme.Color(theme.ColorNameButton)):              "button",
		nrgbaColor(theme.Color(theme.ColorNameDisabledButton)):      "disabled button",
		nrgbaColor(theme.Color(theme.ColorNameDisabled)):            "disabled",
		nrgbaColor(theme.Color(theme.ColorNameError)):               "error",
		nrgbaColor(theme.Color(theme.ColorNameFocus)):               "focus",
		nrgbaColor(theme.Color(theme.ColorNameForeground)):          "foreground",
		nrgbaColor(theme.Color(theme.ColorNameForegroundOnError)):   "foregroundOnError",
		nrgbaColor(theme.Color(theme.ColorNameForegroundOnPrimary)): "foregroundOnPrimary",
		nrgbaColor(theme.Color(theme.ColorNameForegroundOnSuccess)): "foregroundOnSuccess",
		nrgbaColor(theme.Color(theme.ColorNameForegroundOnWarning)): "foregroundOnWarning",
		nrgbaColor(theme.Color(theme.ColorNameHeaderBackground)):    "headerBackground",
		nrgbaColor(theme.Color(theme.ColorNameHover)):               "hover",
		nrgbaColor(theme.Color(theme.ColorNameHyperlink)):           "hyperlink",
		nrgbaColor(theme.Color(theme.ColorNameInputBackground)):     "inputBackground",
		nrgbaColor(theme.Color(theme.ColorNameInputBorder)):         "inputBorder",
		nrgbaColor(theme.Color(theme.ColorNameMenuBackground)):      "menuBackground",
		nrgbaColor(theme.Color(theme.ColorNameOverlayBackground)):   "overlayBackground",
		nrgbaColor(theme.Color(theme.ColorNamePlaceHolder)):         "placeholder",
		nrgbaColor(theme.Color(theme.ColorNamePressed)):             "pressed",
		nrgbaColor(theme.Color(theme.ColorNamePrimary)):             "primary",
		nrgbaColor(theme.Color(theme.ColorNameScrollBar)):           "scrollbar",
		nrgbaColor(theme.Color(theme.ColorNameSelection)):           "selection",
		nrgbaColor(theme.Color(theme.ColorNameSeparator)):           "separator",
		nrgbaColor(theme.Color(theme.ColorNameSuccess)):             "success",
		nrgbaColor(theme.Color(theme.ColorNameShadow)):              "shadow",
		nrgbaColor(theme.Color(theme.ColorNameWarning)):             "warning",
	}[nrgbaColor(c)]
}

func knownResource(rsc fyne.Resource) string {
	return map[fyne.Resource]string{
		theme.CancelIcon():             "cancelIcon",
		theme.CheckButtonCheckedIcon(): "checkButtonCheckedIcon",
		theme.CheckButtonFillIcon():    "checkButtonFillIcon",
		theme.CheckButtonIcon():        "checkButtonIcon",
		theme.ColorAchromaticIcon():    "colorAchromaticIcon",
		theme.ColorChromaticIcon():     "colorChromaticIcon",
		theme.ColorPaletteIcon():       "colorPaletteIcon",
		theme.ComputerIcon():           "computerIcon",
		theme.ConfirmIcon():            "confirmIcon",
		theme.ContentAddIcon():         "contentAddIcon",
		theme.ContentClearIcon():       "contentClearIcon",
		theme.ContentCopyIcon():        "contentCopyIcon",
		theme.ContentCutIcon():         "contentCutIcon",
		theme.ContentPasteIcon():       "contentPasteIcon",
		theme.ContentRedoIcon():        "contentRedoIcon",
		theme.ContentRemoveIcon():      "contentRemoveIcon",
		theme.ContentUndoIcon():        "contentUndoIcon",
		theme.DeleteIcon():             "deleteIcon",
		theme.DesktopIcon():            "desktopIcon",
		theme.DocumentCreateIcon():     "documentCreateIcon",
		theme.DocumentIcon():           "documentIcon",
		theme.DocumentPrintIcon():      "documentPrintIcon",
		theme.DocumentSaveIcon():       "documentSaveIcon",
		theme.DownloadIcon():           "downloadIcon",
		theme.ErrorIcon():              "errorIcon",
		theme.FileApplicationIcon():    "fileApplicationIcon",
		theme.FileAudioIcon():          "fileAudioIcon",
		theme.FileIcon():               "fileIcon",
		theme.FileImageIcon():          "fileImageIcon",
		theme.FileTextIcon():           "fileTextIcon",
		theme.FileVideoIcon():          "fileVideoIcon",
		theme.FolderIcon():             "folderIcon",
		theme.FolderNewIcon():          "folderNewIcon",
		theme.FolderOpenIcon():         "folderOpenIcon",
		theme.FyneLogo():               "fyneLogo", //lint:ignore SA1019 This needs to stay until the API is removed.
		theme.HelpIcon():               "helpIcon",
		theme.HistoryIcon():            "historyIcon",
		theme.HomeIcon():               "homeIcon",
		theme.InfoIcon():               "infoIcon",
		theme.MailAttachmentIcon():     "mailAttachementIcon",
		theme.MailComposeIcon():        "mailComposeIcon",
		theme.MailForwardIcon():        "mailForwardIcon",
		theme.MailReplyAllIcon():       "mailReplyAllIcon",
		theme.MailReplyIcon():          "mailReplyIcon",
		theme.MailSendIcon():           "mailSendIcon",
		theme.MediaFastForwardIcon():   "mediaFastForwardIcon",
		theme.MediaFastRewindIcon():    "mediaFastRewindIcon",
		theme.MediaPauseIcon():         "mediaPauseIcon",
		theme.MediaPlayIcon():          "mediaPlayIcon",
		theme.MediaRecordIcon():        "mediaRecordIcon",
		theme.MediaReplayIcon():        "mediaReplayIcon",
		theme.MediaSkipNextIcon():      "mediaSkipNextIcon",
		theme.MediaSkipPreviousIcon():  "mediaSkipPreviousIcon",
		theme.MenuDropDownIcon():       "menuDropDownIcon",
		theme.MenuDropUpIcon():         "menuDropUpIcon",
		theme.MenuExpandIcon():         "menuExpandIcon",
		theme.MenuIcon():               "menuIcon",
		theme.MoveDownIcon():           "moveDownIcon",
		theme.MoveUpIcon():             "moveUpIcon",
		theme.NavigateBackIcon():       "navigateBackIcon",
		theme.NavigateNextIcon():       "navigateNextIcon",
		theme.QuestionIcon():           "questionIcon",
		theme.RadioButtonCheckedIcon(): "radioButtonCheckedIcon",
		theme.RadioButtonFillIcon():    "radioButtonFillIcon",
		theme.RadioButtonIcon():        "radioButtonIcon",
		theme.SearchIcon():             "searchIcon",
		theme.SearchReplaceIcon():      "searchReplaceIcon",
		theme.SettingsIcon():           "settingsIcon",
		theme.StorageIcon():            "storageIcon",
		theme.ViewFullScreenIcon():     "viewFullScreenIcon",
		theme.ViewRefreshIcon():        "viewRefreshIcon",
		theme.ViewRestoreIcon():        "viewRestoreIcon",
		theme.VisibilityIcon():         "visibilityIcon",
		theme.VisibilityOffIcon():      "visibilityOffIcon",
		theme.VolumeDownIcon():         "volumeDownIcon",
		theme.VolumeMuteIcon():         "volumeMuteIcon",
		theme.VolumeUpIcon():           "volumeUpIcon",
		theme.WarningIcon():            "warningIcon",
		theme.ZoomFitIcon():            "zoomFitIcon",
		theme.ZoomInIcon():             "zoomInIcon",
		theme.ZoomOutIcon():            "zoomOutIcon",
	}[rsc]
}

func sortedKeys(m map[string]*string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
