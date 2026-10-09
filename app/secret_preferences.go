package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

// secretStore persists the encoded secret preferences.
// Implementations either hand the data to a secure operating system store, or protect it and
// then pass the result on to another secretStore that writes bytes unchanged (a file, or browser local storage).
type secretStore interface {
	// load returns the last saved data, or errEmptyPreferencesStore if nothing has been saved.
	load() ([]byte, error)
	// save replaces the stored data.
	save([]byte) error
}

type secretPreferences struct {
	*internal.InMemoryPreferences

	prefLock            sync.RWMutex
	savedRecently       bool
	changedDuringSaving bool

	saveLock sync.Mutex
	app      *fyneApp
	store    secretStore
}

// Declare conformity with Preferences interface
var _ fyne.Preferences = (*secretPreferences)(nil)

// newSecretPreferences returns the secret preferences loaded from the app's secret store.
// An error is returned if existing data could not be loaded, in which case nothing must be saved over it.
func newSecretPreferences(a *fyneApp, newStore func([]byte) secretStore, password []byte) (*secretPreferences, error) {
	p := &secretPreferences{app: a, InMemoryPreferences: internal.NewInMemoryPreferences()}
	if a.uniqueID == "" && a.Metadata().ID == "" {
		return p, nil
	}

	p.store = newStore(password)
	if err := p.load(); err != nil {
		return nil, err
	}
	p.AddChangeListener(func() {
		p.prefLock.Lock()
		shouldIgnoreChange := p.savedRecently
		if p.savedRecently {
			p.changedDuringSaving = true
		}
		p.prefLock.Unlock()

		if shouldIgnoreChange { // callback after loading from storage, or too many updates in a row
			return
		}

		if err := p.save(); err != nil {
			fyne.LogError("Failed on saving secret preferences", err)
		}
	})
	return p, nil
}

// forceImmediateSave writes secret preferences to storage immediately, ignoring the debouncing
// logic in the change listener. Does nothing if not backed with a persistent store.
func (p *secretPreferences) forceImmediateSave() {
	if p.store == nil {
		return
	}
	if err := p.save(); err != nil {
		fyne.LogError("Failed on force saving secret preferences", err)
	}
}

func (p *secretPreferences) load() error {
	if p.store == nil {
		return nil
	}
	data, err := p.store.load()
	if err == errEmptyPreferencesStore {
		return nil
	}
	if err != nil {
		return err
	}
	return p.loadFromData(data)
}

func (p *secretPreferences) loadFromData(data []byte) (err error) {
	p.WriteValues(func(values map[string]any) {
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&values)
		if err != nil {
			return
		}
		convertLists(values)
	})
	return err
}

func (p *secretPreferences) save() error {
	if p.store == nil {
		return nil
	}
	p.saveLock.Lock()
	defer p.saveLock.Unlock()

	p.prefLock.Lock()
	p.savedRecently = true
	p.prefLock.Unlock()
	defer p.resetSavedRecently()

	data, err := p.encode()
	if err != nil {
		return err
	}
	return p.store.save(data)
}

func (p *secretPreferences) encode() (data []byte, err error) {
	buf := &bytes.Buffer{}
	p.ReadValues(func(values map[string]any) {
		err = json.NewEncoder(buf).Encode(&values)
	})
	return buf.Bytes(), err
}

func (p *secretPreferences) resetSavedRecently() {
	go func() {
		time.Sleep(time.Millisecond * 100) // writes are not always atomic. 10ms worked, 100 is safer.

		fyne.DoAndWait(func() {
			p.prefLock.Lock()
			p.savedRecently = false
			changedDuringSaving := p.changedDuringSaving
			p.changedDuringSaving = false
			p.prefLock.Unlock()

			if changedDuringSaving {
				if err := p.save(); err != nil {
					fyne.LogError("failed on saving secret preferences", err)
				}
			}
		})
	}()
}

// newEncryptedSecretStore is the portable secure store - AES-256-GCM over the platform's plain store,
// keyed from the passed password.
func (a *fyneApp) newEncryptedSecretStore(password []byte) secretStore {
	return &encryptedStore{storage: a.newPlainSecretStore(), password: password}
}

