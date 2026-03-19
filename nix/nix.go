package nix

import (
	"bytes"
	"io"
)

func NewNixShell(packages ShellNixPackages, w io.Writer) error {
	return shellNixTemplate.Execute(w, packages)
}

func NewNixShellString(packages ShellNixPackages) (string, error) {
	var buf bytes.Buffer
	err := NewNixShell(packages, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
