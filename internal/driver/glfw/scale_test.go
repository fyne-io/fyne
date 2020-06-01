// +build !ci
// +build !mobile
// +build !windows

package glfw

import (
	"os"
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDetectedScale(t *testing.T) {
	lowDPI := calculateDetectedScale(300, 1600)
	assert.Equal(t, float32(1.1288888), lowDPI)

	hiDPI := calculateDetectedScale(420, 3800)
	assert.Equal(t, float32(1.9150794), hiDPI)
}

func TestCalculateDetectedScale_Min(t *testing.T) {
	lowDPI := calculateDetectedScale(300, 1280)
	assert.Equal(t, float32(1), lowDPI)
}

func TestCalculateScale(t *testing.T) {
	one := calculateScale(1.0, 1.0, 1.0)
	assert.Equal(t, float32(1.0), one)

	larger := calculateScale(1.5, 1.0, 1.0)
	assert.Equal(t, float32(1.5), larger)

	smaller := calculateScale(0.8, 1.0, 1.0)
	assert.Equal(t, float32(0.8), smaller)

	auto := calculateScale(fyne.SettingsScaleAuto, fyne.SettingsScaleAuto, 1.1)
	assert.Equal(t, float32(1.1), auto)

	autoUser := calculateScale(fyne.SettingsScaleAuto, 1.0, 1.1)
	assert.Equal(t, float32(1.0), autoUser)

	hiDPI := calculateScale(0.8, 2.0, 1.0)
	assert.Equal(t, float32(1.6), hiDPI)

	hiDPIAuto := calculateScale(0.8, fyne.SettingsScaleAuto, 2.0)
	assert.Equal(t, float32(1.6), hiDPIAuto)

	large := calculateScale(1.5, 2.0, 2.0)
	assert.Equal(t, float32(3.0), large)
}

func TestCalculateScale_Round(t *testing.T) {
	trunc := calculateScale(1.04321, 1.0, 1.0)
	assert.Equal(t, float32(1.0), trunc)

	round := calculateScale(1.1, 1.1, 1.0)
	assert.Equal(t, float32(1.2), round)
}

func TestUserScale(t *testing.T) {
	envVal := os.Getenv(scaleEnvKey)
	defer os.Setenv(scaleEnvKey, envVal)

	_ = os.Setenv(scaleEnvKey, "auto")
	set := fyne.CurrentApp().Settings().Scale()
	if set == float32(0.0) { // no config set
		assert.Equal(t, float32(1.0), userScale())
	} else {
		assert.Equal(t, set, userScale())
	}

	_ = os.Setenv(scaleEnvKey, "1.2")
	assert.Equal(t, float32(1.2), userScale())
}
