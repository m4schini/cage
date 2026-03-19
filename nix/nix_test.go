package nix

import "testing"

func TestNewNixShell(t *testing.T) {
	data := ShellNixPackages{
		Shell:    "zsh",
		Packages: []string{"go", "git", "curl"},
	}
	t.Log(NewNixShellString(data))
}
