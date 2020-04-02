// +build !ios,!android

package gomobile

import "fyne.io/fyne"

func (*device) SystemScale(_ fyne.Window) float32 {
	return 2
}
