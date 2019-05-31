// +build !ci

package gl

import (
	"image/color"
	"os"
	"runtime"
	"testing"
	"time"

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
	var menuHeight int
	if hasNativeMenu() {
		menuHeight = 0
	} else {
		menuHeight = 22
	}
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	tests := []struct {
		name               string
		padding            bool
		menu               bool
		expectedPad        int
		expectedMenuHeight int
	}{
		{"window without padding", false, false, 0, 0},
		{"window with padding", true, false, 4, 0},
		{"window with menu without padding", false, true, 0, menuHeight},
		{"window with menu and padding", true, true, 4, menuHeight},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := d.CreateWindow("Test").(*window)
			w.Canvas().SetScale(1)
			w.SetPadded(!tt.padding)
			if tt.menu {
				w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Test", fyne.NewMenuItem("Test", func() {}))))
			}
			content := canvas.NewRectangle(color.White)
			w.SetContent(content)
			oldCanvasSize := fyne.NewSize(100, 100)
			w.Resize(oldCanvasSize)

			// wait for canvas to get its size right
			for s := w.Canvas().Size(); s != oldCanvasSize; s = w.Canvas().Size() {
				time.Sleep(time.Millisecond * 10)
			}
			contentSize := content.Size()

			w.SetPadded(tt.padding)
			// wait (max 0.1s) for canvas resize
			for i := 0; i < 10; i++ {
				if w.Canvas().Size() != oldCanvasSize {
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			assert.Equal(t, contentSize, content.Size())
			assert.Equal(t, fyne.NewPos(tt.expectedPad, tt.expectedPad+tt.expectedMenuHeight), content.Position())
			expectedCanvasSize := contentSize.
				Add(fyne.NewSize(2*tt.expectedPad, 2*tt.expectedPad)).
				Add(fyne.NewSize(0, tt.expectedMenuHeight))
			assert.Equal(t, expectedCanvasSize, w.Canvas().Size())
		})
	}
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
