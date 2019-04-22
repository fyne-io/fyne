package gl

import (
	"image/color"
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

var d = NewGLDriver()

func init() {
	runtime.LockOSThread()
}

// TestMain makes sure that our driver is running on the main thread.
// This must be done for some of our tests to function correctly.
func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
	d.Run()
}

func TestWindow_SetTitle(t *testing.T) {
	w := d.CreateWindow("Test")

	title := "My title"
	w.SetTitle(title)

	assert.Equal(t, title, w.Title())
}

func TestWindow_SetIcon(t *testing.T) {
	w := d.CreateWindow("Test")
	assert.Equal(t, fyne.CurrentApp().Icon(), w.Icon())

	newIcon := theme.CancelIcon()
	w.SetIcon(newIcon)
	assert.Equal(t, newIcon, w.Icon())
}

func TestWindow_PixelSize(t *testing.T) {
	w := d.CreateWindow("Test")
	w.SetPadded(false)

	rect := &canvas.Rectangle{}
	rect.SetMinSize(fyne.NewSize(100, 100))
	w.SetContent(fyne.NewContainer(rect))
	w.Canvas().Refresh(w.Content())

	winW, winH := w.(*window).minSizeOnScreen()
	assert.Equal(t, scaleInt(w.Canvas(), 100), winW)
	assert.Equal(t, scaleInt(w.Canvas(), 100), winH)
}

func TestWindow_Padded(t *testing.T) {
	w := d.CreateWindow("Test")
	content := canvas.NewRectangle(color.White)
	w.Canvas().SetScale(1.0)
	w.SetContent(content)

	width, _ := w.(*window).minSizeOnScreen()
	assert.Equal(t, theme.Padding()*2+content.MinSize().Width, width)
	assert.Equal(t, theme.Padding(), content.Position().X)
}

func TestWindow_SetPadded(t *testing.T) {
	w := d.CreateWindow("Test")
	w.SetPadded(false)

	content := canvas.NewRectangle(color.White)
	w.Canvas().SetScale(1.0)
	w.SetContent(content)

	width, _ := w.(*window).minSizeOnScreen()
	assert.Equal(t, content.MinSize().Width, width)
	assert.Equal(t, 0, content.Position().X)

	w.SetPadded(true)
	width, _ = w.(*window).minSizeOnScreen()
	assert.Equal(t, theme.Padding()*2+content.MinSize().Width, width)
	assert.Equal(t, theme.Padding(), content.Position().X)
}

func TestWindow_Clipboard(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	text := "My content from test window"
	cb := w.Clipboard()

	cliboardContent := cb.Content()
	if cliboardContent != "" {
		// Current environment has some content stored in clipboard,
		// set temporary to an empty string to allow test and restore later.
		cb.SetContent("")
	}

	assert.Empty(t, cb.Content())

	cb.SetContent(text)
	assert.Equal(t, text, cb.Content())

	// Restore cliboardContent, if any
	cb.SetContent(cliboardContent)
}

func TestWindow_Shortcut(t *testing.T) {
	d := NewGLDriver()
	w := d.CreateWindow("Test")

	shortcutFullScreenWindow := &desktop.CustomShortcut{
		KeyName: fyne.KeyF12,
	}

	w.Canvas().AddShortcut(shortcutFullScreenWindow, func(sc fyne.Shortcut) {
		w.SetFullScreen(true)
	})

	assert.False(t, w.FullScreen())

	w.Canvas().(*glCanvas).shortcut.TypedShortcut(shortcutFullScreenWindow)
	assert.True(t, w.FullScreen())
}
