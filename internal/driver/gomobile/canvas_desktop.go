// +build !ios,!android

package gomobile

import (
	"os"
	"strconv"

	"fyne.io/fyne"
)

func devicePadding() (topLeft, bottomRight fyne.Size) {
	return fyne.NewSize(0, 0), fyne.NewSize(0, 0)
}

func deviceScale() float32 {
	env := os.Getenv("FYNE_SCALE")
	if env != "" && env != "auto" {
		scale, err := strconv.ParseFloat(env, 32)
		if err == nil && scale != 0 {
			return float32(scale)
		}
		fyne.LogError("Error reading scale", err)
	}

	return 1
}
