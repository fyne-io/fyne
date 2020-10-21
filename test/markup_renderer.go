package test

import (
	"fmt"
	"image/color"
	"math"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
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
	var value string
	switch c {
	case d, nil:
		return
	case theme.BackgroundColor():
		value = "background"
	case theme.ButtonColor():
		value = "button"
	case theme.DisabledButtonColor():
		value = "disabled button"
	case theme.DisabledTextColor():
		value = "disabled text"
	case theme.FocusColor():
		value = "focus"
	case theme.HoverColor():
		value = "hover"
	case theme.PlaceHolderColor():
		value = "placeholder"
	case theme.PrimaryColor():
		value = "primary"
	case theme.ScrollBarColor():
		value = "scrollbar"
	case theme.ShadowColor():
		value = "shadow"
	case theme.TextColor():
		value = "text"
	case theme.DisabledIconColor():
		value = "disabled icon"
	case theme.HyperlinkColor():
		value = "hyperlink"
	case theme.IconColor():
		value = "icon"
	default:
		for _, n := range theme.PrimaryColorNames() {
			if c == theme.PrimaryColorNamed(n) {
				value = n
				break
			}
		}
		if value == "" {
			rd, g, b, a := c.RGBA()
			value = fmt.Sprintf("rgba(%d,%d,%d,%d)", uint8(rd), uint8(g), uint8(b), uint8(a))
		}
	}
	attrs[name] = &value
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

func (r *markupRenderer) setIntAttr(attrs map[string]*string, name string, i int) {
	r.setIntAttrWithDefault(attrs, name, i, 0)
}

func (r *markupRenderer) setIntAttrWithDefault(attrs map[string]*string, name string, i int, d int) {
	if i == d {
		return
	}
	value := fmt.Sprintf("%d", i)
	attrs[name] = &value
}

func (r *markupRenderer) setPosAttr(attrs map[string]*string, name string, pos fyne.Position) {
	if pos.X == 0 && pos.Y == 0 {
		return
	}
	value := fmt.Sprintf("%d,%d", pos.X, pos.Y)
	attrs[name] = &value
}

func (r *markupRenderer) setRectAttr(attrs map[string]*string, name string, pos fyne.Position, size fyne.Size) {
	value := fmt.Sprintf("%dx%d@%d,%d", size.Width, size.Height, pos.X, pos.Y)
	attrs[name] = &value
}

func (r *markupRenderer) setResourceAttr(attrs map[string]*string, name string, rsc fyne.Resource) {
	if rsc == nil {
		return
	}
	var value string
	switch rsc {
	case theme.CancelIcon():
		value = "cancelIcon"
	case theme.CheckButtonCheckedIcon():
		value = "checkButtonCheckedIcon"
	case theme.CheckButtonIcon():
		value = "checkButtonIcon"
	case theme.ColorAchromaticIcon():
		value = "colorAchromaticIcon"
	case theme.ColorChromaticIcon():
		value = "colorChromaticIcon"
	case theme.ColorPaletteIcon():
		value = "colorPaletteIcon"
	case theme.ComputerIcon():
		value = "computerIcon"
	case theme.ConfirmIcon():
		value = "confirmIcon"
	case theme.ContentAddIcon():
		value = "contentAddIcon"
	case theme.ContentClearIcon():
		value = "contentClearIcon"
	case theme.ContentCopyIcon():
		value = "contentCopyIcon"
	case theme.ContentCutIcon():
		value = "contentCutIcon"
	case theme.ContentPasteIcon():
		value = "contentPasteIcon"
	case theme.ContentRedoIcon():
		value = "contentRedoIcon"
	case theme.ContentRemoveIcon():
		value = "contentRemoveIcon"
	case theme.ContentUndoIcon():
		value = "contentUndoIcon"
	case theme.DeleteIcon():
		value = "deleteIcon"
	case theme.DocumentCreateIcon():
		value = "documentCreateIcon"
	case theme.DocumentIcon():
		value = "documentIcon"
	case theme.DocumentPrintIcon():
		value = "documentPrintIcon"
	case theme.DocumentSaveIcon():
		value = "documentSaveIcon"
	case theme.DownloadIcon():
		value = "downloadIcon"
	case theme.ErrorIcon():
		value = "errorIcon"
	case theme.FileApplicationIcon():
		value = "fileApplicationIcon"
	case theme.FileAudioIcon():
		value = "fileAudioIcon"
	case theme.FileIcon():
		value = "fileIcon"
	case theme.FileImageIcon():
		value = "fileImageIcon"
	case theme.FileTextIcon():
		value = "fileTextIcon"
	case theme.FileVideoIcon():
		value = "fileVideoIcon"
	case theme.FolderIcon():
		value = "folderIcon"
	case theme.FolderNewIcon():
		value = "folderNewIcon"
	case theme.FolderOpenIcon():
		value = "folderOpenIcon"
	case theme.FyneLogo():
		value = "fyneLogo"
	case theme.HelpIcon():
		value = "helpIcon"
	case theme.HistoryIcon():
		value = "historyIcon"
	case theme.HomeIcon():
		value = "homeIcon"
	case theme.InfoIcon():
		value = "infoIcon"
	case theme.MailAttachmentIcon():
		value = "mailAttachementIcon"
	case theme.MailComposeIcon():
		value = "mailComposeIcon"
	case theme.MailForwardIcon():
		value = "mailForwardIcon"
	case theme.MailReplyAllIcon():
		value = "mailReplyAllIcon"
	case theme.MailReplyIcon():
		value = "mailReplyIcon"
	case theme.MailSendIcon():
		value = "mailSendIcon"
	case theme.MediaFastForwardIcon():
		value = "mediaFastForwardIcon"
	case theme.MediaFastRewindIcon():
		value = "mediaFastRewindIcon"
	case theme.MediaPauseIcon():
		value = "mediaPauseIcon"
	case theme.MediaPlayIcon():
		value = "mediaPlayIcon"
	case theme.MediaRecordIcon():
		value = "mediaRecordIcon"
	case theme.MediaReplayIcon():
		value = "mediaReplayIcon"
	case theme.MediaSkipNextIcon():
		value = "mediaSkipNextIcon"
	case theme.MediaSkipPreviousIcon():
		value = "mediaSkipPreviousIcon"
	case theme.MenuDropDownIcon():
		value = "menuDropDownIcon"
	case theme.MenuDropUpIcon():
		value = "menuDropUpIcon"
	case theme.MenuExpandIcon():
		value = "menuExpandIcon"
	case theme.MenuIcon():
		value = "menuIcon"
	case theme.MoveDownIcon():
		value = "moveDownIcon"
	case theme.MoveUpIcon():
		value = "moveUpIcon"
	case theme.NavigateBackIcon():
		value = "navigateBackIcon"
	case theme.NavigateNextIcon():
		value = "navigateNextIcon"
	case theme.QuestionIcon():
		value = "questionIcon"
	case theme.RadioButtonCheckedIcon():
		value = "radioButtonCheckedIcon"
	case theme.RadioButtonIcon():
		value = "radioButtonIcon"
	case theme.SearchIcon():
		value = "searchIcon"
	case theme.SearchReplaceIcon():
		value = "searchReplaceIcon"
	case theme.SettingsIcon():
		value = "settingsIcon"
	case theme.StorageIcon():
		value = "storageIcon"
	case theme.ViewFullScreenIcon():
		value = "viewFullScreenIcon"
	case theme.ViewRefreshIcon():
		value = "viewRefreshIcon"
	case theme.ViewRestoreIcon():
		value = "viewRestoreIcon"
	case theme.VisibilityIcon():
		value = "visibilityIcon"
	case theme.VisibilityOffIcon():
		value = "visibilityOffIcon"
	case theme.VolumeDownIcon():
		value = "volumeDownIcon"
	case theme.VolumeMuteIcon():
		value = "volumeMuteIcon"
	case theme.VolumeUpIcon():
		value = "volumeUpIcon"
	case theme.WarningIcon():
		value = "warningIcon"
	case theme.ZoomFitIcon():
		value = "zoomFitIcon"
	case theme.ZoomInIcon():
		value = "zoomInIcon"
	case theme.ZoomOutIcon():
		value = "zoomOutIcon"
	}
	if value != "" {
		r.setStringAttr(attrs, name, value)
		return
	}

	var variant string
	switch rsc.(type) {
	case *theme.DisabledResource:
		variant = "disabled"
	case *theme.ErrorThemedResource:
		variant = "error"
	case *theme.InvertedThemedResource:
		variant = "inverted"
	case *theme.PrimaryThemedResource:
		variant = "primary"
	case *theme.ThemedResource:
		variant = "default"
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
	}
	r.setStringAttr(attrs, name, scaleMode)
}

func (r *markupRenderer) setSizeAttr(attrs map[string]*string, name string, size fyne.Size) {
	value := fmt.Sprintf("%dx%d", size.Width, size.Height)
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

func (r *markupRenderer) writeCanvasObject(obj fyne.CanvasObject, _, clipPos fyne.Position, clipSize fyne.Size) bool {
	attrs := map[string]*string{}
	r.setPosAttr(attrs, "pos", obj.Position())
	r.setSizeAttr(attrs, "size", obj.Size())
	if clipPos != fyne.NewPos(0, 0) || clipSize != fyne.NewSize(math.MaxInt32, math.MaxInt32) {
		r.setRectAttr(attrs, "clip", clipPos, clipSize)
	}
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

func (r *markupRenderer) writeCloseCanvasObject(o, _ fyne.CanvasObject) {
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
	r.setBoolAttr(attrs, "img", i.Image != nil)
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
	r.setColorAttrWithDefault(attrs, "color", t.Color, theme.TextColor())
	r.setAlignmentAttr(attrs, "alignment", t.Alignment)
	r.setIntAttrWithDefault(attrs, "textSize", t.TextSize, theme.TextSize())
	r.setBoolAttr(attrs, "bold", t.TextStyle.Bold)
	r.setBoolAttr(attrs, "italic", t.TextStyle.Italic)
	r.setBoolAttr(attrs, "monospace", t.TextStyle.Monospace)
	r.writeTag("text", false, attrs)
	r.w.WriteString(t.Text)
	r.writeCloseTag("text")
}

func (r *markupRenderer) writeWidget(w fyne.Widget, attrs map[string]*string) {
	r.setStringAttr(attrs, "type", reflect.TypeOf(w).String())
	bgColor := cache.Renderer(w).BackgroundColor()
	if bgColor != color.Transparent { // transparent widgets are the future
		r.setColorAttrWithDefault(attrs, "backgroundColor", bgColor, theme.BackgroundColor())
	}
	r.writeTag("widget", false, attrs)
	r.w.WriteRune('\n')
	r.indentation++
}

func sortedKeys(m map[string]*string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
