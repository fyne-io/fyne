// +build hints

package internal

import (
	"log"
	"runtime"
)

// HintsEnabled returns true to indicate that hints are currently switched on.
func HintsEnabled() bool {
	return true
}

// LogHint reports a developer hint that should be followed to improve their app.
func LogHint(reason string) {
	log.Println("Fyne hint: ", reason)

	_, file, line, ok := runtime.Caller(2)
	if ok {
		log.Printf("  Created at: %s:%d", file, line)
	}
}
