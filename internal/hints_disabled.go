//go:build !hints
// +build !hints

package internal

// HintsEnabled is false to indicate that hints are not currently switched on.
// To enable please rebuild with "-tags hints" parameters.
const HintsEnabled = false

// LogHint reports a developer hint that should be followed to improve their app.
// This does nothing unless the "hints" build flag is used.
func LogHint(reason string) {
	// no-op when hints not enabled
}
