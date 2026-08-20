//go:build !android

package app

// StartForegroundService does nothing, as only Android can elevate a background process in this way.
func (a *fyneApp) StartForegroundService(_, _ string) {
}

// StopForegroundService does nothing, as only Android can elevate a background process in this way.
func (a *fyneApp) StopForegroundService() {
}

// RequestNotificationPermission does nothing, as only Android requires the user to grant permission before
// notifications can be posted.
func (a *fyneApp) RequestNotificationPermission() {
}
