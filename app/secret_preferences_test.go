package app

import (
	"encoding/hex"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/test"
)

// memorySecretStore keeps the protected bytes in memory so tests never touch the user's storage.
type memorySecretStore struct {
	data []byte
}

func (m *memorySecretStore) load() ([]byte, error) {
	if m.data == nil {
		return nil, errEmptyPreferencesStore
	}
	return m.data, nil
}

func (m *memorySecretStore) save(data []byte) error {
	m.data = append([]byte{}, data...)
	return nil
}

// newTestSecretApp returns an app whose secret store encrypts into the given store using the password.
// The error is that of loading any existing data from the store.
func newTestSecretApp(id string, store secretStore, password []byte) (*fyneApp, error) {
	a := &fyneApp{uniqueID: id}
	a.prefs = newPreferences(&fyneApp{}) // in-memory only, no ID so nothing is written to disk
	a.secretPrefs = &secretPreferences{
		app: a, InMemoryPreferences: internal.NewInMemoryPreferences(),
		store: &encryptedStore{storage: store, password: password},
	}
	return a, a.secretPrefs.load()
}

func fixedPassword(password string) []byte {
	return []byte(password)
}

func TestSecretPreferences_EncryptRoundTrip(t *testing.T) {
	key := make([]byte, secretKeySize)
	for i := range key {
		key[i] = byte(i)
	}
	salt := make([]byte, secretSaltSize)

	sealed, err := encryptSecret(key, salt, []byte(`{"token":"abc"}`))
	require.NoError(t, err)
	assert.NotContains(t, string(sealed), "token")
	assert.Equal(t, byte(secretFormatVersion), sealed[0])
	assert.Equal(t, salt, sealed[1:1+secretSaltSize])

	plain, err := decryptSecret(key, sealed)
	require.NoError(t, err)
	assert.Equal(t, `{"token":"abc"}`, string(plain))

	// each save uses a fresh nonce so identical content does not produce identical output
	again, err := encryptSecret(key, salt, []byte(`{"token":"abc"}`))
	require.NoError(t, err)
	assert.NotEqual(t, sealed, again)
}

func TestSecretPreferences_DecryptRejectsTamperingAndWrongKey(t *testing.T) {
	key := make([]byte, secretKeySize)
	sealed, err := encryptSecret(key, make([]byte, secretSaltSize), []byte("secret"))
	require.NoError(t, err)

	tampered := append([]byte{}, sealed...)
	tampered[len(tampered)-1] ^= 0xff
	_, err = decryptSecret(key, tampered)
	assert.Error(t, err)

	tampered = append([]byte{}, sealed...)
	tampered[1] ^= 0xff // the salt is authenticated too
	_, err = decryptSecret(key, tampered)
	assert.Error(t, err)

	other := make([]byte, secretKeySize)
	other[0] = 1
	_, err = decryptSecret(other, sealed)
	assert.Error(t, err)

	_, err = decryptSecret(key, []byte{secretFormatVersion + 1, 0, 0})
	assert.ErrorIs(t, err, errSecretFormat)
	_, err = decryptSecret(key, nil)
	assert.ErrorIs(t, err, errSecretFormat)
}

func TestSecretPreferences_DeriveKey(t *testing.T) {
	salt := []byte("0123456789abcdef")
	key := deriveSecretKey([]byte("password"), salt)
	assert.Len(t, key, secretKeySize)
	// PBKDF2-HMAC-SHA256, 600000 iterations, verified against Go's crypto/pbkdf2
	assert.Equal(t, "996d7c90f74a4a169cf7adef42b06848f7d1bac3e568d1cc94d4f79be1ee0263", hex.EncodeToString(key))

	assert.NotEqual(t, key, deriveSecretKey([]byte("other"), salt))
	assert.NotEqual(t, key, deriveSecretKey([]byte("password"), []byte("fedcba9876543210")))
}

func TestSecretPreferences_SaveAndLoad(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	a, err := newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("pass"))
	require.NoError(t, err)

	p := a.secretPrefs
	p.SetString("keyString", "value")
	p.SetStringList("keyStringList", []string{"1", "2", "3"})
	p.SetInt("keyInt", 4)
	p.SetIntList("keyIntList", []int{1, 2, 3})
	p.SetFloat("keyFloat", 3.5)
	p.SetFloatList("keyFloatList", []float64{1.1, 2.2, 3.3})
	p.SetBool("keyBool", true)
	p.SetBoolList("keyBoolList", []bool{true, false, true})
	a.secretPrefs.forceImmediateSave()

	require.NotNil(t, blob.data)
	assert.NotContains(t, string(blob.data), "value", "secret values must not be stored in plain text")
	assert.NotContains(t, string(blob.data), "keyString")
	assert.Equal(t, "", a.prefs.String("fyne.secret.key"), "no key material may be kept in the app preferences")

	// a fresh app given the same password and blob reads the values back
	b, err := newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("pass"))
	require.NoError(t, err)
	loaded := b.secretPrefs
	assert.Equal(t, "value", loaded.String("keyString"))
	assert.Equal(t, []string{"1", "2", "3"}, loaded.StringList("keyStringList"))
	assert.Equal(t, 4, loaded.Int("keyInt"))
	assert.Equal(t, []int{1, 2, 3}, loaded.IntList("keyIntList"))
	assert.Equal(t, 3.5, loaded.Float("keyFloat"))
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, loaded.FloatList("keyFloatList"))
	assert.True(t, loaded.Bool("keyBool"))
	assert.Equal(t, []bool{true, false, true}, loaded.BoolList("keyBoolList"))
}

