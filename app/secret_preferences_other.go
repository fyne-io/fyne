//go:build ci || wasm || test_web_driver || tinygo || (!darwin && !android && !windows) || (darwin && mobile)

package app

// newSecretStore returns the store used for secret preferences.
// This platform has no secure storage that we can use, so the values are encrypted before
// they are written using a key kept in the app preferences.
func (a *fyneApp) newSecretStore() secretStore {
	return a.newEncryptedSecretStore()
}
