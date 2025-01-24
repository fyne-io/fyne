package binding

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
)

func queueItem(f func()) {
	if async.IsMainGoroutine() {
		f()
		return
	}

	fyne.DoAndWait(f)
}
