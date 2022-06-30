package app

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"golang.org/x/sys/execabs"
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

	meta.ID = "fakedInject"
	app = New()

	assert.Equal(t, meta.ID, app.UniqueID())
}

func TestFyneApp_OpenURL(t *testing.T) {
	opened := ""
	app := NewWithID("io.fyne.test")
	app.(*fyneApp).exec = func(cmd string, arg ...string) *execabs.Cmd {
		opened = arg[len(arg)-1]
		return execabs.Command("")
	}

	urlStr := "https://fyne.io"
	u, _ := url.Parse(urlStr)
	err := app.OpenURL(u)

	if err != nil && strings.Contains(err.Error(), "unknown operating system") {
		return // when running in CI mode we don't actually open URLs...
	}

	assert.Equal(t, urlStr, opened)
}

func TestFyneApp_SetCloudProvider(t *testing.T) {
	a := NewWithID("io.fyne.test")
	p := &mockCloud{}
	a.SetCloudProvider(p)

	assert.Equal(t, p, a.CloudProvider())
	assert.True(t, p.configured)
}

func TestFyneApp_transitionCloud(t *testing.T) {
	a := NewWithID("io.fyne.test")
	p := &mockCloud{}
	preferenceChanged := false
	settingsChan := make(chan fyne.Settings)
	a.Preferences().AddChangeListener(func() {
		preferenceChanged = true
	})
	a.Settings().AddChangeListener(settingsChan)
	a.SetCloudProvider(p)

	<-settingsChan // settings were updated
	assert.True(t, preferenceChanged)
}

type mockCloud struct {
	configured bool
}

func (c *mockCloud) ProviderDescription() string {
	return "Mock cloud implementation"
}

func (c *mockCloud) ProviderIcon() fyne.Resource {
	return theme.FyneLogo()
}

func (c *mockCloud) ProviderName() string {
	return "mock"
}

func (c *mockCloud) Setup(fyne.App) error {
	c.configured = true
	return nil
}
