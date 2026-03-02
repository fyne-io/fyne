package mobile

import (
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/widget"
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
		os.Exit(ret)
	}()

	close(waitForStart) // Signal that execution can continue.
	tester.Run()
}

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

type mobileApp struct {
	fyne.App
	driver fyne.Driver
}

func (a *mobileApp) Driver() fyne.Driver {
	return a.driver
}

func (a *mobileApp) Run() {
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
