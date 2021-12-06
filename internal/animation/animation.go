package animation

import (
	"sync/atomic"
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
	stopped     uint32 // atomic, 0 == false 1 == true
}

func newAnim(a *fyne.Animation) *anim {
	animate := &anim{a: a, start: time.Now(), end: time.Now().Add(a.Duration)}
	animate.total = animate.end.Sub(animate.start).Milliseconds()
	animate.repeatsLeft = a.RepeatCount
	return animate
}

func (a *anim) setStopped() {
	atomic.StoreUint32(&a.stopped, 1)
}

func (a *anim) isStopped() bool {
	return atomic.LoadUint32(&a.stopped) == 1
}
