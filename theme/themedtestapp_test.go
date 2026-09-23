package theme_test

import (
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type themedApp struct {
	primaryColor string
	theme        fyne.Theme
	variant      fyne.ThemeVariant
}

func (*themedApp) CloudProvider() fyne.CloudProvider {
	return nil
}

func (*themedApp) BuildType() fyne.BuildType {
	return fyne.BuildStandard
}

func (*themedApp) NewWindow(_ string) fyne.Window {
	return nil
}

func (*themedApp) OpenURL(_ *url.URL) error {
	return nil
}

func (*themedApp) Icon() fyne.Resource {
	return nil
}

func (*themedApp) SetIcon(fyne.Resource) {
}

func (*themedApp) Run() {
}

func (*themedApp) Quit() {
}

func (*themedApp) Driver() fyne.Driver {
	return nil
}

func (*themedApp) UniqueID() string {
	return ""
}

func (*themedApp) SendNotification(_ *fyne.Notification) {
}

func (*themedApp) ScheduleNotification(_ *fyne.Notification, _ time.Time) (*fyne.ScheduledNotification, error) {
	return nil, nil
}

func (*themedApp) CancelScheduledNotification(_ string) error {
	return nil
}

func (a *themedApp) Settings() fyne.Settings {
	return a
}

func (*themedApp) Storage() fyne.Storage {
	return nil
}

func (*themedApp) Preferences() fyne.Preferences {
	return nil
}

func (*themedApp) Lifecycle() fyne.Lifecycle {
	return nil
}

func (*themedApp) Metadata() fyne.AppMetadata {
	return fyne.AppMetadata{}
}

func (a *themedApp) PrimaryColor() string {
	if a.primaryColor != "" {
		return a.primaryColor
	}

	return theme.ColorBlue
}

func (a *themedApp) Theme() fyne.Theme {
	return a.theme
}

func (a *themedApp) SetTheme(t fyne.Theme) {
	a.theme = t
}

func (a *themedApp) ThemeVariant() fyne.ThemeVariant {
	return a.variant // The null value is theme.VariantDark
}

func (*themedApp) SetCloudProvider(fyne.CloudProvider) {
}

func (*themedApp) Scale() float32 {
	return 1.0
}

func (*themedApp) ShowAnimations() bool {
	return true
}

func (*themedApp) AddChangeListener(chan fyne.Settings) {
}

func (*themedApp) AddListener(func(fyne.Settings)) {
}

func (*themedApp) Cache() fyne.Cache {
	return nil
}

func (*themedApp) Clipboard() fyne.Clipboard {
	return nil
}
