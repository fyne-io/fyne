// +build !ci,efl

package efl

// #cgo pkg-config: ecore evas
// #include <Ecore.h>
// #include <Evas.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

func updateFont(obj *C.Evas_Object, c *eflCanvas, size int, style fyne.TextStyle) {
	font := theme.TextFont()

	if style.Monospace {
		font = theme.TextMonospaceFont()
	} else {
		if style.Bold {
			if style.Italic {
				font = theme.TextBoldItalicFont()
			} else {
				font = theme.TextBoldFont()
			}
		} else if style.Italic {
			font = theme.TextItalicFont()
		}
	}

	cstr := C.CString(font.CachePath())
	C.evas_object_text_font_set(obj, cstr, C.Evas_Font_Size(scaleInt(c, size)))
	C.free(unsafe.Pointer(cstr))
}

func getTextPosition(t *canvas.Text, pos fyne.Position, size fyne.Size) fyne.Position {
	min := t.MinSize()

	switch t.Alignment {
	case fyne.TextAlignCenter:
		return fyne.NewPos(pos.X+(size.Width-min.Width)/2, pos.Y+(size.Height-min.Height)/2)
	case fyne.TextAlignTrailing:
		return fyne.NewPos(pos.X+size.Width-min.Width, pos.Y+(size.Height-min.Height)/2)
	default:
		return fyne.NewPos(pos.X, pos.Y+(size.Height-min.Height)/2)
	}
}

func nativeTextBounds(obj *C.Evas_Object) fyne.Size {
	var x, w, h C.Evas_Coord
	length := len(C.GoString(C.evas_object_text_text_get(obj)))
	height := 0

	for i := 0; i < length; i++ {
		C.evas_object_text_char_pos_get(obj, C.int(i), &x, nil, &w, &h)
		if int(h) > height {
			height = int(h)
		}
	}

	return fyne.NewSize(int(x+w), height)
}

func (d *eFLDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
	if len(d.AllWindows()) == 0 {
		return fyne.NewSize(0, 0)
	}
	c := d.AllWindows()[0].Canvas().(*eflCanvas)
	var native fyne.Size

	runOnMain(func() {
		textObj := C.evas_object_text_add(c.evas)

		cstr := C.CString(text)
		C.evas_object_text_text_set(textObj, cstr)
		updateFont(textObj, c, size, style)
		native = nativeTextBounds(textObj)
		C.free(unsafe.Pointer(cstr))

		C.evas_object_del(textObj)
	})

	return fyne.NewSize(unscaleInt(c, native.Width), unscaleInt(c, native.Height))
}
