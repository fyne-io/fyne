package test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"

	"github.com/stretchr/testify/assert"
)

func TestAssertNotificationSent(t *testing.T) {
	n := fyne.NewNotification("Test Title", "Some content")
	myApp := fyne.CurrentApp()

	AssertNotificationSent(t, n, func() {
		fyne.CurrentApp().SendNotification(n)
	})
	assert.Equal(t, myApp, fyne.CurrentApp())
}

func TestAssertNotificationSent_Nil(t *testing.T) {
	AssertNotificationSent(t, nil, func() {
		// don't send anything
	})
}

func TestAssertNotificationSent_NotSent(t *testing.T) {
	tt := &testing.T{}

	AssertNotificationSent(tt, &fyne.Notification{}, func() {
		// don't send anything
	})
	assert.True(t, tt.Failed(), "notification assert should fail if no notification was sent")
}

func TestAssertNotificationScheduled(t *testing.T) {
	n := fyne.NewNotification("Reminder", "Take a break")
	myApp := fyne.CurrentApp()

	var scheduled *fyne.ScheduledNotification
	AssertNotificationScheduled(t, n, func() {
		s, err := fyne.CurrentApp().ScheduleNotification(n, time.Now().Add(time.Minute))
		assert.NoError(t, err)
		scheduled = s
	})
	assert.NotNil(t, scheduled)
	assert.NotEmpty(t, scheduled.ID())
	assert.Equal(t, myApp, fyne.CurrentApp())
}

func TestAssertNotificationScheduled_CancelClearsState(t *testing.T) {
	n := fyne.NewNotification("Reminder", "Take a break")
	s, err := fyne.CurrentApp().ScheduleNotification(n, time.Now().Add(time.Minute))
	assert.NoError(t, err)

	err = fyne.CurrentApp().CancelScheduledNotification(s.ID())
	assert.NoError(t, err)
}
