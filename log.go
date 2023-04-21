package fyne

import (
	"log"
	"runtime"
)

// LogError reports an error to the command line with the specified err cause,
// if not nil.
// The function also reports basic information about the code location.
func LogError(reason string, err error) {
	log.Println("Fyne error: ", reason)
	if err != nil {
		log.Println("  Cause:", err)
	}

	for i := 1; i < 3; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			log.Printf("  At: %s:%d", file, line)
		}
	}
}
