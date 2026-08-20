//go:build !mobile

package test

func (*device) IsMobile() bool {
	return false
}
