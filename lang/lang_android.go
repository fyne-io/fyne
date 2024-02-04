package lang

import (
	"fyne.io/fyne/v2/internal/driver/mobile/app"

	"github.com/jeandeaual/go-locale"
)

func init() {
	locale.SetRunOnJVM(app.RunOnJVM)
}
