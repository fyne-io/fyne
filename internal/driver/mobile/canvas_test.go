//go:build !windows || !ci

package mobile

import (
	"image/color"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal/async"
	intdriver "fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	intwidget "fyne.io/fyne/v2/internal/widget"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

var d *driver

func TestMain(m *testing.M) {
	currentApp := fyne.CurrentApp()
	tester := newTestMobileApp()
	d = tester.Driver().(*driver)
	d.queuedFuncs = async.NewUnboundedChan[func()]()
	fyne.SetCurrentApp(tester)

	waitForStart := make(chan struct{})
	go func() {
		// Wait for app loop to be running (plus a moment in case of scheduling switches).
		<-waitForStart

		// Just like the GLFW tests, wait a short while for the driver to start
		time.Sleep(time.Millisecond * 100)

		ret := m.Run()
		fyne.SetCurrentApp(currentApp)
		os.Exit(ret) //revive:disable-line:redundant-test-main-exit
	}()

	close(waitForStart) // Signal that execution can continue.
	tester.Run()
}

func Test_canvas_ChildMinSizeChangeAffectsAncestorsUpToRoot(t *testing.T) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	leftObj1 := fynecanvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(100, 50))
	leftObj2 := fynecanvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(100, 50))
	leftCol := container.NewVBox(leftObj1, leftObj2)
	rightObj1 := fynecanvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(100, 50))
	rightObj2 := fynecanvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(100, 50))
	rightCol := container.NewVBox(rightObj1, rightObj2)
	content := container.NewHBox(leftCol, rightCol)
	c.SetContent(container.NewCenter(content))
	c.Resize(fyne.NewSize(300, 300))

	oldContentSize := fyne.NewSize(200+theme.Padding(), 100+theme.Padding())
	assert.Equal(t, oldContentSize, content.Size())

	leftObj1.SetMinSize(fyne.NewSize(110, 60))
	c.EnsureMinSize()

	expectedContentSize := oldContentSize.Add(fyne.NewSize(10, 10))
	assert.Equal(t, expectedContentSize, content.Size())
}

func Test_canvas_DraggingOutOfWidget(t *testing.T) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	slider := widget.NewSlider(0.0, 100.0)
	c.SetContent(container.NewGridWithRows(2, slider, widget.NewLabel("Outside")))
	c.Resize(fyne.NewSize(100, 200))

	assert.Zero(t, slider.Value)
	lastValue := slider.Value

	dragged := false
	c.tapDown(fyne.NewPos(23, 13), 0)
	c.tapMove(fyne.NewPos(30, 13), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		assert.Equal(t, slider, wid)
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Greater(t, slider.Value, lastValue)
	lastValue = slider.Value

	dragged = false
	c.tapMove(fyne.NewPos(40, 120), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		assert.Equal(t, slider, wid)
		assert.Equal(t, fyne.NewPos(40, 120), ev.AbsolutePosition)
		wid.Dragged(ev)
		dragged = true
	})
	assert.True(t, dragged)
	assert.Greater(t, slider.Value, lastValue)
}

func Test_canvas_Focusable(t *testing.T) {
	// Discarding log output for tests
	// The following method logs an error:
	// wid.Tapped(ev) on lines 148 and 159
	// c.Focus(content) on line 168
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	content := newFocusableEntry(c)
	c.SetContent(content)
	content.Resize(fyne.NewSize(25, 25))

	pos := fyne.NewPos(10, 10)
	d.DoFromGoroutine(func() {
		c.tapDown(pos, 0)
		c.tapUp(pos, 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
			wid.Tapped(ev)
		}, nil, nil, nil)
	}, true)
	waitAndCheck(tapDoubleDelay/time.Millisecond+150, func() {
		assert.Equal(t, 1, content.focusedTimes)
		assert.Equal(t, 0, content.unfocusedTimes)
	})

	d.DoFromGoroutine(func() {
		c.tapDown(pos, 1)
		c.tapUp(pos, 1, func(wid fyne.Tappable, ev *fyne.PointEvent) {
			wid.Tapped(ev)
		}, nil, nil, nil)
	}, true)
	waitAndCheck(tapDoubleDelay/time.Millisecond+150, func() {
		assert.Equal(t, 1, content.focusedTimes)
		assert.Equal(t, 0, content.unfocusedTimes)
	})

	d.DoFromGoroutine(func() {
		c.Focus(content)
		assert.Equal(t, 1, content.focusedTimes)
		assert.Equal(t, 0, content.unfocusedTimes)

		c.Unfocus()
		assert.Equal(t, 1, content.focusedTimes)
		assert.Equal(t, 1, content.unfocusedTimes)

		content.Disable()
		c.Focus(content)
		assert.Equal(t, 1, content.focusedTimes)
		assert.Equal(t, 1, content.unfocusedTimes)

		c.tapDown(fyne.NewPos(10, 10), 2)
		assert.Equal(t, 1, content.focusedTimes)
		assert.Equal(t, 1, content.unfocusedTimes)
	}, true)
}

