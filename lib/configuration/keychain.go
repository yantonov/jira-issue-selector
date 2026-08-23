package configuration

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

// KeychainService is the service name used for every entry stored by this tool
// in the OS keychain (macOS Keychain, Windows Credential Manager,
// Linux Secret Service).
const KeychainService = "jira-issue-selector"

// Accounts used inside KeychainService, one entry per setting.
const (
	KeychainUserKey     = "user"
	KeychainHostNameKey = "hostname"
	KeychainApiKeyKey   = "apikey"
)

// KeychainKeys is ordered, the settings are asked and displayed in this order.
var KeychainKeys = []string{KeychainUserKey, KeychainHostNameKey, KeychainApiKeyKey}

type Keychain interface {
	// Get returns an empty string when the entry does not exist.
	Get(key string) (string, error)
	Set(key string, value string) error
	// Delete is a no-op when the entry does not exist.
	Delete(key string) error
}

type SystemKeychain struct{}

func (e SystemKeychain) Get(key string) (string, error) {
	value, err := keyring.Get(KeychainService, key)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", keychainError(err)
	}
	return value, nil
}

func (e SystemKeychain) Set(key string, value string) error {
	if err := keyring.Set(KeychainService, key, value); err != nil {
		return keychainError(err)
	}
	return nil
}

func (e SystemKeychain) Delete(key string) error {
	err := keyring.Delete(KeychainService, key)
	if err == nil || errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return keychainError(err)
}

func keychainError(err error) error {
	return fmt.Errorf("keychain (service %s) is not available: %w", KeychainService, err)
}
