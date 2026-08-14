package systray

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	sys "fyne.io/systray"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// newTestTray returns a Tray that converts icons by passing the bytes straight
// through, so tests see exactly what the shared code produced.
func newTestTray(t *testing.T) *Tray {
	t.Helper()
	test.NewTempApp(t)

	return &Tray{
		IconSize:  32,
		RunOnMain: func(f func()) { f() },
		Quit:      func() {},
		ToOSIcon:  func(icon []byte) ([]byte, error) { return icon, nil },
	}
}

func TestTray_addMissingQuit(t *testing.T) {
	quitted := false
	tray := newTestTray(t)
	tray.Quit = func() { quitted = true }

	t.Run("appends a quit item when there is none", func(t *testing.T) {
		m := fyne.NewMenu("tray", fyne.NewMenuItem("Open", nil))
		tray.addMissingQuit(m)

		assert.Equal(t, 3, len(m.Items)) // Open, separator, Quit
		assert.True(t, m.Items[1].IsSeparator)
		assert.True(t, m.Items[2].IsQuit)

		m.Items[2].Action()
		assert.True(t, quitted)
	})

	t.Run("adopts a trailing item labelled Quit", func(t *testing.T) {
		mine := fyne.NewMenuItem("Quit", func() {})
		m := fyne.NewMenu("tray", fyne.NewMenuItem("Open", nil), mine)
		tray.addMissingQuit(m)

		assert.Equal(t, 2, len(m.Items))
		assert.True(t, mine.IsQuit)
	})

	t.Run("keeps an existing quit item and its action", func(t *testing.T) {
		called := false
		mine := fyne.NewMenuItem("Exit app", func() { called = true })
		mine.IsQuit = true
		m := fyne.NewMenu("tray", mine)
		tray.addMissingQuit(m)

		assert.Equal(t, 1, len(m.Items))
		m.Items[0].Action()
		assert.True(t, called)
	})

	t.Run("adds a quit item to an empty menu", func(t *testing.T) {
		m := fyne.NewMenu("tray")
		tray.addMissingQuit(m)

		assert.Equal(t, 2, len(m.Items))
		assert.True(t, m.Items[1].IsQuit)
	})
}

func TestShortcutKey(t *testing.T) {
	assert.Equal(t, "Enter", shortcutKey(fyne.KeyEnter))
	assert.Equal(t, "PageUp", shortcutKey(fyne.KeyPageUp))
	assert.Equal(t, "S", shortcutKey(fyne.KeyS)) // not remapped, passed straight through
}

func TestShortcutModifiers(t *testing.T) {
	assert.Equal(t, sys.KeyModifier(0), shortcutModifiers(0))
	assert.Equal(t, sys.KeyModifierControl, shortcutModifiers(fyne.KeyModifierControl))
	assert.Equal(t, sys.KeyModifierShift|sys.KeyModifierAlt,
		shortcutModifiers(fyne.KeyModifierShift|fyne.KeyModifierAlt))
}

func TestTray_trayIcon(t *testing.T) {
	tray := newTestTray(t)

	t.Run("rasterises SVG at the icon size", func(t *testing.T) {
		data, err := tray.trayIcon(theme.BrokenImageIcon())
		assert.NoError(t, err)

		img, err := png.Decode(bytes.NewReader(data))
		assert.NoError(t, err)
		assert.Equal(t, tray.IconSize, img.Bounds().Dx())
		assert.Equal(t, tray.IconSize, img.Bounds().Dy())
	})

	t.Run("passes SVG through where the platform supports it", func(t *testing.T) {
		tray.SupportsSVG = true
		defer func() { tray.SupportsSVG = false }()

		res := theme.BrokenImageIcon()
		data, err := tray.trayIcon(res)
		assert.NoError(t, err)
		assert.Equal(t, res.Content(), data)
	})

	t.Run("passes bitmaps through untouched", func(t *testing.T) {
		res := pngResource(t)
		data, err := tray.trayIcon(res)
		assert.NoError(t, err)
		assert.Equal(t, res.Content(), data)
	})
}

func TestTray_menuIcon(t *testing.T) {
	tray := newTestTray(t)

	t.Run("rasterises SVG at the icon size", func(t *testing.T) {
		data, err := tray.menuIcon(theme.BrokenImageIcon())
		assert.NoError(t, err)

		img, err := png.Decode(bytes.NewReader(data))
		assert.NoError(t, err)
		assert.Equal(t, tray.IconSize, img.Bounds().Dx())
	})

	t.Run("inverts SVG when the platform asks for it", func(t *testing.T) {
		plain, err := tray.menuIcon(theme.BrokenImageIcon())
		assert.NoError(t, err)

		asked := false
		tray.InvertMenuIcons = func() bool { asked = true; return true }
		defer func() { tray.InvertMenuIcons = nil }()

		inverted, err := tray.menuIcon(theme.BrokenImageIcon())
		assert.NoError(t, err)
		assert.True(t, asked)
		assert.NotEqual(t, plain, inverted)
	})

	t.Run("passes bitmaps through untouched", func(t *testing.T) {
		res := pngResource(t)
		data, err := tray.menuIcon(res)
		assert.NoError(t, err)
		assert.Equal(t, res.Content(), data)
	})
}

func pngResource(t *testing.T) fyne.Resource {
	t.Helper()

	img := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	img.Set(1, 1, color.NRGBA{R: 0xff, A: 0xff})

	buf := &bytes.Buffer{}
	assert.NoError(t, png.Encode(buf, img))
	return fyne.NewStaticResource("icon.png", buf.Bytes())
}
