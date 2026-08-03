//go:build windows

// Turning a fyne.Resource into the HICON that WM_SETICON wants.

package directx

import (
	"image"
	imgdraw "image/draw"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
)

// rasterise draws a resource into a square image of the given pixel size.
//
// Going through a software canvas rather than image.Decode covers SVG and bitmap
// resources with one path, applies the current theme to themed resources, and
// scales to the exact size Windows asked for instead of leaving the OS to stretch
// whatever the source happened to be.
func rasterise(res fyne.Resource, size int) image.Image {
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain

	c := software.NewTransparentCanvas()
	c.SetPadded(false)
	c.SetContent(img)
	c.Resize(fyne.NewSquareSize(float32(size)))

	return c.Capture()
}

// iconFromResource rasterises a resource at the given square pixel size and wraps
// it in a Win32 icon. The caller owns the handle and must DestroyIcon it.
func iconFromResource(res fyne.Resource, size int) windows.Handle {
	return hIconFor(rasterise(res, size))
}

// menuBitmapFromResource rasterises a resource into the 32 bit per pixel bitmap a
// menu item wants for its icon column. The caller owns the handle and must
// DeleteObject it.
//
// Unlike an icon this has to be a DIB section: menus composite hbmpItem with
// AlphaBlend, which reads the alpha channel only from a device independent
// bitmap. A device dependent bitmap from CreateBitmap draws as a black box.
func menuBitmapFromResource(res fyne.Resource, size int) windows.Handle {
	src := rasterise(res, size)
	b := src.Bounds()
	width, height := b.Dx(), b.Dy()
	if width <= 0 || height <= 0 {
		return 0
	}

	header := bitmapInfoHeader{
		biSize:     uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		biWidth:    int32(width),
		biHeight:   int32(-height), // negative for a top-down row order
		biPlanes:   1,
		biBitCount: 32,
	}
	var bitsAddr uintptr
	h, _, _ := procCreateDIBSection.Call(0, uintptr(unsafe.Pointer(&header)), dibRGBColors,
		uintptr(unsafe.Pointer(&bitsAddr)), 0, 0)
	if h == 0 || bitsAddr == 0 {
		return 0
	}

	// color.Color.RGBA is already alpha premultiplied, which is exactly what
	// AlphaBlend expects, so the channels only need reordering to BGRA.
	bits := unsafe.Slice((*byte)(pointerFromAddr(bitsAddr)), width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, bl, a := src.At(b.Min.X+x, b.Min.Y+y).RGBA()
			i := (y*width + x) * 4
			bits[i], bits[i+1], bits[i+2], bits[i+3] = byte(bl>>8), byte(g>>8), byte(r>>8), byte(a>>8)
		}
	}
	return windows.Handle(h)
}

// menuIconSize reports the icon column size for a window's DPI.
// GetSystemMetricsForDpi is Windows 10 1607+; older builds get the system-wide
// metric, which is right at 96 DPI and merely stretched above it.
func menuIconSize(hwnd windows.Handle) int {
	if procGetSystemMetricsForDpi.Find() == nil {
		if n, _, _ := procGetSystemMetricsForDpi.Call(smCxMenuCheck,
			uintptr(getDpiForWindow(hwnd))); n != 0 {
			return int(n)
		}
	}
	n, _, _ := procGetSystemMetrics.Call(smCxMenuCheck)
	if n == 0 {
		return 16
	}
	return int(n)
}

// hIconFor builds an icon from a rasterised image. The colour bitmap carries the
// alpha channel, so the mask is left blank - Windows only consults it for the
// 1-bit icons that predate alpha.
func hIconFor(src image.Image) windows.Handle {
	return createIconIndirect(src, iconInfo{fIcon: 1})
}

// hCursorFor builds a cursor from a desktop.Cursor image and its hotspot. The
// caller owns the handle and must DestroyCursor it.
func hCursorFor(src image.Image, hotX, hotY int) windows.Handle {
	return createIconIndirect(src, iconInfo{xHotspot: uint32(hotX), yHotspot: uint32(hotY)})
}

// createIconIndirect uploads the image as the colour bitmap of an icon or
// cursor; info carries the fIcon flag and hotspot, the bitmaps are filled here.
//
// The colour bitmap must be a 32bpp DIB section: a device-dependent bitmap from
// CreateBitmap has no guaranteed alpha channel, and when Windows finds no alpha
// it falls back to the AND mask - which must itself be defined, because
// CreateBitmap with no bits leaves undefined content that renders as garbage.
func createIconIndirect(src image.Image, info iconInfo) windows.Handle {
	b := src.Bounds()
	width, height := b.Dx(), b.Dy()
	if width <= 0 || height <= 0 {
		return 0
	}

	header := bitmapInfoHeader{
		biSize:     uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		biWidth:    int32(width),
		biHeight:   int32(-height), // negative for a top-down row order
		biPlanes:   1,
		biBitCount: 32,
	}
	var bitsAddr uintptr
	colour, _, _ := procCreateDIBSection.Call(0, uintptr(unsafe.Pointer(&header)), dibRGBColors,
		uintptr(unsafe.Pointer(&bitsAddr)), 0, 0)
	if colour == 0 || bitsAddr == 0 {
		return 0
	}
	defer procDeleteObject.Call(colour)

	// Icon colour bitmaps carry straight (non-premultiplied) alpha - the same
	// convention .ico files use - so convert via NRGBA rather than RGBA().
	straight := image.NewNRGBA(image.Rect(0, 0, width, height))
	imgdraw.Draw(straight, straight.Rect, src, b.Min, imgdraw.Src)
	bits := unsafe.Slice((*byte)(pointerFromAddr(bitsAddr)), width*height*4)
	for i := 0; i < len(bits); i += 4 {
		bits[i], bits[i+1], bits[i+2], bits[i+3] =
			straight.Pix[i+2], straight.Pix[i+1], straight.Pix[i], straight.Pix[i+3]
	}

	// All-zero mask bits mean fully opaque, letting the alpha channel decide.
	// 1bpp bitmap rows are padded to 16-bit boundaries.
	maskBits := make([]byte, (width+15)/16*2*height)
	mask, _, _ := procCreateBitmap.Call(uintptr(width), uintptr(height), 1, 1,
		uintptr(unsafe.Pointer(&maskBits[0])))
	runtime.KeepAlive(maskBits)
	if mask == 0 {
		return 0
	}
	defer procDeleteObject.Call(mask)

	info.hbmMask = windows.Handle(mask)
	info.hbmColor = windows.Handle(colour)
	h, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&info)))
	return windows.Handle(h)
}
