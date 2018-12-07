package fyne

import "os"
import "time"

import "testing"
import "github.com/stretchr/testify/assert"

func TestThemeDefault(t *testing.T) {
	assert.Equal(t, "dark", GlobalSettings().Theme())
}

func TestThemeFromEnv(t *testing.T) {
	settingsCache = nil
	os.Setenv("FYNE_THEME", "light")

	assert.Equal(t, "light", GlobalSettings().Theme())
}

func TestSetTheme(t *testing.T) {
	GlobalSettings().SetTheme("light")

	assert.Equal(t, "light", GlobalSettings().Theme())
}

func TestListenerCallback(t *testing.T) {
	listener := make(chan Settings)
	GlobalSettings().AddChangeListener(listener)
	GlobalSettings().SetTheme("light")

	func() {
		select {
		case <-listener:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for callback")
		}
	}()
}
