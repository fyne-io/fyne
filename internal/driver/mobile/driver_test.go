package mobile

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"
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

func Test_tapPosition(t *testing.T) {
	const fingerOffset float32 = -8
	unknown := touch.Event{}

	drv := &driver{}
	w := drv.CreateWindow("tap").(*window)
	w.SetPadded(false)
	w.SetContent(container.NewWithoutLayout())
	w.Resize(fyne.NewSize(200, 200))

	tests := []struct {
		name    string
		scale   float32
		x, y    float32
		offset  float32
		precise bool
		want    fyne.Position
	}{
		{name: "finger scale 1", scale: 1, x: 12, y: 40, offset: fingerOffset, precise: false, want: fyne.NewPos(12, 32)},
		{name: "unknown scale 1", scale: 1, x: 12, y: 40, offset: fingerOffset, precise: unknown.Precise, want: fyne.NewPos(12, 32)},
		{name: "precise scale 1", scale: 1, x: 12, y: 40, offset: fingerOffset, precise: true, want: fyne.NewPos(12, 40)},
		{name: "finger scale 2", scale: 2, x: 20, y: 40, offset: fingerOffset, precise: false, want: fyne.NewPos(10, 12)},
		{name: "precise scale 2", scale: 2, x: 20, y: 40, offset: fingerOffset, precise: true, want: fyne.NewPos(10, 20)},
		{name: "zero offset finger", scale: 1, x: 12, y: 40, offset: 0, precise: false, want: fyne.NewPos(12, 40)},
		{name: "zero offset precise", scale: 1, x: 12, y: 40, offset: 0, precise: true, want: fyne.NewPos(12, 40)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w.canvas.scale = tt.scale
			got := tapPosition(w, tt.x, tt.y, tt.offset, tt.precise)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_tapCanvas_precisePointer(t *testing.T) {
	const (
		objX float32 = 20
		objY float32 = 50
		objW float32 = 120
		objH float32 = 36
	)
	edgeX := objX + 8
	edgeY := objY
	belowY := objY + objH

	setup := func(t *testing.T) (*driver, *window, *pointerTarget) {
		t.Helper()
		drv := &driver{}
		w := drv.CreateWindow("precise").(*window)
		w.SetPadded(false)
		target := &pointerTarget{}
		w.SetContent(container.NewWithoutLayout(target))
		w.Resize(fyne.NewSize(240, 240))
		target.Move(fyne.NewPos(objX, objY))
		target.Resize(fyne.NewSize(objW, objH))
		w.canvas.scale = 1
		assert.Equal(t, fyne.NewPos(0, 0), w.canvas.content.Position())
		assert.Equal(t, fyne.NewPos(objX, objY), target.Position())
		return drv, w, target
	}

	t.Run("precise edge hits and below misses", func(t *testing.T) {
		drv, w, target := setup(t)
		drv.tapDownCanvas(w, edgeX, edgeY, 0, true)
		drv.tapUpCanvas(w, edgeX, edgeY, 0, true)
		want := tapPosition(w, edgeX, edgeY, tapYOffset, true)
		if assert.Len(t, target.downs, 1) {
			assert.Equal(t, want, target.downs[0].absolute)
			assert.Equal(t, 0, target.downs[0].id)
		}
		if assert.Len(t, target.ups, 1) {
			assert.Equal(t, want, target.ups[0].absolute)
		}
		assert.Equal(t, []fyne.Position{want}, target.taps)

		drv, w, target = setup(t)
		drv.tapDownCanvas(w, edgeX, belowY, 0, true)
		drv.tapUpCanvas(w, edgeX, belowY, 0, true)
		assert.Empty(t, target.downs)
		assert.Empty(t, target.ups)
		assert.Empty(t, target.taps)
	})

	t.Run("move and release follow the press", func(t *testing.T) {
		drv, w, target := setup(t)
		drv.tapDownCanvas(w, edgeX, edgeY, 0, true)
		drv.tapMoveCanvas(w, edgeX+8, edgeY, 0, true)
		move := tapPosition(w, edgeX+8, edgeY, tapYOffset, true)
		if assert.Len(t, target.moves, 1) {
			assert.Equal(t, move.Subtract(target.Position()), target.moves[0].position)
			assert.Equal(t, 0, target.moves[0].id)
		}
		if assert.Len(t, target.drags, 1) {
			assert.Equal(t, float32(8), target.drags[0].dx)
			assert.Equal(t, float32(0), target.drags[0].dy)
		}

		// The last delta has to stay under the inertia threshold.
		drv.tapMoveCanvas(w, edgeX+9, edgeY, 0, true)
		drv.tapUpCanvas(w, edgeX+9, edgeY, 0, true)
		assert.Equal(t, 1, target.ends)
		assert.Empty(t, target.taps)

		drv, w, target = setup(t)
		drv.tapDownCanvas(w, edgeX, belowY, 0, true)
		drv.tapMoveCanvas(w, edgeX+8, belowY, 0, true)
		drv.tapUpCanvas(w, edgeX+8, belowY, 0, true)
		assert.Empty(t, target.moves)
		assert.Empty(t, target.drags)
		assert.Empty(t, target.taps)
	})

	t.Run("finger compensation and legacy events", func(t *testing.T) {
		drv, w, target := setup(t)
		fingerY := edgeY - tapYOffset
		want := tapPosition(w, edgeX, fingerY, tapYOffset, false)
		assert.Equal(t, fyne.NewPos(edgeX, edgeY), want)

		legacy := touch.Event{X: edgeX, Y: fingerY, Type: touch.TypeBegin, Sequence: 4}
		assert.False(t, legacy.Precise)
		drv.tapDownCanvas(w, legacy.X, legacy.Y, legacy.Sequence, legacy.Precise)
		if assert.Len(t, target.downs, 1) {
			assert.Equal(t, want, target.downs[0].absolute)
			assert.Equal(t, int(legacy.Sequence), target.downs[0].id)
		}

		if tapYOffset == 0 {
			assert.Equal(t, fyne.NewPos(edgeX, edgeY), tapPosition(w, edgeX, edgeY, tapYOffset, true))
		}
		if tapYOffset < 0 {
			fingerBelow := tapPosition(w, edgeX, belowY, tapYOffset, false)
			assert.True(t, posInTarget(target, fingerBelow))
			assert.False(t, posInTarget(target, tapPosition(w, edgeX, belowY, tapYOffset, true)))

			drv, w, target = setup(t)
			drv.tapDownCanvas(w, edgeX, belowY, 0, false)
			if assert.Len(t, target.downs, 1) {
				assert.Equal(t, fingerBelow, target.downs[0].absolute)
			}
			drv, w, target = setup(t)
			drv.tapDownCanvas(w, edgeX, belowY, 0, true)
			assert.Empty(t, target.downs)
		}
	})

	t.Run("concurrent finger and stylus", func(t *testing.T) {
		drv, w, target := setup(t)
		const (
			fingerID = touch.Sequence(0)
			stylusID = touch.Sequence(1)
		)
		drv.tapDownCanvas(w, edgeX, edgeY, stylusID, true)
		drv.tapDownCanvas(w, edgeX, belowY, fingerID, false)
		drv.tapMoveCanvas(w, edgeX+2, edgeY, stylusID, true)
		drv.tapMoveCanvas(w, edgeX+2, belowY, fingerID, false)
		drv.tapUpCanvas(w, edgeX+2, edgeY, stylusID, true)
		drv.tapUpCanvas(w, edgeX+2, belowY, fingerID, false)

		stylusDown := tapPosition(w, edgeX, edgeY, tapYOffset, true)
		fingerDown := tapPosition(w, edgeX, belowY, tapYOffset, false)
		stylusMove := tapPosition(w, edgeX+2, edgeY, tapYOffset, true)
		fingerMove := tapPosition(w, edgeX+2, belowY, tapYOffset, false)

		down, ok := hitByID(target.downs, int(stylusID))
		if assert.True(t, ok) {
			assert.Equal(t, stylusDown, down.absolute)
		}
		if posInTarget(target, fingerDown) {
			down, ok = hitByID(target.downs, int(fingerID))
			if assert.True(t, ok) {
				assert.Equal(t, fingerDown, down.absolute)
			}
		} else {
			_, ok = hitByID(target.downs, int(fingerID))
			assert.False(t, ok)
		}

		move, ok := hitByID(target.moves, int(stylusID))
		if assert.True(t, ok) {
			assert.Equal(t, stylusMove.Subtract(target.Position()), move.position)
		}
		if posInTarget(target, fingerMove) {
			move, ok = hitByID(target.moves, int(fingerID))
			if assert.True(t, ok) {
				assert.Equal(t, fingerMove.Subtract(target.Position()), move.position)
			}
		} else {
			_, ok = hitByID(target.moves, int(fingerID))
			assert.False(t, ok)
		}

		up, ok := hitByID(target.ups, int(stylusID))
		if assert.True(t, ok) {
			assert.Equal(t, stylusMove, up.absolute)
		}
		wantTaps := []fyne.Position{stylusMove}
		if posInTarget(target, fingerMove) {
			up, ok = hitByID(target.ups, int(fingerID))
			if assert.True(t, ok) {
				assert.Equal(t, fingerMove, up.absolute)
			}
			wantTaps = append(wantTaps, fingerMove)
		} else {
			_, ok = hitByID(target.ups, int(fingerID))
			assert.False(t, ok)
		}
		assert.Equal(t, wantTaps, target.taps)
	})
}

func posInTarget(o fyne.CanvasObject, pos fyne.Position) bool {
	top := o.Position()
	size := o.Size()
	return pos.X >= top.X && pos.Y >= top.Y && pos.X < top.X+size.Width && pos.Y < top.Y+size.Height
}

func hitByID(hits []pointerHit, id int) (pointerHit, bool) {
	for _, hit := range hits {
		if hit.id == id {
			return hit, true
		}
	}
	return pointerHit{}, false
}

type pointerHit struct {
	id       int
	absolute fyne.Position
	position fyne.Position
}

type pointerDrag struct {
	dx, dy float32
}

type pointerTarget struct {
	fynecanvas.Rectangle

	downs []pointerHit
	ups   []pointerHit
	moves []pointerHit
	drags []pointerDrag
	taps  []fyne.Position
	ends  int
}

var (
	_ mobile.Touchable = (*pointerTarget)(nil)
	_ mobile.Movable   = (*pointerTarget)(nil)
	_ fyne.Draggable   = (*pointerTarget)(nil)
	_ fyne.Tappable    = (*pointerTarget)(nil)
)

func (p *pointerTarget) TouchDown(ev *mobile.TouchEvent) {
	p.downs = append(p.downs, pointerHit{id: ev.ID, absolute: ev.AbsolutePosition, position: ev.Position})
}

func (p *pointerTarget) TouchUp(ev *mobile.TouchEvent) {
	p.ups = append(p.ups, pointerHit{id: ev.ID, absolute: ev.AbsolutePosition, position: ev.Position})
}

func (*pointerTarget) TouchCancel(*mobile.TouchEvent) {}

func (p *pointerTarget) TouchMoved(ev *mobile.TouchEvent) {
	p.moves = append(p.moves, pointerHit{id: ev.ID, absolute: ev.AbsolutePosition, position: ev.Position})
}

func (p *pointerTarget) Dragged(ev *fyne.DragEvent) {
	p.drags = append(p.drags, pointerDrag{dx: ev.Dragged.DX, dy: ev.Dragged.DY})
}

func (p *pointerTarget) DragEnd() {
	p.ends++
}

func (p *pointerTarget) Tapped(ev *fyne.PointEvent) {
	p.taps = append(p.taps, ev.AbsolutePosition)
}
