// +build !ios,!android

package gomobile

func (*device) SystemScale() float32 {
	return 2
}
