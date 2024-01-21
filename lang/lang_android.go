package lang

import (
	"fyne.io/fyne/v2/internal/driver/mobile/app"

	"github.com/fyne-io/go-locale"
)

func init() {
	locale.SetRunOnJVM(app.RunOnJVM)
}
