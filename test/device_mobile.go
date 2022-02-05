//go:build mobile
// +build mobile

package test

func (d *device) IsMobile() bool {
	return true
}
