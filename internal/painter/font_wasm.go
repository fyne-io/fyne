//go:build wasm

package painter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
	"syscall/js"

	"golang.org/x/image/math/fixed"
)

// In the browser there are no system fonts to fall back on: text that no
// font Fyne has can draw (Hebrew, Chinese, …) is measured and drawn by
// the browser, with its own fonts.
func init() {
	missingAdvance = browserAdvance
	drawMissing = browserDraw
}

const browserFont = "sans-serif"

var (
	browserMu      sync.Mutex
	browserCanvas  js.Value
	browserCtx     js.Value
	browserWidths  = map[browserKey]fixed.Int26_6{}
	browserImages  = map[browserKey]*image.NRGBA{}
	browserTooMany = 4096 // cached glyphs before the caches start again
)

type browserKey struct {
	text string
	size float32
	c    color.NRGBA
}

func context2D() js.Value {
	if !browserCtx.Truthy() {
		browserCanvas = js.Global().Get("document").Call("createElement", "canvas")
		browserCtx = browserCanvas.Call("getContext", "2d", map[string]any{"willReadFrequently": true})
	}
	return browserCtx
}

// browserAdvance is how wide the browser draws text at size.
func browserAdvance(text string, size fixed.Int26_6) fixed.Int26_6 {
	browserMu.Lock()
	defer browserMu.Unlock()
	key := browserKey{text: text, size: fixed266ToFloat32(size)}
	if w, ok := browserWidths[key]; ok {
		return w
	}
	ctx := context2D()
	ctx.Set("font", fmt.Sprintf("%.2fpx %s", key.size, browserFont))
	w := float32ToFixed266(float32(ctx.Call("measureText", text).Get("width").Float()))
	if len(browserWidths) > browserTooMany {
		clear(browserWidths)
	}
	browserWidths[key] = w
	return w
}

// browserDraw draws text as the browser does, its baseline at baseline.
func browserDraw(dst draw.Image, text string, c color.Color, size, scale float32, x, baseline int) {
	browserMu.Lock()
	defer browserMu.Unlock()
	px := size * scale
	key := browserKey{text: text, size: px, c: color.NRGBAModel.Convert(c).(color.NRGBA)}
	img, ok := browserImages[key]
	ascent := int(math.Ceil(float64(px) * 1.2))
	if !ok {
		ctx := context2D()
		font := fmt.Sprintf("%.2fpx %s", px, browserFont)
		ctx.Set("font", font)
		w := int(math.Ceil(ctx.Call("measureText", text).Get("width").Float())) + 2
		h := int(math.Ceil(float64(px) * 1.6))
		browserCanvas.Set("width", w) // resets the context
		browserCanvas.Set("height", h)
		ctx.Set("font", font)
		ctx.Set("textBaseline", "alphabetic")
		k := key.c
		ctx.Set("fillStyle", fmt.Sprintf("rgba(%d,%d,%d,%.3f)", k.R, k.G, k.B, float64(k.A)/255))
		ctx.Call("fillText", text, 0, ascent)
		data := ctx.Call("getImageData", 0, 0, w, h).Get("data")
		img = image.NewNRGBA(image.Rect(0, 0, w, h))
		js.CopyBytesToGo(img.Pix, js.Global().Get("Uint8Array").New(data.Get("buffer")))
		if len(browserImages) > browserTooMany {
			clear(browserImages)
		}
		browserImages[key] = img
	}
	at := image.Pt(x, baseline-ascent)
	draw.Draw(dst, img.Bounds().Add(at), img, image.Point{}, draw.Over)
}
