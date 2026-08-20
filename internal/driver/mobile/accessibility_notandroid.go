//go:build !accessibility || (!android && !ios)

package mobile

// Stub implementations for platforms without accessibility bridges.

func (*window) updateAccessibility() {
}

func (*window) initAccessibilityForWindow() {
}

func (*window) cleanupAccessibilityForWindow() {
}
