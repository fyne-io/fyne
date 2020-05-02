package glfw

import (
	"math"
	"os"
	"strconv"

	"fyne.io/fyne"
)

const (
	baselineDPI = 120.0
	scaleEnvKey = "FYNE_SCALE"
)

func calculateDetectedScale(widthMm, widthPx int) float32 {
	dpi := float32(widthPx) / (float32(widthMm) / 25.4)
	if dpi > 1000 || dpi < 10 {
		dpi = baselineDPI
	}

	scale := float32(float64(dpi) / baselineDPI)
	if scale < 1.0 {
		return 1.0
	}
	return scale
}

func calculateScale(user, system, detected float32) float32 {
	if user == fyne.SettingsScaleAuto {
		user = 1.0
	}

	if system == fyne.SettingsScaleAuto {
		system = detected
	}

	raw := system * user
	return float32(math.Round(float64(raw*10.0))) / 10.0
}

func userScale() float32 {
	env := os.Getenv(scaleEnvKey)

	if env != "" && env != "auto" {
		scale, err := strconv.ParseFloat(env, 32)
		if err == nil && scale != 0 {
			return float32(scale)
		}
		fyne.LogError("Error reading scale", err)
	}

	if env != "auto" {
		setting := fyne.CurrentApp().Settings().Scale()
		if setting != fyne.SettingsScaleAuto && setting != 0.0 {
			return setting
		}
	}

	return 1.0 // user preference for auto is now passed as 1 so the system auto is picked up
}
