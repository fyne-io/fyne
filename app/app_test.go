package app

import (
	"net/url"
	"os/exec"
	"testing"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestDummyApp(t *testing.T) {
	app := New()

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := New()

	assert.Equal(t, app, fyne.CurrentApp())
}

func TestFyneApp_OpenURL(t *testing.T) {
	opened := ""
	app := New()
	app.(*fyneApp).exec = func(cmd string, arg ...string) *exec.Cmd {
		opened = arg[len(arg)-1]
		return exec.Command("")
	}

	urlStr := "https://fyne.io"
	u, _ := url.Parse(urlStr)
	_ = app.OpenURL(u)

	// wait for the command to execute then check the URL we loaded
	for opened == "" {
	}
	assert.Equal(t, urlStr, opened)
}
