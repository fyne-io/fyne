package fyne

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummyApp struct{}

func (dummyApp) CloudProvider() CloudProvider {
	return nil
}

func (dummyApp) NewWindow(_ string) Window {
	return nil
}

func (dummyApp) OpenURL(_ *url.URL) error {
	return nil
}

func (dummyApp) Icon() Resource {
	return nil
}

func (dummyApp) SetIcon(Resource) {
}

func (dummyApp) Run() {
}

func (dummyApp) Quit() {
}

func (dummyApp) Driver() Driver {
	return nil
}

func (dummyApp) UniqueID() string {
	return "dummy"
}

func (dummyApp) SendNotification(*Notification) {
}

func (dummyApp) ScheduleNotification(*Notification, time.Time) (*ScheduledNotification, error) {
	return nil, errors.New("unimplemented")
}

func (dummyApp) CancelScheduledNotification(string) error {
	return nil
}

func (dummyApp) SetCloudProvider(CloudProvider) {
}

func (dummyApp) Settings() Settings {
	return nil
}

func (dummyApp) Storage() Storage {
	return nil
}

func (dummyApp) Preferences() Preferences {
	return nil
}

func (dummyApp) Lifecycle() Lifecycle {
	return nil
}

func (dummyApp) Metadata() AppMetadata {
	return AppMetadata{}
}

func (dummyApp) Cache() Cache {
	return nil
}

func (dummyApp) Clipboard() Clipboard {
	return nil
}

func TestSetCurrentApp(t *testing.T) {
	a := &dummyApp{}
	SetCurrentApp(a)

	assert.Equal(t, a, CurrentApp())
}
