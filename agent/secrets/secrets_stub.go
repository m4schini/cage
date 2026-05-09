//go:build !darwin && !ios && !linux

// Package secrets manages multiple named Anthropic API keys in the
// native secret store of the host operating system.
// This file is the fallback stub for platforms that go-keychain does not support.
package secrets

func store(label, apiKey string) error      { return ErrNotSupported }
func retrieve(label string) (string, error) { return "", ErrNotSupported }
func list() ([]string, error)               { return nil, ErrNotSupported }
func del(label string) error                { return ErrNotSupported }
