//go:build windows

package directx

// System tray support, ported from the GLFW driver's driver_desktop.go.
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
	"image/png"

	"fyne.io/systray"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/sys/windows/registry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

// systrayIconSize matches what fyne.io/systray's LoadImage call produces: it
// passes LR_DEFAULTSIZE with zero extents, so Windows always materialises the
// icon at SM_CXICON (32px at 96 DPI). Handing it exactly that size avoids a
// lossy GDI stretch - Wine's in particular eats the anti-aliased edges.
// ponytail: fixed 32px; loses sharpness on high-DPI trays where SM_CXICON is
// larger. Fix by reading GetSystemMetrics(SM_CXICON) if that ever matters.
const systrayIconSize = 32

var (
	systrayIcon    fyne.Resource
	systrayRunning bool
)

func (d *dxDriver) SetSystemTrayMenu(m *fyne.Menu) {
	if !systrayRunning {
		systrayRunning = true
		d.runSystray(m)
	}
	d.refreshSystray(m)
}

func (d *dxDriver) SystemTrayMenu() *fyne.Menu {
	return d.systrayMenu
}

func (d *dxDriver) runSystray(m *fyne.Menu) {
	d.trayStart, d.trayStop = systray.RunWithExternalLoop(func() {
		switch {
		case systrayIcon != nil:
			d.SetSystemTrayIcon(systrayIcon)
		case fyne.CurrentApp().Icon() != nil:
			d.SetSystemTrayIcon(fyne.CurrentApp().Icon())
		default:
			d.SetSystemTrayIcon(theme.BrokenImageIcon())
		}

		if m != nil {
			// The menu has to be rebuilt after init; doing it earlier has no effect.
			runOnMain(func() {
				d.refreshSystray(m)
			})
		}
	}, func() {
		// nothing to tear down
	})
}

func (d *dxDriver) refreshSystray(m *fyne.Menu) {
	d.systrayMenu = m

	systray.ResetMenu()
	d.refreshSystrayMenu(m, nil)

	addMissingQuitForMenu(m, d)
}

func (d *dxDriver) refreshSystrayMenu(m *fyne.Menu, parent *systray.MenuItem) {
	if m == nil {
		return
	}
	for _, i := range m.Items {
		item := itemForMenuItem(i, parent)
		if item == nil {
			continue // separator
		}
		if i.ChildMenu != nil {
			d.refreshSystrayMenu(i.ChildMenu, item)
		}

		fn := i.Action
		go func() {
			for range item.ClickedCh {
				if fn != nil {
					runOnMain(fn)
				}
			}
		}()
	}
}

func itemForMenuItem(i *fyne.MenuItem, parent *systray.MenuItem) *systray.MenuItem {
	if i.IsSeparator {
		if parent != nil {
			parent.AddSeparator()
		} else {
			systray.AddSeparator()
		}
		return nil
	}

	var item *systray.MenuItem
	switch {
	case i.Checked && parent != nil:
		item = parent.AddSubMenuItemCheckbox(i.Label, "", true)
	case i.Checked:
		item = systray.AddMenuItemCheckbox(i.Label, "", true)
	case parent != nil:
		item = parent.AddSubMenuItem(i.Label, "")
	default:
		item = systray.AddMenuItem(i.Label, "")
	}

	if i.Disabled {
		item.Disable()
	}
	if i.Icon == nil {
		return item
	}

	data := i.Icon.Content()
	if svg.IsResourceSVG(i.Icon) {
		res := i.Icon
		if isDark() { // Windows menus do not follow dark mode, so invert the icon
			res = theme.NewInvertedThemedResource(i.Icon)
		}
		b := &bytes.Buffer{}
		img := paint.PaintImage(canvas.NewImageFromResource(res), nil, systrayIconSize, systrayIconSize)
		if err := png.Encode(b, img); err != nil {
			fyne.LogError("directx: encode SVG icon for menu", err)
		} else {
			data = b.Bytes()
		}
	}

	img, err := toOSIcon(data)
	if err != nil {
		fyne.LogError("directx: convert systray menu icon", err)
		return item
	}
	if _, ok := i.Icon.(*theme.ThemedResource); ok {
		item.SetTemplateIcon(img, img)
	} else {
		item.SetIcon(img)
	}
	return item
}

func (*dxDriver) SetSystemTrayIcon(resource fyne.Resource) {
	systrayIcon = resource // kept in case the tray is (re)started later

	// Windows has no SVG tray icon support, so rasterise first.
	if svg.IsResourceSVG(resource) {
		img := canvas.NewImageFromResource(resource)
		c := software.NewTransparentCanvas()
		c.SetContent(img)
		c.SetPadded(false)
		c.Resize(fyne.NewSquareSize(systrayIconSize))

		buf := &bytes.Buffer{}
		if err := png.Encode(buf, c.Capture()); err != nil {
			fyne.LogError("directx: encode SVG system tray icon", err)
			return
		}
		resource = fyne.NewStaticResource(resource.Name()+".png", buf.Bytes())
	}

	img, err := toOSIcon(resource.Content())
	if err != nil {
		fyne.LogError("directx: convert system tray icon", err)
		return
	}

	if _, ok := resource.(*theme.ThemedResource); ok {
		systray.SetTemplateIcon(img, img)
	} else {
		systray.SetIcon(img)
	}
}

func (d *dxDriver) SetSystemTrayWindow(w fyne.Window) {
	if !systrayRunning {
		systrayRunning = true
		d.runSystray(nil)
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

// addMissingQuitForMenu guarantees the tray menu can always quit the app, which
// matters here because a driver with a tray menu deliberately keeps running after
// its last window closes.
func addMissingQuitForMenu(menu *fyne.Menu, d *dxDriver) {
	localQuit := lang.L("Quit")

	var lastItem *fyne.MenuItem
	if len(menu.Items) > 0 {
		lastItem = menu.Items[len(menu.Items)-1]
		if lastItem.Label == localQuit {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit {
		quitItem := fyne.NewMenuItem(localQuit, nil)
		quitItem.IsQuit = true
		menu.Items = append(menu.Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menu.Items {
		if item.IsQuit && item.Action == nil {
			item.Action = d.Quit
		}
	}
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
