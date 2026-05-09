// Package secrets manages multiple named API keys in the
// native secret store of the host operating system.
//
// Supported platforms:
//
//	macOS  – Security.framework (Keychain)
//	Linux  – D-Bus SecretService (gnome-keyring / ksecretservice)
//
// On every other platform the four functions return ErrNotSupported.
//
// Keys are namespaced under the "cage:" prefix so that List() only surfaces
// items written by this package even if other applications share the same
// keychain service name.
package secrets

import (
	"errors"
	"fmt"
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

// Store saves apiKey in the system keychain under label, replacing any
// previously stored value.
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
	return store(label, apiKey)
}

// Retrieve returns the API key stored under label.
// Returns ErrNotFound if no key exists for that label.
func Retrieve(label string) (string, error) {
	if err := validateLabel(label); err != nil {
		return "", err
	}
	return retrieve(label)
}

// List returns the labels of all stored keys (never the key values).
// The slice is empty (not nil) when no keys have been stored yet.
func List() ([]string, error) {
	labels, err := list()
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
	return del(label)
}

func validateLabel(label string) error {
	if label == "" {
		return fmt.Errorf("%w", ErrLabelEmpty)
	}
	return nil
}
