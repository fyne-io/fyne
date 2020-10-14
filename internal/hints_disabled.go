// +build !hints

package internal

// HintsEnabled returns false to indicate that hints are not currently switched on.
// To enable please rebuild with "-tags hints" parameters.
func HintsEnabled() bool {
	return false
}

// LogHint reports a developer hint that should be followed to improve their app.
// This does nothing unless the "hints" build flag is used.
func LogHint(reason string) {
	// no-op when hints not enabled
}
