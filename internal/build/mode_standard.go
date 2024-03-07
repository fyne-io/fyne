//go:build !debug && !release

package build

import "fyne.io/fyne/v2"

// Mode is the application's build mode.
const Mode = fyne.BuildStandard
