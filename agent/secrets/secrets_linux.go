//go:build linux

// Package secrets manages multiple Anthropic API keys in the system keychain.
// On Linux this communicates with a D-Bus SecretService provider such as
// gnome-keyring or KSecretService.
package secrets

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/keybase/go-keychain/secretservice"
)

const (
	// collection is the D-Bus path for the default secrets collection.
	collection = secretservice.DefaultCollection

	// accountPrefix namespaces every label stored in the attributes map so
	// that List() can distinguish our items from other apps' items.
	accountPrefix = "cage:"
)

// attrService is the D-Bus attribute key used to identify our service.
const attrService = "service"
const attrAccount = "account"
const serviceName = "cage-secrets"

func (KeychainBackend) Store(label, apiKey string) error {
	srv, err := secretservice.NewService()
	if err != nil {
		return fmt.Errorf("opening SecretService: %w", err)
	}

	session, err := srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		return fmt.Errorf("opening session: %w", err)
	}

	secret, err := session.NewSecret([]byte(apiKey))
	if err != nil {
		return fmt.Errorf("creating secret: %w", err)
	}

	attrs := attributes(label)
	props := secretservice.NewSecretProperties(label, attrs)

	// Replace = true performs an upsert: updates the secret if the item
	// already exists, or creates it if not.
	_, err = srv.CreateItem(collection, props, secret, secretservice.ReplaceBehaviorReplace)
	if err != nil {
		return fmt.Errorf("storing key %q: %w", label, err)
	}
	return nil
}

func (KeychainBackend) Retrieve(label string) (string, error) {
	srv, err := secretservice.NewService()
	if err != nil {
		return "", fmt.Errorf("opening SecretService: %w", err)
	}

	session, err := srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		return "", fmt.Errorf("opening session: %w", err)
	}

	items, err := srv.SearchCollection(collection, attributes(label))
	if err != nil {
		return "", fmt.Errorf("searching key %q: %w", label, err)
	}
	if len(items) == 0 {
		return "", ErrNotFound
	}

	secret, err := srv.GetSecret(items[0], *session)
	if err != nil {
		return "", fmt.Errorf("getting secret %q: %w", label, err)
	}
	return string(secret), nil
}

func (KeychainBackend) List() ([]string, error) {
	srv, err := secretservice.NewService()
	if err != nil {
		return nil, fmt.Errorf("opening SecretService: %w", err)
	}

	// Search with only the service attribute to find all our items.
	items, err := srv.SearchCollection(collection, map[string]string{
		attrService: serviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("listing keys: %w", err)
	}

	labels := make([]string, 0, len(items))
	for _, item := range items {
		attrs, err := srv.GetAttributes(item)
		if err != nil {
			continue
		}
		acct, ok := attrs[attrAccount]
		if !ok {
			continue
		}
		if lbl, ok := labelFromAccount(acct); ok {
			labels = append(labels, lbl)
		}
	}
	return labels, nil
}

func (KeychainBackend) Delete(label string) error {
	srv, err := secretservice.NewService()
	if err != nil {
		return fmt.Errorf("opening SecretService: %w", err)
	}

	items, err := srv.SearchCollection(collection, attributes(label))
	if err != nil {
		return fmt.Errorf("searching key %q: %w", label, err)
	}
	if len(items) == 0 {
		return ErrNotFound
	}

	for _, item := range items {
		if err := srv.DeleteItem(item); err != nil {
			return fmt.Errorf("deleting key %q: %w", label, err)
		}
	}
	return nil
}

// attributes builds the D-Bus attribute map for a given label.
func attributes(label string) map[string]string {
	return map[string]string{
		attrService: serviceName,
		attrAccount: accountPrefix + strings.ToLower(strings.TrimSpace(label)),
	}
}

// labelFromAccount is the inverse of the account field encoding.
func labelFromAccount(acct string) (string, bool) {
	after, found := strings.CutPrefix(acct, accountPrefix)
	return after, found
}

// Silence "imported and not used" if dbus ends up indirect only.
var _ dbus.ObjectPath
