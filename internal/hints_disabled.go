// +build !hints

package internal

// LogHint reports a developer hint that should be followed to improve their app.
// This does nothing unless the "hints" build flag is used.
func LogHint(reason string) {
	// no-op when hints not enabled
}
