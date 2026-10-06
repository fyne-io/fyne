package lang

import (
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/internal/driver/mobile/mobileinit"

	"github.com/jeandeaual/go-locale"
)

func initRuntime() {
	locale.SetRunOnJVM(app.RunOnJVM)
}

// platformLocaleReady reports whether we can call into the JVM yet.
func platformLocaleReady() bool {
	return mobileinit.HasContext()
}