func Test_canvas_InteractiveArea(t *testing.T) {
	dev := &device{
		safeTop:    17,
		safeLeft:   42,
		safeBottom: 71,
		safeRight:  24,
	}
	scale := dev.SystemScaleForWindow(nil)

	c := newCanvas(dev)
	c.SetContent(fynecanvas.NewRectangle(color.Black))
	canvasSize := fyne.NewSize(600, 800)
	c.(*canvas).Resize(canvasSize)

	t.Run("for canvas with size", func(t *testing.T) {
		pos, size := c.InteractiveArea()
		assert.Equal(t, fyne.NewPos(float32(dev.safeLeft)/scale, float32(dev.safeTop)/scale), pos)
		assert.Equal(t, canvasSize.SubtractWidthHeight(float32(dev.safeLeft+dev.safeRight)/scale, float32(dev.safeTop+dev.safeBottom)/scale), size)
	})

	t.Run("when canvas size has changed", func(t *testing.T) {
		changedCanvasSize := fyne.NewSize(800, 600)
		c.(*canvas).Resize(changedCanvasSize)
		pos, size := c.InteractiveArea()
		assert.Equal(t, fyne.NewPos(float32(dev.safeLeft)/scale, float32(dev.safeTop)/scale), pos)
		assert.Equal(t, changedCanvasSize.SubtractWidthHeight(float32(dev.safeLeft+dev.safeRight)/scale, float32(dev.safeTop+dev.safeBottom)/scale), size)
	})

	t.Run("when canvas got hovering window head", func(t *testing.T) {
		c.(*canvas).Resize(canvasSize)
		hoveringMenu := fynecanvas.NewRectangle(color.Black)
		hoveringMenu.SetMinSize(fyne.NewSize(17, 17))
		windowHead := container.NewHBox(hoveringMenu) // a single object in window head is considered to be the hovering menu
		c.(*canvas).setWindowHead(windowHead)
		pos, size := c.InteractiveArea()
		assert.Equal(t, fyne.NewPos(float32(dev.safeLeft)/scale, float32(dev.safeTop)/scale), pos)
		assert.Equal(t, canvasSize.SubtractWidthHeight(float32(dev.safeLeft+dev.safeRight)/scale, float32(dev.safeTop+dev.safeBottom)/scale), size)
	})

	t.Run("when canvas got displacing window head", func(t *testing.T) {
		c.(*canvas).Resize(canvasSize)
		menu := fynecanvas.NewRectangle(color.Black)
		menu.SetMinSize(fyne.NewSize(17, 17))
		title := fynecanvas.NewRectangle(color.Black)
		title.SetMinSize(fyne.NewSize(42, 20))
		expectedOffset := 20 + 2*theme.Padding()
		windowHead := container.NewHBox(menu, title) // two or more objects in window head are considered to be the window title bar
		c.(*canvas).setWindowHead(windowHead)
		pos, size := c.InteractiveArea()
		assert.Equal(t, fyne.NewPos(float32(dev.safeLeft)/scale, float32(dev.safeTop)/scale+expectedOffset), pos)
		expectedWidth := canvasSize.Width - float32(dev.safeLeft+dev.safeRight)/scale
		expectedHeight := canvasSize.Height - float32(dev.safeTop+dev.safeBottom)/scale - expectedOffset
		assert.InDelta(t, expectedWidth, size.Width, 0.001)
		assert.InDelta(t, expectedHeight, size.Height, 0.001)
	})
}

func Test_canvas_PixelCoordinateAtPosition(t *testing.T) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)

	pos := fyne.NewPos(4, 4)
	c.scale = 2.5
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, 10, x)
	assert.Equal(t, 10, y)
}

