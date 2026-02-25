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
	case *fynecanvas.Polygon:
		r.writePolygon(o, attrs)
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
	case *fynecanvas.Arc:
		r.writeArc(o, attrs)
	default:
		panic(fmt.Sprint("please add support for", reflect.TypeOf(o)))
	}

	return false
}

func (r *markupRenderer) writeArc(a *fynecanvas.Arc, attrs map[string]*string) {
	r.setColorAttr(attrs, "fillColor", a.FillColor)
	r.setFloatAttr(attrs, "cutoutRatio", float64(a.CutoutRatio))
	r.setFloatAttr(attrs, "startAngle", float64(a.StartAngle))
	r.setFloatAttr(attrs, "endAngle", float64(a.EndAngle))
	r.setFloatAttr(attrs, "radius", float64(a.CornerRadius))
	r.setColorAttr(attrs, "strokeColor", a.StrokeColor)
	r.setFloatAttr(attrs, "strokeWidth", float64(a.StrokeWidth))
	r.writeTag("arc", true, attrs)
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

func (r *markupRenderer) writePolygon(rct *fynecanvas.Polygon, attrs map[string]*string) {
	r.setColorAttr(attrs, "fillColor", rct.FillColor)
	r.setColorAttr(attrs, "strokeColor", rct.StrokeColor)
	r.setFloatAttr(attrs, "strokeWidth", float64(rct.StrokeWidth))
	r.setFloatAttr(attrs, "radius", float64(rct.CornerRadius))
	r.setFloatAttr(attrs, "angle", float64(rct.Angle))
	r.setFloatAttr(attrs, "sides", float64(rct.Sides))
	r.writeTag("polygon", true, attrs)
}

func (r *markupRenderer) writeRectangle(rct *fynecanvas.Rectangle, attrs map[string]*string) {
	r.setColorAttr(attrs, "fillColor", rct.FillColor)
	r.setColorAttr(attrs, "strokeColor", rct.StrokeColor)
	r.setFloatAttr(attrs, "strokeWidth", float64(rct.StrokeWidth))
	r.setFloatAttr(attrs, "radius", float64(rct.CornerRadius))
	r.setFloatAttr(attrs, "aspect", float64(rct.Aspect))
	r.setFloatAttr(attrs, "topRightRadius", float64(rct.TopRightCornerRadius))
	r.setFloatAttr(attrs, "topLeftRadius", float64(rct.TopLeftCornerRadius))
	r.setFloatAttr(attrs, "bottomRightRadius", float64(rct.BottomRightCornerRadius))
	r.setFloatAttr(attrs, "bottomLeftRadius", float64(rct.BottomLeftCornerRadius))
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

//gocyclo:ignore
func knownColor(c color.Color) string {
	switch nrgbaColor(c) {
	case nrgbaColor(theme.Color(theme.ColorNameBackground)):
		return "background"
	case nrgbaColor(theme.Color(theme.ColorNameButton)):
		return "button"
	case nrgbaColor(theme.Color(theme.ColorNameDisabledButton)):
		return "disabled button"
	case nrgbaColor(theme.Color(theme.ColorNameDisabled)):
		return "disabled"
	case nrgbaColor(theme.Color(theme.ColorNameError)):
		return "error"
	case nrgbaColor(theme.Color(theme.ColorNameFocus)):
		return "focus"
	case nrgbaColor(theme.Color(theme.ColorNameForeground)):
		return "foreground"
	case nrgbaColor(theme.Color(theme.ColorNameForegroundOnError)):
		return "foregroundOnError"
	case nrgbaColor(theme.Color(theme.ColorNameForegroundOnPrimary)):
		return "foregroundOnPrimary"
	case nrgbaColor(theme.Color(theme.ColorNameForegroundOnSuccess)):
		return "foregroundOnSuccess"
	case nrgbaColor(theme.Color(theme.ColorNameForegroundOnWarning)):
		return "foregroundOnWarning"
	case nrgbaColor(theme.Color(theme.ColorNameHeaderBackground)):
		return "headerBackground"
	case nrgbaColor(theme.Color(theme.ColorNameHover)):
		return "hover"
	case nrgbaColor(theme.Color(theme.ColorNameHyperlink)):
		return "hyperlink"
	case nrgbaColor(theme.Color(theme.ColorNameInputBackground)):
		return "inputBackground"
	case nrgbaColor(theme.Color(theme.ColorNameInputBorder)):
		return "inputBorder"
	case nrgbaColor(theme.Color(theme.ColorNameMenuBackground)):
		return "menuBackground"
	case nrgbaColor(theme.Color(theme.ColorNameOverlayBackground)):
		return "overlayBackground"
	case nrgbaColor(theme.Color(theme.ColorNamePlaceHolder)):
		return "placeholder"
	case nrgbaColor(theme.Color(theme.ColorNamePressed)):
		return "pressed"
	case nrgbaColor(theme.Color(theme.ColorNamePrimary)):
		return "primary"
	case nrgbaColor(theme.Color(theme.ColorNameScrollBar)):
		return "scrollbar"
	case nrgbaColor(theme.Color(theme.ColorNameScrollBarBackground)):
		return "scrollbarBackground"
	case nrgbaColor(theme.Color(theme.ColorNameSelection)):
		return "selection"
	case nrgbaColor(theme.Color(theme.ColorNameSeparator)):
		return "separator"
	case nrgbaColor(theme.Color(theme.ColorNameSuccess)):
		return "success"
	case nrgbaColor(theme.Color(theme.ColorNameShadow)):
		return "shadow"
	case nrgbaColor(theme.Color(theme.ColorNameWarning)):
		return "warning"
	default:
		return ""
	}
}

//gocyclo:ignore
func knownResource(rsc fyne.Resource) string {
	switch rsc {
	case theme.CancelIcon():
		return "cancelIcon"
	case theme.CheckButtonCheckedIcon():
		return "checkButtonCheckedIcon"
	case theme.CheckButtonFillIcon():
		return "checkButtonFillIcon"
	case theme.CheckButtonIcon():
		return "checkButtonIcon"
	case theme.ColorAchromaticIcon():
		return "colorAchromaticIcon"
	case theme.ColorChromaticIcon():
		return "colorChromaticIcon"
	case theme.ColorPaletteIcon():
		return "colorPaletteIcon"
	case theme.ComputerIcon():
		return "computerIcon"
	case theme.ConfirmIcon():
		return "confirmIcon"
	case theme.ContentAddIcon():
		return "contentAddIcon"
	case theme.ContentClearIcon():
		return "contentClearIcon"
	case theme.ContentCopyIcon():
		return "contentCopyIcon"
	case theme.ContentCutIcon():
		return "contentCutIcon"
	case theme.ContentPasteIcon():
		return "contentPasteIcon"
	case theme.ContentRedoIcon():
		return "contentRedoIcon"
	case theme.ContentRemoveIcon():
		return "contentRemoveIcon"
	case theme.ContentUndoIcon():
		return "contentUndoIcon"
	case theme.DeleteIcon():
		return "deleteIcon"
	case theme.DesktopIcon():
		return "desktopIcon"
	case theme.DocumentCreateIcon():
		return "documentCreateIcon"
	case theme.DocumentIcon():
		return "documentIcon"
	case theme.DocumentPrintIcon():
		return "documentPrintIcon"
	case theme.DocumentSaveIcon():
		return "documentSaveIcon"
	case theme.DownloadIcon():
		return "downloadIcon"
	case theme.ErrorIcon():
		return "errorIcon"
	case theme.FileApplicationIcon():
		return "fileApplicationIcon"
	case theme.FileAudioIcon():
		return "fileAudioIcon"
	case theme.FileIcon():
		return "fileIcon"
	case theme.FileImageIcon():
		return "fileImageIcon"
	case theme.FileTextIcon():
		return "fileTextIcon"
	case theme.FileVideoIcon():
		return "fileVideoIcon"
	case theme.FolderIcon():
		return "folderIcon"
	case theme.FolderNewIcon():
		return "folderNewIcon"
	case theme.FolderOpenIcon():
		return "folderOpenIcon"
	case theme.FyneLogo():
		return "fyneLogo" //lint:ignore SA1019 This needs to stay until the API is removed.
	case theme.HelpIcon():
		return "helpIcon"
	case theme.HistoryIcon():
		return "historyIcon"
	case theme.HomeIcon():
		return "homeIcon"
	case theme.InfoIcon():
		return "infoIcon"
	case theme.MailAttachmentIcon():
		return "mailAttachementIcon"
	case theme.MailComposeIcon():
		return "mailComposeIcon"
	case theme.MailForwardIcon():
		return "mailForwardIcon"
	case theme.MailReplyAllIcon():
		return "mailReplyAllIcon"
	case theme.MailReplyIcon():
		return "mailReplyIcon"
	case theme.MailSendIcon():
		return "mailSendIcon"
	case theme.MediaFastForwardIcon():
		return "mediaFastForwardIcon"
	case theme.MediaFastRewindIcon():
		return "mediaFastRewindIcon"
	case theme.MediaPauseIcon():
		return "mediaPauseIcon"
	case theme.MediaPlayIcon():
		return "mediaPlayIcon"
	case theme.MediaRecordIcon():
		return "mediaRecordIcon"
	case theme.MediaReplayIcon():
		return "mediaReplayIcon"
	case theme.MediaSkipNextIcon():
		return "mediaSkipNextIcon"
	case theme.MediaSkipPreviousIcon():
		return "mediaSkipPreviousIcon"
	case theme.MenuDropDownIcon():
		return "menuDropDownIcon"
	case theme.MenuDropUpIcon():
		return "menuDropUpIcon"
	case theme.MenuExpandIcon():
		return "menuExpandIcon"
	case theme.MenuIcon():
		return "menuIcon"
	case theme.MoveDownIcon():
		return "moveDownIcon"
	case theme.MoveUpIcon():
		return "moveUpIcon"
	case theme.NavigateBackIcon():
		return "navigateBackIcon"
	case theme.NavigateNextIcon():
		return "navigateNextIcon"
	case theme.QuestionIcon():
		return "questionIcon"
	case theme.RadioButtonCheckedIcon():
		return "radioButtonCheckedIcon"
	case theme.RadioButtonFillIcon():
		return "radioButtonFillIcon"
	case theme.RadioButtonIcon():
		return "radioButtonIcon"
	case theme.SearchIcon():
		return "searchIcon"
	case theme.SearchReplaceIcon():
		return "searchReplaceIcon"
	case theme.SettingsIcon():
		return "settingsIcon"
	case theme.StorageIcon():
		return "storageIcon"
	case theme.ViewFullScreenIcon():
		return "viewFullScreenIcon"
	case theme.ViewRefreshIcon():
		return "viewRefreshIcon"
	case theme.ViewRestoreIcon():
		return "viewRestoreIcon"
	case theme.VisibilityIcon():
		return "visibilityIcon"
	case theme.VisibilityOffIcon():
		return "visibilityOffIcon"
	case theme.VolumeDownIcon():
		return "volumeDownIcon"
	case theme.VolumeMuteIcon():
		return "volumeMuteIcon"
	case theme.VolumeUpIcon():
		return "volumeUpIcon"
	case theme.WarningIcon():
		return "warningIcon"
	case theme.ZoomFitIcon():
		return "zoomFitIcon"
	case theme.ZoomInIcon():
		return "zoomInIcon"
	case theme.ZoomOutIcon():
		return "zoomOutIcon"
	default:
		return ""
	}
}

func sortedKeys(m map[string]*string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
