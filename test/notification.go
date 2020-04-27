package test

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

// ExpectNotification allows an app developer to assert that a notification was sent.
// After the content of f has executed this utility will check that the specified notification was sent.
func ExpectNotification(t *testing.T, n *fyne.Notification, f func(a fyne.App)) {
	a := NewApp().(*testApp)
	fyne.SetCurrentApp(a)

	f(a)
	if n == nil {
		return
	}

	assert.NotNil(t, a.lastNotification)
	assert.Equal(t, n.Title, a.lastNotification.Title)
	assert.Equal(t, n.Content, a.lastNotification.Content)
}