func Test_canvas_ResizeWithModalPopUpOverlay(t *testing.T) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)

	c.SetContent(widget.NewLabel("Content"))

	popup := widget.NewModalPopUp(widget.NewLabel("PopUp"), c)
	popupBgSize := fyne.NewSize(200, 200)
	popup.Show()
	popup.Resize(popupBgSize)

	canvasSize := fyne.NewSize(600, 700)
	c.Resize(canvasSize)

	// get popup content padding dynamically
	popupContentPadding := popup.MinSize().Subtract(popup.Content.MinSize())

	// TODO modal non-zero again
	assert.Equal(t, popupBgSize.Subtract(popupContentPadding), popup.Content.Size())
	assert.Equal(t, canvasSize, c.Overlays().Top().Size())
}

func Test_canvas_Tappable(t *testing.T) {
	content := &touchableLabel{Label: widget.NewLabel("Hi\nHi\nHi"), ids: make(map[int]bool)}
	content.ExtendBaseWidget(content)
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(content)
	c.Resize(fyne.NewSize(36, 24))
	content.Resize(fyne.NewSize(24, 24))

	c.tapDown(fyne.NewPos(15, 15), 0)
	assert.True(t, content.down)

	c.tapUp(fyne.NewPos(15, 15), 0, func(fyne.Tappable, *fyne.PointEvent) {
	}, func(fyne.SecondaryTappable, *fyne.PointEvent) {
	}, func(fyne.DoubleTappable, *fyne.PointEvent) {
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})
	assert.True(t, content.up)

	c.tapDown(fyne.NewPos(15, 15), 0)
	c.tapMove(fyne.NewPos(35, 15), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
	})
	assert.True(t, content.cancel)
}

func Test_canvas_TouchID(t *testing.T) {
	content := &touchableLabel{Label: widget.NewLabel("Hi\nHi\nHi"), ids: make(map[int]bool)}
	content.ExtendBaseWidget(content)
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(content)
	c.Resize(fyne.NewSize(36, 24))
	content.Resize(fyne.NewSize(24, 24))

	c.tapDown(fyne.NewPos(15, 15), 0)
	c.tapUp(fyne.NewPos(15, 15), 0, func(fyne.Tappable, *fyne.PointEvent) {
	}, func(fyne.SecondaryTappable, *fyne.PointEvent) {
	}, func(fyne.DoubleTappable, *fyne.PointEvent) {
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})
	assert.True(t, content.ids[0])

	c.tapDown(fyne.NewPos(15, 15), 1)
	c.tapUp(fyne.NewPos(15, 15), 1, func(fyne.Tappable, *fyne.PointEvent) {
	}, func(fyne.SecondaryTappable, *fyne.PointEvent) {
	}, func(fyne.DoubleTappable, *fyne.PointEvent) {
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})
	assert.True(t, content.ids[1])
}

func Test_canvas_Tapped(t *testing.T) {
	tapped := false
	altTapped := false
	buttonTap := false
	var pointEvent *fyne.PointEvent
	var tappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(button)
	c.Resize(fyne.NewSize(36, 24))
	button.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	c.tapUp(tapPos, 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		tappedObj = wid
		pointEvent = ev
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		altTapped = true
		wid.TappedSecondary(ev)
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})

	assert.True(t, tapped, "tap primary")
	assert.False(t, altTapped, "don't tap secondary")
	assert.True(t, buttonTap, "button should be tapped")
	assert.Equal(t, button, tappedObj)
	if assert.NotNil(t, pointEvent) {
		assert.Equal(t, fyne.NewPos(6, 6), pointEvent.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(3, 3), pointEvent.Position)
	}
}

func Test_canvas_TappedAndDoubleTapped(t *testing.T) {
	tapped := 0
	but := newDoubleTappableButton()
	but.OnTapped = func() {
		tapped = 1
	}
	but.onDoubleTap = func() {
		tapped = 2
	}

	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(container.NewStack(but))
	c.Resize(fyne.NewSize(36, 24))

	simulateTap(c)
	waitAndCheck(700, func() {
		assert.Equal(t, 1, tapped)
	})

	simulateTap(c)
	simulateTap(c)
	waitAndCheck(700, func() {
		assert.Equal(t, 2, tapped)
	})
}

