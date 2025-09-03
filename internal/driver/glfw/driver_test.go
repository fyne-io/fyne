//go:build !no_glfw && !mobile

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func Test_gLDriver_AbsolutePositionForObject(t *testing.T) {
	w := createWindow("Test")

	cr1c1 := widget.NewLabel("row 1 col 1")
	cr1c2 := widget.NewLabel("row 1 col 2")
	cr1c3 := widget.NewLabel("row 1 col 3")
	cr2c1 := widget.NewLabel("row 2 col 1")
	cr2c2 := widget.NewLabel("row 2 col 2")
	cr2c3 := widget.NewLabel("row 2 col 3")
	cr3c1 := widget.NewLabel("row 3 col 1")
	cr3c2 := widget.NewLabel("row 3 col 2")
	cr3c3 := widget.NewLabel("row 3 col 3")
	cr1 := container.NewHBox(cr1c1, cr1c2, cr1c3)
	cr2 := container.NewHBox(cr2c1, cr2c2, cr2c3)
	cr3 := container.NewHBox(cr3c1, cr3c2, cr3c3)
	content := container.NewVBox(cr1, cr2, cr3)
	cr2c2.Hide()

	mm := fyne.NewMainMenu(
		fyne.NewMenu("Menu 1", fyne.NewMenuItem("Menu 1 Item", nil)),
		fyne.NewMenu("Menu 2", fyne.NewMenuItem("Menu 2 Item", nil)),
	)
	// We want to test the handling of the canvas' Fyne menu here.
	// We work around w.SetMainMenu because on MacOS the main menu is a native menu.
	c := w.Canvas()
	movl := buildMenuOverlay(mm, w.window)
	c.setMenuOverlay(movl)
	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 200))
	ensureCanvasSize(t, w, fyne.NewSize(300, 200))

	ovli1 := widget.NewLabel("Overlay Item 1")
	ovli2 := widget.NewLabel("Overlay Item 2")
	ovli3 := widget.NewLabel("Overlay Item 3")
	ovlContent := container.NewVBox(ovli1, ovli2, ovli3)
	// use the unsafe canvas so we can wrap the whole show in safety without deadlock
	ovl := widget.NewModalPopUp(ovlContent, w.window.Canvas())
	runOnMain(ovl.Show)

	// This helps to detect size changes which might happen when the font size or rendering are changed.
	// It gives also a hint on the expected offset for the overlay components.
	var size fyne.Size
	runOnMain(func() {
		size = ovlContent.Size()
	})
	assert.InDelta(t, 112, size.Width, 1)
	assert.InDelta(t, 113, size.Height, 1)

	repaintWindow(w)
	ensureCanvasSize(t, w, fyne.NewSize(300, 200))
	// accessing the menu bar's actual CanvasObjects isn't straight forward
	// 0 is the shadow
	// 1 is the menu bar’s underlay
	// 2 is the menu bar's background
	// 3 is the container holding the items
	var m2 fyne.CanvasObject
	runOnMain(func() {
		mbarCont := cache.Renderer(movl.(fyne.Widget)).Objects()[3].(*fyne.Container)
		m2 = mbarCont.Objects[1]
	})

	tests := map[string]struct {
		object       fyne.CanvasObject
		wantX, wantY int
	}{
		"a cell": {
			object: cr1c3,
			wantX:  184,
			wantY:  31,
		},
		"a row": {
			object: cr2,
			wantX:  4,
			wantY:  70,
		},
		"the window content": {
			object: content,
			wantX:  4,
			wantY:  31,
		},
		"a hidden element": {
			object: cr2c2,
			wantX:  0,
			wantY:  0,
		},

		"a menu": {
			object: m2,
			wantX:  77,
			wantY:  0,
		},

		"an overlay item": {
			object: ovli2,
			wantX:  93,
			wantY:  82,
		},
		"the overlay content": {
			object: ovlContent,
			wantX:  93,
			wantY:  43,
		},
		"the overlay": {
			object: ovl,
			wantX:  0,
			wantY:  0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runOnMain(func() {
				pos := d.AbsolutePositionForObject(tt.object)
				assert.Equal(t, tt.wantX, int(pos.X))
				assert.Equal(t, tt.wantY, int(pos.Y))
			})
		})
	}
}
