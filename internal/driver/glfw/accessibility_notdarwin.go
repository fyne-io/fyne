//go:build !accessibility || (!darwin && !windows)

package glfw

func (*window) updateAccessibility() {
}

func (*window) initAccessibilityForWindow() {
}

func (*window) cleanupAccessibilityForWindow() {
}
