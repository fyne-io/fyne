// +build android ios mobile nacl

package app

func (s *settings) watchSettings() {
	// no-op on mobile
}

func (s *settings) stopWatching() {
	// no-op on mobile
}