const (
	secretKeySize       = 32 // AES-256
	secretSaltSize      = 16
	secretKeyIterations = 600000
	secretFormatVersion = 2
)

var errSecretFormat = errors.New("secret preferences data is not in a recognised format")

// encryptedStore encrypts data with AES-GCM before handing it to the storage store.
// The key is derived, using PBKDF2, from a password that the app provides.
// The stored layout is: version byte, key salt, GCM nonce, ciphertext with authentication tag.
type encryptedStore struct {
	storage  secretStore
	password []byte

	lock      sync.Mutex
	salt, key []byte
}

func (e *encryptedStore) load() ([]byte, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	data, err := e.storage.load()
	if err == errEmptyPreferencesStore { // prepare the key now so a missing password is reported on load
		if keyErr := e.newKey(); keyErr != nil {
			return nil, keyErr
		}
	}
	if err != nil {
		return nil, err
	}
	if len(data) < 1+secretSaltSize || data[0] != secretFormatVersion {
		return nil, errSecretFormat
	}

	key, err := e.keyForSalt(data[1 : 1+secretSaltSize])
	if err != nil {
		return nil, err
	}
	plain, err := decryptSecret(key, data)
	if err != nil && err != errSecretFormat {
		return nil, fyne.ErrPasswordIncorrect
	}
	return plain, err
}

// newKey derives a key for new data, which needs a new salt.
func (e *encryptedStore) newKey() error {
	salt := make([]byte, secretSaltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	_, err := e.keyForSalt(salt)
	return err
}

func (e *encryptedStore) save(data []byte) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.key == nil { // nothing was loaded
		if err := e.newKey(); err != nil {
			return err
		}
	}
	sealed, err := encryptSecret(e.key, e.salt, data)
	if err != nil {
		return err
	}
	return e.storage.save(sealed)
}

// keyForSalt returns the key for the given salt, deriving it from the password if that was not already done.
// The key is kept, rather than the password, which is not retained once used.
func (e *encryptedStore) keyForSalt(salt []byte) ([]byte, error) {
	if e.key != nil && bytes.Equal(salt, e.salt) {
		return e.key, nil
	}
	if len(e.password) == 0 {
		return nil, fyne.ErrPasswordRequired
	}

	e.salt = append([]byte{}, salt...)
	e.key = deriveSecretKey(e.password, e.salt)
	e.password = nil
	return e.key, nil
}

// deriveSecretKey is PBKDF2 (RFC 8018) with HMAC-SHA256, producing a single block which is the size of our key.
// We can move this to stdlib once upgraded to Go 1.24
func deriveSecretKey(password []byte, salt []byte) []byte {
	prf := hmac.New(sha256.New, password)
	prf.Write(salt)
	prf.Write([]byte{0, 0, 0, 1}) // block index
	u := prf.Sum(nil)

	key := append([]byte{}, u...)
	for i := 1; i < secretKeyIterations; i++ {
		prf.Reset()
		prf.Write(u)
		u = prf.Sum(u[:0])
		for j := range key {
			key[j] ^= u[j]
		}
	}
	return key
}

func encryptSecret(key, salt, plain []byte) ([]byte, error) {
	gcm, err := newSecretCipher(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	header := 1 + len(salt)
	out := make([]byte, 0, header+len(nonce)+len(plain)+gcm.Overhead())
	out = append(out, secretFormatVersion)
	out = append(out, salt...)
	out = append(out, nonce...)
	return gcm.Seal(out, nonce, plain, out[:header]), nil
}

func decryptSecret(key, data []byte) ([]byte, error) {
	gcm, err := newSecretCipher(key)
	if err != nil {
		return nil, err
	}
	header := 1 + secretSaltSize
	if len(data) < header+gcm.NonceSize() || data[0] != secretFormatVersion {
		return nil, errSecretFormat
	}
	nonce := data[header : header+gcm.NonceSize()]
	return gcm.Open(nil, nonce, data[header+gcm.NonceSize():], data[:header])
}

func newSecretCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
