// Package systray holds the system tray support that is shared between Fyne's
// desktop drivers - menu construction, icon conversion and the automatic quit
// item. Only the parts that need a driver's own windowing, rendering or main
// thread handling are left to the drivers themselves, supplied through the
// fields of Tray.
package systray

import (
	"bytes"
	_ "image/jpeg" // allow JPEGs as icon sources
	"image/png"

	sys "fyne.io/systray"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

// Tray drives the system tray for one driver. The function fields are required
// and must be set before any other method is called.
type Tray struct {
	// IconSize is the pixel size that SVG resources are rasterised to.
	IconSize int
	// RunOnMain runs f on the driver's main goroutine.
	RunOnMain func(f func())
	// Quit shuts the app down, for the tray's automatic quit item.
	Quit func()
	// ToOSIcon converts image bytes to the icon format the platform tray wants.
	ToOSIcon func(icon []byte) ([]byte, error)
	// SupportsSVG reports whether the platform tray takes SVG resources as they
	// are; when false they are rasterised to PNG at IconSize first.
	SupportsSVG bool
	// InvertMenuIcons reports whether themed SVG menu icons should be inverted,
	// which Windows needs because its menus do not follow the dark theme.
	InvertMenuIcons func() bool

	icon    fyne.Resource
	running bool
}

// Running reports whether the tray has been started.
func (t *Tray) Running() bool { return t.running }

// Start brings the tray up with the given menu, which may be nil, and returns
// the start and stop functions that the driver's run loop must call. The ready
// callback, if set, runs once the tray is initialised.
func (t *Tray) Start(m *fyne.Menu, ready func()) (start, stop func()) {
	t.running = true
	return sys.RunWithExternalLoop(func() {
		switch {
		case t.icon != nil:
			t.SetIcon(t.icon)
		case fyne.CurrentApp().Icon() != nil:
			t.SetIcon(fyne.CurrentApp().Icon())
		default:
			t.SetIcon(theme.BrokenImageIcon())
		}

		if ready != nil {
			ready()
		}

		if m != nil {
			// the menu has to be rebuilt after init, doing it earlier has no effect
			t.RunOnMain(func() {
				t.Refresh(m)
			})
		}
	}, func() {
		// nothing to tear down
	})
}

// Refresh replaces the tray menu with m, adding a quit item if it has none.
func (t *Tray) Refresh(m *fyne.Menu) {
	sys.ResetMenu()
	t.refreshMenu(m, nil)

	t.addMissingQuit(m)
}

func (t *Tray) refreshMenu(m *fyne.Menu, parent *sys.MenuItem) {
	if m == nil {
		return
	}

	for _, i := range m.Items {
		item := t.itemForMenuItem(i, parent)
		if item == nil {
			continue // separator
		}
		if i.ChildMenu != nil {
			t.refreshMenu(i.ChildMenu, item)
		}

		fn := i.Action
		go func() {
			for range item.ClickedCh {
				if fn != nil {
					t.RunOnMain(fn)
				}
			}
		}()
	}
}

func (t *Tray) itemForMenuItem(i *fyne.MenuItem, parent *sys.MenuItem) *sys.MenuItem {
	if i.IsSeparator {
		if parent != nil {
			parent.AddSeparator()
		} else {
			sys.AddSeparator()
		}
		return nil
	}

	var item *sys.MenuItem
	switch {
	case i.Checked && parent != nil:
		item = parent.AddSubMenuItemCheckbox(i.Label, "", true)
	case i.Checked:
		item = sys.AddMenuItemCheckbox(i.Label, "", true)
	case parent != nil:
		item = parent.AddSubMenuItem(i.Label, "")
	default:
		item = sys.AddMenuItem(i.Label, "")
	}

	if i.Disabled {
		item.Disable()
	}
	if s, ok := i.Shortcut.(fyne.KeyboardShortcut); ok {
		item.SetShortcut(shortcutModifiers(s.Mod()), shortcutKey(s.Key()))
	}
	if i.Icon == nil {
		return item
	}

	img, err := t.menuIcon(i.Icon)
	if err != nil {
		fyne.LogError("Failed to convert systray icon", err)
		return item
	}
	if _, ok := i.Icon.(*theme.ThemedResource); ok {
		item.SetTemplateIcon(img, img)
	} else {
		item.SetIcon(img)
	}
	return item
}

// shortcutKeys maps the few key names that Fyne spells differently to the
// platform neutral names that the systray package understands.
var shortcutKeys = map[fyne.KeyName]string{
	fyne.KeyEnter:    "Enter",
	fyne.KeyPageDown: "PageDown",
	fyne.KeyPageUp:   "PageUp",
}

func shortcutKey(key fyne.KeyName) string {
	if name, ok := shortcutKeys[key]; ok {
		return name
	}
	return string(key)
}

func shortcutModifiers(mod fyne.KeyModifier) (mods sys.KeyModifier) {
	if mod&fyne.KeyModifierShift != 0 {
		mods |= sys.KeyModifierShift
	}
	if mod&fyne.KeyModifierControl != 0 {
		mods |= sys.KeyModifierControl
	}
	if mod&fyne.KeyModifierAlt != 0 {
		mods |= sys.KeyModifierAlt
	}
	if mod&fyne.KeyModifierSuper != 0 {
		mods |= sys.KeyModifierSuper
	}
	return mods
}

// menuIcon converts a menu item icon to the platform icon format, rasterising
// SVG resources at IconSize and inverting them where the platform asks for it.
func (t *Tray) menuIcon(res fyne.Resource) ([]byte, error) {
	data := res.Content()
	if svg.IsResourceSVG(res) {
		src := res
		if t.InvertMenuIcons != nil && t.InvertMenuIcons() {
			src = theme.NewInvertedThemedResource(res)
		}
		b := &bytes.Buffer{}
		img := painter.PaintImage(canvas.NewImageFromResource(src), nil, t.IconSize, t.IconSize)
		if err := png.Encode(b, img); err != nil {
			fyne.LogError("Failed to encode SVG icon for menu", err)
		} else {
			data = b.Bytes()
		}
	}
	return t.ToOSIcon(data)
}

// trayIcon converts the tray icon resource to the platform icon format,
// rasterising SVG resources at IconSize unless the platform takes them as-is.
func (t *Tray) trayIcon(res fyne.Resource) ([]byte, error) {
	data := res.Content()
	if !t.SupportsSVG && svg.IsResourceSVG(res) {
		c := software.NewTransparentCanvas()
		c.SetContent(canvas.NewImageFromResource(res))
		c.SetPadded(false)
		c.Resize(fyne.NewSquareSize(float32(t.IconSize)))

		buf := &bytes.Buffer{}
		if err := png.Encode(buf, c.Capture()); err != nil {
			return nil, err
		}
		data = buf.Bytes()
	}
	return t.ToOSIcon(data)
}

// SetIcon shows resource as the tray icon.
func (t *Tray) SetIcon(resource fyne.Resource) {
	t.icon = resource // in case the tray is (re)started later

	img, err := t.trayIcon(resource)
	if err != nil {
		fyne.LogError("Failed to convert systray icon", err)
		return
	}

	if _, ok := resource.(*theme.ThemedResource); ok {
		sys.SetTemplateIcon(img, img)
	} else {
		sys.SetIcon(img)
	}
}

// addMissingQuit guarantees the tray menu can always quit the app, which matters
// because a driver with a tray menu keeps running after its last window closes.
func (t *Tray) addMissingQuit(menu *fyne.Menu) {
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
			item.Action = t.Quit
		}
	}
}
