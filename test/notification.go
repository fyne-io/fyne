package test

import (
	"testing"

	"fyne.io/fyne/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertNotificationSent allows an app developer to assert that a notification was sent.
// After the content of f has executed this utility will check that the specified notification was sent.
func AssertNotificationSent(t *testing.T, n *fyne.Notification, f func()) {
	require.NotNil(t, f, "function has to be specified")
	require.IsType(t, &testApp{}, fyne.CurrentApp())
	a := fyne.CurrentApp().(*testApp)
	a.lastNotification = nil

	f()
	if n == nil {
		assert.Nil(t, a.lastNotification)
		return
	} else if a.lastNotification == nil {
		t.Error("No notification sent")
		return
	}

	assert.Equal(t, n.Title, a.lastNotification.Title)
	assert.Equal(t, n.Content, a.lastNotification.Content)
}
