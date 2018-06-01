package efl

// #cgo pkg-config: ecore-evas
// #include <Evas.h>
import "C"

import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/canvas"
import "github.com/fyne-io/fyne/api/ui/theme"

func updateFont(obj *C.Evas_Object, c *eflCanvas, t *canvas.Text) {
	font := theme.TextFont()

	if t.Bold {
		if t.Italic {
			font = theme.TextBoldItalicFont()
		} else {
			font = theme.TextBoldFont()
		}
	} else if t.Italic {
		font = theme.TextItalicFont()
	}

	C.evas_object_text_font_set(obj, C.CString(font), C.Evas_Font_Size(scaleInt(c, t.FontSize)))
}

func getTextPosition(t *canvas.Text, pos ui.Position, size ui.Size) ui.Position {
	min := t.MinSize()

	switch t.Alignment {
	case ui.TextAlignCenter:
		return ui.NewPos(pos.X+(size.Width-min.Width)/2, pos.Y+(size.Height-min.Height)/2)
	case ui.TextAlignTrailing:
		return ui.NewPos(pos.X+size.Width-min.Width, pos.Y+(size.Height-min.Height)/2)
	default:
		return ui.NewPos(pos.X, pos.Y+(size.Height-min.Height)/2)
	}
}

func nativeTextBounds(obj *C.Evas_Object) ui.Size {
	width, height := 0, 0
	var w, h C.Evas_Coord
	length := int(C.strlen(C.evas_object_text_text_get(obj)))

	for i := 0; i < length; i++ {
		C.evas_object_text_char_pos_get(obj, C.int(i), nil, nil, &w, &h)
		width += int(w) + 2
		if int(h) > height {
			height = int(h)
		}
	}

	return ui.Size{width, height}
}

func (d *eFLDriver) RenderedTextSize(text string, size int) ui.Size {
	c := ui.GetDriver().AllWindows()[0].Canvas().(*eflCanvas)

	textObj := C.evas_object_text_add(c.evas)
	C.evas_object_text_text_set(textObj, C.CString(text))
	C.evas_object_text_font_set(textObj, C.CString(theme.TextFont()), C.Evas_Font_Size(scaleInt(c, size)))

	native := nativeTextBounds(textObj)
	return ui.NewSize(unscaleInt(c, native.Width), unscaleInt(c, native.Height))
}
