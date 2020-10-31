// !build test

package theme

import (
	"net/url"

	"fyne.io/fyne"
)

type themedApp struct {
	theme fyne.Theme
}

func (t *themedApp) BuildType() fyne.BuildType {
	return fyne.BuildStandard
}

func (t *themedApp) NewWindow(title string) fyne.Window {
	return nil
}

func (t *themedApp) OpenURL(url *url.URL) error {
	return nil
}

func (t *themedApp) Icon() fyne.Resource {
	return nil
}

func (t *themedApp) SetIcon(fyne.Resource) {
}

func (t *themedApp) Run() {
}

func (t *themedApp) Quit() {
}

func (t *themedApp) Driver() fyne.Driver {
	return nil
}

func (t *themedApp) UniqueID() string {
	return ""
}

func (t *themedApp) SendNotification(notification *fyne.Notification) {
}

func (t *themedApp) Settings() fyne.Settings {
	return t
}

func (t *themedApp) Storage() fyne.Storage {
	return nil
}

func (t *themedApp) Preferences() fyne.Preferences {
	return nil
}

func (t *themedApp) PrimaryColor() string {
	return ColorBlue
}

func (t *themedApp) Theme() fyne.Theme {
	return t.theme
}

func (t *themedApp) SetTheme(theme fyne.Theme) {
	t.theme = theme
}

func (t *themedApp) Scale() float32 {
	return 1.0
}

func (t *themedApp) AddChangeListener(chan fyne.Settings) {
}
