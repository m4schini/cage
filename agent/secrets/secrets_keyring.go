package secrets

import (
	"github.com/99designs/keyring"
)

type Keyring99Backend struct {
	cfg keyring.Config
}

func (k *Keyring99Backend) Store(label, secret string) error {
	keyring.AvailableBackends()

	//TODO implement me
	panic("implement me")
}

func (k *Keyring99Backend) Retrieve(label string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Keyring99Backend) List() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (k *Keyring99Backend) Delete(label string) error {
	//TODO implement me
	panic("implement me")
}
