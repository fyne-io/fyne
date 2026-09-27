package app

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

// newTestSecretApp returns an app with in-memory preferences whose secret store encrypts into the given store.
func newTestSecretApp(id string, next secretStore) *fyneApp {
	a := &fyneApp{uniqueID: id}
	a.prefs = newPreferences(&fyneApp{}) // in-memory only, no ID so nothing is written to disk
	a.secretPrefs = newSecretPreferences(a)
	a.secretPrefs.store = &encryptedStore{next: next, key: a.secretKey}
	return a
}

func TestSecretPreferences_EncryptRoundTrip(t *testing.T) {
	key := make([]byte, secretKeySize)
	for i := range key {
		key[i] = byte(i)
	}

	sealed, err := encryptSecret(key, []byte(`{"token":"abc"}`))
	require.NoError(t, err)
	assert.NotContains(t, string(sealed), "token")
	assert.Equal(t, byte(secretFormatVersion), sealed[0])

	plain, err := decryptSecret(key, sealed)
	require.NoError(t, err)
	assert.Equal(t, `{"token":"abc"}`, string(plain))

	// each save uses a fresh nonce so identical content does not produce identical output
	again, err := encryptSecret(key, []byte(`{"token":"abc"}`))
	require.NoError(t, err)
	assert.NotEqual(t, sealed, again)
}

func TestSecretPreferences_DecryptRejectsTamperingAndWrongKey(t *testing.T) {
	key := make([]byte, secretKeySize)
	sealed, err := encryptSecret(key, []byte("secret"))
	require.NoError(t, err)

	tampered := append([]byte{}, sealed...)
	tampered[len(tampered)-1] ^= 0xff
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

func TestSecretPreferences_SaveAndLoad(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	a := newTestSecretApp("io.fyne.test.secret", blob)

	p := a.SecretPreferences()
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

	// a fresh app using the same preferences (and so the same key) and blob reads the values back
	b := &fyneApp{uniqueID: "io.fyne.test.secret", prefs: a.prefs}
	b.secretPrefs = &secretPreferences{app: b, InMemoryPreferences: internal.NewInMemoryPreferences(),
		store: &encryptedStore{next: blob, key: b.secretKey}}
	b.secretPrefs.load()

	loaded := b.SecretPreferences()
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
	a := newTestSecretApp("io.fyne.test.secret", blob)

	a.SecretPreferences().SetString("keep", "yes")
	a.SecretPreferences().SetString("drop", "no")
	a.SecretPreferences().RemoveValue("drop")
	a.secretPrefs.forceImmediateSave()

	b := newTestSecretApp("io.fyne.test.secret", blob)
	b.prefs = a.prefs
	b.secretPrefs.store = &encryptedStore{next: blob, key: b.secretKey}
	b.secretPrefs.load()
	assert.Equal(t, "yes", b.SecretPreferences().String("keep"))
	assert.Equal(t, "missing", b.SecretPreferences().StringWithFallback("drop", "missing"))
}

func TestSecretPreferences_WrongKeyLoadsEmpty(t *testing.T) {
	test.NewTempApp(t)
	blob := &memorySecretStore{}
	a := newTestSecretApp("io.fyne.test.secret", blob)
	a.SecretPreferences().SetString("token", "abc")
	a.secretPrefs.forceImmediateSave()

	// an app whose preferences hold a different key cannot read the data, but must not crash
	b := newTestSecretApp("io.fyne.test.secret", blob)
	b.secretPrefs.load()
	assert.Equal(t, "", b.SecretPreferences().String("token"))
}

func TestSecretPreferences_KeyStoredInPreferences(t *testing.T) {
	a := newTestSecretApp("io.fyne.test.secret", &memorySecretStore{})
	assert.Equal(t, "", a.prefs.String(secretKeyPreference), "no key should exist before first use")

	key, err := a.secretKey()
	require.NoError(t, err)
	assert.Len(t, key, secretKeySize)

	stored := a.prefs.String(secretKeyPreference)
	decoded, err := base64.StdEncoding.DecodeString(stored)
	require.NoError(t, err)
	assert.Equal(t, key, decoded, "the key is kept in the app preferences")

	same, err := a.secretKey()
	require.NoError(t, err)
	assert.Equal(t, key, same)

	// if the preferences are replaced (e.g. by a cloud provider) the key in use is written to the new store
	a.prefs = newPreferences(&fyneApp{})
	same, err = a.secretKey()
	require.NoError(t, err)
	assert.Equal(t, key, same)
	assert.Equal(t, stored, a.prefs.String(secretKeyPreference))

	// a new app loading the same preferences derives the same key
	b := &fyneApp{uniqueID: "io.fyne.test.secret", prefs: a.prefs}
	loaded, err := b.secretKey()
	require.NoError(t, err)
	assert.Equal(t, key, loaded)
}

func TestSecretPreferences_InvalidStoredKeyIsReplaced(t *testing.T) {
	a := newTestSecretApp("io.fyne.test.secret", &memorySecretStore{})
	a.prefs.SetString(secretKeyPreference, "not base64!")

	key, err := a.secretKey()
	require.NoError(t, err)
	assert.Len(t, key, secretKeySize)
	assert.NotEqual(t, "not base64!", a.prefs.String(secretKeyPreference))
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
	a := NewWithID("io.fyne.test")
	assert.NotNil(t, a.SecretPreferences())
	assert.NotSame(t, a.Preferences(), a.SecretPreferences())
}
