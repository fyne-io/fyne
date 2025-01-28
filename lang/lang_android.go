package lang

import (
	"fyne.io/fyne/v2/internal/driver/mobile/app"

	"github.com/jeandeaual/go-locale"
)

func initRuntime() {
	locale.SetRunOnJVM(app.RunOnJVM)
}
