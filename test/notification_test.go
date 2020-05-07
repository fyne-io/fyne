package test_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestAssertNotificationSent(t *testing.T) {
	n := fyne.NewNotification("Test Title", "Some content")
	myApp := fyne.CurrentApp()

	test.AssertNotificationSent(t, n, func() {
		fyne.CurrentApp().SendNotification(n)
	})
	assert.Equal(t, myApp, fyne.CurrentApp())
}

func TestAssertNotificationSent_Nil(t *testing.T) {
	test.AssertNotificationSent(t, nil, func() {
		// don't send anything
	})
}

func TestAssertNotificationSent_NotSent(t *testing.T) {
	tt := &testing.T{}

	test.AssertNotificationSent(tt, &fyne.Notification{}, func() {
		// don't send anything
	})
	assert.True(t, tt.Failed(), "notification assert should fail if no notification was sent")
}
