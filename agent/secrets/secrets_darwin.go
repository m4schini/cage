//go:build darwin || ios

// Package secrets manages multiple API keys in the system keychain.
// On macOS and iOS this uses the Security framework via go-keychain.
package secrets

import (
	"errors"
	"fmt"
	"strings"

	keychain "github.com/keybase/go-keychain"
)

const (
	service     = "cage-secrets"
	accessGroup = "" // empty = default keychain access group

	// accountPrefix is prepended to every label before it is written to the
	// keychain account field.  The "cage:" namespace prevents collisions with
	// any other application that might share the same service string.
	accountPrefix = "cage:"
)

func (k KeychainBackend) Store(label, apiKey string) error {
	// Delete any pre-existing item so we can do a clean add (no upsert in
	// the Security framework).
	if err := k.Delete(label); err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("replacing existing key: %w", err)
	}

	item := keychain.NewGenericPassword(
		service,
		account(label),
		label, // human-readable label shown in Keychain Access.app
		[]byte(apiKey),
		accessGroup,
	)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	if err := keychain.AddItem(item); err != nil {
		return fmt.Errorf("storing key %q: %w", label, err)
	}
	return nil
}

func (KeychainBackend) Retrieve(label string) (string, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(account(label))
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		if err == keychain.ErrorItemNotFound {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("querying key %q: %w", label, err)
	}
	if len(results) == 0 {
		return "", ErrNotFound
	}
	return string(results[0].Data), nil
}

func (KeychainBackend) List() ([]string, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetMatchLimit(keychain.MatchLimitAll)
	query.SetReturnAttributes(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		if err == keychain.ErrorItemNotFound {
			return nil, nil // service exists but has no items
		}
		return nil, fmt.Errorf("listing keys: %w", err)
	}

	labels := make([]string, 0, len(results))
	for _, r := range results {
		if lbl, ok := labelFromAccount(r.Account); ok {
			labels = append(labels, lbl)
		}
	}
	return labels, nil
}

func (KeychainBackend) Delete(label string) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(account(label))

	if err := keychain.DeleteItem(item); err != nil {
		if err == keychain.ErrorItemNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// account converts a user-facing label into the stable keychain account field.
func account(label string) string {
	return accountPrefix + strings.ToLower(strings.TrimSpace(label))
}

// labelFromAccount is the inverse of account.  It returns (label, true) when
// the account string carries our prefix, and ("", false) otherwise — allowing
// List() to skip items written by other apps under the same service name.
func labelFromAccount(acct string) (string, bool) {
	after, found := strings.CutPrefix(acct, accountPrefix)
	return after, found
}
