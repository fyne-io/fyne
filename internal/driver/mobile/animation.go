package mobile

import "fyne.io/fyne/v2"

func (d *mobileDriver) StartAnimation(a *fyne.Animation) {
	a.Start()
}

func (d *mobileDriver) StartAnimationInternal(a *fyne.Animation) {
	d.animation.Start(a)
}

func (d *mobileDriver) StopAnimation(a *fyne.Animation) {
	a.Stop()
}
