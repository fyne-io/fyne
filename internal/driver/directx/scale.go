//go:build windows && directx

// Scale handling: the user preference and the per-monitor DPI that combine into
// a canvas scale.

package directx

import (
	"math"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
)

const (
	fallbackWidth  = 800
	fallbackHeight = 600
	scaleEnvKey    = "FYNE_SCALE"
)

// userScale reads the user's scale preference: the FYNE_SCALE environment
// variable when set to anything but "auto", otherwise the app's scale setting.
func userScale() float32 {
	env := os.Getenv(scaleEnvKey)
	if env != "" && env != "auto" {
		s, err := strconv.ParseFloat(env, 32)
		if err == nil && s != 0 {
			return float32(s)
		}
		fyne.LogError("Error reading scale", err)
	}
	if env != "auto" {
		if setting := fyne.CurrentApp().Settings().Scale(); setting > 0 {
			return setting
		}
	}
	return 1.0
}

// calculatedScale combines the monitor DPI with the user's scale preference.
func (w *window) calculatedScale() float32 {
	dpi := float32(getDpiForWindow(w.hwnd))
	system := dpi / 96.0
	raw := system * userScale()
	return float32(math.Round(float64(raw*10.0))) / 10.0
}
