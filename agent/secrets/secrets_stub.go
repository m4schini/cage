//go:build !darwin && !ios && !linux

// Package secrets manages multiple named Anthropic API keys in the
// native secret store of the host operating system.
// This file is the fallback stub for platforms that go-keychain does not support.
package secrets

func (KeychainBackend) Store(label, apiKey string) error      { return ErrNotSupported }
func (KeychainBackend) Retrieve(label string) (string, error) { return "", ErrNotSupported }
func (KeychainBackend) List() ([]string, error)               { return nil, ErrNotSupported }
func (KeychainBackend) Delete(label string) error             { return ErrNotSupported }
