//go:build android || ios || mobile
// +build android ios mobile

package app

func (s *settings) watchSettings() {
	// no-op on mobile
}

func (s *settings) stopWatching() {
	// no-op on mobile
}
