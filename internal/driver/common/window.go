package common

import (
	"fyne.io/fyne/v2/internal/async"
)

var DonePool = async.Pool[chan struct{}]{
	New: func() chan struct{} {
		return make(chan struct{})
	},
}
