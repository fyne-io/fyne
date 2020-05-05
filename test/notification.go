package test

import (
	"testing"

	"fyne.io/fyne"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertNotificationSent allows an app developer to assert that a notification was sent.
// After the content of f has executed this utility will check that the specified notification was sent.
func AssertNotificationSent(t *testing.T, n *fyne.Notification, f func()) {
	if f == nil {
		return
	}
	require.IsType(t, &testApp{}, fyne.CurrentApp())
	a := fyne.CurrentApp().(*testApp)
	a.lastNotification = nil

	f()
	if n == nil {
		assert.Nil(t, a.lastNotification)
		return
	}

	assert.NotNil(t, a.lastNotification)
	assert.Equal(t, n.Title, a.lastNotification.Title)
	assert.Equal(t, n.Content, a.lastNotification.Content)
}
