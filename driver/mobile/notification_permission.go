package mobile

// NotificationPermission is implemented by a [fyne.io/fyne/v2.App] on
// platforms where the user has to grant permission before notifications are
// shown. Currently only Android apps implement it.
//
// Since: 2.9
type NotificationPermission interface {
	// RequestNotificationPermission asks the user to allow this app to post
	// notifications, which Android 13+ requires before any notification is
	// shown. The system dialog is skipped when the permission is already
	// granted or the device runs an older Android version.
	//
	// The request is only issued when the app targets SDK 33 or later, which
	// for Fyne means a release build. Apps with a lower target SDK cannot
	// request this permission themselves; Android prompts on its own when the
	// app first creates a notification channel, so this call does nothing
	// for them.
	RequestNotificationPermission()
}
