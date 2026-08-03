//go:build windows

// Animation control, delegated to the shared runner.

package directx

import (
	"fyne.io/fyne/v2"
)

func (d *dxDriver) StartAnimation(a *fyne.Animation) { d.animation.Start(a) }

func (d *dxDriver) StopAnimation(a *fyne.Animation) { d.animation.Stop(a) }

// DoubleTapDelay reports the system double-click time so taps match the rest of
// the desktop rather than a hard coded constant.