func TestSecretPreferences_Remove(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	a, err := newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("pass"))
	require.NoError(t, err)

	a.secretPrefs.SetString("keep", "yes")
	a.secretPrefs.SetString("drop", "no")
	a.secretPrefs.RemoveValue("drop")
	a.secretPrefs.forceImmediateSave()

	b, err := newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("pass"))
	require.NoError(t, err)
	assert.Equal(t, "yes", b.secretPrefs.String("keep"))
	assert.Equal(t, "missing", b.secretPrefs.StringWithFallback("drop", "missing"))
}

func TestSecretPreferences_WrongPasswordIsAnError(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	a, err := newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("pass"))
	require.NoError(t, err)
	a.secretPrefs.SetString("token", "abc")
	a.secretPrefs.forceImmediateSave()

	_, err = newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("wrong"))
	assert.ErrorIs(t, err, fyne.ErrPasswordIncorrect)
}

func TestSecretPreferences_PasswordNotRetained(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	password := fixedPassword("pass")

	a, err := newTestSecretApp("io.fyne.test.secret", blob, password)
	require.NoError(t, err)
	assert.Equal(t, []byte("pass"), password, "the caller's password is not modified")
	assert.Nil(t, a.secretPrefs.store.(*encryptedStore).password)

	// the derived key is kept, so saving does not need the password again
	a.secretPrefs.SetString("token", "abc")
	a.secretPrefs.forceImmediateSave()
	assert.NotNil(t, blob.data)
}

func TestSecretPreferences_MissingPasswordIsAnError(t *testing.T) {
	test.NewTempApp(t)
	for name, password := range map[string][]byte{"nil": nil, "empty": {}} {
		t.Run(name, func(t *testing.T) {
			blob := &memorySecretStore{}
			_, err := newTestSecretApp("io.fyne.test.secret", blob, password)
			assert.ErrorIs(t, err, fyne.ErrPasswordRequired)

			store := &encryptedStore{storage: blob, password: password}
			assert.ErrorIs(t, store.save([]byte("secret")), fyne.ErrPasswordRequired)
			assert.Nil(t, blob.data)
		})
	}
}

func TestFileSecretStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "preferences.secret")
	store := &fileSecretStore{path: path}

	_, err := store.load()
	assert.ErrorIs(t, err, errEmptyPreferencesStore)

	require.NoError(t, store.save([]byte{1, 2, 3}))
	data, err := store.load()
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, data)

	require.NoError(t, store.save([]byte{9}))
	data, err = store.load()
	require.NoError(t, err)
	assert.Equal(t, []byte{9}, data)
	assert.NoFileExists(t, path+".tmp")
}

func TestFyneApp_SecretPreferences(t *testing.T) {
	test.NewTempApp(t) // restore the test app when done, so later tests are not left with a real app
	a := &fyneApp{uniqueID: "io.fyne.test.secret"}
	a.prefs = newPreferences(&fyneApp{})

	newStore := func(password []byte) secretStore {
		return &encryptedStore{storage: &memorySecretStore{}, password: password}
	}

	secret, err := a.secretPreferencesFrom(newStore, fixedPassword("pass"))
	require.NoError(t, err)
	assert.NotNil(t, secret)
	assert.NotSame(t, a.Preferences(), secret)

	again, err := a.SecretPreferences(nil)
	require.NoError(t, err)
	assert.Same(t, secret, again, "later calls return the same store")
}

func TestFyneApp_SecretPreferences_ErrorLeavesDataAndCanRetry(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	a, err := newTestSecretApp("io.fyne.test.secret", blob, fixedPassword("pass"))
	require.NoError(t, err)
	a.secretPrefs.SetString("token", "abc")
	a.secretPrefs.forceImmediateSave()
	saved := append([]byte{}, blob.data...)

	b := &fyneApp{uniqueID: "io.fyne.test.secret"}
	b.prefs = newPreferences(&fyneApp{})
	newStore := func(password []byte) secretStore {
		return &encryptedStore{storage: blob, password: password}
	}

	secret, err := b.secretPreferencesFrom(newStore, nil)
	assert.ErrorIs(t, err, fyne.ErrPasswordRequired)
	assert.Nil(t, secret)

	secret, err = b.secretPreferencesFrom(newStore, fixedPassword("wrong"))
	assert.ErrorIs(t, err, fyne.ErrPasswordIncorrect)
	assert.Nil(t, secret)
	assert.Equal(t, saved, blob.data, "a failed load must not change the stored data")

	secret, err = b.secretPreferencesFrom(newStore, fixedPassword("pass"))
	require.NoError(t, err)
	assert.Equal(t, "abc", secret.String("token"))
}
