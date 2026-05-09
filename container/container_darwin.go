package container

import (
	"os"

	"golang.org/x/term"
)

func PrepareTTY() (cleanup func(), err error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return func() {}, err
	}
	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}, nil
}
