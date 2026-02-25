//go:build tamago || noos || tinygo

package app

func (s *settings) load() {
	s.schema.Scale = 1
}

func (s *settings) loadFromFile(_ string) error {
	// not supported
	return nil
}

func watchFile(_ string, _ func()) {
	// not supported
}

func (s *settings) watchSettings() {
	// not supported
}

func (s *settings) stopWatching() {
	// not supported
}