func Test_canvas_TappedMulti(t *testing.T) {
	buttonTap := false
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(button)
	c.Resize(fyne.NewSize(36, 24))
	button.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	c.tapUp(tapPos, 1, func(wid fyne.Tappable, ev *fyne.PointEvent) { // different tapID
		wid.Tapped(ev)
	}, func(fyne.SecondaryTappable, *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})

	assert.False(t, buttonTap, "button should not be tapped")
}

func Test_canvas_TappedSecondary(t *testing.T) {
	var pointEvent *fyne.PointEvent
	var altTappedObj fyne.SecondaryTappable
	obj := &tappableLabel{}
	obj.ExtendBaseWidget(obj)
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(obj)
	c.Resize(fyne.NewSize(36, 24))
	obj.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	time.Sleep(310 * time.Millisecond)
	c.tapUp(tapPos, 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		obj.tap = true
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		obj.altTap = true
		altTappedObj = wid
		pointEvent = ev
		wid.TappedSecondary(ev)
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})

	assert.False(t, obj.tap, "don't tap primary")
	assert.True(t, obj.altTap, "tap secondary")
	assert.Equal(t, obj, altTappedObj)
	if assert.NotNil(t, pointEvent) {
		assert.Equal(t, fyne.NewPos(6, 6), pointEvent.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(3, 3), pointEvent.Position)
	}
}

func Test_canvas_TapSurvivesDialogRelayout(t *testing.T) {
	t.Run("growth keeps the pressed button", func(t *testing.T) {
		var topTaps, bottomTaps int
		top := widget.NewButton("Top", func() { topTaps++ })
		bottom := widget.NewButton("Bottom", func() { bottomTaps++ })
		c, entry := newModalTapCanvas(top, bottom)
		roots := overlayRoots(c)
		tapPos := pressPos(bottom, roots)
		local := tapPos.Subtract(objectOrigin(bottom, roots))
		deltaY := objectOrigin(bottom, roots).Y - objectOrigin(top, roots).Y
		if !assert.Greater(t, deltaY, float32(8)) || !assert.True(t, objectContains(bottom, roots, tapPos)) {
			return
		}
		assert.False(t, objectContains(top, roots, tapPos))
		assert.Equal(t, entry, c.Focused())
		assertOverlayCentered(t, c)

		pop := c.Overlays().Top().(*intwidget.OverlayContainer).Content
		before := pop.Position()
		resized := resizeOnUnfocus(c, fyne.NewSize(360, 640+deltaY*2))
		c.tapDown(tapPos, 0)
		assert.True(t, resized())
		assert.Nil(t, c.Focused())
		assert.True(t, objectContains(top, roots, tapPos), "button under the finger changed")
		assert.False(t, objectContains(bottom, roots, tapPos))
		assert.Greater(t, pop.Position().Y, before.Y)
		assertOverlayCentered(t, c)

		got, ev := dispatchTap(c, tapPos, 0)
		assert.Equal(t, bottom, got)
		assert.Equal(t, 1, bottomTaps)
		assert.Zero(t, topTaps)
		if assert.NotNil(t, ev) {
			assert.Equal(t, local, ev.Position)
			assert.Equal(t, tapPos, ev.AbsolutePosition)
		}
		assertNoTapState(t, c, 0)
	})

	t.Run("shrink keeps the pressed button", func(t *testing.T) {
		var topTaps, bottomTaps int
		top := widget.NewButton("Top", func() { topTaps++ })
		bottom := widget.NewButton("Bottom", func() { bottomTaps++ })
		c, _ := newModalTapCanvas(top, bottom)
		roots := overlayRoots(c)
		tapPos := pressPos(bottom, roots)
		deltaY := objectOrigin(bottom, roots).Y - objectOrigin(top, roots).Y
		shrunkH := float32(640) - deltaY*2
		if !assert.Greater(t, deltaY, float32(8)) || !assert.Greater(t, shrunkH, float32(200)) {
			return
		}

		resized := resizeOnUnfocus(c, fyne.NewSize(360, shrunkH))
		c.tapDown(tapPos, 0)
		assert.True(t, resized())
		assert.False(t, objectContains(bottom, roots, tapPos), "pressed button moved off the finger")
		assert.False(t, objectContains(top, roots, tapPos))
		assertOverlayCentered(t, c)

		got, _ := dispatchTap(c, tapPos, 0)
		assert.Equal(t, bottom, got)
		assert.Equal(t, 1, bottomTaps)
		assert.Zero(t, topTaps)
	})

	t.Run("touchable keeps down up pairing", func(t *testing.T) {
		top := newRecordingTouch("Top")
		bottom := newRecordingTouch("Bottom")
		c, _ := newModalTapCanvas(top, bottom)
		roots := overlayRoots(c)
		tapPos := pressPos(bottom, roots)
		local := tapPos.Subtract(objectOrigin(bottom, roots))
		deltaY := objectOrigin(bottom, roots).Y - objectOrigin(top, roots).Y
		if !assert.Greater(t, deltaY, float32(8)) || !assert.True(t, objectContains(bottom, roots, tapPos)) {
			return
		}

		resized := resizeOnUnfocus(c, fyne.NewSize(360, 640+deltaY*2))
		c.tapDown(tapPos, 3)
		assert.True(t, resized())
		assert.True(t, objectContains(top, roots, tapPos))
		assert.False(t, objectContains(bottom, roots, tapPos))

		got, ev := dispatchTap(c, tapPos, 3)
		assert.Equal(t, bottom, got)
		assert.Nil(t, top.down)
		assert.Nil(t, top.up)
		if assert.NotNil(t, bottom.down) && assert.NotNil(t, bottom.up) {
			assert.Equal(t, 3, bottom.down.ID)
			assert.Equal(t, local, bottom.down.Position)
			assert.Equal(t, tapPos, bottom.down.AbsolutePosition)
			assert.Equal(t, 3, bottom.up.ID)
			assert.Equal(t, local, bottom.up.Position)
			assert.Equal(t, tapPos, bottom.up.AbsolutePosition)
		}
		if assert.NotNil(t, ev) {
			assert.Equal(t, local, ev.Position)
			assert.Equal(t, tapPos, ev.AbsolutePosition)
		}
		assert.Zero(t, bottom.cancels)
	})
}

