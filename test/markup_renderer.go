package test

import (
	"fmt"
	"image/color"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	col "fyne.io/fyne/v2/internal/color"
	"fyne.io/fyne/v2/internal/driver"
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

	for _, n := range theme.PrimaryColorNames() {
		if c == theme.PrimaryColorNamed(n) {
			r.setStringAttr(attrs, name, n)
			return
		}
	}

	rd, g, b, a := col.ToNRGBA(c)
	r.setStringAttr(attrs, name, fmt.Sprintf("rgba(%d,%d,%d,%d)", uint8(rd), uint8(g), uint8(b), uint8(a)))
}

func (r *markupRenderer) setFillModeAttr(attrs map[string]*string, name string, m canvas.ImageFill) {
	var fillMode string
	switch m {
	case canvas.ImageFillStretch:
		// default mode, don’t add an attr
	case canvas.ImageFillContain:
		fillMode = "contain"
	case canvas.ImageFillOriginal:
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

	if value := knownResource(rsc); value != "" {
		r.setStringAttr(attrs, name, value)
		return
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
			variant = "default"
		}
	default:
		r.setStringAttr(attrs, name, rsc.Name())
		return
	}

	// That’s some magic to access the private `source` field of the themed resource.
	v := reflect.ValueOf(rsc).Elem().Field(0)
	src := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface().(fyne.Resource)
	r.setResourceAttr(attrs, name, src)
	r.setStringAttr(attrs, "themed", variant)
}

func (r *markupRenderer) setScaleModeAttr(attrs map[string]*string, name string, m canvas.ImageScale) {
	var scaleMode string
	switch m {
	case canvas.ImageScaleSmooth:
		// default mode, don’t add an attr
	case canvas.ImageScalePixels:
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
	driver.WalkVisibleObjectTree(c.Content(), r.writeCanvasObject, r.writeCloseCanvasObject)
	r.indentation--
	r.writeIndent()
	r.writeCloseTag("content")
	for _, o := range c.Overlays().List() {
		r.writeTag("overlay", false, nil)
		r.w.WriteRune('\n')
		r.indentation++
		driver.WalkVisibleObjectTree(o, r.writeCanvasObject, r.writeCloseCanvasObject)
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
	case *canvas.Circle:
		r.writeCircle(o, attrs)
	case *canvas.Image:
		r.writeImage(o, attrs)
	case *canvas.Line:
		r.writeLine(o, attrs)
	case *canvas.LinearGradient:
		r.writeLinearGradient(o, attrs)
	case *canvas.RadialGradient:
		r.writeRadialGradient(o, attrs)
	case *canvas.Raster:
		r.writeRaster(o, attrs)
	case *canvas.Rectangle:
		r.writeRectangle(o, attrs)
	case *canvas.Text:
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

func (r *markupRenderer) writeCircle(c *canvas.Circle, attrs map[string]*string) {
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

func (r *markupRenderer) writeContainer(c *fyne.Container, attrs map[string]*string) {
	r.writeTag("container", false, attrs)
	r.w.WriteRune('\n')
	r.indentation++
}

func (r *markupRenderer) writeIndent() {
	for i := 0; i < r.indentation; i++ {
		r.w.WriteRune('\t')
	}
}

func (r *markupRenderer) writeImage(i *canvas.Image, attrs map[string]*string) {
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

func (r *markupRenderer) writeLine(l *canvas.Line, attrs map[string]*string) {
	r.setColorAttr(attrs, "strokeColor", l.StrokeColor)
	r.setFloatAttrWithDefault(attrs, "strokeWidth", float64(l.StrokeWidth), 1)
	r.writeTag("line", true, attrs)
}

func (r *markupRenderer) writeLinearGradient(g *canvas.LinearGradient, attrs map[string]*string) {
	r.setColorAttr(attrs, "startColor", g.StartColor)
	r.setColorAttr(attrs, "endColor", g.EndColor)
	r.setFloatAttr(attrs, "angle", g.Angle)
	r.writeTag("linearGradient", true, attrs)
}

func (r *markupRenderer) writeRadialGradient(g *canvas.RadialGradient, attrs map[string]*string) {
	r.setColorAttr(attrs, "startColor", g.StartColor)
	r.setColorAttr(attrs, "endColor", g.EndColor)
	r.setFloatPosAttr(attrs, "centerOffset", g.CenterOffsetX, g.CenterOffsetY)
	r.writeTag("radialGradient", true, attrs)
}

func (r *markupRenderer) writeRaster(rst *canvas.Raster, attrs map[string]*string) {
	r.setFloatAttr(attrs, "translucency", rst.Translucency)
	r.writeTag("raster", true, attrs)
}

func (r *markupRenderer) writeRectangle(rct *canvas.Rectangle, attrs map[string]*string) {
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

func (r *markupRenderer) writeText(t *canvas.Text, attrs map[string]*string) {
	r.setColorAttrWithDefault(attrs, "color", t.Color, theme.ForegroundColor())
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
		nrgbaColor(theme.BackgroundColor()):        "background",
		nrgbaColor(theme.ButtonColor()):            "button",
		nrgbaColor(theme.DisabledButtonColor()):    "disabled button",
		nrgbaColor(theme.DisabledColor()):          "disabled",
		nrgbaColor(theme.ErrorColor()):             "error",
		nrgbaColor(theme.FocusColor()):             "focus",
		nrgbaColor(theme.ForegroundColor()):        "foreground",
		nrgbaColor(theme.HoverColor()):             "hover",
		nrgbaColor(theme.InputBackgroundColor()):   "inputBackground",
		nrgbaColor(theme.InputBorderColor()):       "inputBorder",
		nrgbaColor(theme.MenuBackgroundColor()):    "menuBackground",
		nrgbaColor(theme.OverlayBackgroundColor()): "overlayBackground",
		nrgbaColor(theme.PlaceHolderColor()):       "placeholder",
		nrgbaColor(theme.PrimaryColor()):           "primary",
		nrgbaColor(theme.ScrollBarColor()):         "scrollbar",
		nrgbaColor(theme.SelectionColor()):         "selection",
		nrgbaColor(theme.ShadowColor()):            "shadow",
	}[nrgbaColor(c)]
}

func knownResource(rsc fyne.Resource) string {
	return map[fyne.Resource]string{
		theme.CancelIcon():             "cancelIcon",
		theme.CheckButtonCheckedIcon(): "checkButtonCheckedIcon",
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
		theme.FyneLogo():               "fyneLogo",
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
