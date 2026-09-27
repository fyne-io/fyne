package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

// secretKeyPreference is the key used to keep the random encryption key in the app preferences.
const secretKeyPreference = "fyne.secret.key"

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

func newSecretPreferences(a *fyneApp) *secretPreferences {
	p := &secretPreferences{app: a, InMemoryPreferences: internal.NewInMemoryPreferences()}
	if a.uniqueID == "" && a.Metadata().ID == "" {
		return p
	}

	p.store = a.newSecretStore()
	p.load()
	p.AddChangeListener(func() {
		if p != a.secretPrefs {
			return
		}
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
	return p
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

func (p *secretPreferences) load() {
	if p.store == nil {
		return
	}
	data, err := p.store.load()
	if err == nil {
		err = p.loadFromData(data)
	}
	if err != nil && err != errEmptyPreferencesStore {
		fyne.LogError("Secret preferences load error:", err)
	}
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

// secretKey returns the random key material used by the encrypting stores.
// It is generated on first use and kept in the app preferences so that it survives an app ID change.
func (a *fyneApp) secretKey() ([]byte, error) {
	a.secretKeyLock.Lock()
	defer a.secretKeyLock.Unlock()

	stored := a.prefs.String(secretKeyPreference)
	if a.secretKeyCache != nil {
		if stored == "" { // preferences were replaced (e.g. cloud provider), make sure the key follows
			a.prefs.SetString(secretKeyPreference, base64.StdEncoding.EncodeToString(a.secretKeyCache))
		}
		return a.secretKeyCache, nil
	}

	if stored != "" {
		key, err := base64.StdEncoding.DecodeString(stored)
		if err == nil && len(key) == secretKeySize {
			a.secretKeyCache = key
			return key, nil
		}
		fyne.LogError("Stored secret preferences key is invalid, generating a new one", err)
	}

	key := make([]byte, secretKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	a.prefs.SetString(secretKeyPreference, base64.StdEncoding.EncodeToString(key))
	a.secretKeyCache = key
	return key, nil
}

// newEncryptedSecretStore is the portable secure store - AES-256-GCM over the platform's plain store,
// keyed by the random key that lives in the app preferences.
func (a *fyneApp) newEncryptedSecretStore() secretStore {
	return &encryptedStore{next: a.newPlainSecretStore(), key: a.secretKey}
}

const (
	secretKeySize       = 32 // AES-256
	secretFormatVersion = 1
)

var errSecretFormat = errors.New("secret preferences data is not in a recognised format")

// encryptedStore encrypts data with AES-GCM before handing it to the next store.
// The stored layout is: version byte, GCM nonce, ciphertext with authentication tag.
type encryptedStore struct {
	next secretStore
	key  func() ([]byte, error)
}

func (e *encryptedStore) load() ([]byte, error) {
	data, err := e.next.load()
	if err != nil {
		return nil, err
	}
	key, err := e.key()
	if err != nil {
		return nil, err
	}
	return decryptSecret(key, data)
}

func (e *encryptedStore) save(data []byte) error {
	key, err := e.key()
	if err != nil {
		return err
	}
	sealed, err := encryptSecret(key, data)
	if err != nil {
		return err
	}
	return e.next.save(sealed)
}

func encryptSecret(key, plain []byte) ([]byte, error) {
	gcm, err := newSecretCipher(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	out := make([]byte, 0, 1+len(nonce)+len(plain)+gcm.Overhead())
	out = append(out, secretFormatVersion)
	out = append(out, nonce...)
	return gcm.Seal(out, nonce, plain, out[:1]), nil
}

func decryptSecret(key, data []byte) ([]byte, error) {
	gcm, err := newSecretCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < 1+gcm.NonceSize() || data[0] != secretFormatVersion {
		return nil, errSecretFormat
	}
	nonce := data[1 : 1+gcm.NonceSize()]
	return gcm.Open(nil, nonce, data[1+gcm.NonceSize():], data[:1])
}

func newSecretCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
