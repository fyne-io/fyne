package fyne

import "time"

// Notification represents a user notification that can be sent to the operating system.
type Notification struct {
	Title, Content string
}

// NewNotification creates a notification that can be passed to [App.SendNotification].
func NewNotification(title, content string) *Notification {
	return &Notification{Title: title, Content: content}
}

// ScheduledNotification represents a notification that has been queued for delivery at a
// future time. Instances are returned from [App.ScheduleNotification] and can be cancelled
// using the ID with [App.CancelScheduledNotification].
//
// Since: 2.8
type ScheduledNotification struct {
	*Notification

	// DeliveryTime is the time at which the notification is scheduled to be delivered.
	DeliveryTime time.Time

	id string
}

// ID returns the unique identifier for this scheduled notification.
// Pass this value to [App.CancelScheduledNotification] to cancel a pending delivery.
func (s *ScheduledNotification) ID() string {
	return s.id
}

// Cancel will remove this scheduled notification from future posting.
// If your application might want to cancel a future notification after it has been restarted
// you should persist the `ID` value and then use `CancelScheduledNotification`.
func (s *ScheduledNotification) Cancel() error {
	return CurrentApp().CancelScheduledNotification(s.id)
}

// NewScheduledNotification builds a scheduled notification record with a known ID.
// Application code should not normally call this directly - prefer [App.ScheduleNotification]
// which assigns an ID and arranges delivery.
//
// This is exposed for use by app implementations and advanced testing.
//
// Since: 2.8
func NewScheduledNotification(id string, n *Notification, deliverAt time.Time) *ScheduledNotification {
	return &ScheduledNotification{
		Notification: n,
		DeliveryTime: deliverAt,
		id:           id,
	}
}
