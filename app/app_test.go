package app

import (
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	_ "fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestDummyApp(t *testing.T) {
	app := NewWithID("io.fyne.test")

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := NewWithID("io.fyne.test")

	assert.Equal(t, app, fyne.CurrentApp())
}

func TestFyneApp_UniqueID(t *testing.T) {
	appID := "io.fyne.test"
	app := NewWithID(appID)

	assert.Equal(t, appID, app.UniqueID())
}

func TestFyneApp_OpenURL(t *testing.T) {
	opened := ""
	app := NewWithID("io.fyne.test")
	app.(*fyneApp).exec = func(cmd string, arg ...string) *exec.Cmd {
		opened = arg[len(arg)-1]
		return exec.Command("")
	}

	urlStr := "https://fyne.io"
	u, _ := url.Parse(urlStr)
	err := app.OpenURL(u)

	if err != nil && strings.Contains(err.Error(), "unknown operating system") {
		return // when running in CI mode we don't actually open URLs...
	}

	// wait for the command to execute then check the URL we loaded
	for opened == "" {
	}
	assert.Equal(t, urlStr, opened)
}

func TestFyneApp_SendNotification(t *testing.T) {
	n := fyne.NewNotification("Test Title", "Some content")
	testApp := test.NewApp()
	oldApp := fyne.CurrentApp() // we are testing the app package so don't use the test app
	fyne.SetCurrentApp(testApp)

	test.AssertNotificationSent(t, n, func() {
		fyne.CurrentApp().SendNotification(n)
	})
	fyne.SetCurrentApp(oldApp)
}
