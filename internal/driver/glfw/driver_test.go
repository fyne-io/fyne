// +build !ci
// +build !mobile

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
	c := w.Canvas().(*glCanvas)
	movl := buildMenuOverlay(mm, c)
	c.setMenuOverlay(movl)
	w.SetContent(content)
	w.Resize(fyne.NewSize(200, 199))

	ovli1 := widget.NewLabel("Overlay Item 1")
	ovli2 := widget.NewLabel("Overlay Item 2")
	ovli3 := widget.NewLabel("Overlay Item 3")
	ovlContent := container.NewVBox(ovli1, ovli2, ovli3)
	ovl := widget.NewModalPopUp(ovlContent, c)
	ovl.Show()

	repaintWindow(w.(*window))
	// accessing the menu bar's actual CanvasObjects isn't straight forward
	// 0 is the shadow
	// 1 is the menu bar’s underlay
	// 2 is the menu bar's background
	// 3 is the container holding the items
	mbarCont := cache.Renderer(movl.(fyne.Widget)).Objects()[3].(*fyne.Container)
	m2 := mbarCont.Objects[1]

	tests := map[string]struct {
		object fyne.CanvasObject
		want   fyne.Position
	}{
		"a cell": {
			object: cr1c3,
			want:   fyne.NewPos(182, 33),
		},
		"a row": {
			object: cr2,
			want:   fyne.NewPos(4, 74),
		},
		"the window content": {
			object: content,
			want:   fyne.NewPos(4, 33),
		},
		"a hidden element": {
			object: cr2c2,
			want:   fyne.NewPos(0, 0),
		},

		"a menu": {
			object: m2,
			want:   fyne.NewPos(78, 0),
		},

		"an overlay item": {
			object: ovli2,
			want:   fyne.NewPos(79, 81),
		},
		"the overlay content": {
			object: ovlContent,
			want:   fyne.NewPos(79, 40),
		},
		"the overlay": {
			object: ovl,
			want:   fyne.NewPos(0, 0),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, d.AbsolutePositionForObject(tt.object))
		})
	}
}

var mainRoutineID int

func init() {
	mainRoutineID = goroutineID()
}

func TestGoroutineID(t *testing.T) {
	assert.Equal(t, 1, mainRoutineID)

	var childID1, childID2 int
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
	assert.Greater(t, childID1, 0)
	assert.NotEqual(t, testID1, childID1)
	assert.Greater(t, childID2, 0)
	assert.NotEqual(t, childID1, childID2)
}