func Test_canvas_TapCaptureDoesNotRetarget(t *testing.T) {
	t.Run("move away then tap again", func(t *testing.T) {
		var aTaps, bTaps int
		a := widget.NewButton("A", func() { aTaps++ })
		b := widget.NewButton("B", func() { bTaps++ })
		a.Resize(fyne.NewSize(40, 40))
		b.Resize(fyne.NewSize(40, 40))
		a.Move(fyne.NewPos(10, 10))
		b.Move(fyne.NewPos(100, 10))
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(container.NewWithoutLayout(a, b))
		c.Resize(fyne.NewSize(220, 120))
		roots := []fyne.CanvasObject{c.content}
		tapA := objectOrigin(a, roots).AddXY(20, 20)
		tapB := objectOrigin(b, roots).AddXY(20, 20)

		c.tapDown(tapA, 0)
		c.tapMove(tapB, 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
			wid.Dragged(ev)
		})
		got, _ := dispatchTap(c, tapB, 0)
		assert.Equal(t, b, got)
		assert.Zero(t, aTaps)
		assert.Equal(t, 1, bTaps)

		c.tapDown(tapA, 0)
		got, _ = dispatchTap(c, tapA, 0)
		assert.Equal(t, a, got)
		assert.Equal(t, 1, aTaps)
		assertNoTapState(t, c, 0)
	})

	t.Run("drag does not fire captured button", func(t *testing.T) {
		taps := 0
		button := widget.NewButton("A", func() { taps++ })
		scroll := container.NewVScroll(container.NewVBox(button, widget.NewLabel("1\n2\n3\n4\n5\n6\n7\n8\n9\n10")))
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(scroll)
		c.Resize(fyne.NewSize(120, 80))
		roots := []fyne.CanvasObject{c.content}
		pos := pressPos(button, roots)

		c.tapDown(pos, 0)
		c.tapMove(pos.AddXY(0, 70), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
			wid.Dragged(ev)
		})
		dispatchTap(c, pos.AddXY(0, 70), 0)
		assert.Zero(t, taps)
		assertNoTapState(t, c, 0)

		scroll.ScrollToTop()
		pos = pressPos(button, roots)
		c.tapDown(pos, 1)
		got, _ := dispatchTap(c, pos, 1)
		assert.Equal(t, button, got)
		assert.Equal(t, 1, taps)
	})

	t.Run("hidden pressed widget", func(t *testing.T) {
		var aTaps, bTaps int
		a := widget.NewButton("A", func() { aTaps++ })
		b := widget.NewButton("B", func() { bTaps++ })
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(container.NewStack(b, a))
		c.Resize(fyne.NewSize(120, 80))
		roots := []fyne.CanvasObject{c.content}
		pos := objectOrigin(a, roots).AddXY(a.Size().Width/2, a.Size().Height/2)

		c.tapDown(pos, 0)
		a.Hide()
		dispatchTap(c, pos, 0)
		assert.Zero(t, aTaps)
		assert.Zero(t, bTaps)
		assert.False(t, a.Visible())
		assert.True(t, objectContains(b, roots, pos))

		c.tapDown(pos, 0)
		got, _ := dispatchTap(c, pos, 0)
		assert.Equal(t, b, got)
		assert.Equal(t, 1, bTaps)
	})

	t.Run("removed pressed widget", func(t *testing.T) {
		var aTaps, bTaps int
		a := widget.NewButton("A", func() { aTaps++ })
		b := widget.NewButton("B", func() { bTaps++ })
		stack := container.NewStack(b, a)
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(stack)
		c.Resize(fyne.NewSize(120, 80))
		roots := []fyne.CanvasObject{c.content}
		pos := objectOrigin(a, roots).AddXY(a.Size().Width/2, a.Size().Height/2)

		c.tapDown(pos, 0)
		stack.Remove(a)
		dispatchTap(c, pos, 0)
		assert.Zero(t, aTaps)
		assert.Zero(t, bTaps)

		c.tapDown(pos, 0)
		got, _ := dispatchTap(c, pos, 0)
		assert.Equal(t, b, got)
		assert.Equal(t, 1, bTaps)
	})

	t.Run("replaced overlay", func(t *testing.T) {
		var aTaps, underTaps int
		under := widget.NewButton("Under", func() { underTaps++ })
		a := widget.NewButton("A", func() { aTaps++ })
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(under)
		c.Resize(fyne.NewSize(300, 400))
		popup := widget.NewModalPopUp(a, c)
		popup.Show()
		roots := overlayRoots(c)
		pos := objectOrigin(a, roots).AddXY(a.Size().Width/2, a.Size().Height/2)
		if !assert.True(t, objectContains(a, roots, pos)) {
			return
		}

		c.tapDown(pos, 0)
		popup.Hide()
		dispatchTap(c, pos, 0)
		assert.Zero(t, aTaps)
		assert.Zero(t, underTaps)
		assert.Nil(t, c.Overlays().Top())

		c.tapDown(pos, 0)
		got, _ := dispatchTap(c, pos, 0)
		assert.Equal(t, under, got)
		assert.Equal(t, 1, underTaps)
	})

	t.Run("content sliding under an empty press", func(t *testing.T) {
		taps := 0
		button := widget.NewButton("A", func() { taps++ })
		button.Resize(fyne.NewSize(40, 40))
		button.Move(fyne.NewPos(10, 10))
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(container.NewWithoutLayout(button))
		c.Resize(fyne.NewSize(200, 200))
		roots := []fyne.CanvasObject{c.content}
		pos := fyne.NewPos(150, 150)
		assert.False(t, objectContains(button, roots, pos))

		c.tapDown(pos, 0)
		button.Move(pos.Subtract(objectOrigin(c.content, roots)))
		assert.True(t, objectContains(button, roots, pos))
		dispatchTap(c, pos, 0)
		assert.Zero(t, taps)
		assertNoTapState(t, c, 0)
	})

	t.Run("unmatched reuse and interleaved ids", func(t *testing.T) {
		label := &tappableLabel{}
		label.ExtendBaseWidget(label)
		c := newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(label)
		c.Resize(fyne.NewSize(80, 80))
		pos := fyne.NewPos(20, 20)

		dispatchTap(c, pos, 9)
		assert.False(t, label.tap)
		assert.False(t, label.altTap)

		c.tapDown(pos, 4)
		dispatchTap(c, pos, 4)
		assert.True(t, label.tap)
		assert.False(t, label.altTap)
		label.tap = false
		dispatchTap(c, pos, 4)
		assert.False(t, label.tap)
		assert.False(t, label.altTap)
		assertNoTapState(t, c, 4)

		c.tapDown(pos, 4)
		dispatchTap(c, pos, 4)
		assert.True(t, label.tap)
		assertNoTapState(t, c, 4)

		var aTaps, bTaps int
		a := widget.NewButton("A", func() { aTaps++ })
		b := widget.NewButton("B", func() { bTaps++ })
		a.Resize(fyne.NewSize(40, 40))
		b.Resize(fyne.NewSize(40, 40))
		a.Move(fyne.NewPos(10, 10))
		b.Move(fyne.NewPos(100, 10))
		c = newCanvas(fyne.CurrentDevice()).(*canvas)
		c.SetContent(container.NewWithoutLayout(a, b))
		c.Resize(fyne.NewSize(220, 120))
		roots := []fyne.CanvasObject{c.content}
		tapA := objectOrigin(a, roots).AddXY(20, 20)
		tapB := objectOrigin(b, roots).AddXY(20, 20)

		c.tapDown(tapA, 0)
		c.tapDown(tapB, 1)
		got, _ := dispatchTap(c, tapB, 1)
		assert.Equal(t, b, got)
		got, _ = dispatchTap(c, tapA, 0)
		assert.Equal(t, a, got)
		assert.Equal(t, 1, aTaps)
		assert.Equal(t, 1, bTaps)
		dispatchTap(c, tapA, 2)
		assert.Equal(t, 1, aTaps)
		assert.Equal(t, 1, bTaps)
		assertNoTapState(t, c, 0)
		assertNoTapState(t, c, 1)
		assert.Empty(t, c.tapCaptured)
		assert.Empty(t, c.touched)
	})
}

