package state

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

var NixStorePath = filepath.Join(xdg.StateHome, "cage", "nix")

func InitNixStore() error {
	return os.MkdirAll(NixStorePath, 0750)
}
