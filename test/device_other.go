// +build !mobile

package test

func (d *device) IsMobile() bool {
	return false
}
