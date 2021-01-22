package glfw

import "fyne.io/fyne/v2"

func (g *gLDriver) StartAnimation(a *fyne.Animation) {
	g.animation.Start(a)
}

func (g *gLDriver) StopAnimation(a *fyne.Animation) {
	g.animation.Stop(a)
}
