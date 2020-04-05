// +build android

package gomobile

func (*device) SystemScale() float32 {
	if currentDPI >= 600 {
		return 4
	} else if currentDPI >= 405 {
		return 3
	} else if currentDPI >= 270 {
		return 2
	} else if currentDPI >= 180 {
		return 1.5
	}
	return 1
}
