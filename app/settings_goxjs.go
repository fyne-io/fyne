// +build js wasm web

package app

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
