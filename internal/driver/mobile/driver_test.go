package mobile

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Test_mobileDriver_AbsolutePositionForObject(t *testing.T) {
	for name, tt := range map[string]struct {
		want          fyne.Position
		windowIsChild bool
		windowPadded  bool
	}{
		"for an unpadded primary (non-child) window it is (0,0)": {
			want:          fyne.NewPos(0, 0),
			windowIsChild: false,
			windowPadded:  false,
		},
		"for a padded primary (non-child) window it is (padding,padding)": {
			want:          fyne.NewPos(4, 4),
			windowIsChild: false,
			windowPadded:  true,
		},
		"for an unpadded child window it is (0,0)": {
			want:          fyne.NewPos(0, 0),
			windowIsChild: true,
			windowPadded:  false,
		},
		"for a padded child window it is (padding,padding)": {
			want:          fyne.NewPos(4, 4),
			windowIsChild: true,
			windowPadded:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			var o fyne.CanvasObject
			size := fyne.NewSize(100, 100)
			d := &driver{}
			w := d.CreateWindow("main")
			w.SetPadded(tt.windowPadded)
			l := widget.NewLabel("main window")
			if !tt.windowIsChild {
				o = l
			}
			w.SetContent(l)
			w.Show()
			w.Resize(size)
			w = d.CreateWindow("child1")
			w.SetContent(widget.NewLabel("first child"))
			if tt.windowIsChild {
				w.Show()
			}
			w.Resize(size)
			w = d.CreateWindow("child2 - hidden")
			w.SetContent(widget.NewLabel("second child"))
			w.Resize(size)
			w = d.CreateWindow("child3")
			r := fynecanvas.NewRectangle(color.White)
			r.SetMinSize(fyne.NewSize(42, 17))
			w.SetPadded(tt.windowPadded)
			w.SetContent(container.NewVBox(r))
			if tt.windowIsChild {
				w.Show()
				o = r
			}
			w.Resize(size)
			w = d.CreateWindow("child4 - hidden")
			w.SetContent(widget.NewLabel("fourth child"))
			w.Resize(size)

			got := d.AbsolutePositionForObject(o)
			assert.Equal(t, tt.want, got)
		})
	}
}
