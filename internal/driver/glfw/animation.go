package glfw

import "fyne.io/fyne/v2"

func (d *gLDriver) StartAnimation(a *fyne.Animation) {
	a.Start()
}

func (d *gLDriver) StartAnimationInternal(a *fyne.Animation) {
	d.animation.Start(a)
}

func (d *gLDriver) StopAnimation(a *fyne.Animation) {
	a.Stop()
}
