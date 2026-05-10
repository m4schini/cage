package secrets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/99designs/keyring"
)

const keyringServiceName = "cage-secrets"
const labelPrefix = "cage:"

// Keyring99Backend stores secrets via 99designs/keyring, which abstracts the
// macOS Keychain, freedesktop secret-service, Windows credential manager, KDE
// kwallet, and an encrypted-file fallback behind a single API.
type Keyring99Backend struct{}

func (Keyring99Backend) open() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: keyringServiceName,
	})
}

func (k Keyring99Backend) Store(label, secret string) error {
	ring, err := k.open()
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:   fmt.Sprintf("%v%v", labelPrefix, label),
		Data:  []byte(secret),
		Label: label,
	})
}

func (k Keyring99Backend) Retrieve(label string) (string, error) {
	ring, err := k.open()
	if err != nil {
		return "", err
	}
	item, err := ring.Get(fmt.Sprintf("%v%v", labelPrefix, label))
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func (k Keyring99Backend) List() ([]string, error) {
	ring, err := k.open()
	if err != nil {
		return nil, err
	}
	list, err := ring.Keys()
	if err != nil {
		return nil, err
	}
	out := make([]string, len(list))
	for i, s := range list {
		out[i] = strings.TrimPrefix(s, labelPrefix)
	}
	return out, nil
}

func (k Keyring99Backend) Delete(label string) error {
	ring, err := k.open()
	if err != nil {
		return err
	}
	if err := ring.Remove(fmt.Sprintf("%v%v", labelPrefix, label)); err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