type recordingTouch struct {
	*widget.Label
	down, up *mobile.TouchEvent
	cancels  int
}

func newRecordingTouch(text string) *recordingTouch {
	r := &recordingTouch{Label: widget.NewLabel(text)}
	r.ExtendBaseWidget(r)
	return r
}

func (r *recordingTouch) Tapped(*fyne.PointEvent) {
}

func (r *recordingTouch) TouchDown(ev *mobile.TouchEvent) {
	cp := *ev
	r.down = &cp
}

func (r *recordingTouch) TouchUp(ev *mobile.TouchEvent) {
	cp := *ev
	r.up = &cp
}

func (r *recordingTouch) TouchCancel(*mobile.TouchEvent) {
	r.cancels++
}

func newModalTapCanvas(top, bottom fyne.CanvasObject) (*canvas, *widget.Entry) {
	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.OnFocus = nil
	c.SetContent(widget.NewLabel("content"))
	c.Resize(fyne.NewSize(360, 640))

	entry := widget.NewEntry()
	widget.NewModalPopUp(container.NewVBox(entry, top, bottom), c).Show()
	c.Focus(entry)
	return c, entry
}

func resizeOnUnfocus(c *canvas, size fyne.Size) func() bool {
	done := false
	c.OnUnfocus = func() {
		if done {
			return
		}
		done = true
		c.Resize(size)
	}
	return func() bool { return done }
}

