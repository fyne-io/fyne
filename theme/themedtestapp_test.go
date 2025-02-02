package theme_test

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type themedApp struct {
	primaryColor string
	theme        fyne.Theme
	variant      fyne.ThemeVariant
}

func (t *themedApp) CloudProvider() fyne.CloudProvider {
	return nil
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

func (t *themedApp) Lifecycle() fyne.Lifecycle {
	return nil
}

func (t *themedApp) Metadata() fyne.AppMetadata {
	return fyne.AppMetadata{}
}

func (t *themedApp) PrimaryColor() string {
	if t.primaryColor != "" {
		return t.primaryColor
	}

	return theme.ColorBlue
}

func (t *themedApp) Theme() fyne.Theme {
	return t.theme
}

func (t *themedApp) SetTheme(theme fyne.Theme) {
	t.theme = theme
}

func (t *themedApp) ThemeVariant() fyne.ThemeVariant {
	return t.variant // The null value is theme.VariantDark
}

func (t *themedApp) SetCloudProvider(fyne.CloudProvider) {
}

func (t *themedApp) Scale() float32 {
	return 1.0
}

func (t *themedApp) ShowAnimations() bool {
	return true
}

func (t *themedApp) AddChangeListener(chan fyne.Settings) {
}

func (t *themedApp) AddListener(func(fyne.Settings)) {
}

func (t *themedApp) Clipboard() fyne.Clipboard {
	return nil
}
