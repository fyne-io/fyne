package gomobile

import (
	"fyne.io/fyne/driver/mobile"
	"github.com/fyne-io/mobile/event/size"

	"fyne.io/fyne"
)

type device struct {
	safeTop, safeLeft, safeWidth, safeHeight int
}

//lint:file-ignore U1000 Var currentDPI is used in other files, but not here
var (
	currentOrientation size.Orientation
	currentDPI         float32
)

// Declare conformity with Device
var _ fyne.Device = (*device)(nil)

func (*device) ScreenInteractiveArea() (fyne.Position, fyne.Size) {
	scale := fyne.CurrentDevice().SystemScaleForWindow(nil) // we don't need a window parameter on mobile

	dev, ok := fyne.CurrentDevice().(*device)
	if !ok {
		wins := fyne.CurrentApp().Driver().AllWindows()
		size := fyne.Size{}
		if len(wins) > 0 {
			size = wins[0].Canvas().Size()
		}
		return fyne.NewPos(0, 0), size // running in test mode
	}

	return fyne.NewPos(int(float32(dev.safeLeft)/scale), int(float32(dev.safeTop)/scale)),
		fyne.NewSize(int(float32(dev.safeWidth)/scale), int(float32(dev.safeHeight)/scale))
}

func (*device) Orientation() fyne.DeviceOrientation {
	switch currentOrientation {
	case size.OrientationLandscape:
		return fyne.OrientationHorizontalLeft
	default:
		return fyne.OrientationVertical
	}
}

func (*device) IsMobile() bool {
	return true
}

func (*device) HasKeyboard() bool {
	return false
}

func (d *device) SystemScale() float32 {
	return d.SystemScaleForWindow(nil)
}

func (*device) ShowVirtualKeyboard() {
	showVirtualKeyboard(mobile.DefaultKeyboard)
}

func (*device) ShowVirtualKeyboardType(keyboard mobile.KeyboardType) {
	showVirtualKeyboard(keyboard)
}

func (*device) HideVirtualKeyboard() {
	hideVirtualKeyboard()
}
