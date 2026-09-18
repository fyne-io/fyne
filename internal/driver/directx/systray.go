//go:build windows && directx

package directx

// System tray support. The menu handling and icon conversion shared with the
// other desktop drivers lives in internal/driver/systray; what stays here is the
// Win32 specifics - the classic ICO encoder and the dark mode lookup.
//
// fyne.io/systray only needs cgo on darwin; its Windows implementation is pure
// Go over the same Win32 calls this package already uses, so pulling it in keeps
// the driver cgo free.

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	_ "image/jpeg" // allow JPEG icon sources

	"fyne.io/systray"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/sys/windows/registry"

	"fyne.io/fyne/v2"
	intsystray "fyne.io/fyne/v2/internal/driver/systray"
)

// systrayIconSize matches what fyne.io/systray's LoadImage call produces: it
// passes LR_DEFAULTSIZE with zero extents, so Windows always materialises the
// icon at SM_CXICON (32px at 96 DPI). Handing it exactly that size avoids a
// lossy GDI stretch - Wine's in particular eats the anti-aliased edges.
// ponytail: fixed 32px; loses sharpness on high-DPI trays where SM_CXICON is
// larger. Fix by reading GetSystemMetrics(SM_CXICON) if that ever matters.
const systrayIconSize = 32

var systrayTray *intsystray.Tray

func (d *dxDriver) systray() *intsystray.Tray {
	if systrayTray == nil {
		systrayTray = &intsystray.Tray{
			IconSize:        systrayIconSize,
			RunOnMain:       runOnMain,
			Quit:            d.Quit,
			ToOSIcon:        toOSIcon,
			InvertMenuIcons: isDark, // Windows menus do not follow dark mode
		}
	}
	return systrayTray
}

func (d *dxDriver) SetSystemTrayMenu(m *fyne.Menu) {
	d.systrayMenu = m
	t := d.systray()
	if !t.Running() {
		d.trayStart, d.trayStop = t.Start(m, nil)
	}
	t.Refresh(m)
}

func (d *dxDriver) SystemTrayMenu() *fyne.Menu {
	return d.systrayMenu
}

func (d *dxDriver) SetSystemTrayIcon(resource fyne.Resource) {
	d.systray().SetIcon(resource)
}

func (d *dxDriver) SetSystemTrayWindow(w fyne.Window) {
	t := d.systray()
	if !t.Running() {
		d.trayStart, d.trayStop = t.Start(nil, nil)
	}

	w.SetCloseIntercept(w.Hide)
	win, ok := w.(*window)
	if !ok {
		return
	}
	systray.SetOnTapped(func() { fyne.Do(win.toggleVisible) })
}

// toOSIcon converts image bytes to a classic uncompressed ICO. The PNG-in-ICO
// format that github.com/fyne-io/image/ico produces is valid on Windows Vista
// and later, but Wine's icon loader cannot decode PNG entries ("NtUserDrawIconEx
// Error retrieving icon frame 0") and draws garbage - the 32bpp BMP entry
// written here is the original format every icon consumer understands.
func toOSIcon(icon []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(icon))
	if err != nil {
		return nil, err
	}
	// Scale to the size LoadImage will materialise anyway (LR_DEFAULTSIZE, see
	// systrayIconSize), with a real scaler rather than GDI's.
	const w, h = systrayIconSize, systrayIconSize
	straight := image.NewNRGBA(image.Rect(0, 0, w, h))
	if b := img.Bounds(); b.Dx() == w && b.Dy() == h {
		draw.Draw(straight, straight.Rect, img, b.Min, draw.Src)
	} else {
		xdraw.CatmullRom.Scale(straight, straight.Rect, img, b, xdraw.Src, nil)
	}

	xorSize := w * h * 4
	andStride := (w + 31) / 32 * 4 // 1bpp mask rows pad to 32-bit boundaries
	andSize := andStride * h
	entrySize := 40 + xorSize + andSize

	buf := &bytes.Buffer{}
	// ICONDIR (reserved, type icon, one image) + ICONDIRENTRY. The width and
	// height bytes hold 0 for 256px, which w%256 produces naturally.
	binary.Write(buf, binary.LittleEndian, [3]uint16{0, 1, 1})
	buf.Write([]byte{byte(w % 256), byte(h % 256), 0, 0}) // size, colours, reserved
	binary.Write(buf, binary.LittleEndian, [2]uint16{1, 32})
	binary.Write(buf, binary.LittleEndian, [2]uint32{uint32(entrySize), 6 + 16})

	// BITMAPINFOHEADER: the height counts the XOR and AND blocks together, and
	// icon bitmap rows are stored bottom-up.
	binary.Write(buf, binary.LittleEndian, [3]uint32{40, uint32(w), uint32(h * 2)})
	binary.Write(buf, binary.LittleEndian, [2]uint16{1, 32})
	binary.Write(buf, binary.LittleEndian, [2]uint32{0 /* BI_RGB */, uint32(xorSize + andSize)})
	buf.Write(make([]byte, 16)) // resolutions and colour counts, all zero

	for y := h - 1; y >= 0; y-- { // XOR data: bottom-up BGRA, straight alpha
		row := straight.Pix[y*straight.Stride : y*straight.Stride+w*4]
		for x := 0; x < len(row); x += 4 {
			buf.Write([]byte{row[x+2], row[x+1], row[x], row[x+3]})
		}
	}

	// AND mask, also bottom-up: renderers without per-pixel alpha (Wine's tray)
	// composite with this instead, so transparent pixels must set their bit or
	// the background draws as solid black.
	mask := make([]byte, andSize)
	for y := 0; y < h; y++ {
		maskRow := mask[(h-1-y)*andStride:]
		for x := 0; x < w; x++ {
			if straight.Pix[y*straight.Stride+x*4+3] < 0x80 {
				maskRow[x/8] |= 0x80 >> (x % 8)
			}
		}
	}
	buf.Write(mask)

	return buf.Bytes(), nil
}

// isDark reports whether Windows is in dark mode, so tray icons can be inverted.
func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil { // older Windows versions do not have this key
		return false
	}
	defer k.Close()

	useLight, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return false
	}
	return useLight == 0
}
