package fyne

import "os"
import "time"

import "testing"
import "github.com/stretchr/testify/assert"

func TestThemeDefault(t *testing.T) {
	assert.Equal(t, "dark", GetSettings().Theme())
}

func TestThemeFromEnv(t *testing.T) {
	settingsCache = nil
	os.Setenv("FYNE_THEME", "light")

	assert.Equal(t, "light", GetSettings().Theme())
}

func TestSetTheme(t *testing.T) {
	GetSettings().SetTheme("light")

	assert.Equal(t, "light", GetSettings().Theme())
}

func TestListenerCallback(t *testing.T) {
	listener := make(chan Settings)
	GetSettings().AddChangeListener(listener)
	GetSettings().SetTheme("light")

	func() {
		select {
		case <-listener:
		case <-time.After(1 * time.Second):
			assert.Fail(t, "Timed out waiting for callback")
		}
	}()
}
