//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package app

// TODO: #2734

func (s *settings) load() {
	s.setupTheme()
	s.schema.Scale = 1
}

func (s *settings) loadFromFile(path string) error {
	return nil
}

func watchFile(path string, callback func()) {
}

func (s *settings) watchSettings() {
}

func (s *settings) stopWatching() {
}
