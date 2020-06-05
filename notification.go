package fyne

// Notification represents a user notification that can be sent to the operating system.
type Notification struct {
	Title, Content string
}

// NewNotification creates a notification that can be passed to App.SendNotification.
func NewNotification(title, content string) *Notification {
	return &Notification{Title: title, Content: content}
}
