// Package secrets manages multiple named API keys in a pluggable secret store.
//
// Backends are selected at runtime via the viper key "secrets.backend":
//
//	keyring (default) – via 99designs/keyring; transparently uses the macOS
//	                    Keychain, freedesktop secret-service, Windows
//	                    credential manager, KDE kwallet, or an encrypted-file
//	                    fallback depending on the host.
//	bitwarden         – the `bw` CLI, compatible with Bitwarden and Vaultwarden.
//
// Items are scoped to the service name "cage", so List() only surfaces items
// written by this package.
package secrets

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Sentinel errors.
var (
	// ErrNotFound is returned when no key is stored under the requested label.
	ErrNotFound = errors.New("api key not found")

	// ErrLabelEmpty is returned when the caller passes an empty label.
	ErrLabelEmpty = errors.New("label must not be empty")

	// ErrNotSupported is returned on platforms without a keychain backend.
	ErrNotSupported = errors.New("keychain is not supported on this platform")
)

// Backend is a pluggable secret store. Implementations must namespace items
// themselves so List() does not leak items written by other applications.
type Backend interface {
	Store(label, secret string) error
	Retrieve(label string) (string, error)
	List() ([]string, error)
	Delete(label string) error
}

// currentBackend resolves the active backend from viper. The default is the
// 99designs/keyring backend, which works across macOS, Linux, and Windows
// without per-platform code.
func currentBackend() Backend {
	switch viper.GetString("secrets.backend") {
	case "bitwarden":
		return BitwardenBackend{}
	case "keyring":
		return Keyring99Backend{}
	default:
		return Keyring99Backend{}
	}
}

// Store saves apiKey under label, replacing any previously stored value.
//
//	label  – a short identifier such as "work", "personal", or "ci-prod"
//	apiKey – the raw "sk-ant-…" value
func Store(label, apiKey string) error {
	if err := validateLabel(label); err != nil {
		return err
	}
	if apiKey == "" {
		return errors.New("apiKey must not be empty")
	}
	return currentBackend().Store(label, apiKey)
}

// Retrieve returns the API key stored under label.
// Returns ErrNotFound if no key exists for that label.
func Retrieve(label string) (string, error) {
	if err := validateLabel(label); err != nil {
		return "", err
	}
	return currentBackend().Retrieve(label)
}

// List returns the labels of all stored keys (never the key values).
// The slice is empty (not nil) when no keys have been stored yet.
func List() ([]string, error) {
	labels, err := currentBackend().List()
	if err != nil {
		return nil, err
	}
	if labels == nil {
		return []string{}, nil
	}
	return labels, nil
}

// Delete removes the key stored under label.
// Returns ErrNotFound if no such key exists.
func Delete(label string) error {
	if err := validateLabel(label); err != nil {
		return err
	}
	return currentBackend().Delete(label)
}

func validateLabel(label string) error {
	if label == "" {
		return fmt.Errorf("%w", ErrLabelEmpty)
	}
	return nil
}