func overlayRoots(c *canvas) []fyne.CanvasObject {
	return []fyne.CanvasObject{c.Overlays().Top()}
}

func objectOrigin(obj fyne.CanvasObject, roots []fyne.CanvasObject) fyne.Position {
	return intdriver.AbsolutePositionForObject(obj, roots)
}

func objectContains(obj fyne.CanvasObject, roots []fyne.CanvasObject, p fyne.Position) bool {
	origin := objectOrigin(obj, roots)
	size := obj.Size()
	return p.X >= origin.X && p.Y >= origin.Y && p.X < origin.X+size.Width && p.Y < origin.Y+size.Height
}

func pressPos(obj fyne.CanvasObject, roots []fyne.CanvasObject) fyne.Position {
	return objectOrigin(obj, roots).AddXY(obj.Size().Width/2, 2)
}

func assertOverlayCentered(t *testing.T, c *canvas) {
	t.Helper()
	ov, ok := c.Overlays().Top().(*intwidget.OverlayContainer)
	if !assert.True(t, ok) {
		return
	}
	pop := ov.Content
	assert.InDelta(t, float64((ov.Size().Width-pop.Size().Width)/2), float64(pop.Position().X), 0.01)
	assert.InDelta(t, float64((ov.Size().Height-pop.Size().Height)/2), float64(pop.Position().Y), 0.01)
}

