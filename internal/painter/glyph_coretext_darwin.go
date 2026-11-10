//go:build darwin

package painter

/*
#cgo LDFLAGS: -framework CoreText -framework CoreGraphics -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <CoreGraphics/CoreGraphics.h>
#include <CoreText/CoreText.h>
#include <stdlib.h>

// fyne_ct_descriptor returns the descriptor of face index of a font file or
// its bytes, retained, or NULL. A collection has one descriptor per face.
static CTFontDescriptorRef fyne_ct_descriptor(const char *path, const void *data, long len, int index) {
	CFArrayRef all = NULL;
	if (path) {
		CFURLRef url = CFURLCreateFromFileSystemRepresentation(NULL, (const UInt8 *)path, strlen(path), false);
		if (!url) {
			return NULL;
		}
		all = CTFontManagerCreateFontDescriptorsFromURL(url);
		CFRelease(url);
	} else {
		CFDataRef bytes = CFDataCreate(NULL, data, len);
		if (!bytes) {
			return NULL;
		}
		all = CTFontManagerCreateFontDescriptorsFromData(bytes);
		CFRelease(bytes);
	}
	if (!all) {
		return NULL;
	}
	CTFontDescriptorRef out = NULL;
	if (index >= 0 && index < CFArrayGetCount(all)) {
		out = (CTFontDescriptorRef)CFRetain(CFArrayGetValueAtIndex(all, index));
	}
	CFRelease(all);
	return out;
}

static CTFontRef fyne_ct_font(CTFontDescriptorRef desc, double size) {
	return CTFontCreateWithFontDescriptor(desc, size, NULL);
}

static void fyne_ct_release(CFTypeRef ref) {
	if (ref) {
		CFRelease(ref);
	}
}

// fyne_ct_draw_glyph draws one glyph white on black into an 8-bit gray
// bitmap of w by h, its pen position x, y pixels from the top left. Font
// smoothing thickens strokes only for a glyph drawn opaque over an opaque
// background, so the gray value is the coverage.
static int fyne_ct_draw_glyph(CTFontRef font, unsigned short glyph, unsigned char *buf, int w, int h, double x, double y) {
	CGColorSpaceRef cs = CGColorSpaceCreateDeviceGray();
	CGContextRef ctx = CGBitmapContextCreate(buf, w, h, 8, w, cs, kCGImageAlphaNone);
	CGColorSpaceRelease(cs);
	if (!ctx) {
		return 0;
	}
	CGContextSetGrayFillColor(ctx, 0, 1);
	CGContextFillRect(ctx, CGRectMake(0, 0, w, h));
	CGContextSetGrayFillColor(ctx, 1, 1);
	CGContextSetAllowsAntialiasing(ctx, true);
	CGContextSetShouldAntialias(ctx, true);
	CGContextSetAllowsFontSmoothing(ctx, true);
	CGContextSetShouldSmoothFonts(ctx, true);
	// The pen position is already where the shaper put it, so CoreText
	// must not snap it to a pixel or a quarter of one.
	CGContextSetAllowsFontSubpixelPositioning(ctx, true);
	CGContextSetShouldSubpixelPositionFonts(ctx, true);
	CGContextSetAllowsFontSubpixelQuantization(ctx, false);
	CGContextSetShouldSubpixelQuantizeFonts(ctx, false);
	CGGlyph g = glyph;
	CGPoint at = {x, h - y};
	CTFontDrawGlyphs(font, &g, &at, 1, ctx);
	CGContextRelease(ctx);
	return 1;
}
*/
import "C"

import (
	"image"
	"math"
	"sync"
	"unsafe"

	"github.com/go-text/typesetting/font"
)

// nativeGlyphs is true where rasteriseGlyphNative draws, so the whole-run
// path draws a glyph at a time to match the atlas.
const nativeGlyphs = true

// glyphMargin is the room the atlas leaves around a glyph's ink. Smoothing
// thickens a stroke by about a pixel, which needs a second one.
const glyphMargin = 2

type coreTextKey struct {
	font *font.Font
	size float32
}

// coreText holds what CoreText draws with. Its lock is held while a glyph is
// drawn, so a font is never released under a draw.
var coreText struct {
	sync.Mutex
	// sources holds the bytes each theme font was parsed from, so CoreText
	// opens exactly the font the shaper picked glyph IDs from.
	sources map[*font.Font][]byte
	// descriptors and fonts hold each font CoreText opened, and nil for one
	// it cannot, so a miss is asked once.
	descriptors map[*font.Font]C.CTFontDescriptorRef
	fonts       map[coreTextKey]C.CTFontRef
}

