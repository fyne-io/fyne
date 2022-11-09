//go:build !ci && !mobile
// +build !ci,!mobile

package glfw

import (
	"sync"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func Test_gLDriver_AbsolutePositionForObject(t *testing.T) {
	w := createWindow("Test").(*window)

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
	c := w.Canvas().(*glCanvas)
	movl := buildMenuOverlay(mm, w)
	c.Lock()
	c.setMenuOverlay(movl)
	c.Unlock()
	w.SetContent(content)
	w.Resize(fyne.NewSize(200, 199))

	ovli1 := widget.NewLabel("Overlay Item 1")
	ovli2 := widget.NewLabel("Overlay Item 2")
	ovli3 := widget.NewLabel("Overlay Item 3")
	ovlContent := container.NewVBox(ovli1, ovli2, ovli3)
	ovl := widget.NewModalPopUp(ovlContent, c)
	ovl.Show()

	repaintWindow(w)
	// accessing the menu bar's actual CanvasObjects isn't straight forward
	// 0 is the shadow
	// 1 is the menu bar’s underlay
	// 2 is the menu bar's background
	// 3 is the container holding the items
	mbarCont := cache.Renderer(movl.(fyne.Widget)).Objects()[3].(*fyne.Container)
	m2 := mbarCont.Objects[1]

	tests := map[string]struct {
		object       fyne.CanvasObject
		wantX, wantY int
	}{
		"a cell": {
			object: cr1c3,
			wantX:  180,
			wantY:  35,
		},
		"a row": {
			object: cr2,
			wantX:  6,
			wantY:  75,
		},
		"the window content": {
			object: content,
			wantX:  6,
			wantY:  35,
		},
		"a hidden element": {
			object: cr2c2,
			wantX:  0,
			wantY:  0,
		},

		"a menu": {
			object: m2,
			wantX:  84,
			wantY:  0,
		},

		"an overlay item": {
			object: ovli2,
			wantX:  81,
			wantY:  60,
		},
		"the overlay content": {
			object: ovlContent,
			wantX:  81,
			wantY:  28,
		},
		"the overlay": {
			object: ovl,
			wantX:  0,
			wantY:  0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			pos := d.AbsolutePositionForObject(tt.object)
			assert.Equal(t, tt.wantX, int(pos.X))
			assert.Equal(t, tt.wantY, int(pos.Y))
		})
	}
}

var mainRoutineID uint64

func init() {
	mainRoutineID = goroutineID()
}

func TestGoroutineID(t *testing.T) {
	assert.Equal(t, uint64(1), mainRoutineID)

	var childID1, childID2 uint64
	testID1 := goroutineID()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		childID1 = goroutineID()
		wg.Done()
	}()
	go func() {
		childID2 = goroutineID()
		wg.Done()
	}()
	wg.Wait()
	testID2 := goroutineID()

	assert.Equal(t, testID1, testID2)
	assert.Greater(t, childID1, uint64(0))
	assert.NotEqual(t, testID1, childID1)
	assert.Greater(t, childID2, uint64(0))
	assert.NotEqual(t, childID1, childID2)
}