func assertNoTapState(t *testing.T, c *canvas, id int) {
	t.Helper()
	_, ok := c.tapCaptured[id]
	assert.False(t, ok)
	_, ok = c.touched[id]
	assert.False(t, ok)
	_, ok = c.lastTapDown[id]
	assert.False(t, ok)
}

func dispatchTap(c *canvas, pos fyne.Position, id int) (fyne.CanvasObject, *fyne.PointEvent) {
	var obj fyne.CanvasObject
	var ev *fyne.PointEvent
	c.tapUp(pos, id, func(wid fyne.Tappable, e *fyne.PointEvent) {
		obj, _ = wid.(fyne.CanvasObject)
		cp := *e
		ev = &cp
		wid.Tapped(e)
	}, func(wid fyne.SecondaryTappable, e *fyne.PointEvent) {
		wid.TappedSecondary(e)
	}, func(wid fyne.DoubleTappable, e *fyne.PointEvent) {
		wid.DoubleTapped(e)
	}, func(wid fyne.Draggable, _ *fyne.DragEvent) {
		wid.DragEnd()
	})
	return obj, ev
}

type touchableLabel struct {
	*widget.Label
	down, up, cancel bool
	ids              map[int]bool
}

func (t *touchableLabel) TouchDown(ev *mobile.TouchEvent) {
	t.down = true
	t.ids[ev.ID] = true
}

func (t *touchableLabel) TouchUp(ev *mobile.TouchEvent) {
	t.up = true
	t.ids[ev.ID] = true
}

func (t *touchableLabel) TouchCancel(ev *mobile.TouchEvent) {
	t.cancel = true
	t.ids[ev.ID] = true
}

type tappableLabel struct {
	widget.Label
	tap, altTap bool
}

func (t *tappableLabel) Tapped(_ *fyne.PointEvent) {
	t.tap = true
}

func (t *tappableLabel) TappedSecondary(_ *fyne.PointEvent) {
	t.altTap = true
}

type focusableEntry struct {
	widget.Entry
	c fyne.Canvas

	focusedTimes   int
	unfocusedTimes int
}

func newFocusableEntry(c fyne.Canvas) *focusableEntry {
	entry := &focusableEntry{c: c}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (f *focusableEntry) FocusGained() {
	f.focusedTimes++
	f.Entry.FocusGained()
}

func (f *focusableEntry) FocusLost() {
	f.unfocusedTimes++
	f.Entry.FocusLost()
}

func (f *focusableEntry) Tapped(*fyne.PointEvent) {
	f.c.Focus(f)
}

type doubleTappableButton struct {
	widget.Button

	onDoubleTap func()
}

func (t *doubleTappableButton) DoubleTapped(_ *fyne.PointEvent) {
	t.onDoubleTap()
}

func newDoubleTappableButton() *doubleTappableButton {
	but := &doubleTappableButton{}
	but.ExtendBaseWidget(but)

	return but
}

func simulateTap(c *canvas) {
	c.tapDown(fyne.NewPos(15, 15), 0)
	time.Sleep(50 * time.Millisecond)
	c.tapUp(fyne.NewPos(15, 15), 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		wid.Tapped(ev)
	}, func(fyne.SecondaryTappable, *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(fyne.Draggable, *fyne.DragEvent) {
	})
}

type mobileApp struct {
	fyne.App
	driver fyne.Driver
}

func (a *mobileApp) Driver() fyne.Driver {
	return a.driver
}

func (*mobileApp) Run() {
	// This is an incomplete driver loop - our CI does not currently support booting the mobile graphics
	// TODO replace with a full mobileApp.Run() once that is resolved
	async.SetMainGoroutine()

	for fn := range d.queuedFuncs.Out() {
		fn()
	}
}

func newTestMobileApp() fyne.App {
	return &mobileApp{
		App:    fyne.CurrentApp(),
		driver: NewGoMobileDriver(),
	}
}

func waitAndCheck(msWait time.Duration, fn func()) {
	waitForCheck := common.DonePool.Get()

	go func() {
		defer common.DonePool.Put(waitForCheck)

		time.Sleep(msWait * time.Millisecond)
		d.DoFromGoroutine(func() {
			fn()

			waitForCheck <- struct{}{}
		}, true)
	}()
	<-waitForCheck
}