// coreTextFontsMax bounds the fonts kept open, one per font and size.
const coreTextFontsMax = 256

func rememberFontSource(f *font.Font, data []byte) {
	coreText.Lock()
	defer coreText.Unlock()
	if coreText.sources == nil {
		coreText.sources = make(map[*font.Font][]byte)
	}
	coreText.sources[f] = data
}

// forgetNativeFonts releases everything CoreText opened, for ClearFontCache:
// the fonts it was opened from are no longer used.
func forgetNativeFonts() {
	coreText.Lock()
	defer coreText.Unlock()
	releaseFonts()
	for _, desc := range coreText.descriptors {
		C.fyne_ct_release(C.CFTypeRef(desc))
	}
	coreText.sources, coreText.descriptors = nil, nil
}

func releaseFonts() {
	for _, ct := range coreText.fonts {
		C.fyne_ct_release(C.CFTypeRef(ct))
	}
	coreText.fonts = nil
}

// rasteriseGlyphNative draws a glyph with CoreText, as the rest of macOS
// draws text: with font smoothing, where go-text draws strokes noticeably
// lighter. Shaping, measuring and placement are go-text's either way.
// It reports false when CoreText cannot draw this face, and the caller
// rasterises it in Go.
func rasteriseGlyphNative(dst *image.RGBA, g PlacedGlyph, fontSize, scale float32, x, y int, subX float32) bool {
	// A variable font's instance is not something the descriptor carries.
	if len(g.Face.Coords()) != 0 {
		return false
	}
	b := dst.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return false
	}
	gray := make([]byte, w*h)

	coreText.Lock()
	ct := coreTextFont(g.Face.Font, fontSize*scale)
	drawn := ct != 0 && C.fyne_ct_draw_glyph(ct, C.ushort(g.Glyph.GlyphID), (*C.uchar)(unsafe.Pointer(&gray[0])),
		C.int(w), C.int(h), C.double(float32(x)+subX), C.double(y)) != 0
	coreText.Unlock()
	if !drawn {
		return false
	}

	// Premultiplied white, as go-text draws it: every channel is coverage.
	for row := 0; row < h; row++ {
		o := dst.PixOffset(b.Min.X, b.Min.Y+row)
		for col, v := range gray[row*w : row*w+w] {
			p := dst.Pix[o+4*col : o+4*col+4 : o+4*col+4]
			p[0], p[1], p[2], p[3] = v, v, v, v
		}
	}
	return true
}

// coreTextFont is called with coreText locked.
func coreTextFont(f *font.Font, size float32) C.CTFontRef {
	key := coreTextKey{font: f, size: size}
	if ct, ok := coreText.fonts[key]; ok {
		return ct
	}
	if len(coreText.fonts) >= coreTextFontsMax {
		releaseFonts()
	}
	if coreText.fonts == nil {
		coreText.fonts = make(map[coreTextKey]C.CTFontRef)
	}
	var ct C.CTFontRef
	if desc := coreTextDescriptor(f); desc != 0 {
		ct = C.fyne_ct_font(desc, C.double(math.Max(float64(size), 1)))
	}
	coreText.fonts[key] = ct
	return ct
}

// coreTextDescriptor is called with coreText locked.
func coreTextDescriptor(f *font.Font) C.CTFontDescriptorRef {
	if desc, ok := coreText.descriptors[f]; ok {
		return desc
	}
	if coreText.descriptors == nil {
		coreText.descriptors = make(map[*font.Font]C.CTFontDescriptorRef)
	}
	var desc C.CTFontDescriptorRef
	if raw, ok := coreText.sources[f]; ok {
		if len(raw) > 0 {
			desc = C.fyne_ct_descriptor(nil, unsafe.Pointer(&raw[0]), C.long(len(raw)), 0)
		}
	} else if loc, ok := systemFontLocation(f); ok && loc.Instance == 0 {
		path := C.CString(loc.File)
		desc = C.fyne_ct_descriptor(path, nil, 0, C.int(loc.Index))
		C.free(unsafe.Pointer(path))
	}
	coreText.descriptors[f] = desc
	return desc
}
