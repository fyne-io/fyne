package animation

import (
	"time"

	"fyne.io/fyne/v2"
)

type anim struct {
	a           *fyne.Animation
	end         time.Time
	repeatsLeft int
	reverse     bool
	start       time.Time
	total       int64
}

func newAnim(a *fyne.Animation) *anim {
	animate := &anim{a: a, start: time.Now(), end: time.Now().Add(a.Duration)}
	animate.total = animate.end.Sub(animate.start).Nanoseconds() / 1000000 // TODO change this to Milliseconds() when we drop Go 1.12
	animate.repeatsLeft = a.RepeatCount
	return animate
}
