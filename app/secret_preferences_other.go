//go:build ci || wasm || test_web_driver || tinygo || (!darwin && !android && !windows) || (darwin && mobile)

package app

// newSecretStore returns the store used for secret preferences.
// This platform has no secure storage that we can use, so the values are encrypted before
// they are written using a key derived from the password that the app provides.
func (a *fyneApp) newSecretStore(password []byte) secretStore {
	return a.newEncryptedSecretStore(password)
}
